package v1

// jd账号创建
type JDAccountCreateReq struct {
	AccountList []BaseJDAccount `json:"accountList" validate:"required"`
	Remark      string          `json:"remark"`
}

type BaseJDAccount struct {
	Account string `json:"account"`
	WsKey   string `json:"wsKey"`
}

type JDAccountCreateResp struct {
}

// 启用或禁用
type JDAccountEnableReq struct {
	Id     uint `json:"id" validate:"required"`
	Status int  `json:"status" validate:"required"`
}

type JDAccountEnableResp struct {
}

// jd账号列表
type ListJDAccountReq struct {
	JDAccountSearchParams
	Pagination
}

type JDAccountSearchParams struct {
	Id      uint   `query:"id" json:"id"`
	Account string `query:"account" json:"account"`
	Status  []int  `query:"status" json:"status"`
}

type ListJDAccountResp struct {
	ListTableData[JDAccount]
}

type JDAccount struct {
	Id                     uint   `json:"id"`
	Account                string `json:"account"`
	RealNameStatus         int    `json:"realNameStatus"`
	TotalOrderCount        int    `json:"totalOrderCount"`
	TodayOrderCount        int    `json:"todayOrderCount"`
	TotalSuccessOrderCount int    `json:"totalSuccessOrderCount"`
	OnlineStatus           int    `json:"onlineStatus"`
	Status                 int    `json:"status"`
	Remark                 string `json:"remark"`
	CreateAt               int64  `json:"createAt"`
	UpdateAt               int64  `json:"updateAt"`
}

// 删除
type JDAccountDeleteReq struct {
	JDAccountSearchParams
	IsAll bool `json:"isAll"`
}

type JDAccountDeleteResp struct {
}

// 重置异常状态
type JDAccountResetReq struct {
	JDAccountSearchParams
}

type JDAccountResetResp struct {
}
