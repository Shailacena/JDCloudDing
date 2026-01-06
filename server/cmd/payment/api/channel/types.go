package channel

type ValidCookieResp struct {
	Code    int    `json:"code"`
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    string `json:"data"`
}

type AddAddressResp struct {
	Code                      string     `json:"code"`
	Flag                      bool       `json:"Flag"`
	AddressCheckStructUpgrade string     `json:"addressCheckStructUpgrade"`
	AddAddress                AddAddress `json:"addAddress"`
}

type AddAddress struct {
	AddressUUID string `json:"addressUUID"`
	Flag        bool   `json:"Flag"`
	AddressId   bool   `json:"AddressId"`
}

type SubmitOrderResp struct {
	SubmitOrderExtMap         SubmitOrderExtMap  `json:"submitOrderExtMap"`
	CoMsg                     string             `json:"coMsg"`
	Code                      string             `json:"code"`
	SyncCartNoResponse        int                `json:"syncCartNoResponse"`
	JSONParams                JSONParams         `json:"jsonParams"`
	SubmitOrder               SubmitOrder        `json:"submitOrder"`
	SubmitOrderExtFlag        SubmitOrderExtFlag `json:"submitOrderExtFlag"`
	SwitchControl             SwitchControl      `json:"switchControl"`
	OnlinePay                 bool               `json:"onlinePay"`
	UnableUpdateGlobalAddress bool               `json:"unableUpdateGlobalAddress"`
	Status                    string             `json:"status"`
	Message                   string             `json:"message"`
}
type SubmitOrderExtMap struct {
}
type JSONParams struct {
	CheckoutSource string `json:"checkoutSource"`
	Source         string `json:"source"`
}
type CashierPaymentParam struct {
	OrderType     string `json:"orderType"`
	OrderTypeCode string `json:"orderTypeCode"`
	OrderID       string `json:"orderId"`
	OrderPrice    string `json:"orderPrice"`
}
type SubmitOrder struct {
	IsPayForMeSubmit       bool                `json:"isPayForMeSubmit"`
	OrderTypeCode          string              `json:"orderTypeCode"`
	Message                string              `json:"Message"`
	JpsStatus              int                 `json:"jpsStatus"`
	SubmitSkuNum           int                 `json:"SubmitSkuNum"`
	Source                 string              `json:"source"`
	SubmitWithWmCard       string              `json:"submitWithWmCard"`
	UseBalance             int                 `json:"UseBalance"`
	OrderType              int                 `json:"OrderType"`
	IsZeroOrder            bool                `json:"isZeroOrder"`
	JumpCashierPayPage     bool                `json:"jumpCashierPayPage"`
	IDCompanyBranch        int                 `json:"IdCompanyBranch"`
	JpsNum                 string              `json:"jpsNum"`
	MessageType            int                 `json:"MessageType"`
	UseScore               int                 `json:"UseScore"`
	CashierPaymentParam    CashierPaymentParam `json:"cashierPaymentParam"`
	OrderID                int64               `json:"OrderId"`
	Flag                   bool                `json:"Flag"`
	TradeResultCode        string              `json:"TradeResultCode"`
	Price                  float64             `json:"Price"`
	Appid                  string              `json:"appid"`
	FactPrice              float64             `json:"FactPrice"`
	IsFriendPayForMeSubmit bool                `json:"isFriendPayForMeSubmit"`
	ErrMessage             string              `json:"errMessage"`
	IDPaymentType          int                 `json:"IdPaymentType"`
}
type SubmitOrderExtFlag struct {
	PayPath                 int    `json:"payPath"`
	OptimalPlan             string `json:"optimalPlan"`
	OptimalPlanFeeFlag      string `json:"optimalPlanFeeFlag"`
	IsGoodsDetailJinCaiFlag string `json:"isGoodsDetailJinCaiFlag"`
}
type SwitchControl struct {
	TraceSubmitOrderPathSwitch string `json:"traceSubmitOrderPathSwitch"`
}

type AppOrderWxPayResp struct {
	Code    int    `json:"code"`
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    string `json:"data"`
}
