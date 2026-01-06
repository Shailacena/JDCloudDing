package channel

import (
	"apollo/server/cmd/payment/api/common"
	"apollo/server/cmd/payment/pkg/errorx"
	"apollo/server/cmd/payment/pkg/types"
	"apollo/server/internal/model"
	"apollo/server/pkg/config"
	"apollo/server/pkg/contextx"
	"apollo/server/pkg/data"
	"apollo/server/pkg/util"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/spf13/cast"
	"gorm.io/gorm"
)

type TBPay struct {
	MerchantId int32
	Merchant   *model.Merchant
	db         *gorm.DB
	CD         int

	jsonData types.TBJsonData
}

func NewTBPay(merchantId int32) (*TBPay, error) {
	db := data.Instance()
	merchant, err := common.FindMerchant(db, merchantId)
	if err != nil {
		return nil, err
	}

	cd := 5
	if !config.IsProd() {
		cd = 2
	}

	return &TBPay{
		db:         db,
		MerchantId: merchantId,
		Merchant:   merchant,
		CD:         cd,
	}, nil
}

func (tb *TBPay) GetMerchant(c echo.Context) *model.Merchant {
	return tb.Merchant
}

func (tb *TBPay) FindGoods(c echo.Context, channelId string, merchantId int32, amount float64) (*common.FindGoodsResult, error) {
	db := tb.db

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

		// 天猫直付，商品是否在使用中
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

func (tb *TBPay) GenOrder(c echo.Context, req types.CreateOrderReq, goods *common.FindGoodsResult, orderId string) (model.Order, error) {
	skuId := goods.SkuId
	partnerId := goods.PartnerId
	shop := goods.ShopName
	payType := goods.PayType
	merchant := tb.Merchant

	cd := tb.CD
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

func (tb *TBPay) GenPayUrl(c echo.Context, baseUrl string, o model.Order, numId string) string {
	return fmt.Sprintf("%s?merchantOrderId=%s&orderId=%s&price=%f&sku=%s&time=%d&ts=%d&numId=%s", baseUrl, o.MerchantOrderId, o.OrderId, o.Amount, o.SkuId, tb.CD*60, time.Now().Unix(), numId)
}

func (tb *TBPay) UnmarshalNotifyJsonData(c echo.Context, jsonStr string) error {
	jsonData := types.TBJsonData{}
	err := json.Unmarshal([]byte(jsonStr), &jsonData)
	if err != nil {
		return err
	}
	tb.jsonData = jsonData
	return nil
}

type TBPayNotify struct {
	db *gorm.DB

	fromPlatform string
	aopic        string
	sign         string
	jsonStr      string
	jsonData     types.TBJsonData
	timestamp    string
}

func NewTBPayNotify(fromPlatform, timestamp, aopic, sign, jsonStr string) (*TBPayNotify, error) {
	db := data.Instance()

	jsonData := types.TBJsonData{}
	err := json.Unmarshal([]byte(jsonStr), &jsonData)
	if err != nil {
		return nil, err
	}

	return &TBPayNotify{
		db:           db,
		fromPlatform: fromPlatform,
		aopic:        aopic,
		sign:         sign,
		jsonStr:      jsonStr,
		timestamp:    timestamp,
		jsonData:     jsonData,
	}, nil
}

func (tbn *TBPayNotify) Handle(c echo.Context) error {
	db := tbn.db
	jsonData := tbn.jsonData
	jsonStr := tbn.jsonStr
	timestamp := tbn.timestamp
	sign := tbn.sign

	if jsonData.Status != types.TRADE_FINISHED {
		return errors.New("jsonData.Status error")
	}

	for _, v := range jsonData.Orders {
		skuId := cast.ToString(v.NumIid)
		if len(v.SkuID) > 0 {
			skuId = v.SkuID
		}

		c.Logger().Infof("skuId: %s", skuId)

		o, err := common.FindOrderBySkuId(c, db, skuId, jsonData.Created)
		if err != nil {
			return fmt.Errorf("findOrderBySkuId: %s", err.Error())
		}

		partner, err := common.FindPartner(db, o.PartnerId)
		if err != nil {
			return fmt.Errorf("findPartner: %s", err.Error())
		}

		secret := partner.AqsAppSecret
		if partner.Type == model.PartnerTypeAnssy {
			secret = partner.AnssyAppSecret
		}

		agisoSign := common.NewAgisoSign(secret, jsonStr, timestamp)
		if !agisoSign.Check(c, sign) {
			return errorx.ErrInvalidSign
		}

		err = tbn.Transaction(c, db, jsonData, skuId)
		if err != nil {
			c.Logger().Errorf("AgisoNotifyHandler Transaction error= %s", err.Error())
			continue
		}
	}

	return nil
}

func (tbn *TBPayNotify) Transaction(c echo.Context, db *gorm.DB, jsonData types.TBJsonData, skuId string) error {
	var o *model.Order
	var err error

	err = db.Transaction(func(tx *gorm.DB) error {
		o, err = common.FindOrderBySkuId(c, tx, skuId, jsonData.Created)
		if err != nil {
			return fmt.Errorf("findOrderBySkuId, skuId=%s, error: %s", skuId, err.Error())
		}

		if o.Status == model.OrderStatusFinish {
			return fmt.Errorf("订单已完成, orderId=%s, skuId=%s, order.Status=%d", o.OrderId, skuId, o.Status)
		}

		o.ReceivedAmount = util.ToDecimal(cast.ToFloat64(jsonData.Payment))
		var payTime time.Time
		if len(jsonData.PayTime) > 0 {
			payTime, err = time.Parse(time.DateTime, jsonData.PayTime)
			if err != nil {
				return err
			}
		}
		params := common.UpdateOrderParams{
			PartnerOrderId: fmt.Sprintf("%d", jsonData.Tid),
			OrderId:        o.OrderId,
			PayAccount:     jsonData.BuyerNick,
			PayTime:        payTime,
			ShopName:       jsonData.SellerNick,
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

	return nil
}
