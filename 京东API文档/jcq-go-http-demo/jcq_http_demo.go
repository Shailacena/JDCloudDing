package main

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"
)

func md5f(content string) string {
	h := md5.New()
	h.Write([]byte(content))
	return hex.EncodeToString(h.Sum(nil))
}

func message2str(message map[string]interface{}) string {
	m := make(map[string]interface{}, len(message))
	for k, v := range message {
		m[k] = v
	}
	if props, ok := m["properties"]; ok {
		for k, v := range props.(map[string]interface{}) {
			m[k] = v
		}
		delete(m, "properties")
	}
	sortedKeys := make([]string, 0, len(m))
	for k := range m {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Strings(sortedKeys)
	var ms strings.Builder
	for _, k := range sortedKeys {
		ms.WriteString(k)
		ms.WriteString("=")
		ms.WriteString(fmt.Sprintf("%v", m[k]))
		ms.WriteString("&")
	}
	return md5f(ms.String()[:len(ms.String())-1])
}

func getSignSource(headers map[string]string, params map[string]interface{}) string {
	d := make(map[string]interface{}, 2+len(params))
	d["accessKey"] = headers["accessKey"]
	d["dateTime"] = headers["dateTime"]
	for k, v := range params {
		d[k] = v
	}
	if messages, ok := d["messages"]; ok && reflect.TypeOf(messages).Kind() == reflect.Slice {
		var messageStrings []string
		if messages != nil {
			if m, ok := messages.([]map[string]interface{}); ok {
				for _, message := range m {
					messageStr := message2str(message)
					messageStrings = append(messageStrings, messageStr)
				}
			}
		}
		d["messages"] = strings.Join(messageStrings, ",")
	}
	var signSource strings.Builder
	sortedKeys := make([]string, 0, len(d))
	for k := range d {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Strings(sortedKeys)
	for _, k := range sortedKeys {
		signSource.WriteString(k)
		signSource.WriteString("=")
		signSource.WriteString(fmt.Sprintf("%v", d[k]))
		signSource.WriteString("&")
	}
	return signSource.String()[:len(signSource.String())-1]
}

func getSignature(source string, secretKey string) (string, error) {
	key := []byte(secretKey)
	sourceBytes := []byte(source)
	mac := hmac.New(sha1.New, key)
	_, err := mac.Write(sourceBytes)
	if err != nil {
		return "", err
	}
	signature := mac.Sum(nil)
	encodedSignature := base64.StdEncoding.EncodeToString(signature)
	return strings.TrimSpace(encodedSignature), nil
}

func getHeaders(accessKey string) map[string]string {
	return map[string]string{
		"Content-Type": "application/json",
		"accessKey":    accessKey,
		"dateTime":     time.Now().UTC().Format(time.RFC3339),
	}
}

func ackMessage(accessKey, secretKey, topic, consumerGroupId, ackAction, ackIndex, httpUrl string) {
	headers := getHeaders(accessKey)
	body := map[string]interface{}{
		"topic":           topic,
		"consumerGroupId": consumerGroupId,
		"ackAction":       ackAction,
		"ackIndex":        ackIndex,
	}
	signSource := getSignSource(headers, body)
	signature, _ := getSignature(signSource, secretKey)
	headers["signature"] = signature
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return
	}
	req, err := http.NewRequest("POST", httpUrl+"/v2/ack", strings.NewReader(string(jsonBody)))
	if err != nil {
		return
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	fmt.Println(string(respBody))
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(resp.Body)
}

func getMessages(accessKey, secretKey, topic, consumerGroupId, httpUrl string, size int, autoAck bool) (string, error) {
	headers := getHeaders(accessKey)
	params := map[string]interface{}{
		"topic":           topic,
		"consumerGroupId": consumerGroupId,
		"size":            strconv.Itoa(size),
		"ack":             strconv.FormatBool(autoAck),
	}
	signSource := getSignSource(headers, params)
	headers["signature"], _ = getSignature(signSource, secretKey)
	req, err := http.NewRequest("GET", httpUrl+"/v2/messages", nil)
	if err != nil {
		return "", err
	}
	q := req.URL.Query()
	for k, v := range params {
		switch val := v.(type) {
		case string:
			q.Add(k, val)
		}
	}
	req.URL.RawQuery = q.Encode()
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(resp.Body)
	return string(respBody), nil
}

func testSendMessage(accessKey, secretKey, topic, httpUrl string) {
	headers := getHeaders(accessKey)
	messages := make([]map[string]interface{}, 1)
	for i := 0; i < 1; i++ {
		messages[i] = map[string]interface{}{
			"body":         fmt.Sprintf("message-%d", i),
			"delaySeconds": 0,
			"tag":          "mytag",
			"properties":   map[string]interface{}{fmt.Sprintf("%d", 1): "test"},
		}
	}
	body := map[string]interface{}{
		"topic":    topic,
		"type":     "NORMAL",
		"messages": messages,
	}
	signSource := getSignSource(headers, body)
	signature, _ := getSignature(signSource, secretKey)
	headers["signature"] = signature
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return
	}
	req, err := http.NewRequest("POST", httpUrl+"/v2/messages", strings.NewReader(string(jsonBody)))
	if err != nil {
		return
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	fmt.Println("Send result: ", string(respBody))
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(resp.Body)
}

func testReceiveAutoAck(accessKey, secretKey, topic, consumerGroupId, httpUrl string) {
	startTime := int(time.Now().Unix())
	for {
		if time.Now().Unix()-int64(startTime) >= 30 {
			break
		}
		response, err := getMessages(accessKey, secretKey, topic, consumerGroupId, httpUrl, 5, true)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("Receive response: ", response)
		var result map[string]interface{}
		err = json.Unmarshal([]byte(response), &result)
		if err != nil {
			panic(err)
		}
		messages, ok := result["result"].(map[string]interface{})["messages"]
		if !ok || messages == nil {
			fmt.Println("No messages received")
			time.Sleep(1 * time.Second)
			continue
		}
		for _, m := range messages.([]interface{}) {
			fmt.Println("Receive message: ", m)
		}
	}
}

func testReceiveManualAck(accessKey, secretKey, topic, consumerGroupId, httpUrl string) {
	startTime := int(time.Now().Unix())
	for {
		if time.Now().Unix()-int64(startTime) >= 30 {
			break
		}
		response, err := getMessages(accessKey, secretKey, topic, consumerGroupId, httpUrl, 5, false)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("Receive response: ", response)
		var result map[string]interface{}
		err = json.Unmarshal([]byte(response), &result)
		if err != nil {
			panic(err)
		}
		messages, ok := result["result"].(map[string]interface{})["messages"]
		if !ok || messages == nil {
			fmt.Println("No messages received")
			time.Sleep(1 * time.Second)
			continue
		}
		ackIndex, ok := result["result"].(map[string]interface{})["ackIndex"].(string)
		if !ok {
			fmt.Println("Invalid ackIndex")
			time.Sleep(1 * time.Second)
			continue
		}
		fmt.Println("ackIndex = ", ackIndex)
		for _, m := range messages.([]interface{}) {
			fmt.Println("Receive message: ", m)
		}
		ackMessage(accessKey, secretKey, topic, consumerGroupId, "SUCCESS", ackIndex, httpUrl)
	}
}

/*
 * 详细接口文档请参考: https://docs.jdcloud.com/cn/message-queue/consume-message
 */
func main() {
	// 京东云 access key
	accessKey := "xxx"
	// 京东云 secret key
	secretKey := "xxx"
	// topic名称, 被授权订阅的topic需要填写topic全称。云鼎集群全称和简称相同。
	topic := "xxx"
	// http 或者 https 接入点， 详情参考 https://docs.jdcloud.com/cn/message-queue/faq
	httpUrl := "https://jcq-shared-004.cn-north-1.jdcloud.com"
	// 订阅组名称
	consumerGroupId := "xxx"
	// 测试发送，用户有生产权限的topic才能发送
	//testSendMessage(accessKey, secretKey, topic, httpUrl)
	// 测试接收，用户有订阅权限的topic才能接收， 自动ack
	//testReceiveAutoAck(accessKey, secretKey, topic, consumerGroupId, httpUrl)

	// 测试接收，用户有订阅权限的topic才能接收， 手动ack
	testReceiveManualAck(accessKey, secretKey, topic, consumerGroupId, httpUrl)

}
