package v1

// 订单列表
type ListOrderReq struct {
	PartnerId       uint   `query:"partnerId"`
	MerchantId      uint   `query:"merchantId"`
	OrderId         string `query:"orderId"`
	PartnerOrderId  string `query:"partnerOrderId"`
	MerchantOrderId string `query:"merchantOrderId"`
	PayType         int    `query:"payType"`
	OrderStatus     uint   `query:"orderStatus"`
	NotifyStatus    uint   `query:"notifyStatus"`
	StartAt         string `query:"startAt"`
	EndAt           string `query:"endAt"`
	Pagination
}

type ListOrderResp struct {
	ListTableData[Order]
}

type Order struct {
	OrderId         string  `json:"orderId"`
	MerchantId      uint    `json:"merchantId"`
	PartnerId       uint    `json:"partnerId"`
	MerchantOrderId string  `json:"merchantOrderId"`
	PartnerOrderId  string  `json:"partnerOrderId"`
	Amount          float64 `json:"amount"`
	ReceivedAmount  float64 `json:"receivedAmount"`
	ChannelId       string  `json:"channelId"`
	PayType         int     `json:"payType"`
	PayAccount      string  `json:"payAccount"`
	Status          uint    `json:"status"`
	SkuId           string  `json:"skuId"`
	Shop            string  `json:"shop"`
	NotifyStatus    uint    `json:"notifyStatus"`
	IP              string  `json:"ip"`
	Device          string  `json:"device"`
	Remark          string  `json:"remark"`
	CreateAt        int64   `json:"createAt"`
	PayAt           int64   `json:"payAt"`
	MerchantName    string  `json:"merchantName"`
	PartnerName     string  `json:"partnerName"`
}

type ConfirmOrderReq struct {
	OrderId string `json:"orderId"`
}

type ConfirmOrderResp struct {
}

// 订单汇总数据
type GetOrderSummaryReq struct {
}

type GetOrderSummaryesp struct {
    TotalAmount float64 `json:"totalAmount"`
}

type ArchiveOrdersReq struct {
    AdminId uint `json:"adminId"`
}

type ArchiveOrdersResp struct {
    ArchiveDate string  `json:"archiveDate"`
    TotalAmount float64 `json:"totalAmount"`
    OrderCount  int64   `json:"orderCount"`
}
