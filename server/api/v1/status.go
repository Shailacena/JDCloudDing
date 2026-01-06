package v1

// 服务器与数据库状态
type GetServerStatusReq struct{}

type GetServerStatusResp struct {
	DbVersion          string  `json:"dbVersion"`
	OrderTotal         int64   `json:"orderTotal"`
	OrderToday         int64   `json:"orderToday"`
	SuccessTodayCount  int64   `json:"successTodayCount"`
	SuccessTodayAmount float64 `json:"successTodayAmount"`
	PartnerTotal       int64   `json:"partnerTotal"`
	MerchantTotal      int64   `json:"merchantTotal"`
}

// 当日订单趋势（5分钟间隔）
type GetTodayTrendReq struct {
	PartnerId  uint `query:"partnerId"`
	MerchantId uint `query:"merchantId"`
}

type TrendPoint struct {
	Time          string  `json:"time"` // HH:mm
	OrderCount    int64   `json:"orderCount"`
	OrderAmount   float64 `json:"orderAmount"`
	SuccessCount  int64   `json:"successCount"`
	SuccessAmount float64 `json:"successAmount"`
}

type GetTodayTrendResp struct {
	Points []*TrendPoint `json:"points"`
}
