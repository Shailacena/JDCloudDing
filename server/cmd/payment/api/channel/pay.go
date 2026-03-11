package channel

import (
	"apollo/server/cmd/payment/api/iface"
	"apollo/server/cmd/payment/pkg/errorx"
	"apollo/server/cmd/payment/pkg/types"
	"apollo/server/internal/model"

	"gorm.io/gorm"
)

func GetPayByChannelId(channelId string, merchantId int32) (iface.IOrder, error) {
	var (
		instance iface.IOrder
		err      error
	)

	switch model.ChannelId(channelId) {
	case model.ChannelTBPay:
		instance, err = NewTBPay(merchantId)
	case model.ChannelJDPay:
		instance, err = NewJDPay(merchantId)
	case model.ChannelJDCk:
		instance, err = NewJDCKPay(merchantId)
	case model.ChannelJDCard:
		instance, err = NewJDCard(merchantId)
	default:
		err = errorx.ErrInvalidChannelId
	}
	if err != nil {
		return nil, err
	}

	return instance, nil
}

func GetNotifyByPlatform(fromPlatform, timestamp, aopic, sign, jsonStr string) (iface.IHander, error) {
	var (
		instance iface.IHander
		err      error
	)

	switch fromPlatform {
	case types.TbAlds:
		instance, err = NewTBPayNotify(fromPlatform, timestamp, aopic, sign, jsonStr)
	case types.AldsJd:
		instance, err = NewJDPayNotify(fromPlatform, timestamp, aopic, sign, jsonStr)
	default:
		err = errorx.ErrInvalidPlatform
	}
	if err != nil {
		return nil, err
	}

	return instance, nil
}

func GetJDCardNotify(db *gorm.DB, orderId int64, skuId string, orderCreateTime string) *JDCardNotify {
	return NewJDCardNotify(db, orderId, skuId, orderCreateTime)
}
