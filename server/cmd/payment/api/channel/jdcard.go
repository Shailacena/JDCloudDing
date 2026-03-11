package channel

import (
	"apollo/server/cmd/payment/api/common"
	"apollo/server/cmd/payment/pkg/errorx"
	"apollo/server/cmd/payment/pkg/types"
	"apollo/server/internal/model"
	"apollo/server/internal/repository"
	"apollo/server/pkg/config"
	"apollo/server/pkg/contextx"
	"apollo/server/pkg/data"
	"apollo/server/pkg/util"
	"crypto/md5"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type JDCard struct {
	MerchantId int32
	Merchant   *model.Merchant
	db         *gorm.DB
	CD         int
}

func NewJDCard(merchantId int32) (*JDCard, error) {
	db := data.Instance()
	merchant, err := common.FindMerchant(db, merchantId)
	if err != nil {
		return nil, err
	}

	cd := 5
	if !config.IsProd() {
		cd = 2
	}

	return &JDCard{
		db:         db,
		MerchantId: merchantId,
		Merchant:   merchant,
		CD:         cd,
	}, nil
}

func (jd *JDCard) GetMerchant(c echo.Context) *model.Merchant {
	return jd.Merchant
}

func (jd *JDCard) FindGoods(c echo.Context, channelId string, merchantId int32, amount float64) (*common.FindGoodsResult, error) {
	db := jd.db

	rows, err := common.FindGoods(c, db, channelId, merchantId, amount)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var res common.FindGoodsResult
		db.ScanRows(rows, &res)

		if len(res.SkuId) == 0 {
			continue
		}

		// 京东入鼎LOC商品
		order := model.Order{}
		err = db.Where("sku_id = ? AND end_lock_at >= ?", res.SkuId, time.Now()).Order("created_at desc").First(&order).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			continue
		}
		if err == nil {
			continue
		}

		err = db.Model(&model.Goods{}).Where("id = ?", res.Id).Update("weight", gorm.Expr("weight + ?", 1)).Error
		if err != nil {
			c.Logger().Error("update goods weight error", err.Error())
		}
		return &res, nil
	}

	return nil, errorx.ErrGoodsNotFound
}

func (jd *JDCard) GenOrder(c echo.Context, req types.CreateOrderReq, goods *common.FindGoodsResult, orderId string) (model.Order, error) {
	skuId := goods.SkuId
	partnerId := goods.PartnerId
	shop := goods.ShopName
	payType := goods.PayType
	merchant := jd.Merchant

	cd := jd.CD
	o := model.Order{
		OrderId:         orderId,
		ChannelId:       model.ChannelId(req.ChannelId),
		MerchantId:      uint(req.MerchantId),
		MerchantName:    merchant.Nickname,
		MerchantOrderId: req.MerchantTradeNo,
		Amount:          util.ToDecimal(req.Amount),
		SkuId:           skuId,
		IP:              c.RealIP(),
		NotifyUrl:       req.NotifyUrl,
		EndLockAt:       time.Now().Add(time.Minute * time.Duration(cd)),
		Shop:            shop,
		PartnerId:       partnerId,
		PartnerName:     goods.PartnerNickname,
		PayType:         payType,
	}

	return o, nil
}

func (jd *JDCard) GenPayUrl(c echo.Context, baseUrl string, o model.Order, numId string) string {
	return fmt.Sprintf("%s?merchantOrderId=%s&orderId=%s&price=%f&sku=%s&time=%d&ts=%d&numId=%s", baseUrl, o.MerchantOrderId, o.OrderId, o.Amount, o.SkuId, jd.CD*60, time.Now().Unix(), numId)
}

type JDCardNotify struct {
	db *gorm.DB

	jsonData types.JDCloudOrderChangeJsonData
	skuId    string
}

func NewJDCardNotify(db *gorm.DB, orderId int64, skuId string, orderCreateTime string) *JDCardNotify {
	return &JDCardNotify{
		db: db,
		jsonData: types.JDCloudOrderChangeJsonData{
			OrderId:         orderId,
			OrderCreateTime: orderCreateTime,
		},
		skuId: skuId,
	}
}

func (jdc *JDCardNotify) Handle(c echo.Context) error {
	db := jdc.db
	jsonData := jdc.jsonData
	skuId := jdc.skuId

	c.Logger().Infof("JDCloudOrderChangeJsonData skuId: %s, orderId: %d", skuId, jsonData.OrderId)

	createdStr := jsonData.OrderCreateTime
	if len(createdStr) == 0 {
		createdStr = time.Now().Format(time.DateTime)
	}
	if strings.Contains(createdStr, "T") {
		createdStr = strings.Replace(createdStr, "T", " ", 1)
	}
	o, err := common.FindOrderBySkuId(c, db, skuId, createdStr)
	if err != nil {
		return fmt.Errorf("findOrderBySkuId: %s", err.Error())
	}

	if o == nil {
		return fmt.Errorf("订单未找到, skuId=%s, orderId=%d", skuId, jsonData.OrderId)
	}

	partner, err := common.FindPartner(db, o.PartnerId)
	if err != nil {
		return fmt.Errorf("findPartner: %s", err.Error())
	}

	err = jdc.Transaction(c, db, jsonData, skuId, partner)
	if err != nil {
		c.Logger().Errorf("JDCardNotify Transaction error= %s", err.Error())
		return err
	}

	return nil
}

func (jdc *JDCardNotify) Transaction(c echo.Context, db *gorm.DB, jsonData types.JDCloudOrderChangeJsonData, skuId string, partner *model.Partner) error {
	var o *model.Order
	var err error

	err = db.Transaction(func(tx *gorm.DB) error {
		createdStr := jsonData.OrderCreateTime
		if len(createdStr) == 0 {
			createdStr = time.Now().Format(time.DateTime)
		}
		if strings.Contains(createdStr, "T") {
			createdStr = strings.Replace(createdStr, "T", " ", 1)
		}
		o, err = common.FindOrderBySkuId(c, tx, skuId, createdStr)
		if err != nil {
			return fmt.Errorf("findOrderBySkuId, skuId=%s, error: %s", skuId, err.Error())
		}

		if o == nil {
			return fmt.Errorf("订单未找到, skuId=%s, orderId=%d", skuId, jsonData.OrderId)
		}
		if o.SkuId != skuId {
			return fmt.Errorf("订单sku不匹配, expect=%s, actual=%s", skuId, o.SkuId)
		}

		if o.Status == model.OrderStatusFinish {
			return fmt.Errorf("订单已完成, orderId=%s, skuId=%s, order.Status=%d", o.OrderId, skuId, o.Status)
		}

		o.Status = model.OrderStatusPaid

		var payTime time.Time
		if len(jsonData.OrderCreateTime) > 0 {
			createdStr := jsonData.OrderCreateTime
			if strings.Contains(createdStr, "T") {
				createdStr = strings.Replace(createdStr, "T", " ", 1)
			}
			payTime, err = time.Parse(time.DateTime, createdStr)
			if err != nil {
				return err
			}
		}
		o.PayAt = payTime

		params := common.UpdateOrderParams{
			PartnerOrderId: fmt.Sprintf("%d", jsonData.OrderId),
			OrderId:        o.OrderId,
			PayAccount:     "",
			PayTime:        payTime,
			ReceivedAmount: o.ReceivedAmount,
		}
		err = common.UpdateOrder(c, tx, params)
		if err != nil {
			return fmt.Errorf("updateOrder error= %s", err.Error())
		}

		err = common.UpdatePartnerBalance(c, tx, o)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	if o == nil {
		return fmt.Errorf("订单是nil, skuId=%s", skuId)
	}

	err = common.NotifyMerchant(contextx.NewContextFromEcho(c), db, o)
	if err != nil {
		return fmt.Errorf("NotifyMerchant error= %s", err.Error())
	}

	err = jdc.SendCard(c, db, o, partner)
	if err != nil {
		c.Logger().Error("发放卡密失败:", err)
		return err
	}

	return nil
}

func (jdc *JDCardNotify) SendCard(c echo.Context, db *gorm.DB, o *model.Order, partner *model.Partner) error {
	c.Logger().Info("JDCard 开始发放卡密, orderId:", o.OrderId)

	amount := o.Amount

	var exist model.PriceCard
	err := db.Where("order_id = ?", o.OrderId).First(&exist).Error
	if err == nil {
		c.Logger().Info("订单已有卡密，跳过发放, orderId:", o.OrderId, ", cardNo:", exist.CardNo)
		return nil
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("查询已发卡密失败: %s", err.Error())
	}

	cards, err := repository.PriceCard.GetAvailableCards(db, amount, model.CardTypeReal, 1)
	if err != nil {
		return fmt.Errorf("查询卡密失败: %s", err.Error())
	}

	if len(cards) == 0 {
		return fmt.Errorf("库存不足, partnerId=%d, amount=%.2f", partner.ID, amount)
	}

	card := cards[0]
	card.OrderId = o.OrderId
	card.UseIP = o.IP

	jdOrderId := jdc.jsonData.OrderId
	err = uploadCardToJD(c, jdOrderId, card.CardNo, card.Password)
	if err != nil {
		card.CardStatus = model.CardStatusFailed
		card.Remark = fmt.Sprintf("京东回传卡密失败: %s", err.Error())
		db.Save(card)
		return fmt.Errorf("京东回传卡密失败: %s", err.Error())
	}

	card.CardStatus = model.CardStatusSent
	err = db.Save(card).Error
	if err != nil {
		return fmt.Errorf("更新卡密状态失败: %s", err.Error())
	}

	c.Logger().Info("JDCard 发放卡密成功, cardNo:", card.CardNo, ", orderId:", o.OrderId)

	go func() {
		time.Sleep(10 * time.Second)
		consumeCardFromJD(c, jdOrderId, card.CardNo, card.Password, o.ID)
	}()

	return nil
}

func uploadCardToJD(c echo.Context, orderId int64, cardNo, password string) error {
	conf := config.Get()
	if conf == nil {
		conf = config.New("configs/config.yaml")
	}
	jdConf := conf.JDCloudConfig

	if jdConf.AppKey == "" || jdConf.AppSecret == "" || jdConf.Token == "" {
		return errors.New("京东配置缺失")
	}

	paramJson := fmt.Sprintf(`{"orderId":%d,"cardNumber":"%s","pwdNumber":"%s"}`, orderId, cardNo, password)
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	params := map[string]string{
		"method":       "jingdong.pop.oto.checkNumbers.upload",
		"access_token": jdConf.Token,
		"app_key":      jdConf.AppKey,
		"format":       "json",
		"v":            "2.0",
		"timestamp":    timestamp,
		"param_json":   paramJson,
	}
	sign := generateJDCloudSign(jdConf.AppSecret, params)
	params["sign"] = sign

	client := resty.New()
	url := "https://api.jd.com/routerjson"

	var result struct {
		jingdong_pop_oto_checkNumbers_upload_responce struct {
			Result struct {
				ResultMessage string `json:"result_message"`
				ResultCode    string `json:"result_code"`
				IsSuccess     string `json:"is_success"`
			} `json:"result"`
		} `json:"jingdong_pop_oto_checkNumbers_upload_responce"`
	}

	resp, err := client.SetTimeout(10 * time.Second).R().SetFormData(params).SetResult(&result).Post(url)

	if err != nil {
		c.Logger().Error("调用京东回传卡密API失败:", err)
		return err
	}

	c.Logger().Info("京东回传卡密API响应:", string(resp.Body()))

	respResult := result.jingdong_pop_oto_checkNumbers_upload_responce.Result
	if respResult.IsSuccess != "true" {
		return fmt.Errorf("京东回传失败: %s", respResult.ResultMessage)
	}

	return nil
}

func consumeCardFromJD(c echo.Context, orderId int64, cardNo, password string, orderDbId uint) {
	db := data.Instance()

	var card model.PriceCard
	err := db.Where("card_no = ?", cardNo).First(&card).Error
	if err != nil {
		c.Logger().Error("核销查找卡密失败:", err)
		return
	}

	conf := config.Get()
	if conf == nil {
		conf = config.New("configs/config.yaml")
	}
	jdConf := conf.JDCloudConfig

	paramJson := fmt.Sprintf(`{"codeNum":"%s","pwdNumber":"%s"}`, cardNo, password)
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	params := map[string]string{
		"method":       "jingdong.loc.code.consume",
		"access_token": jdConf.Token,
		"app_key":      jdConf.AppKey,
		"format":       "json",
		"v":            "2.0",
		"timestamp":    timestamp,
		"param_json":   paramJson,
	}
	sign := generateJDCloudSign(jdConf.AppSecret, params)
	params["sign"] = sign

	client := resty.New()
	url := "https://api.jd.com/routerjson"

	var result struct {
		jingdong_loc_code_consume_responce struct {
			ReturnType struct {
				Result        string `json:"result"`
				Success       string `json:"success"`
				ResultCode    string `json:"resultCode"`
				ResultMessage string `json:"resultMessage"`
			} `json:"returnType"`
		} `json:"jingdong_loc_code_consume_responce"`
	}

	resp, err := client.SetTimeout(10 * time.Second).R().SetFormData(params).SetResult(&result).Post(url)

	if err != nil {
		c.Logger().Error("调用京东核销API失败:", err)
		card.CardStatus = model.CardStatusFailed
		card.Remark = fmt.Sprintf("京东核销API调用失败: %s", err.Error())
		db.Save(card)
		return
	}

	c.Logger().Info("京东核销API响应:", string(resp.Body()))

	respResult := result.jingdong_loc_code_consume_responce.ReturnType
	if respResult.Success != "true" {
		c.Logger().Error("京东核销失败:", respResult.ResultMessage)
		card.CardStatus = model.CardStatusFailed
		card.Remark = fmt.Sprintf("京东核销失败: %s", respResult.ResultMessage)
		db.Save(card)
		return
	}

	card.CardStatus = model.CardStatusSuccess
	now := time.Now()
	card.UsedAt = &now
	db.Save(card)

	c.Logger().Info("JDCard 核销成功, cardNo:", cardNo, ", orderId:", orderId)
}

func generateJDCloudSign(appSecret string, params map[string]string) string {
	keys := make([]string, 0, len(params))
	for k, v := range params {
		if k == "sign" || len(v) == 0 {
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var b strings.Builder
	b.WriteString(appSecret)
	for _, k := range keys {
		b.WriteString(k)
		b.WriteString(params[k])
	}
	b.WriteString(appSecret)
	sum := md5.Sum([]byte(b.String()))
	return fmt.Sprintf("%X", sum)
}
