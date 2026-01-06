package v1

// 每日流水账单
type ListDailyBillReq struct {
	PartnerId  uint `query:"partnerId"`
	MerchantId uint `query:"merchantId"`
}

type ListDailyBillResp struct {
	List []*BaseDailyBill `json:"list"`
}

type BaseDailyBill struct {
	Date                 string  `json:"date"`
	TotalOrderAmount     float64 `json:"totalOrderAmount"`
	TotalOrderNum        float64 `json:"totalOrderNum"`
	TotalSuccessAmount   float64 `json:"totalSuccessAmount"`
	TotalSuccessOrderNum float64 `json:"totalSuccessOrderNum"`
}

type DailyBill struct {
	BaseDailyBill
	Id       uint    `json:"id"`
	Nickname string  `json:"nickname"`
	Balance  float64 `json:"balance"`
	Time     int64   `json:"time"`
}

// 合作商每日流水账单
type ListDailyBillByPartnerReq struct {
	PartnerId uint   `query:"partnerId"`
	StartAt   string `query:"startAt"`
	EndAt     string `query:"endAt"`
}

type ListDailyBillByPartnerResp struct {
	ListTableData[DailyBill]
}

// 商户每日流水账单
type ListDailyBillByMerchantReq struct {
	MerchantId uint   `query:"merchantId"`
	StartAt    string `query:"startAt"`
	EndAt      string `query:"endAt"`
}

type ListDailyBillByMerchantResp struct {
	ListTableData[DailyBill]
}
