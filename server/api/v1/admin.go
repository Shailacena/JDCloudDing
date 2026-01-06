package v1

// 管理员登录
type AdminLoginReq struct {
	Username   string `json:"username" validate:"required"`
	Password   string `json:"password" validate:"required"`
	VerifiCode string `json:"verifiCode" validate:"required"`
}

type AdminLoginResp struct {
	Id       uint   `json:"id"`
	Token    string `json:"token"`
	Nickname string `json:"nickname"`
	Role     int    `json:"role"`
}

// 管理员登出
type AdminLogoutReq struct {
}

type AdminLogoutResp struct {
}

// 管理员注册
type AdminRegister11Req struct {
	AdminRegisterReq
	S string `json:"s" validate:"required"`
}

// 管理员注册
type AdminRegisterReq struct {
	Username string `json:"username" validate:"required"`
	Nickname string `json:"nickname" validate:"required"`
	Role     int    `json:"role" validate:"required"`
	Remark   string `json:"remark"`
}

type AdminRegisterResp struct {
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Password string `json:"password"`
}

// 管理员列表
type ListAdminReq struct {
}

type ListAdminResp struct {
	ListTableData[Admin]
}

// 密码修改
type AdminSetPasswordReq struct {
	OldPassword string `json:"oldPassword" validate:"required"`
	NewPassword string `json:"newPassword" validate:"required"`
}

type AdminSetPasswordResp struct {
}

// 密码重置
type AdminResetPasswordReq struct {
	Username string `json:"username" validate:"required"`
}

type AdminResetPasswordResp struct {
	Password string `json:"password"`
}

// 删除管理员
type AdminDeleteReq struct {
	Username string `json:"username" validate:"required"`
}

type AdminDeleteResp struct {
}

// 更新信息
type AdminUpdateReq struct {
	Username string `json:"username" validate:"required"`
	Nickname string `json:"nickname" validate:"required"`
	Remark   string `json:"remark"`
}

type AdminUpdateResp struct {
}

// 启用或禁用
type AdminEnableReq struct {
	Username string `json:"username" validate:"required"`
	Enable   int    `json:"enable"`
}

type AdminEnableResp struct {
	Enable int `json:"enable"`
}

type Admin struct {
	Id       uint   `json:"id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Remark   string `json:"remark"`
	Enable   int    `json:"enable"`
	Role     int    `json:"role"`
	UrlKey   string `json:"urlKey"`
	ParentId uint   `json:"parentId"`
}

// 重置验证码
type AdminResetVerifiCodeReq struct {
	Id uint `json:"id" validate:"required"`
}

type AdminResetVerifiCodeResp struct {
	UrlKey string `json:"urlKey"`
}

// 获取主账号总收入
type GetMasterIncomeReq struct {
	MasterId uint `json:"masterId" validate:"required"`
}

type GetMasterIncomeResp struct {
	TotalIncome float64 `json:"totalIncome"`
}
