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
	"errors"
	"fmt"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/spf13/cast"
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

	jsonData types.JDJsonData
	skuId    string
}

func NewJDCardNotify(db *gorm.DB, orderId int64, skuId string, gameAccount string) *JDCardNotify {
	return &JDCardNotify{
		db: db,
		jsonData: types.JDJsonData{
			OrderId:     orderId,
			SkuId:       cast.ToInt64(skuId),
			GameAccount: gameAccount,
		},
		skuId: skuId,
	}
}

func (jdc *JDCardNotify) Handle(c echo.Context) error {
	db := jdc.db
	jsonData := jdc.jsonData
	skuId := jdc.skuId

	c.Logger().Infof("JDCardNotify skuId: %s, orderId: %d", skuId, jsonData.OrderId)

	o, err := repository.Order.GetBySkuIdDarkNumber(c, db, skuId, fmt.Sprintf("%d", jsonData.OrderId))
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

func (jdc *JDCardNotify) Transaction(c echo.Context, db *gorm.DB, jsonData types.JDJsonData, skuId string, partner *model.Partner) error {
	var o *model.Order
	var err error

	err = db.Transaction(func(tx *gorm.DB) error {
		o, err = repository.Order.GetBySkuIdDarkNumber(c, tx, skuId, fmt.Sprintf("%d", jsonData.OrderId))
		if err != nil {
			return fmt.Errorf("findOrderBySkuId, skuId=%s, error: %s", skuId, err.Error())
		}

		if o == nil {
			return fmt.Errorf("订单未找到, skuId=%s, orderId=%d", skuId, jsonData.OrderId)
		}

		if o.Status == model.OrderStatusFinish {
			return fmt.Errorf("订单已完成, orderId=%s, skuId=%s, order.Status=%d", o.OrderId, skuId, o.Status)
		}

		o.Status = model.OrderStatusPaid
		o.ReceivedAmount = util.ToDecimal(cast.ToFloat64(jsonData.TotalPrice))

		var payTime time.Time
		if len(jsonData.CreateTime) > 0 {
			payTime, err = time.Parse("2006-01-02T15:04:05", jsonData.CreateTime)
			if err != nil {
				return err
			}
		}
		o.PayAt = payTime

		params := common.UpdateOrderParams{
			PartnerOrderId: fmt.Sprintf("%d", jsonData.OrderId),
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

	cards, err := repository.PriceCard.GetAvailableCards(db, amount, model.CardTypeReal, 1)
	if err != nil {
		return fmt.Errorf("查询卡密失败: %s", err.Error())
	}

	if len(cards) == 0 {
		return fmt.Errorf("库存不足, partnerId=%d, amount=%.2f", partner.ID, amount)
	}

	card := cards[0]
	card.UsedStatus = true
	card.UsedAt = func() *time.Time { t := time.Now(); return &t }()
	card.OrderId = o.OrderId
	card.UseIP = o.IP

	err = db.Save(card).Error
	if err != nil {
		return fmt.Errorf("更新卡密状态失败: %s", err.Error())
	}

	c.Logger().Info("JDCard 发放卡密成功, cardNo:", card.CardNo, ", orderId:", o.OrderId)

	return nil
}
