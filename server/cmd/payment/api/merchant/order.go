package merchant

// import (
// 	"apollo/server/internal/model"
// 	"apollo/server/internal/repository"
// 	"apollo/server/pkg/config"
// 	"apollo/server/pkg/data"
// 	"apollo/server/pkg/response"
// 	"apollo/server/pkg/util"
// 	"crypto/md5"
// 	"encoding/hex"
// 	"encoding/json"
// 	"errors"
// 	"fmt"
// 	"net/http"
// 	"reflect"
// 	"sort"
// 	"strconv"
// 	"strings"
// 	"time"

// 	"github.com/go-resty/resty/v2"
// 	"github.com/labstack/echo/v4"
// 	"github.com/samber/lo"
// 	"github.com/spf13/cast"
// 	"golang.org/x/exp/rand"
// 	"gorm.io/gorm"
// )

// const (
// 	Success = "success"
// )

// type CreateOrderReq struct {
// 	ChannelId       string  `json:"channelId"`
// 	MerchantId      int32   `json:"merchantId"`
// 	MerchantTradeNo string  `json:"merchantTradeNo"`
// 	Amount          float64 `json:"amount"`
// 	NotifyUrl       string  `json:"notifyUrl"`
// 	Timestamp       string  `json:"timestamp"`
// 	Sign            string  `json:"sign"`
// }

// type CreateOrderResp struct {
// 	MerchantTradeNo string  `json:"merchantTradeNo"`
// 	Amount          float64 `json:"amount"`
// 	TradeNo         string  `json:"tradeNo"`
// 	PayPageUrl      string  `json:"payPageUrl"`
// 	Sign            string  `json:"sign"`
// }

// // 校验商户id、校验sign
// func CreateOrder(c echo.Context) error {
// 	req := CreateOrderReq{}
// 	err := c.Bind(&req)
// 	if err != nil {
// 		response.ResponseError(c, http.StatusBadRequest, "下单失败, 参数错误", nil)
// 		return nil
// 	}

// 	db := data.Instance()
// 	merchant, err := findMerchant(db, req.MerchantId)
// 	if err != nil {
// 		return err
// 	}

// 	privateKey := merchant.PrivateKey
// 	if !checkSign(c, privateKey, req) {
// 		return errors.New("参数错误")
// 	}

// 	goods, err := findGoods(c, req.ChannelId, req.Amount)
// 	if err != nil {
// 		return err
// 	}

// 	skuId := goods.SkuId
// 	partnerId := goods.PartnerId
// 	shop := goods.ShopName
// 	payType := goods.PayType
// 	orderId, err := genOrderId()
// 	if err != nil {
// 		return err
// 	}

// 	cd := 5
// 	if !config.IsProd() {
// 		cd = 2
// 	}
// 	o := model.Order{
// 		OrderId:         orderId,
// 		ChannelId:       model.ChannelId(req.ChannelId),
// 		MerchantId:      uint(req.MerchantId),
// 		MerchantName:    merchant.Nickname,
// 		MerchantOrderId: req.MerchantTradeNo,
// 		Amount:          util.ToDecimal(req.Amount),
// 		SkuId:           skuId,
// 		IP:              c.RealIP(),
// 		NotifyUrl:       req.NotifyUrl,
// 		EndLockAt:       time.Now().Add(time.Minute * time.Duration(cd)),
// 		Shop:            shop,
// 		PartnerId:       partnerId,
// 		PartnerName:     goods.PartnerNickname,
// 		PayType:         payType,
// 		DarkNumber:      RandomDarkNumber(),
// 	}

// 	err = repository.Order.Create(c, o)
// 	if err != nil {
// 		return err
// 	}

// 	baseUrl, err := getPayBaseUrl(partnerId, req.ChannelId)

// 	if err != nil {
// 		return err
// 	}

// 	url := fmt.Sprintf(baseUrl+"?merchantId=%s&orderId=%s&price=%f&sku=%s&time=%d&ts=%d&code=%s", o.MerchantOrderId, orderId, req.Amount, skuId, cd*60, time.Now().Unix(), o.DarkNumber)
// 	// conf := config.Get()
// 	// if config.IsProd() {
// 	// 	url = fmt.Sprintf(baseUrl + "?orderid=%s&price=%f&sku=%s&time=%d", conf.PaymentHttpConfig.Host, orderId, req.Amount, skuId, cd*60)
// 	// }

// 	resp := CreateOrderResp{
// 		MerchantTradeNo: req.MerchantTradeNo,
// 		TradeNo:         orderId,
// 		Amount:          req.Amount,
// 		PayPageUrl:      url,
// 	}
// 	sign, _ := makeSign(privateKey, resp)
// 	resp.Sign = sign
// 	response.ResponseSuccess(c, resp)
// 	return nil
// }

// type Filed struct {
// 	Name  string
// 	Value any
// }

// const (
// 	FiledSign = "sign" // 签名字段名
// )

// func getPayBaseUrl(partnerId uint, channel string) (string, error) {
// 	db := data.Instance()

// 	partner, err := findPartner(db, partnerId)
// 	if err != nil {
// 		return "", err
// 	}

// 	switch partner.PayType {
// 	case model.WebPay:
// 		return "https://takeorder-apollo.s3.ap-southeast-1.amazonaws.com/html/pay-" + channel + "-web.html", nil
// 	case model.AliPay:
// 		return "https://takeorder-apollo.s3.ap-southeast-1.amazonaws.com/html/pay-" + channel + "-alipay.html", nil
// 	case model.AppPay:
// 		return "https://takeorder-apollo.s3.ap-southeast-1.amazonaws.com/html/pay-" + channel + "-app.html", nil
// 	default:
// 		return "https://takeorder-apollo.s3.ap-southeast-1.amazonaws.com/html/pay-" + channel + "-web.html", nil
// 	}
// }

// // 参数键值对
// func findFiledOfParams(params any) ([]Filed, string) {
// 	v := reflect.ValueOf(params)
// 	t := v.Type()
// 	fields := make([]Filed, 0, v.NumField())

// 	var sign string
// 	for i := 0; i < v.NumField(); i++ {
// 		name := t.Field(i).Tag.Get("json")
// 		if name == FiledSign {
// 			sign = v.Field(i).String()
// 			continue
// 		}

// 		keys := strings.Split(name, ",")
// 		if len(keys) > 0 {
// 			name = keys[0]
// 		}

// 		fields = append(fields, Filed{
// 			Name:  name,
// 			Value: v.Field(i).Interface(),
// 		})
// 	}

// 	return fields, sign
// }

// func makeSign(privateKey string, params any) (string, string) {
// 	fields, sign := findFiledOfParams(params)

// 	sort.Slice(fields, func(i, j int) bool {
// 		return fields[i].Name < fields[j].Name
// 	})

// 	fieldList := lo.Map(fields, func(f Filed, index int) string {
// 		return fmt.Sprintf("%s=%s", f.Name, cast.ToString(f.Value))
// 	})

// 	// 加入SignKey
// 	fieldList = append(fieldList, fmt.Sprintf("key=%s", privateKey))
// 	// 组装
// 	fieldStr := strings.Join(fieldList, "&")

// 	// md5
// 	firstHash := md5.Sum([]byte(fieldStr))
// 	hashString := hex.EncodeToString(firstHash[:])
// 	hashString = strings.ToUpper(hashString)

// 	return hashString, sign
// }

// // 校验签名
// func checkSign(c echo.Context, privateKey string, params any) bool {
// 	hashString, sign := makeSign(privateKey, params)

// 	if !strings.EqualFold(hashString, sign) {
// 		c.Logger().Infof("checkSign: %s--%s", hashString, sign)
// 	}

// 	return strings.EqualFold(hashString, sign)
// }

// func findMerchant(db *gorm.DB, merchantId int32) (*model.Merchant, error) {
// 	merchant := model.Merchant{}
// 	err := db.Where("id = ? AND enable = ?", merchantId, model.Enabled).First(&merchant).Error
// 	if err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return nil, errors.New("非法商户")
// 		}
// 		return nil, err
// 	}
// 	return &merchant, nil
// }

// // 生成订单号
// func genOrderId() (string, error) {
// 	id, err := util.SFlake.GenString()
// 	if err != nil {
// 		return "", err
// 	}
// 	return id, nil
// }

// type simpleGoods struct {
// 	Id              uint
// 	SkuId           string
// 	ShopName        string
// 	PartnerId       uint
// 	PartnerNickname string
// 	PayType         model.PayType
// }

// // 查找商品
// func findGoods(c echo.Context, channelId string, amount float64) (*simpleGoods, error) {
// 	db := data.Instance()

// 	rows, err := db.Model(model.Goods{}).
// 		Select("goods.id, goods.sku_id, goods.shop_name, partner.id as partner_id, partner.nickname as partner_nickname, partner.pay_type").
// 		Joins("left join partner on partner.id = goods.partner_id").
// 		Where("partner.channel_id = ? AND amount = ?", channelId, amount).
// 		Order("partner.priority desc, goods.weight").
// 		Rows()
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		var g simpleGoods
// 		db.ScanRows(rows, &g)

// 		if len(g.SkuId) == 0 {
// 			continue
// 		}

// 		// 天猫直付，商品是否在使用中
// 		if channelId == string(model.ChannelTBPay) {
// 			order := model.Order{}
// 			err = db.Where("sku_id = ? AND end_lock_at >= ?", g.SkuId, time.Now()).Order("created_at desc").First(&order).Error
// 			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
// 				continue
// 			}
// 			if err == nil {
// 				continue
// 			}
// 		}

// 		err = db.Model(&model.Goods{}).Where("id = ?", g.Id).Update("weight", gorm.Expr("weight + ?", 1)).Error
// 		if err != nil {
// 			c.Logger().Error("update goods weight error", err.Error())
// 		}
// 		return &g, nil
// 	}

// 	return nil, errors.New("没有可用的商品")
// }

// type QueryOrderReq struct {
// 	MerchantId      int32  `json:"merchantId"`
// 	MerchantTradeNo string `json:"merchantTradeNo"`
// 	Timestamp       string `json:"timestamp"`
// 	Sign            string `json:"sign"`
// }

// type QueryOrderResp struct {
// 	MerchantId      int32   `json:"merchantId"`
// 	MerchantTradeNo string  `json:"merchantTradeNo"`
// 	TradeNo         string  `json:"tradeNo"`
// 	Amount          float64 `json:"amount"`
// 	ActualAmount    float64 `json:"actualAmount"`
// 	Status          int32   `json:"status"`
// 	PayAt           string  `json:"payAt,omitempty"`
// 	Sign            string  `json:"sign"`
// }

// func QueryOrder(c echo.Context) error {
// 	req := QueryOrderReq{}
// 	err := c.Bind(&req)
// 	if err != nil {
// 		response.ResponseError(c, http.StatusBadRequest, "查询订单失败，参数错误", nil)
// 		return nil
// 	}

// 	db := data.Instance()
// 	merchant, err := findMerchant(db, req.MerchantId)
// 	if err != nil {
// 		return err
// 	}

// 	privateKey := merchant.PrivateKey
// 	if !checkSign(c, privateKey, req) {
// 		return errors.New("参数错误")
// 	}

// 	o, err := findOrderById(c, req.MerchantTradeNo)
// 	if err != nil {
// 		return err
// 	}

// 	var payAt string
// 	if !o.PayAt.IsZero() {
// 		payAt = cast.ToString(o.PayAt.Unix())
// 	}

// 	resp := QueryOrderResp{
// 		MerchantId:      req.MerchantId,
// 		MerchantTradeNo: req.MerchantTradeNo,
// 		Amount:          util.ToDecimal(o.Amount),
// 		ActualAmount:    util.ToDecimal(o.ReceivedAmount),
// 		TradeNo:         o.OrderId,
// 		Status:          int32(o.Status),
// 		PayAt:           payAt,
// 	}
// 	sign, _ := makeSign(privateKey, resp)
// 	resp.Sign = sign
// 	response.ResponseSuccess(c, resp)
// 	return nil
// }

// func findOrderById(c echo.Context, merchantOrderId string) (*model.Order, error) {
// 	o, err := repository.Order.GetByMerchantOrderId(c, merchantOrderId)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return o, nil
// }

// type AgisoNotifyReq struct {
// 	FromPlatform string          `json:"fromPlatform"`
// 	Timestamp    int64           `json:"timestamp"`
// 	Aopic        int32           `json:"aopic"`
// 	Json         AgisoNotifyJson `json:"json"`
// 	Sign         string          `json:"sign"`
// }

// type AgisoNotifyJson struct {
// 	Tid           string `json:"tid"`           // 订单编号
// 	Status        string `json:"status"`        // 订单状态
// 	SellerNick    string `json:"sellerNick"`    // 卖家昵称
// 	SellerOpenUid string `json:"sellerOpenUid"` // 卖家ID
// 	BuyerNick     string `json:"buyerNick"`     // 买家昵称
// 	BuyerOpenUid  string `json:"buyerOpenUid"`  // 买家ID
// 	Payment       string `json:"payment"`       // 支付金额
// 	Type          string `json:"type"`          // 交易类型
// }

// type AgisoNotifyResp struct {
// }

// type JSONDataOrder struct {
// 	Oid                               int64  `json:"Oid"`
// 	OidStr                            string `json:"OidStr"`
// 	NumIid                            int64  `json:"NumIid"`
// 	OuterIid                          string `json:"OuterIid"`
// 	OuterSkuID                        string `json:"OuterSkuId"`
// 	Title                             string `json:"Title"`
// 	Price                             string `json:"Price"`
// 	Num                               int    `json:"Num"`
// 	TotalFee                          string `json:"TotalFee"`
// 	Payment                           string `json:"Payment"`
// 	PicPath                           string `json:"PicPath"`
// 	SkuID                             string `json:"SkuId"`
// 	SkuPropertiesName                 string `json:"SkuPropertiesName"`
// 	DivideOrderFee                    string `json:"DivideOrderFee"`
// 	PartMjzDiscount                   string `json:"PartMjzDiscount"`
// 	ExpandCardExpandPriceUsedSuborder any    `json:"ExpandCardExpandPriceUsedSuborder"`
// 	Customization                     any    `json:"Customization"`
// }

// type JSONData struct {
// 	Platform                  string          `json:"Platform"`
// 	PlatformUserID            string          `json:"PlatformUserId"`
// 	ReceiverName              string          `json:"ReceiverName"`
// 	ReceiverMobile            string          `json:"ReceiverMobile"`
// 	ReceiverPhone             string          `json:"ReceiverPhone"`
// 	ReceiverAddress           string          `json:"ReceiverAddress"`
// 	BuyerArea                 string          `json:"BuyerArea"`
// 	SellerOpenUID             string          `json:"SellerOpenUid"`
// 	Tid                       int64           `json:"Tid"`
// 	TidStr                    string          `json:"TidStr"`
// 	Status                    string          `json:"Status"`
// 	SellerNick                string          `json:"SellerNick"`
// 	BuyerNick                 string          `json:"BuyerNick"`
// 	BuyerOpenUID              string          `json:"BuyerOpenUid"`
// 	Type                      string          `json:"Type"`
// 	BuyerMessage              string          `json:"BuyerMessage"`
// 	Price                     string          `json:"Price"`
// 	Num                       int32           `json:"Num"`
// 	TotalFee                  string          `json:"TotalFee"`
// 	Payment                   string          `json:"Payment"`
// 	PayTime                   string          `json:"PayTime"`
// 	PicPath                   string          `json:"PicPath"`
// 	PostFee                   string          `json:"PostFee"`
// 	Created                   string          `json:"Created"`
// 	TradeFrom                 string          `json:"TradeFrom"`
// 	Orders                    []JSONDataOrder `json:"Orders"`
// 	SellerMemo                string          `json:"SellerMemo"`
// 	SellerFlag                int             `json:"SellerFlag"`
// 	CreditCardFee             string          `json:"CreditCardFee"`
// 	ExpandCardExpandPriceUsed string          `json:"ExpandCardExpandPriceUsed"`
// }

// type JSONJDGAMEData struct {
// 	Id   string `json:"Id"`
// 	Name string `json:"Name"`
// }

// type JSONJDData struct {
// 	PlatformShopId  int64          `json:"PlatformShopId"`
// 	CustomerId      int64          `json:"CustomerId"`
// 	OrderId         int64          `json:"OrderId"`
// 	OrderType       int            `json:"OrderType"`
// 	Pin             string         `json:"Pin"`
// 	BuyNum          int            `json:"BuyNum"`
// 	SkuId           int64          `json:"SkuId"`
// 	BrandId         int            `json:"BrandId"`
// 	UserIp          string         `json:"UserIp"`
// 	TotalPrice      float64        `json:"TotalPrice"`
// 	CreateTime      string         `json:"CreateTime"`
// 	Features        any            `json:"Features"`
// 	SourceType      int            `json:"SourceType"`
// 	FacePrice       int            `json:"FacePrice"`
// 	GameAccount     string         `json:"GameAccount"`
// 	Permit          string         `json:"Permit"`
// 	GameAccountType JSONJDGAMEData `json:"GameAccountType"`
// 	ChargeType      JSONJDGAMEData `json:"ChargeType"`
// 	GameArea        JSONJDGAMEData `json:"GameArea"`
// 	GameServer      JSONJDGAMEData `json:"GameServer"`
// }

// func AgisoNotify(c echo.Context) error {
// 	err := AgisoNotifyHandler(c)
// 	if err != nil {
// 		c.Logger().Error(err)
// 		response.ResponseSuccess(c, nil)
// 	}
// 	return nil
// }

// const (
// 	BuyerConfirms = "4"       // 买家确认收货
// 	MockNotify    = "21"      // 模拟推送
// 	BuyerPay      = "2097152" // 买家付款
// 	JDGAMECARD    = "8"       // 京东游戏点卡订单支付成功后;
// )

// const (
// 	TRADE_FINISHED         = "TRADE_FINISHED"         // 交易已完成
// 	WAIT_SELLER_SEND_GOODS = "WAIT_SELLER_SEND_GOODS" // 等待卖方发货
// )

// func AgisoNotifyHandler(c echo.Context) error {
// 	fromPlatform := c.QueryParam("fromPlatform")
// 	timestamp := c.QueryParam("timestamp")
// 	aopic := c.QueryParam("aopic")
// 	sign := c.QueryParam("sign")
// 	jsonStr := c.FormValue("json")

// 	c.Logger().Info("fromPlatform:", fromPlatform)
// 	c.Logger().Info("timestamp:", timestamp)
// 	c.Logger().Info("aopic:", aopic)
// 	c.Logger().Info("sign:", sign)
// 	c.Logger().Info("json:", jsonStr)

// 	if fromPlatform == "AldsJd" {
// 		var err error
// 		jsonData := JSONJDData{}
// 		err = json.Unmarshal([]byte(jsonStr), &jsonData)

// 		if err != nil {
// 			return err
// 		}

// 		db := data.Instance()

// 		skuId := cast.ToString(jsonData.SkuId)
// 		c.Logger().Infof("skuId: %s", skuId)

// 		o, err := findOrderByDarkNumber(c, db, skuId, jsonData.GameAccount)
// 		if err != nil {
// 			return fmt.Errorf("findOrderByDarkNumber: %s", err.Error())
// 		}

// 		partner, err := findPartner(db, o.PartnerId)
// 		if err != nil {
// 			return fmt.Errorf("findPartner: %s", err.Error())
// 		}

// 		if !checkSignFromAgiso(c, sign, partner.AqsAppSecret, jsonStr, timestamp) {
// 			return errors.New("checkSignFromAgiso invalid parameter")
// 		}

// 		err = TransactionJD(c, db, jsonData, skuId)
// 		if err != nil {
// 			c.Logger().Errorf("AgisoNotifyHandler Transaction error= %s", err.Error())
// 			return nil
// 		}
// 		return nil
// 	} else {
// 		var err error
// 		jsonData := JSONData{}
// 		err = json.Unmarshal([]byte(jsonStr), &jsonData)

// 		switch aopic {
// 		case BuyerConfirms:
// 			if jsonData.Status != TRADE_FINISHED {
// 				return errors.New("jsonData.Status error")
// 			}
// 		case MockNotify:

// 		default:
// 			return errors.New("aopic not found")
// 		}
// 		if err != nil {
// 			return err
// 		}

// 		if len(jsonData.Orders) == 0 {
// 			return errors.New("jsonData.Orders is empty")
// 		}

// 		db := data.Instance()

// 		for _, v := range jsonData.Orders {
// 			skuId := cast.ToString(v.NumIid)
// 			c.Logger().Infof("skuId: %s", skuId)

// 			o, err := findOrderBySkuId(c, db, skuId, jsonData)
// 			if err != nil {
// 				return fmt.Errorf("findOrderBySkuId: %s", err.Error())
// 			}

// 			partner, err := findPartner(db, o.PartnerId)
// 			if err != nil {
// 				return fmt.Errorf("findPartner: %s", err.Error())
// 			}

// 			if !checkSignFromAgiso(c, sign, partner.AqsAppSecret, jsonStr, timestamp) {
// 				return errors.New("checkSignFromAgiso invalid parameter")
// 			}

// 			err = Transaction(c, db, jsonData, skuId)
// 			if err != nil {
// 				c.Logger().Errorf("AgisoNotifyHandler Transaction error= %s", err.Error())
// 				continue
// 			}
// 		}
// 		return nil
// 	}
// }

// func findPartner(db *gorm.DB, partnerId uint) (*model.Partner, error) {
// 	partner := model.Partner{}
// 	err := db.Where("id = ? AND enable = ?", partnerId, model.Enabled).First(&partner).Error
// 	if err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return nil, errors.New("非法合作商")
// 		}
// 		return nil, err
// 	}
// 	return &partner, nil
// }

// func checkSignFromAgiso(c echo.Context, sign, appSecret, jsonStr, timestamp string) bool {
// 	signStr := fmt.Sprintf("%sjson%stimestamp%s%s", appSecret, jsonStr, timestamp, appSecret)

// 	hash := md5.Sum([]byte(signStr))
// 	s := fmt.Sprintf("%X", hash)

// 	if !strings.EqualFold(sign, s) {
// 		c.Logger().Infof("checkSignFromAgiso: %s--%s", s, sign)
// 	}

// 	return strings.EqualFold(sign, s)
// }

// func Transaction(c echo.Context, db *gorm.DB, jsonData JSONData, skuId string) error {
// 	var o *model.Order
// 	var err error

// 	err = db.Transaction(func(tx *gorm.DB) error {
// 		o, err = findOrderBySkuId(c, tx, skuId, jsonData)
// 		if err != nil {
// 			return fmt.Errorf("findOrderBySkuId error: %s", err.Error())
// 		}

// 		if o.Status == model.OrderStatusFinish {
// 			return fmt.Errorf("订单已完成, skuId=%s, order.Status=%d", skuId, o.Status)
// 		}

// 		err = updateOrder(c, tx, o.OrderId, jsonData)
// 		if err != nil {
// 			return fmt.Errorf("updateOrder error= %s", err.Error())
// 		}

// 		o.ReceivedAmount = util.ToDecimal(float64(jsonData.Num) * cast.ToFloat64(jsonData.Price))
// 		err = UpdatePartnerBalance(c, tx, o)
// 		if err != nil {
// 			return err
// 		}

// 		return nil
// 	})
// 	if err != nil {
// 		return err
// 	}

// 	if o == nil {
// 		return fmt.Errorf("订单是nil, skuId=%s", skuId)
// 	}

// 	err = NotifyMerchant(c, db, o)
// 	if err != nil {
// 		return fmt.Errorf("NotifyMerchant error= %s", err.Error())
// 	}

// 	return nil
// }

// func UpdatePartnerBalance(c echo.Context, db *gorm.DB, o *model.Order) error {
// 	c.Logger().Infof("Update Partner Balance=%d", -o.ReceivedAmount)
// 	err := repository.Partner.UpdateBalance(c, db, o.PartnerId, o.OrderId, -o.ReceivedAmount, model.BalanceFromTypeOrderDeduct)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func TransactionJD(c echo.Context, db *gorm.DB, jsonData JSONJDData, skuId string) error {
// 	var o *model.Order
// 	var err error

// 	err = db.Transaction(func(tx *gorm.DB) error {
// 		o, err = findOrderByDarkNumber(c, tx, skuId, jsonData.GameAccount)
// 		if err != nil {
// 			return fmt.Errorf("findOrderByDarkNumber error: %s", err.Error())
// 		}

// 		if o.Status == model.OrderStatusFinish {
// 			return fmt.Errorf("订单已完成, skuId=%s, order.Status=%d", skuId, o.Status)
// 		}

// 		err = updateOrderJD(c, tx, o.OrderId, jsonData)
// 		if err != nil {
// 			return fmt.Errorf("updateOrder error= %s", err.Error())
// 		}

// 		o.ReceivedAmount = util.ToDecimal(cast.ToFloat64(jsonData.TotalPrice))
// 		err = UpdatePartnerBalance(c, tx, o)
// 		if err != nil {
// 			return err
// 		}

// 		return nil
// 	})
// 	if err != nil {
// 		return err
// 	}

// 	if o == nil {
// 		return fmt.Errorf("订单是nil, skuId=%s", skuId)
// 	}

// 	err = NotifyMerchant(c, db, o)
// 	if err != nil {
// 		return fmt.Errorf("NotifyMerchant error= %s", err.Error())
// 	}

// 	return nil
// }

// func UpdateMerchantBalance(c echo.Context, db *gorm.DB, o *model.Order) error {
// 	c.Logger().Infof("Update Merchant Balance=%d", -o.ReceivedAmount)
// 	err := repository.Merchant.UpdateBalance(c, db, o.MerchantId, o.OrderId, -o.ReceivedAmount, model.BalanceFromTypeOrderDeduct)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func findOrderBySkuId(c echo.Context, db *gorm.DB, skuId string, jsonData JSONData) (*model.Order, error) {
// 	createdTime, err := time.Parse(time.DateTime, jsonData.Created)
// 	if err != nil {
// 		return nil, err
// 	}

// 	o, err := repository.Order.GetBySkuId(c, db, skuId, createdTime)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return o, nil
// }

// func findOrderByDarkNumber(c echo.Context, db *gorm.DB, skuId string, gameAccount string) (*model.Order, error) {
// 	o, err := repository.Order.GetBySkuIdDarkNumber(c, db, skuId, gameAccount)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return o, nil
// }

// func updateOrder(c echo.Context, db *gorm.DB, orderId string, jsonData JSONData) error {
// 	var err error
// 	var parsedTime time.Time
// 	if len(jsonData.PayTime) > 0 {
// 		parsedTime, err = time.Parse(time.DateTime, jsonData.PayTime)
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	err = repository.Order.Update(c, db, orderId, model.Order{
// 		Status:         model.OrderStatusPaid,
// 		PayAt:          parsedTime,
// 		PartnerOrderId: jsonData.TidStr,
// 		PayAccount:     jsonData.BuyerNick,
// 		ReceivedAmount: util.ToDecimal(float64(jsonData.Num) * cast.ToFloat64(jsonData.Price)),
// 	})
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func updateOrderJD(c echo.Context, db *gorm.DB, orderId string, jsonData JSONJDData) error {
// 	var err error
// 	var parsedTime time.Time
// 	if len(jsonData.CreateTime) > 0 {
// 		parsedTime, err = time.Parse(time.DateTime, jsonData.CreateTime)
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	err = repository.Order.Update(c, db, orderId, model.Order{
// 		Status:         model.OrderStatusPaid,
// 		PayAt:          parsedTime,
// 		PartnerOrderId: strconv.FormatInt(jsonData.OrderId, 10),
// 		PayAccount:     jsonData.Pin,
// 		ReceivedAmount: util.ToDecimal(cast.ToFloat64(jsonData.TotalPrice)),
// 	})
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// type NotifyData struct {
// 	MerchantId      int32   `json:"merchantId"`
// 	MerchantTradeNo string  `json:"merchantTradeNo"`
// 	TradeNo         string  `json:"tradeNo"`
// 	Amount          float64 `json:"amount"`
// 	ActualAmount    float64 `json:"actualAmount"`
// 	Timestamp       string  `json:"timestamp"`
// 	Sign            string  `json:"sign"`
// }

// func NotifyMerchant(c echo.Context, db *gorm.DB, o *model.Order) error {
// 	if len(o.NotifyUrl) == 0 || o == nil {
// 		return nil
// 	}

// 	merchant, err := findMerchant(db, int32(o.MerchantId))
// 	if err != nil {
// 		return err
// 	}
// 	privateKey := merchant.PrivateKey

// 	notifyData := NotifyData{
// 		MerchantId:      int32(o.MerchantId),
// 		MerchantTradeNo: o.MerchantOrderId,
// 		TradeNo:         o.OrderId,
// 		Amount:          util.ToDecimal(o.Amount),
// 		ActualAmount:    util.ToDecimal(o.ReceivedAmount),
// 		Timestamp:       cast.ToString(time.Now().Unix()),
// 	}
// 	sign, _ := makeSign(privateKey, notifyData)
// 	notifyData.Sign = sign

// 	c.Logger().Infof("NotifyUrl=%s, NotifyData=%+v", o.NotifyUrl, notifyData)

// 	if config.IsProd() {
// 		var result string
// 		client := resty.New()
// 		resp, err := client.SetTimeout(5 * time.Second).R().SetBody(notifyData).SetResult(&result).Post(o.NotifyUrl)
// 		if err != nil {
// 			return err
// 		}
// 		rawBody := resp.Body()
// 		c.Logger().Infof("response code=%d,  body=%s", resp.StatusCode(), string(rawBody))

// 		if resp.StatusCode() != http.StatusOK {
// 			return fmt.Errorf("通知商户失败, http.StatusCode=%d", resp.StatusCode())
// 		}
// 		if strings.ToLower(string(rawBody)) != Success {
// 			return fmt.Errorf("通知商户失败, 商户回复 result=%s", string(rawBody))
// 		}
// 	}

// 	err = db.Transaction(func(tx *gorm.DB) error {
// 		// 支付状态判断
// 		o, err := repository.Order.GetByOrderId(c, tx, o.OrderId)
// 		if err != nil {
// 			return err
// 		}

// 		c.Logger().Infof("order status=%d", o.Status)

// 		err = repository.Order.Update(c, tx, o.OrderId, model.Order{
// 			Status:       model.OrderStatusFinish,
// 			NotifyStatus: model.Notify,
// 			NotifyAt:     time.Now(),
// 		})
// 		if err != nil {
// 			return err
// 		}

// 		if o.Status == model.OrderStatusFinish {
// 			return nil
// 		}

// 		err = UpdateMerchantBalance(c, tx, o)
// 		if err != nil {
// 			return err
// 		}

// 		return nil
// 	})
// 	if err != nil {
// 		return err
// 	}

// 	c.Logger().Infof("通知商户成功, NotifyData=%+v", notifyData)
// 	return nil
// }

// type QueryBalanceReq struct {
// 	MerchantId int32  `json:"merchantId"`
// 	Timestamp  string `json:"timestamp"`
// 	Sign       string `json:"sign"`
// }

// type QueryBalanceResp struct {
// 	MerchantId int32   `json:"merchantId"`
// 	Balance    float64 `json:"balance"`
// 	Sign       string  `json:"sign"`
// }

// func QueryBalance(c echo.Context) error {
// 	req := QueryOrderReq{}
// 	err := c.Bind(&req)
// 	if err != nil {
// 		response.ResponseError(c, http.StatusBadRequest, "查询余额失败, 参数错误", nil)
// 		return nil
// 	}

// 	db := data.Instance()
// 	merchant, err := findMerchant(db, req.MerchantId)
// 	if err != nil {
// 		return err
// 	}
// 	privateKey := merchant.PrivateKey

// 	if !checkSign(c, privateKey, req) {
// 		return errors.New("参数错误")
// 	}

// 	resp := QueryBalanceResp{
// 		MerchantId: req.MerchantId,
// 		Balance:    util.ToDecimal(merchant.Balance),
// 	}
// 	sign, _ := makeSign(privateKey, resp)
// 	resp.Sign = sign
// 	response.ResponseSuccess(c, resp)
// 	return nil
// }

// func RandomDarkNumber() string {
// 	rand.Seed(uint64(time.Now().UnixNano()))
// 	now := time.Now()
// 	millis := now.UnixNano() / 1e6
// 	timestampPart := fmt.Sprintf("%06d", millis%1e6)
// 	randomPart := fmt.Sprintf("%04d", rand.Intn(10000))
// 	result := randomPart + timestampPart
// 	return result
// }
