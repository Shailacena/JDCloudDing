package types

// 创建订单
type CreateOrderReq struct {
	ChannelId       string  `json:"channelId"`
	MerchantId      int32   `json:"merchantId"`
	MerchantTradeNo string  `json:"merchantTradeNo"`
	Amount          float64 `json:"amount"`
	NotifyUrl       string  `json:"notifyUrl"`
	Timestamp       string  `json:"timestamp"`
	Sign            string  `json:"sign"`
}

type CreateOrderResp struct {
	MerchantTradeNo string  `json:"merchantTradeNo"`
	Amount          float64 `json:"amount"`
	TradeNo         string  `json:"tradeNo"`
	PayPageUrl      string  `json:"payPageUrl"`
	Sign            string  `json:"sign"`
}

// 查询订单
type QueryOrderReq struct {
	MerchantId      int32  `json:"merchantId"`
	MerchantTradeNo string `json:"merchantTradeNo"`
	Timestamp       string `json:"timestamp"`
	Sign            string `json:"sign"`
}

type QueryOrderResp struct {
	MerchantId      int32   `json:"merchantId"`
	MerchantTradeNo string  `json:"merchantTradeNo"`
	TradeNo         string  `json:"tradeNo"`
	Amount          float64 `json:"amount"`
	ActualAmount    float64 `json:"actualAmount"`
	Status          int32   `json:"status"`
	PayAt           string  `json:"payAt,omitempty"`
	Sign            string  `json:"sign"`
}

// 查询余额
type QueryBalanceReq struct {
	MerchantId int32  `json:"merchantId"`
	Timestamp  string `json:"timestamp"`
	Sign       string `json:"sign"`
}

type QueryBalanceResp struct {
	MerchantId int32   `json:"merchantId"`
	Balance    float64 `json:"balance"`
	Sign       string  `json:"sign"`
}