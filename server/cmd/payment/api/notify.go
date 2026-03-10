package api

import (
	"apollo/server/cmd/payment/api/channel"
	"apollo/server/cmd/payment/api/common"
	"apollo/server/cmd/payment/pkg/errorx"
	"apollo/server/cmd/payment/pkg/types"
	"apollo/server/internal/model"
	"apollo/server/internal/repository"
	"apollo/server/pkg/config"
	"apollo/server/pkg/contextx"
	"apollo/server/pkg/data"
	"apollo/server/pkg/response"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/labstack/echo/v4"
	"github.com/spf13/cast"
	"gorm.io/gorm"
)

func AgisoNotify(c echo.Context) error {
	fromPlatform := c.QueryParam("fromPlatform")
	timestamp := c.QueryParam("timestamp")
	aopic := c.QueryParam("aopic")
	sign := c.QueryParam("sign")
	jsonStr := c.FormValue("json")

	err := agisoNotifyHandler(c, fromPlatform, timestamp, aopic, sign, jsonStr)
	if err != nil {
		c.Logger().Error(err)
		response.ResponseSuccess(c, nil)
	}
	return nil
}

// 阿奇索通知
func agisoNotifyHandler(c echo.Context, fromPlatform, timestamp, aopic, sign, jsonStr string) error {
	c.Logger().Info("fromPlatform:", fromPlatform)
	c.Logger().Info("timestamp:", timestamp)
	c.Logger().Info("aopic:", aopic)
	c.Logger().Info("sign:", sign)
	c.Logger().Info("json:", jsonStr)

	notifyHandler, err := channel.GetNotifyByPlatform(fromPlatform, timestamp, aopic, sign, jsonStr)
	if err != nil {
		return err
	}

	switch aopic {
	case types.BuyerConfirms:
		fallthrough
	case types.JDGameCard:
		err = notifyHandler.Handle(c)
	case types.MockNotify:
		c.Logger().Info("AgisoNotifyHandler MockNotify")
		if !config.IsProd() {
			err = notifyHandler.Handle(c)
		}
	default:
		return errors.New("aopic not found")
	}
	if err != nil {
		return err
	}

	return nil
}

func NotifySuccess(c echo.Context) error {
	req := types.NotifyData{}
	err := c.Bind(&req)
	if err != nil {
		response.ResponseError(c, http.StatusBadRequest, errorx.ErrTxtNotifySuccessInvalidParams, nil)
		return nil
	}

	db := data.Instance()
	merchant, err := common.FindMerchant(db, req.MerchantId)
	if err != nil {
		return err
	}

	queryReq := types.QueryOrderReq{
		MerchantId:      req.MerchantId,
		MerchantTradeNo: req.MerchantTradeNo,
		Timestamp:       cast.ToString(time.Now()),
	}

	sign := common.NewSign(merchant.PrivateKey)
	queryReq.Sign = sign.Generate(contextx.NewContextFromEcho(c), queryReq)
	url := "http://127.0.0.1:9000/api/order/query"
	var result struct {
		Code    int                  `json:"code"`
		Message string               `json:"message"`
		Data    types.QueryOrderResp `json:"data"`
	}

	client := resty.New()
	resp, err := client.SetTimeout(3 * time.Second).R().SetBody(queryReq).SetResult(&result).Post(url)
	if err != nil {
		return err
	}

	c.Logger().Infof("response code=%d, QueryOrderResp body=%s", resp.StatusCode(), string(resp.Body()))
	c.Logger().Infof("response QueryOrderResp result=%+v", result)
	c.String(http.StatusOK, types.Success)

	return nil
}

func AnssyAuthNotify(c echo.Context) error {
	req := types.AnssyAuthNotifyReq{}
	err := c.Bind(&req)
	if err != nil {
		response.ResponseError(c, http.StatusBadRequest, "参数解析错误", nil)
		return nil
	}

	id := cast.ToUint(req.State)
	// 更新合作商token
	err = common.UpdateAnssyPartner(c, id, req)
	if err != nil {
		return err
	}

	response.ResponseSuccess(c, nil)

	return nil
}

func AnssyNotify(c echo.Context) error {
	fromPlatform := c.QueryParam("fromPlatform")
	timestamp := c.QueryParam("timestamp")
	sign := c.QueryParam("sign")
	aopic := types.BuyerConfirms

	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.String(http.StatusBadRequest, "读取请求体失败")
	}
	defer c.Request().Body.Close()

	// req := types.TBJsonData{}
	// err = json.Unmarshal(body, &req)
	// if err != nil {
	// 	return c.String(http.StatusBadRequest, "解析请求体失败")
	// }

	err = agisoNotifyHandler(c, fromPlatform, timestamp, aopic, sign, string(body))
	if err != nil {
		c.Logger().Error(err)
		response.ResponseSuccess(c, nil)
	}
	return nil
}

type JDCloudNotifyReq struct {
	Token             string `json:"token"`
	AppKey            string `json:"app_key"`
	Sign              string `json:"sign"`
	Timestamp         string `json:"timestamp"`
	Format            string `json:"format"`
	V                 string `json:"v"`
	JdParamJson       string `json:"jd_param_json"`
	EncryptJdParamJson string `json:"encrypt_jd_param_json"`
}

type JDCloudOrderPay struct {
	OrderId   int64  `json:"orderId"`
	OrderType int    `json:"orderType"`
	StatusId  int    `json:"statusId"`
	Timestamp string `json:"timestamp"`
}

type JDCloudOrderDetailReq struct {
	OrderId int64  `json:"orderId"`
	Token   string `json:"token"`
	AppKey  string `json:"app_key"`
	Sign    string `json:"sign"`
	Format  string `json:"format"`
	Timestamp string `json:"timestamp"`
	V       string `json:"v"`
}

type JDCloudOrderDetailResp struct {
	Code      int    `json:"code"`
	Message   string `json:"message"`
	OrderInfo any    `json:"orderInfo"`
}

type JDCloudOrderInfo struct {
	OrderId          int64       `json:"orderId"`
	OrderType        int         `json:"orderType"`
	OrderState       string      `json:"orderState"`
	OrderSellerPrice float64     `json:"orderSellerPrice"`
	ItemInfoList     []JDCloudItemInfo `json:"itemInfoList"`
}

type JDCloudItemInfo struct {
	SkuId   string `json:"skuId"`
	SkuName string `json:"skuName"`
	ItemTotal int   `json:"itemTotal"`
}

func generateJDCloudSign(appSecret, timestamp, paramJson string) string {
	signStr := fmt.Sprintf("app_key%sformatjsonparam_json%stimestamp%s%s%s",
		"", timestamp, paramJson, timestamp, appSecret)
	hash := md5.Sum([]byte(signStr))
	return fmt.Sprintf("%X", hash)
}

func getJDCloudOrderDetail(c echo.Context, orderId int64) (*JDCloudOrderInfo, error) {
	conf := config.Get()
	if conf == nil {
		conf = config.New("configs/config.yaml")
	}
	jdConf := conf.JDCloudConfig

	if jdConf.AppKey == "" || jdConf.AppSecret == "" {
		return nil, errors.New("京东配置缺失app_key或app_secret")
	}

	paramJson := fmt.Sprintf(`{"orderId":%d}`, orderId)
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	sign := generateJDCloudSign(jdConf.AppSecret, timestamp, paramJson)

	client := resty.New()
	url := "https://api.jd.com/routerjson"
	var result struct {
		jingdong_pop_order_get_responce struct {
			Code    string `json:"code"`
			Message string `json:"message"`
			OrderDetailInfo struct {
				OrderInfo JDCloudOrderInfo `json:"orderInfo"`
			} `json:"orderDetailInfo"`
		} `json:"jingdong_pop_order_get_responce"`
	}

	resp, err := client.SetTimeout(10*time.Second).R().SetFormData(map[string]string{
		"method":        "jingdong.pop.order.get",
		"access_token": jdConf.Token,
		"app_key":      jdConf.AppKey,
		"sign":         sign,
		"format":       "json",
		"v":            "1.0",
		"timestamp":    timestamp,
		"param_json":   paramJson,
	}).SetResult(&result).Post(url)

	if err != nil {
		c.Logger().Error("调用京东API失败:", err)
		return nil, err
	}

	c.Logger().Info("京东订单详情API响应:", string(resp.Body()))

	respData := result.jingdong_pop_order_get_responce
	if respData.Code != "0" {
		return nil, fmt.Errorf("京东API返回错误: %s", respData.Message)
	}

	return &respData.OrderDetailInfo.OrderInfo, nil
}

func JDCloudNotify(c echo.Context) error {
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.String(http.StatusBadRequest, "读取请求体失败")
	}
	defer c.Request().Body.Close()

	c.Logger().Info("JDCloudNotify received:", string(body))

	var req JDCloudNotifyReq
	if err := json.Unmarshal(body, &req); err != nil {
		c.Logger().Error("解析请求体失败:", err)
		return c.String(http.StatusBadRequest, "参数解析错误")
	}

	c.Logger().Info("Token:", req.Token)
	c.Logger().Info("AppKey:", req.AppKey)
	c.Logger().Info("Sign:", req.Sign)
	c.Logger().Info("Timestamp:", req.Timestamp)
	c.Logger().Info("JdParamJson:", req.JdParamJson)

	var orderPay JDCloudOrderPay
	if req.JdParamJson != "" {
		if err := json.Unmarshal([]byte(req.JdParamJson), &orderPay); err != nil {
			c.Logger().Error("解析jd_param_json失败:", err)
			return c.String(http.StatusBadRequest, "参数解析错误")
		}
	} else if req.EncryptJdParamJson != "" {
		c.Logger().Info("收到加密参数，需要解密:", req.EncryptJdParamJson)
		return c.String(http.StatusOK, "{\"code\":0,\"message\":\"success\"}")
	}

	c.Logger().Info("OrderId:", orderPay.OrderId)
	c.Logger().Info("StatusId:", orderPay.StatusId)

	if orderPay.OrderId == 0 {
		c.Logger().Error("OrderId为空")
		return c.String(http.StatusOK, "{\"code\":0,\"message\":\"success\"}")
	}

	orderDetail, err := getJDCloudOrderDetail(c, orderPay.OrderId)
	if err != nil {
		c.Logger().Error("获取订单详情失败:", err)
		return c.String(http.StatusOK, "{\"code\":0,\"message\":\"success\"}")
	}

	c.Logger().Info("订单详情:", orderDetail)

	if len(orderDetail.ItemInfoList) == 0 {
		c.Logger().Error("订单商品为空")
		return c.String(http.StatusOK, "{\"code\":0,\"message\":\"success\"}")
	}

	db := data.Instance()
	skuId := orderDetail.ItemInfoList[0].SkuId

	o, err := repository.Order.GetByOrderId(db, fmt.Sprintf("%d", orderPay.OrderId))
	if err != nil {
		c.Logger().Error("查找订单失败:", err)
		return c.String(http.StatusOK, "{\"code\":0,\"message\":\"success\"}")
	}

	if o != nil && o.Status == 0 {
		o.Status = 1
		err = db.Save(o).Error
		if err != nil {
			c.Logger().Error("更新订单状态失败:", err)
		}
		c.Logger().Info("订单状态已更新为已支付:", orderPay.OrderId)
	}

	partnerId := o.PartnerId
	_, err = common.FindPartner(db, partnerId)
	if err != nil {
		c.Logger().Error("查找合作商失败:", err)
		return c.String(http.StatusOK, "{\"code\":0,\"message\":\"success\"}")
	}

	jdc := channel.GetJDCardNotify(db, orderPay.OrderId, skuId, "")
	if jdc != nil {
		err = jdc.Handle(c)
		if err != nil {
			c.Logger().Error("处理订单事务失败:", err)
		} else {
			c.Logger().Info("订单处理成功:", orderPay.OrderId)
		}
	}

	return c.String(http.StatusOK, "{\"code\":0,\"message\":\"success\"}")
}

func FindOrderByOrderId(c echo.Context, db *gorm.DB, orderId string) (*model.Order, error) {
	return repository.Order.GetByOrderId(db, orderId)
}
