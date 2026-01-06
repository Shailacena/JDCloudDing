package types

// 平台推送给商户的数据
type NotifyData struct {
	MerchantId      int32   `json:"merchantId"`
	MerchantTradeNo string  `json:"merchantTradeNo"`
	TradeNo         string  `json:"tradeNo"`
	Amount          float64 `json:"amount"`
	ActualAmount    float64 `json:"actualAmount"`
	Timestamp       string  `json:"timestamp"`
	Sign            string  `json:"sign"`
}

// 阿奇索推送
type AgisoNotifyReq struct {
	FromPlatform string          `json:"fromPlatform"` // 平台
	Timestamp    int64           `json:"timestamp"`    // 时间戳
	Aopic        int32           `json:"aopic"`        // 推送类型
	Json         AgisoNotifyJson `json:"json"`         // 推送数据
	Sign         string          `json:"sign"`         // 签名
}

type AgisoNotifyJson struct {
	Tid           string `json:"tid"`           // 订单编号
	Status        string `json:"status"`        // 订单状态
	SellerNick    string `json:"sellerNick"`    // 卖家昵称
	SellerOpenUid string `json:"sellerOpenUid"` // 卖家ID
	BuyerNick     string `json:"buyerNick"`     // 买家昵称
	BuyerOpenUid  string `json:"buyerOpenUid"`  // 买家ID
	Payment       string `json:"payment"`       // 支付金额
	Type          string `json:"type"`          // 交易类型
}

// 淘宝
type TBJsonData struct {
	Platform                  string            `json:"Platform"`
	PlatformUserID            string            `json:"PlatformUserId"`
	ReceiverName              string            `json:"ReceiverName"`
	ReceiverMobile            string            `json:"ReceiverMobile"`
	ReceiverPhone             string            `json:"ReceiverPhone"`
	ReceiverAddress           string            `json:"ReceiverAddress"`
	BuyerArea                 string            `json:"BuyerArea"`
	SellerOpenUID             string            `json:"SellerOpenUid"`
	Tid                       int64             `json:"Tid"`
	TidStr                    string            `json:"TidStr"`
	Status                    string            `json:"Status"`
	SellerNick                string            `json:"SellerNick"`
	BuyerNick                 string            `json:"BuyerNick"`
	BuyerOpenUID              string            `json:"BuyerOpenUid"`
	Type                      string            `json:"Type"`
	BuyerMessage              string            `json:"BuyerMessage"`
	Price                     string            `json:"Price"`
	Num                       int32             `json:"Num"`
	TotalFee                  string            `json:"TotalFee"`
	Payment                   string            `json:"Payment"`
	PayTime                   string            `json:"PayTime"`
	PicPath                   string            `json:"PicPath"`
	PostFee                   string            `json:"PostFee"`
	Created                   string            `json:"Created"`
	TradeFrom                 string            `json:"TradeFrom"`
	Orders                    []TBOrderJsonData `json:"Orders"`
	SellerMemo                string            `json:"SellerMemo"`
	SellerFlag                int               `json:"SellerFlag"`
	CreditCardFee             string            `json:"CreditCardFee"`
	ExpandCardExpandPriceUsed string            `json:"ExpandCardExpandPriceUsed"`
}

type TBOrderJsonData struct {
	Oid                               int64  `json:"Oid"`
	OidStr                            string `json:"OidStr"`
	NumIid                            int64  `json:"NumIid"`
	OuterIid                          string `json:"OuterIid"`
	OuterSkuID                        string `json:"OuterSkuId"`
	Title                             string `json:"Title"`
	Price                             string `json:"Price"`
	Num                               int    `json:"Num"`
	TotalFee                          string `json:"TotalFee"`
	Payment                           string `json:"Payment"`
	PicPath                           string `json:"PicPath"`
	SkuID                             string `json:"SkuId"`
	SkuPropertiesName                 string `json:"SkuPropertiesName"`
	DivideOrderFee                    string `json:"DivideOrderFee"`
	PartMjzDiscount                   string `json:"PartMjzDiscount"`
	ExpandCardExpandPriceUsedSuborder any    `json:"ExpandCardExpandPriceUsedSuborder"`
	Customization                     any    `json:"Customization"`
}

// 京东
type JDJsonData struct {
	PlatformShopId  int64          `json:"PlatformShopId"`
	CustomerId      int64          `json:"CustomerId"`
	OrderId         int64          `json:"OrderId"`
	OrderType       int            `json:"OrderType"`
	Pin             string         `json:"Pin"`
	BuyNum          int            `json:"BuyNum"`
	SkuId           int64          `json:"SkuId"`
	BrandId         int            `json:"BrandId"`
	UserIp          string         `json:"UserIp"`
	TotalPrice      float64        `json:"TotalPrice"`
	CreateTime      string         `json:"CreateTime"`
	Features        any            `json:"Features"`
	SourceType      int            `json:"SourceType"`
	FacePrice       int            `json:"FacePrice"`
	GameAccount     string         `json:"GameAccount"`
	Permit          string         `json:"Permit"`
	GameAccountType JDGameJsonData `json:"GameAccountType"`
	ChargeType      JDGameJsonData `json:"ChargeType"`
	GameArea        JDGameJsonData `json:"GameArea"`
	GameServer      JDGameJsonData `json:"GameServer"`
}

type JDGameJsonData struct {
	Id   string `json:"Id"`
	Name string `json:"Name"`
}

type AnssyAuthNotifyReq struct {
	State       string `query:"state"`            // 合作商id
	AccessToken string `query:"access_token"`     // 授权token
	ExpiresIn   int64  `query:"expires_in"`       // token有效期
	TBUserId    string `query:"taobao_user_id"`   // 淘宝账户id
	TBUserNick  string `query:"taobao_user_nick"` // 店铺昵称
}

type AnssyNotifyReq struct {
	Id          int32  `json:"id"`               // 合作商id
	AccessToken string `json:"access_token"`     // 授权token
	ExpiresIn   int64  `json:"expires_in"`       // token有效期
	TBUserId    int32  `json:"taobao_user_id"`   // 淘宝账户id
	TBUserNick  string `json:"taobao_user_nick"` // 店铺昵称
	Sign        string `json:"sign"`             // 签名
}
