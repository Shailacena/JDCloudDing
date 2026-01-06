package common

import (
	v1 "apollo/server/api/v1"
	"apollo/server/cmd/payment/pkg/errorx"
	"apollo/server/cmd/payment/pkg/types"
	"apollo/server/internal/model"
	"apollo/server/internal/repository"
	"apollo/server/pkg/contextx"
	"apollo/server/pkg/data"
	"apollo/server/pkg/rand"
	"apollo/server/pkg/timex"
	"apollo/server/pkg/util"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/samber/lo"

	"github.com/go-resty/resty/v2"
	"github.com/labstack/echo/v4"
	"github.com/spf13/cast"
	"gorm.io/gorm"
)

// 生成订单号
func GenOrderId() (string, error) {
	id, err := util.SFlake.GenString()
	if err != nil {
		return "", err
	}
	return id, nil
}

// 查找商户
func FindMerchant(db *gorm.DB, merchantId int32) (*model.Merchant, error) {
	merchant := model.Merchant{}
	err := db.Where("id = ? AND enable = ?", merchantId, model.Enabled).First(&merchant).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorx.ErrInvalidMerchant
		}
		return nil, err
	}
	return &merchant, nil
}

type FindGoodsResult struct {
	Id              uint
	SkuId           string
	NumId           string
	ShopName        string
	PartnerId       uint
	PartnerNickname string
	PayType         model.PayType
}

// 查找商品
func FindGoods(c echo.Context, db *gorm.DB, channelId string, merchantId int32, amount float64) (*sql.Rows, error) {
	var admin model.SysUser
	err := db.Model(model.SysUser{}).
		Joins("left join merchant on sys_user.id = merchant.parent_id").
		Where("merchant.id = ?", merchantId).
		Find(&admin).Error
	if err != nil {
		return nil, err
	}

	if admin.ID == 0 {
		return nil, errorx.ErrInvalidMerchant
	}

	parentIds, err := repository.Admin.FindAdminIds(c, admin.ID, admin.Role)
	if err != nil {
		return nil, err
	}

	partners, _, err := repository.Partner.List(c, &v1.ListPartnerReq{}, parentIds)
	if err != nil {
		return nil, err
	}

	partnerIds := lo.Map(partners, func(item *model.Partner, _ int) uint {
		return item.ID
	})

	rows, err := db.Model(model.Goods{}).
		Select("goods.id, goods.sku_id, goods.num_id, goods.shop_name, partner.id as partner_id, partner.nickname as partner_nickname, partner.pay_type").
		Joins("left join partner on partner.id = goods.partner_id").
		Where("partner.channel_id = ? AND goods.amount = ? AND goods.status = ? AND partner.enable = ? AND partner.id IN (?)", channelId, amount, model.GoodsStatusEnabled, model.Enabled, partnerIds).
		Order("partner.priority desc, goods.weight").
		Rows()
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// 查找合作商
func FindPartner(db *gorm.DB, partnerId uint) (*model.Partner, error) {
	partner := model.Partner{}
	err := db.Where("id = ?", partnerId).First(&partner).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("非法合作商 %d", partnerId)
		}
		return nil, err
	}
	return &partner, nil
}

// 根据id查找订单
func FindOrderById(c echo.Context, merchantOrderId string) (*model.Order, error) {
	o, err := repository.Order.GetByMerchantOrderId(c, merchantOrderId)
	if err != nil {
		return nil, err
	}
	return o, nil
}

// 根据skuId查找订单
func FindOrderBySkuId(c echo.Context, db *gorm.DB, skuId string, createdFromJsonData string) (*model.Order, error) {
	// todo 注意时间字符串的时区
	createdTime, err := time.Parse(time.DateTime, createdFromJsonData)
	if err != nil {
		return nil, err
	}

	o, err := repository.Order.GetBySkuId(c, db, skuId, createdTime)
	if err != nil {
		return nil, err
	}
	return o, nil
}

// 根据darkNumber查找订单
func FindOrderByDarkNumber(c echo.Context, db *gorm.DB, skuId string, gameAccount string) (*model.Order, error) {
	o, err := repository.Order.GetBySkuIdDarkNumber(c, db, skuId, gameAccount)
	if err != nil {
		return nil, err
	}
	return o, nil
}

// 根据partnerOrderId查找订单
func FindOrderByPartnerOrderId(c echo.Context, db *gorm.DB, partnerOrderId string) (*model.Order, error) {
	o, err := repository.Order.GetByPartnerOrderId(c, db, partnerOrderId)
	if err != nil {
		return nil, err
	}
	return o, nil
}

// 获得支付页面url
func GetPayBaseUrl(partnerId uint, channel string) (string, error) {
	db := data.Instance()

	partner, err := FindPartner(db, partnerId)
	if err != nil {
		return "", err
	}

	switch partner.PayType {
	case model.WebPay:
		return "https://takeorder-apollo.s3.ap-southeast-1.amazonaws.com/html/pay-" + channel + "-web.html", nil
	case model.AliPay:
		return "https://takeorder-apollo.s3.ap-southeast-1.amazonaws.com/html/pay-" + channel + "-alipay.html", nil
	case model.AppPay:
		return "https://takeorder-apollo.s3.ap-southeast-1.amazonaws.com/html/pay-" + channel + "-app.html", nil
	default:
		return "https://takeorder-apollo.s3.ap-southeast-1.amazonaws.com/html/pay-" + channel + "-web.html", nil
	}
}

// 随机DarkNumber
func RandomDarkNumber() string {
	// 默认生成11位数字
	return RandomDarkNumberWithLength(11)
}

// 根据指定长度生成随机DarkNumber
func RandomDarkNumberWithLength(length int) string {
	if length < 8 {
		length = 8
	} else if length > 15 {
		length = 15
	}

	// 计算最小值和最大值
	low := 1
	for i := 1; i < length; i++ {
		low *= 10
	}
	high := low*10 - 1

	r := rand.Random
	n := r.Int63n(int64(high-low+1)) + int64(low)
	return fmt.Sprintf("%d", n)
}

type UpdateOrderParams struct {
	PartnerOrderId string
	OrderId        string
	PayAccount     string
	PayTime        time.Time
	ReceivedAmount float64
	ShopName       string
}

// 更新订单
func UpdateOrder(c echo.Context, db *gorm.DB, params UpdateOrderParams) error {
	partnerOrderId := params.PartnerOrderId
	orderId := params.OrderId
	payAccount := params.PayAccount
	payTime := params.PayTime
	receivedAmount := params.ReceivedAmount
	shopName := params.ShopName

	o := model.Order{
		Status:         model.OrderStatusPaid,
		PartnerOrderId: partnerOrderId,
		PayAccount:     payAccount,
		Shop:           shopName,
		ReceivedAmount: util.ToDecimal(receivedAmount),
	}
	var err error
	if !payTime.IsZero() {
		o.PayAt = payTime
	}

	err = repository.Order.Update(db, orderId, o)
	if err != nil {
		return err
	}
	return nil
}

// 更新合作商余额
func UpdatePartnerBalance(c echo.Context, db *gorm.DB, o *model.Order) error {
	c.Logger().Infof("OrderId=%s, Update Partner Balance=%.2f", o.OrderId, -o.ReceivedAmount)
	if o.ReceivedAmount == 0 {
		return nil
	}

	err := repository.Partner.UpdateBalance(c, db, o.PartnerId, o.OrderId, -o.ReceivedAmount, model.BalanceFromTypeOrderDeduct)
	if err != nil {
		return err
	}

	return nil
}

func UpdateAnssyPartner(c echo.Context, id uint, req types.AnssyAuthNotifyReq) error {
	var expiredAt time.Time
	if req.ExpiresIn > 0 {
		expiredAt = timex.GetPRCNowTime().Carbon2Time().Add(time.Duration(req.ExpiresIn) * time.Second)
	}

	_, err := repository.Partner.Update(c, id, &v1.PartnerUpdateReq{
		AnssyToken:      req.AccessToken,
		AnssyExpiredAt:  expiredAt,
		AnssyTbUserId:   req.TBUserId,
		AnssyTbUserNick: req.TBUserNick,
	})
	if err != nil {
		return err
	}

	return nil
}

// 更新商户余额
func UpdateMerchantBalance(c contextx.Context, db *gorm.DB, o *model.Order) error {
	c.Logger().Infof("OrderId=%s, Update Merchant Balance=%.2f", o.OrderId, -o.Amount)
	if o.Amount == 0 {
		return nil
	}

	err := repository.Merchant.UpdateBalance(db, o.MerchantId, o.OrderId, -o.Amount, model.BalanceFromTypeOrderDeduct)
	if err != nil {
		return err
	}

	return nil
}

// 通知商户
func NotifyMerchant(c contextx.Context, db *gorm.DB, o *model.Order) error {
	if o == nil || len(o.NotifyUrl) == 0 {
		return nil
	}

	merchant, err := FindMerchant(db, int32(o.MerchantId))
	if err != nil {
		return err
	}

	notifyData := types.NotifyData{
		MerchantId:      int32(o.MerchantId),
		MerchantTradeNo: o.MerchantOrderId,
		TradeNo:         o.OrderId,
		Amount:          util.ToDecimal(o.Amount),
		ActualAmount:    util.ToDecimal(o.ReceivedAmount),
		Timestamp:       cast.ToString(time.Now().Unix()),
	}
	sign := NewSign(merchant.PrivateKey)
	notifyData.Sign = sign.Generate(c, notifyData)

	c.Logger().Infof("NotifyUrl=%s, NotifyData=%+v", o.NotifyUrl, notifyData)

	var result string
	client := resty.New()
	resp, err := client.SetTimeout(10 * time.Second).R().SetBody(notifyData).SetResult(&result).Post(o.NotifyUrl)
	if err != nil {
		c.Logger().Infof("NotifyUrl=%s, NotifyData=%+v, error=%s", o.NotifyUrl, notifyData, err)
		AddFailedOrderNotify(c, o.OrderId)
		return err
	}
	rawBody := resp.Body()
	c.Logger().Infof("response orderId=%s, code=%d,  body=%s", o.OrderId, resp.StatusCode(), string(rawBody))

	if resp.StatusCode() != http.StatusOK {
		AddFailedOrderNotify(c, o.OrderId)
		return fmt.Errorf("通知商户(%d)失败, 商户订单(%s), http.StatusCode=%d", o.MerchantId, o.MerchantOrderId, resp.StatusCode())
	}
	if strings.ToLower(string(rawBody)) != types.Success {
		AddFailedOrderNotify(c, o.OrderId)
		return fmt.Errorf("通知商户(%d)失败, 商户订单(%s), 商户回复 result=%s", o.MerchantId, o.MerchantOrderId, string(rawBody))
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		o, err := repository.Order.GetByOrderId(tx, o.OrderId)
		if err != nil {
			return err
		}

		c.Logger().Infof("orderId=%s, order status=%d", o.OrderId, o.Status)

		err = repository.Order.Update(tx, o.OrderId, model.Order{
			Status:       model.OrderStatusFinish,
			NotifyStatus: model.NotifyDone,
			NotifyAt:     time.Now(),
		})
		if err != nil {
			return err
		}

		if o.Status == model.OrderStatusFinish || o.Status == model.OrderStatusRefundSuccessful {
			c.Logger().Infof("该订单已处理, order=%s, status=%d", o.OrderId, o.Status)
			return nil
		}

		err = UpdateMerchantBalance(c, tx, o)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	c.Logger().Infof("通知商户成功, NotifyData=%+v", notifyData)
	return nil
}

func AddFailedOrderNotify(c contextx.Context, orderId string) {
	go func() {
		total, err := repository.Notify.CountOfToday(c, orderId, model.NotifyBizTypeOrder)
		if err != nil {
			c.Logger().Errorf("FailedNotify CountOfToday error=%s", err)
			return
		}

		if total >= 3 {
			c.Logger().Warnf("FailedNotify CountOfToday total=%d", total)
			return
		}

		now := time.Now()
		_, err = repository.Notify.Create(c, &model.Notify{
			BizId:     orderId,
			BizType:   model.NotifyBizTypeOrder,
			ExpiredAt: now.Add(time.Duration(2) * time.Minute),
		})
		if err != nil {
			c.Logger().Errorf("AddFailedNotify error=%s", err)
		}
	}()
}

// 查找ck
func FindCK(db *gorm.DB, masterId uint, limit int) ([]*model.JDAccount, error) {
	jdAccount := make([]*model.JDAccount, 0)
	err := db.Where("master_id = ? AND enable = ?", masterId, model.Enabled).Limit(limit).Order("weight").Find(&jdAccount).Error
	if err != nil {
		return nil, err
	}
	return jdAccount, nil
}

// 更新ck状态
func UpdateCKStatus(db *gorm.DB, id uint, u map[string]any) error {
	err := db.Model(&model.JDAccount{}).Where("id = ?", id).Updates(u).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	return nil
}

func UpdateCKUseCount(db *gorm.DB, id uint) error {
	err := repository.JDAccount.UseCount(id)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	return nil
}

// 查找地址
func FindAddress(db *gorm.DB, masterId uint) (*model.RealNameAccount, error) {
	realNameAccount := model.RealNameAccount{}
	err := db.Where("master_id = ? AND enable = ?", masterId, model.Enabled).First(&realNameAccount).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("没有可用的地址")
		}
		return nil, err
	}
	return &realNameAccount, nil
}
