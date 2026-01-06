package errorx

import "errors"

var (
	ErrInvalidSign       error = errors.New("非法签名")
	ErrInvalidMerchant   error = errors.New("非法商户")
	ErrInvalidPartner    error = errors.New("非法合作商")
	ErrInvalidChannelId  error = errors.New("非法通道")
	ErrInvalidPlatform   error = errors.New("非法平台")
	ErrGoodsNotFound     error = errors.New("没有可用的商品")
	ErrAutomaticShipment error = errors.New("自动发货失败")
	ErrRefundFailed      error = errors.New("退款失败")
)

var (
	ErrTxtCreateOrderInvalidParams   string = "下单失败, 参数错误"
	ErrTxtNotifySuccessInvalidParams string = "通知失败, 参数错误"
	ErrTxtQueryOrderInvalidParams    string = "查询订单失败, 参数错误"
	ErrTxtQueryBalanceInvalidParams  string = "查询余额失败, 参数错误"
)
