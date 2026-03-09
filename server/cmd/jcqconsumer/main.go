package main

import (
	"apollo/server/cmd/payment/api/common"
	"apollo/server/cmd/payment/pkg/types"
	"apollo/server/internal/model"
	"apollo/server/internal/repository"
	"apollo/server/pkg/config"
	"apollo/server/pkg/contextx"
	"apollo/server/pkg/data"
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	llog "github.com/labstack/gommon/log"
)

type jcqConfig struct {
	AccessKey       string
	SecretKey       string
	Topic           string
	ConsumerGroupID string
	HTTPURL         string
	AutoAck         bool
	Size            int
	Interval        time.Duration
}

func readConfig() jcqConfig {
	conf := config.New("configs/config.yaml")
	sec := conf.JCQConfig
	size := sec.PollSize
	if size <= 0 {
		size = 5
	}
	interval := time.Duration(sec.PollIntervalSeconds) * time.Second
	if interval <= 0 {
		interval = time.Second
	}
	return jcqConfig{
		AccessKey:       sec.AccessKey,
		SecretKey:       sec.SecretKey,
		Topic:           sec.Topic,
		ConsumerGroupID: sec.ConsumerGroupID,
		HTTPURL:         sec.HTTPURL,
		AutoAck:         sec.AutoAck,
		Size:            size,
		Interval:        interval,
	}
}

func getHeaders(accessKey string) map[string]string {
	return map[string]string{
		"Content-Type": "application/json",
		"accessKey":    accessKey,
		"dateTime":     time.Now().UTC().Format(time.RFC3339),
	}
}

func getSignSource(headers map[string]string, params map[string]string) string {
	type kv struct{ k, v string }
	all := make([]kv, 0, len(headers)+len(params))
	for k, v := range headers {
		all = append(all, kv{k: k, v: v})
	}
	for k, v := range params {
		all = append(all, kv{k: k, v: v})
	}
	for i := 0; i < len(all)-1; i++ {
		for j := i + 1; j < len(all); j++ {
			if all[i].k > all[j].k {
				all[i], all[j] = all[j], all[i]
			}
		}
	}
	var b strings.Builder
	for i, p := range all {
		if i > 0 {
			b.WriteString("&")
		}
		b.WriteString(p.k)
		b.WriteString("=")
		b.WriteString(p.v)
	}
	return b.String()
}

func sign(source, secret string) string {
	mac := hmac.New(sha1.New, []byte(secret))
	_, _ = mac.Write([]byte(source))
	return strings.TrimSpace(base64.StdEncoding.EncodeToString(mac.Sum(nil)))
}

func getMessages(conf jcqConfig) (string, error) {
	headers := getHeaders(conf.AccessKey)
	params := map[string]string{
		"topic":           conf.Topic,
		"consumerGroupId": conf.ConsumerGroupID,
		"size":            strconv.Itoa(conf.Size),
		"ack":             strconv.FormatBool(conf.AutoAck),
	}
	source := getSignSource(headers, params)
	headers["signature"] = sign(source, conf.SecretKey)
	req, err := http.NewRequest("GET", conf.HTTPURL+"/v2/messages", nil)
	if err != nil {
		return "", err
	}
	q := req.URL.Query()
	for k, v := range params {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func ackMessage(conf jcqConfig, ackIndex string, action string) error {
	headers := getHeaders(conf.AccessKey)
	body := map[string]any{
		"topic":           conf.Topic,
		"consumerGroupId": conf.ConsumerGroupID,
		"ackAction":       action,
		"ackIndex":        ackIndex,
	}
	jsonBody, _ := json.Marshal(body)
	params := map[string]string{
		"topic":           conf.Topic,
		"consumerGroupId": conf.ConsumerGroupID,
		"ackAction":       action,
		"ackIndex":        ackIndex,
	}
	source := getSignSource(headers, params)
	headers["signature"] = sign(source, conf.SecretKey)
	req, err := http.NewRequest("POST", conf.HTTPURL+"/v2/ack", strings.NewReader(string(jsonBody)))
	if err != nil {
		return err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, _ = io.ReadAll(resp.Body)
	return nil
}

func parseJDMessageBody(body string) (*types.JDJsonData, error) {
	var jd types.JDJsonData
	if err := json.Unmarshal([]byte(body), &jd); err != nil {
		return nil, err
	}
	return &jd, nil
}

func toString(i any) string {
	switch v := i.(type) {
	case string:
		return v
	case float64:
		return strconv.FormatInt(int64(v), 10)
	default:
		return ""
	}
}

func handleMessage(raw any) error {
	msg, ok := raw.(map[string]any)
	if !ok {
		return nil
	}
	bodyStr := toString(msg["body"])
	if len(bodyStr) == 0 {
		return nil
	}
	jd, err := parseJDMessageBody(bodyStr)
	if err != nil {
		return nil
	}
	db := data.Instance()
	skuId := strconv.FormatInt(jd.SkuId, 10)
	kw := strconv.FormatInt(jd.OrderId, 10)
	if len(jd.GameAccount) > 0 {
		kw = jd.GameAccount
	}
	o, err := repository.Order.GetBySkuIdDarkNumber(nil, db, skuId, kw)
	if err != nil || o == nil {
		return nil
	}
	if o.Status == model.OrderStatusFinish {
		return nil
	}
	var payTime time.Time
	if len(jd.CreateTime) > 0 {
		t, err := time.Parse("2006-01-02T15:04:05", jd.CreateTime)
		if err == nil {
			payTime = t
		}
	}
	params := common.UpdateOrderParams{
		PartnerOrderId: strconv.FormatInt(jd.OrderId, 10),
		OrderId:        o.OrderId,
		PayAccount:     jd.Pin,
		PayTime:        payTime,
		ReceivedAmount: jd.TotalPrice,
		ShopName:       o.Shop,
	}
	_ = common.UpdateOrder(nil, db, params)
	_ = repository.Partner.UpdateBalance(nil, db, o.PartnerId, o.OrderId, -params.ReceivedAmount, model.BalanceFromTypeOrderDeduct)
	updated, err := repository.Order.GetByOrderId(db, o.OrderId)
	if err != nil || updated == nil {
		return nil
	}
	logger := llog.New("jcq-consumer")
	_ = common.NotifyMerchant(contextx.NewContext(logger), db, updated)
	return nil
}

func run(ctx context.Context, conf jcqConfig) {
	t := time.NewTicker(conf.Interval)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			respStr, err := getMessages(conf)
			if err != nil {
				continue
			}
			var result map[string]any
			if err := json.Unmarshal([]byte(respStr), &result); err != nil {
				continue
			}
			r, ok := result["result"].(map[string]any)
			if !ok || r == nil {
				continue
			}
			msgs, _ := r["messages"].([]any)
			if len(msgs) == 0 {
				continue
			}
			ackIndex, _ := r["ackIndex"].(string)
			for _, m := range msgs {
				_ = handleMessage(m)
			}
			if !conf.AutoAck && len(ackIndex) > 0 {
				_ = ackMessage(conf, ackIndex, "SUCCESS")
			}
		}
	}
}

func main() {
	conf := readConfig()
	srv := config.Get()
	if srv == nil {
		srv = config.New("configs/config.yaml")
	}
	_ = data.Init(srv.MysqlConfig)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	go run(ctx, conf)
	<-ctx.Done()
}
