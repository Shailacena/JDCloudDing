package v1

// 商户登录
type MerchantLoginReq struct {
	Username   string `json:"username" validate:"required"`
	Password   string `json:"password" validate:"required"`
	VerifiCode string `json:"verifiCode" validate:"required"`
}

type MerchantLoginResp struct {
	Id       uint   `json:"id"`
	Token    string `json:"token"`
	Nickname string `json:"nickname"`
}

// 商户登出
type MerchantLogoutReq struct {
}

type MerchantLogoutResp struct {
}

// 商户注册
type MerchantRegisterReq struct {
	Nickname string `json:"nickname" validate:"required"`
	Remark   string `json:"remark"`
}

type MerchantRegisterResp struct {
	Nickname string `json:"nickname"`
	Password string `json:"password"`
}

// 商户修改
type MerchantUpdateReq struct {
	Username string `json:"username" validate:"required"`
	Nickname string `json:"nickname" validate:"required"`
	Remark   string `json:"remark"`
	IsDel    bool   `json:"isDel"`
}

type MerchantUpdateResp struct {
}

// 商户修改余额
type MerchantUpdateBalanceReq struct {
	AdminId      uint    `json:"adminId" validate:"required"`
	MerchantId   uint    `json:"merchantId" validate:"required"`
	ChangeAmount float64 `json:"changeAmount" validate:"required"`
	Password     string  `json:"password" validate:"required"`
}

type MerchantUpdateBalanceResp struct {
}

type MerchantEnableReq struct {
	Username string `json:"username" validate:"required"`
	Enable   int    `json:"enable"`
}

type MerchantEnableResp struct {
	Enable int `json:"enable"`
}

// 商户列表
type ListMerchantReq struct {
	MerchantId       uint `query:"merchantId"`
	IgnoreStatistics bool `query:"ignoreStatistics"`
	Pagination
}

type ListMerchantResp struct {
	ListTableData[Merchant]
}

// 商户修改密码
type MerchantSetPasswordReq struct {
	// Username    string `json:"username" validate:"required"`
	OldPassword string `json:"oldpassword" validate:"required"`
	NewPassword string `json:"newpassword" validate:"required"`
}

type MerchantSetPasswordResp struct {
	Token    string `json:"token"`
	Nickname string `json:"nickname"`
}

// 密码重置
type MerchantResetPasswordReq struct {
	Id uint `json:"id" validate:"required"`
}

type MerchantResetPasswordResp struct {
	Password string `json:"password"`
}

// 重置验证码
type MerchantResetVerifiCodeReq struct {
	Id uint `json:"id" validate:"required"`
}

type MerchantResetVerifiCodeResp struct {
	UrlKey string `json:"urlKey"`
}

type Merchant struct {
	Id          uint    `json:"id"`
	Username    string  `json:"username"`
	Nickname    string  `json:"nickname"`
	PrivateKey  string  `json:"privateKey"`
	CreateAt    int64   `json:"createAt"`
	TotalAmount float64 `json:"totalAmount"`
	TodayAmount float64 `json:"todayAmount"`
	Enable      int     `json:"enable"`
	Balance     float64 `json:"balance"`
	Remark      string  `json:"remark"`
	UrlKey      string  `json:"urlKey"`
	ParentId    uint    `json:"parentId"`
}

// 商户余额账单
type ListMerchantBalanceBillReq struct {
	MerchantId uint   `query:"merchantId"`
	StartAt    string `query:"startAt"`
	EndAt      string `query:"endAt"`
	Pagination
}

type ListMerchantBalanceBillResp struct {
	ListTableData[MerchantBalanceBill]
}

type MerchantBalanceBill struct {
	Id           uint    `json:"id"`
	MerchantId   uint    `json:"merchantId"`
	Nickname     string  `json:"nickname"`
	OrderId      string  `json:"orderId"`
	From         int     `json:"from"`
	Balance      float64 `json:"balance"`
	ChangeAmount float64 `json:"changeAmount"`
	CreateAt     int64   `json:"createAt"`
}
type MerchantBalanceReq struct {
}

type MerchantBalanceResp struct {
	Balance float64 `json:"balance"`
}
