package v1

// 实名账号创建
type RealNameAccountCreateReq struct {
	AccountList []*BaseRealNameAccount `json:"accountList" validate:"required"`
	Remark      string                 `json:"remark"`
}

type RealNameAccountCreateResp struct {
}

// 实名账号列表
type ListRealNameAccountReq struct {
	Pagination
}

type ListRealNameAccountResp struct {
	ListTableData[RealNameAccount]
}

type BaseRealNameAccount struct {
	IdNumber string `json:"idNumber"`
	Name     string `json:"name"`
	Mobile   string `json:"mobile"`
	Address  string `json:"address"`
}

type RealNameAccount struct {
	BaseRealNameAccount
	RealNameCount int64  `json:"realNameCount"`
	Enable        int    `json:"enable"`
	Remark        string `json:"remark"`
}
