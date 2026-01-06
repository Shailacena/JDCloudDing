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
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/labstack/echo/v4"
	"github.com/spf13/cast"
	"gorm.io/gorm"
)

type JDPay struct {
	MerchantId int32
	Merchant   *model.Merchant
	db         *gorm.DB
	CD         int
}

func NewJDPay(merchantId int32) (*JDPay, error) {
	db := data.Instance()
	merchant, err := common.FindMerchant(db, merchantId)
	if err != nil {
		return nil, err
	}

	cd := 10
	if !config.IsProd() {
		cd = 2
	}

	return &JDPay{
		db:         db,
		MerchantId: merchantId,
		Merchant:   merchant,
		CD:         cd,
	}, nil
}

func (jd *JDPay) GetMerchant(c echo.Context) *model.Merchant {
	return jd.Merchant
}

func (jd *JDPay) FindGoods(c echo.Context, channelId string, merchantId int32, amount float64) (*common.FindGoodsResult, error) {
	db := jd.db

	rows, err := common.FindGoods(c, db, channelId, merchantId, amount)
	if err != nil {
		c.Logger().Errorf("common.FindGoods error=%s", err)
		return nil, errorx.ErrGoodsNotFound
	}
	defer rows.Close()

	for rows.Next() {
		var res common.FindGoodsResult
		db.ScanRows(rows, &res)

		if len(res.SkuId) == 0 {
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

func (jd *JDPay) GenOrder(c echo.Context, req types.CreateOrderReq, goods *common.FindGoodsResult, orderId string) (model.Order, error) {
	skuId := goods.SkuId
	partnerId := goods.PartnerId
	shop := goods.ShopName
	payType := goods.PayType
	merchant := jd.Merchant

	// 获取partner信息以确定darkNumber长度
	partner, err := common.FindPartner(jd.db, partnerId)
	if err != nil {
		c.Logger().Errorf("FindPartner error=%s", err)
		// 如果获取失败，使用默认长度11
		return model.Order{
			OrderId:         orderId,
			ChannelId:       model.ChannelId(req.ChannelId),
			MerchantId:      uint(req.MerchantId),
			MerchantName:    merchant.Nickname,
			MerchantOrderId: req.MerchantTradeNo,
			Amount:          util.ToDecimal(req.Amount),
			SkuId:           skuId,
			IP:              c.RealIP(),
			NotifyUrl:       req.NotifyUrl,
			Shop:            shop,
			PartnerId:       partnerId,
			PartnerName:     goods.PartnerNickname,
			PayType:         payType,
			DarkNumber:      common.RandomDarkNumber(),
		}, nil
	}

	// 使用partner配置的darkNumber长度
	darkNumber := common.RandomDarkNumberWithLength(partner.DarkNumberLength)

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
		Shop:            shop,
		PartnerId:       partnerId,
		PartnerName:     goods.PartnerNickname,
		PayType:         payType,
		DarkNumber:      darkNumber,
	}

	return o, nil
}

func (jd *JDPay) GenPayUrl(c echo.Context, baseUrl string, o model.Order, numId string) string {
	return fmt.Sprintf("%s?merchantOrderId=%s&orderId=%s&price=%f&sku=%s&time=%d&ts=%d&code=%s", baseUrl, o.MerchantOrderId, o.OrderId, o.Amount, o.SkuId, jd.CD*60, time.Now().Unix(), o.DarkNumber)
}

type JDPayNotify struct {
	db *gorm.DB

	fromPlatform string
	aopic        string
	sign         string
	jsonStr      string
	jsonData     types.JDJsonData
	timestamp    string
}

func NewJDPayNotify(fromPlatform, timestamp, aopic, sign, jsonStr string) (*JDPayNotify, error) {
	db := data.Instance()

	jsonData := types.JDJsonData{}
	err := json.Unmarshal([]byte(jsonStr), &jsonData)
	if err != nil {
		return nil, err
	}

	return &JDPayNotify{
		db:           db,
		fromPlatform: fromPlatform,
		aopic:        aopic,
		sign:         sign,
		jsonStr:      jsonStr,
		timestamp:    timestamp,
		jsonData:     jsonData,
	}, nil
}

func (jdn *JDPayNotify) Handle(c echo.Context) error {
	var err error
	db := jdn.db
	jsonData := jdn.jsonData
	jsonStr := jdn.jsonStr
	timestamp := jdn.timestamp
	sign := jdn.sign

	skuId := cast.ToString(jsonData.SkuId)
	gameAccount := jsonData.GameAccount
	c.Logger().Infof("skuId: %s", skuId)

	var partnerId uint
	var orderId string
	o, findOrderByDarkNumberErr := common.FindOrderByDarkNumber(c, db, skuId, gameAccount)
	if findOrderByDarkNumberErr != nil {
		c.Logger().Error("FindOrderByDarkNumber Error", skuId, gameAccount, findOrderByDarkNumberErr)

		goods, err1 := repository.Goods.GetBySkuId(c, db, skuId)
		if err1 != nil {
			return fmt.Errorf("GetBySkuId findOrderByDarkNumber=%s, err=%s", gameAccount, err1)
		}
		partnerId = goods.PartnerId
	} else {
		partnerId = o.PartnerId
		orderId = o.OrderId
	}

	if o != nil && o.Status != model.OrderStatusUnpaid {
		c.Logger().Infof("该订单已处理, order=%s, status=%d", orderId, o.Status)
		return nil
	}
	if o == nil && findOrderByDarkNumberErr == nil {
		findOrderByDarkNumberErr = errors.New("findOrderByDarkNumberErr 失败")
	}

	partner, err := common.FindPartner(db, partnerId)
	if err != nil {
		return fmt.Errorf("findPartner: %s", err.Error())
	}

	agisoSign := common.NewAgisoSign(partner.AqsAppSecret, jsonStr, timestamp)
	if !agisoSign.Check(c, sign) {
		return errorx.ErrInvalidSign
	}

	partnerOrderId := jsonData.OrderId
	var rechargeSendErr error
	if findOrderByDarkNumberErr == nil {
		num := 5
		for i := 0; i < num; i++ {
			rechargeSendErr = jdn.RechargeSend(c, partner, cast.ToString(partnerOrderId))
			if rechargeSendErr == nil {
				break
			} else {
				c.Logger().Error("自动发货失败", orderId, partnerOrderId, rechargeSendErr)
				if i == num-1 {
					return nil
				}
			}
			time.Sleep(5 * time.Second)
		}
	} else {
		err = jdn.Refund(c, partner, cast.ToString(partnerOrderId))
		if err != nil {
			err1 := repository.Order.Update(db, orderId, model.Order{Status: model.OrderStatusRefundFailed})
			if err1 != nil {
				c.Logger().Error("订单退款失败更新失败", orderId, partnerOrderId, err1)
			}
			c.Logger().Error("订单退款失败", orderId, partnerOrderId, err)
			return err
		}

		err1 := repository.Order.Update(db, orderId, model.Order{Status: model.OrderStatusRefundSuccessful})
		if err1 != nil {
			c.Logger().Error("订单退款成功更新失败", orderId, partnerOrderId, err1)
		}
		return nil
	}

	err = jdn.Transaction(c, db, jsonData, skuId)
	if err != nil {
		c.Logger().Errorf("AgisoNotifyHandler Transaction error= %s", err.Error())
		return nil
	}

	return nil
}

func (jdn *JDPayNotify) Transaction(c echo.Context, db *gorm.DB, jsonData types.JDJsonData, skuId string) error {
	var o *model.Order
	var err error

	err = db.Transaction(func(tx *gorm.DB) error {
		kw := cast.ToString(jsonData.OrderId)
		if len(jsonData.GameAccount) > 0 {
			kw = jsonData.GameAccount
		}

		o, err = common.FindOrderByDarkNumber(c, tx, skuId, kw)
		if err != nil {
			return fmt.Errorf("findOrderBySkuId, gameAccount=%s, error: %s", jsonData.GameAccount, err.Error())
		}

		if o.Status == model.OrderStatusFinish {
			return fmt.Errorf("订单已完成, orderId=%s, skuId=%s, order.Status=%d", o.OrderId, skuId, o.Status)
		}

		o.ReceivedAmount = util.ToDecimal(cast.ToFloat64(jsonData.TotalPrice))

		var payTime time.Time
		if len(jsonData.CreateTime) > 0 {
			payTime, err = time.Parse("2006-01-02T15:04:05", jsonData.CreateTime)
			if err != nil {
				return err
			}
		}

		params := common.UpdateOrderParams{
			PartnerOrderId: cast.ToString(jsonData.OrderId),
			OrderId:        o.OrderId,
			PayAccount:     jsonData.Pin,
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
		return fmt.Errorf("订单是nil, gameAccount=%s, skuId=%s", jsonData.GameAccount, skuId)
	}

	err = common.NotifyMerchant(contextx.NewContextFromEcho(c), db, o)
	if err != nil {
		return fmt.Errorf("NotifyMerchant error= %s", err.Error())
	}

	return nil
}

type AutoGenerated struct {
	IsSuccess  bool   `json:"IsSuccess"`
	Data       any    `json:"Data"`
	ErrorCode  int    `json:"Error_Code"`
	ErrorMsg   string `json:"Error_Msg"`
	AllowRetry any    `json:"AllowRetry"`
	RequestID  string `json:"RequestId"`
}

func makeHeaders(aqsToken string) map[string]string {
	headers := map[string]string{
		"Authorization": "Bearer " + aqsToken,
		"ApiVersion":    "1",
	}
	return headers
}

func makeParams(aqsAppSecret, partnerOrderId string) map[string]string {
	timestamp := cast.ToString(time.Now().Unix())
	d := struct {
		Timestamp string `json:"timestamp"`
		Tid       string `json:"tid"`
	}{
		Timestamp: timestamp,
		Tid:       partnerOrderId,
	}

	agisoSign := common.NewAgisoSign(aqsAppSecret, "", timestamp)
	sign := agisoSign.Generate(d)

	params := map[string]string{
		"tid":       cast.ToString(d.Tid),
		"timestamp": cast.ToString(d.Timestamp),
		"sign":      sign,
	}
	return params
}

func isInvalidAqs(aqsAppSecret, aqsToken string) bool {
	return len(aqsAppSecret) == 0 || len(aqsToken) == 0
}

func (jdn *JDPayNotify) RechargeSend(c echo.Context, partner *model.Partner, partnerOrderId string) error {
	if isInvalidAqs(partner.AqsAppSecret, partner.AqsToken) {
		return errors.New("阿奇索缺少Secret或Token")
	}

	url := "http://gw.api.agiso.com/aldsJd/GameCard/RechargeSend"

	headers := makeHeaders(partner.AqsToken)

	params := makeParams(partner.AqsAppSecret, partnerOrderId)

	var result AutoGenerated
	client := resty.New()
	resp, err := client.SetHeaders(headers).SetTimeout(60 * time.Second).R().SetFormData(params).SetResult(&result).Post(url)
	if err != nil {
		return err
	}

	c.Logger().Infof(" partnerOrderId=%s, StatusCode=%d", partnerOrderId, resp.StatusCode())
	c.Logger().Infof(" partnerOrderId=%s, resp.Body=%s", partnerOrderId, string(resp.Body()))
	c.Logger().Infof(" partnerOrderId=%s, result=%+v", partnerOrderId, result)
	if resp.StatusCode() != 200 || !result.IsSuccess {
		return errorx.ErrAutomaticShipment
	}

	c.Logger().Infof("自动发货成功, partnerOrderId=%s", partnerOrderId)
	return nil
}

func (jdn *JDPayNotify) Refund(c echo.Context, partner *model.Partner, partnerOrderId string) error {
	if isInvalidAqs(partner.AqsAppSecret, partner.AqsToken) {
		return errors.New("阿奇索缺少Secret或Token")
	}

	url := "http://gw.api.agiso.com/aldsJd/GameCard/Refund"

	headers := makeHeaders(partner.AqsToken)

	params := makeParams(partner.AqsAppSecret, partnerOrderId)

	var result AutoGenerated
	client := resty.New()
	resp, err := client.SetHeaders(headers).SetTimeout(60 * time.Second).R().SetFormData(params).SetResult(&result).Post(url)
	if err != nil {
		return err
	}

	c.Logger().Infof(" partnerOrderId=%s, StatusCode=%d", partnerOrderId, resp.StatusCode())
	c.Logger().Infof(" partnerOrderId=%s, resp.Body=%s", partnerOrderId, string(resp.Body()))
	c.Logger().Infof(" partnerOrderId=%s, result=%+v", partnerOrderId, result)
	if resp.StatusCode() != 200 || !result.IsSuccess {
		if result.ErrorMsg == "该笔订单已经退款" {
			return nil
		}
		return errorx.ErrRefundFailed
	}

	c.Logger().Infof("订单退款成功, partnerOrderId=%s", partnerOrderId)
	return nil
}
