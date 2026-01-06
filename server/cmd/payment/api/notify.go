package api

import (
	"apollo/server/cmd/payment/api/channel"
	"apollo/server/cmd/payment/api/common"
	"apollo/server/cmd/payment/pkg/errorx"
	"apollo/server/cmd/payment/pkg/types"
	"apollo/server/pkg/config"
	"apollo/server/pkg/contextx"
	"apollo/server/pkg/data"
	"apollo/server/pkg/response"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/labstack/echo/v4"
	"github.com/spf13/cast"
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
