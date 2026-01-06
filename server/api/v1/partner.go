package v1

import (
	"apollo/server/internal/model"
	"time"
)

// 合作商登录
type PartnerLoginReq struct {
	Username   string `json:"username" validate:"required"`
	Password   string `json:"password" validate:"required"`
	VerifiCode string `json:"verifiCode" validate:"required"`
}

type PartnerLoginResp struct {
	Id       uint   `json:"id"`
	Token    string `json:"token"`
	Level    int    `json:"level"`
	Nickname string `json:"nickname"`
}

// 合作商登出
type PartnerLogoutReq struct {
}

type PartnerLogoutResp struct {
}

// 合作商注册
type PartnerRegisterReq struct {
	Nickname          string            `json:"nickname" validate:"required"`
	Priority          int               `json:"priority" validate:"required"`
	AqsAppSecret      string            `json:"aqsAppSecret"`
	AqsToken          string            `json:"aqsToken"`
	PayType           model.PayType     `json:"payType" validate:"required"`
	ChannelId         model.ChannelId   `json:"channelId" validate:"required"`
	Level             int               `json:"level"`
	Remark            string            `json:"remark"`
	Type              model.PartnerType `json:"type" validate:"required"`
	DarkNumberLength  int               `json:"darkNumberLength" validate:"min=8,max=15"`
}

type PartnerRegisterResp struct {
	Nickname string `json:"nickname"`
	Password string `json:"password"`
}

// 合作商更新
type PartnerUpdateReq struct {
	Id               uint      `json:"id" validate:"required"`
	Nickname         string    `json:"nickname" validate:"required"`
	Priority         int       `json:"priority" validate:"required"`
	AqsAppSecret     string    `json:"aqsAppSecret"`
	AqsToken         string    `json:"aqsToken"`
	AnssyAppSecret   string    `json:"anssyAppSecret"`
	AnssyToken       string    `json:"anssyToken"`
	PayType          int       `json:"payType"`
	RechargeTime     int64     `json:"rechargeTime"`
	Remark           string    `json:"remark"`
	Enable           int       `json:"enable"`
	AnssyExpiredAt   time.Time `json:"anssyExpiredAt"`
	AnssyTbUserId    string    `json:"anssyTbUserId"`
	AnssyTbUserNick  string    `json:"anssyTbUserNick"`
	Secret           string    `json:"secret"`
	UrlPath          string    `json:"urlPath"`
	DarkNumberLength int       `json:"darkNumberLength" validate:"min=8,max=15"`
}

type PartnerUpdateResp struct {
}

// 合作商修改余额
type PartnerUpdateBalanceReq struct {
	AdminId      uint    `json:"adminId" validate:"required"`
	PartnerId    uint    `json:"partnerId" validate:"required"`
	ChangeAmount float64 `json:"changeAmount" validate:"required"`
	Password     string  `json:"password" validate:"required"`
}

type PartnerUpdateBalanceResp struct {
}

// 密码重置
type PartnerResetPasswordReq struct {
	Id uint `json:"id" validate:"required"`
}

type PartnerResetPasswordResp struct {
	Password string `json:"password"`
}

// 重置验证码
type PartnerResetVerifiCodeReq struct {
	Id uint `json:"id" validate:"required"`
}

type PartnerResetVerifiCodeResp struct {
	UrlKey string `json:"urlKey"`
}

// 删除合作商
type PartnerDeleteReq struct {
	Id uint `json:"id" validate:"required"`
}

type PartnerDeleteResp struct {
}

// 合作商列表
type ListPartnerReq struct {
	PartnerId        uint `query:"partnerId"`
	IgnoreStatistics bool `query:"ignoreStatistics"`
	Pagination
}

type ListPartnerResp struct {
	ListTableData[Partner]
}

type Partner struct {
	Id               uint              `json:"id"`
	Nickname         string            `json:"nickname"`
	PayType          model.PayType     `json:"payType"`
	ChannelId        model.ChannelId   `json:"channelId"`
	Balance          float64           `json:"balance"`
	Priority         int               `json:"priority"`
	SuperiorAgent    string            `json:"superiorAgent"`
	Level            int               `json:"level"`
	StockAmount      int64             `json:"stockAmount"`
	RechargeTime     int64             `json:"rechargeTime"`
	PrivateKey       string            `json:"privateKey"`
	AqsAppSecret     string            `json:"aqsAppSecret"`
	AqsToken         string            `json:"aqsToken"`
	Enable           int               `json:"enable"`
	Remark           string            `json:"remark"`
	Type             model.PartnerType `json:"type"`
	UrlKey           string            `json:"urlKey"`
	ParentId         uint              `json:"parentId"`
	DarkNumberLength int               `json:"darkNumberLength"`

	AnssyAppSecret string `json:"anssyAppSecret"`
	AnssyToken     string `json:"anssyToken"`
	AnssyExpiredAt int64  `json:"anssyExpiredAt"`

	TodayOrderAmount     float64 `json:"todayOrderAmount"`
	TodayOrderNum        float64 `json:"todayOrderNum"`
	TodaySuccessAmount   float64 `json:"todaySuccessAmount"`
	TodaySuccessOrderNum float64 `json:"todaySuccessOrderNum"`

	Last1HourTotal       int64 `json:"last1HourTotal"`
	Last1HourSuccess     int64 `json:"last1HourSuccess"`
	Last30MinutesTotal   int64 `json:"last30MinutesTotal"`
	Last30MinutesSuccess int64 `json:"last30MinutesSuccess"`
}

// // 合作商流水账单
// type ListPartnerBillReq struct {
// }

// type ListPartnerBillResp struct {
// 	List []*PartnerBill `json:"list"`
// }

// type PartnerBill struct {
// 	PartnerId   uint   `json:"partnerId"`
// 	Type        int    `json:"type"`
// 	ChangeMoney int    `json:"changeMoney"`
// 	Money       int    `json:"money"`
// 	Remark      string `json:"remark"`
// 	CreateAt    int64  `json:"createAt"`
// }

// 密码修改
type PartnerSetPasswordReq struct {
	OldPassword string `json:"oldpassword" validate:"required"`
	NewPassword string `json:"newpassword" validate:"required"`
}

type PartnerSetPasswordResp struct {
}

// 合作商余额账单
type ListPartnerBalanceBillReq struct {
	PartnerId   uint   `query:"partnerId"`
	CurrentPage int    `query:"currentPage"`
	PageSize    int    `query:"pageSize"`
	StartAt     string `query:"startAt"`
	EndAt       string `query:"endAt"`
}

type ListPartnerBalanceBillResp struct {
	ListTableData[PartnerBalanceBill]
}

type PartnerBalanceBill struct {
	Id           uint    `json:"id"`
	PartnerId    uint    `json:"partnerId"`
	Nickname     string  `json:"nickname"`
	OrderId      string  `json:"orderId"`
	From         int     `json:"from"`
	Balance      float64 `json:"balance"`
	ChangeAmount float64 `json:"changeAmount"`
	CreateAt     int64   `json:"createAt"`
}
