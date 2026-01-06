package model

import (
	"regexp"
)

type SysUserRole int

const (
	SuperAdminRole SysUserRole = iota + 1
	NormalAdminRole
	ClonedAdminRole
	AgencyAdminRole
)

type EnableStatus int

const (
	Enabled EnableStatus = iota + 1
	Disabled
)

type OnlineStatus int

const (
	Online OnlineStatus = iota + 1
	Offline
)

type BalanceFromType int

const (
	BalanceFromTypeOrderAdd     BalanceFromType = iota + 1 // 订单收入
	BalanceFromTypeOrderDeduct                             // 订单扣减
	BalanceFromTypeSystemAdd                               // 平台增加
	BalanceFromTypeSystemDeduct                            // 平台扣减
)

type OrderStatus int

const (
	OrderStatusUnpaid           OrderStatus = iota + 1 // 待支付
	OrderStatusPaid                                    // 已支付
	OrderStatusFinish                                  // 已完成
	OrderStatusRefundSuccessful                        // 退款成功
	OrderStatusRefundFailed                            // 退款失败
)

type NotifyStatus int

const (
	NotNotify  = iota + 1 // 未通知
	NotifyDone            // 已通知
)

type PayType int

const (
	AppPay PayType = iota + 1
	WebPay
	AliPay
)

type ChannelId string

const (
	ChannelJDGame          ChannelId = "JD000000" // 京东游戏
	ChannelJDEntity        ChannelId = "JD111111" // 京东实物
	ChannelTBShop          ChannelId = "TB000000" // 淘宝店铺
	ChannelTBQRCode        ChannelId = "TB111111" // 淘宝码上收
	ChannelTBClipboard     ChannelId = "TB222222" // 淘宝复制
	ChannelTBPayForAnother ChannelId = "TB668888" // 淘宝代付
	ChannelTBPay           ChannelId = "TB686088" // 天猫直付
	ChannelTBECoupon       ChannelId = "TB087888" // 淘宝电子券
	ChannelJDPay           ChannelId = "JS000000" // 京东复制
	ChannelJDCk            ChannelId = "JS111111" // 京东ck
)

type ConfirmStatus int

const (
	ConfirmStatusAuto   ConfirmStatus = iota + 1 // 自动
	ConfirmStatusManual                          // 手动
)

type GoodsStatus int

const (
	GoodsStatusEnabled GoodsStatus = iota + 1
	GoodsStatusDisabled
)

type PartnerType int

const (
	PartnerTypeAgiso PartnerType = iota + 1
	PartnerTypeAnssy
)

type NotifyBizType int

const (
	NotifyBizTypeOrder = iota + 1
)

var OperationMap = map[string]string{
	"admin/login":           "管理员登录",
	"admin/register":        "管理员注册",
	"admin/list":            "查看管理员列表",
	"admin/setPassword":     "修改管理员密码",
	"admin/resetPassword":   "重置管理员密码",
	"admin/delete":          "删除管理员",
	"admin/update":          "修改管理员信息",
	"admin/enable":          "修改管理员状态",
	"admin/resetVerifiCode": "重置管理员验证码",
	"admin/logout":          "管理员退出",

	"partner/register":        "合作商注册",
	"partner/list":            "查看合作商列表",
	"partner/resetPassword":   "重置合作商密码",
	"partner/delete":          "删除合作商",
	"partner/update":          "修改合作商信息",
	"partner/resetVerifiCode": "重置合作商验证码",
	"partner/updateBalance":   "调整合作商余额",
	"partner/listBalanceBill": "查看合作商余额列表",
	"partner/syncGoods":       "同步合作商商品",

	"merchant/register":        "商户注册",
	"merchant/list":            "查看商户列表",
	"merchant/resetPassword":   "重置商户密码",
	"merchant/enable":          "修改商户状态",
	"merchant/update":          "修改商户信息",
	"merchant/resetVerifiCode": "重置商户验证码",
	"merchant/updateBalance":   "调整商户余额",
	"merchant/listBalanceBill": "查看商户余额列表",

	"realNameAccount/create": "新建实名账号",
	"realNameAccount/list":   "查看实名账号列表",

	"jdAccount/create": "新建京东账号",
	"jdAccount/list":   "查看京东账号列表",
	"jdAccount/enable": "修改京东账号状态",
	"jdAccount/delete": "删除京东账号",

	"statistics/listDailyBill":           "查看每日流水",
	"statistics/listDailyBillByPartner":  "查看合作商每日流水",
	"statistics/listDailyBillByMerchant": "查看商户每日流水",

	"goods/create": "新建商品",
	"goods/list":   "查看商品列表",
	"goods/update": "更新商品",
	"goods/delete": "删除商品",

	"order/list":    "查看订单列表",
	"order/summary": "查看首页汇总",
	"order/confirm": "确认商品",
	"order/archive": "归档订单",
}

func GetOperationStr(uri string) string {
	var pathPart string
	re := regexp.MustCompile(`^/web_api/([^?]+)`)
	matches := re.FindStringSubmatch(uri)

	if len(matches) > 1 {
		pathPart = matches[1]
	}

	return OperationMap[pathPart]
}

type JDAccountStatus int

const (
	JDAccountStatusNormal JDAccountStatus = iota + 1
	JDAccountStatusInvalid
	JDAccountStatusHot
	JDAccountStatusAddAddressErr
	JDAccountStatusSubmitOrderErr
	JDAccountStatusGetWxPayErr
)
