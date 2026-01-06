package api

import (
	"apollo/server/cmd/payment/api/channel"
	"apollo/server/cmd/payment/api/common"
	"apollo/server/cmd/payment/pkg/errorx"
	"apollo/server/cmd/payment/pkg/types"
	"apollo/server/internal/model"
	"apollo/server/internal/repository"
	"apollo/server/pkg/contextx"
	"apollo/server/pkg/data"
	"apollo/server/pkg/response"
	"apollo/server/pkg/util"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/spf13/cast"
)

// 创建订单
func CreateOrder(c echo.Context) error {
	req := types.CreateOrderReq{}
	err := c.Bind(&req)
	if err != nil {
		response.ResponseError(c, http.StatusBadRequest, errorx.ErrTxtCreateOrderInvalidParams, nil)
		return nil
	}

	payer, err := channel.GetPayByChannelId(req.ChannelId, req.MerchantId)
	if err != nil {
		return err
	}
	merchant := payer.GetMerchant(c)

	sign := common.NewSign(merchant.PrivateKey)

	if !sign.Check(c, req) {
		return errorx.ErrInvalidSign
	}

	goods, err := payer.FindGoods(c, req.ChannelId, req.MerchantId, req.Amount)
	if err != nil {
		return err
	}

	orderId, err := common.GenOrderId()
	if err != nil {
		return err
	}

	partnerId := goods.PartnerId
	o, err := payer.GenOrder(c, req, goods, orderId)
	if err != nil {
		return err
	}

	err = repository.Order.Create(c, o)
	if err != nil {
		return err
	}

	baseUrl, err := common.GetPayBaseUrl(partnerId, req.ChannelId)
	if err != nil {
		return err
	}

	url := payer.GenPayUrl(c, baseUrl, o, goods.NumId)

	resp := types.CreateOrderResp{
		MerchantTradeNo: req.MerchantTradeNo,
		TradeNo:         orderId,
		Amount:          req.Amount,
		PayPageUrl:      url,
	}

	resp.Sign = sign.Generate(contextx.NewContextFromEcho(c), resp)
	response.ResponseSuccess(c, resp)

	return nil
}

// 查询订单
func QueryOrder(c echo.Context) error {
	req := types.QueryOrderReq{}
	err := c.Bind(&req)
	if err != nil {
		response.ResponseError(c, http.StatusBadRequest, errorx.ErrTxtQueryOrderInvalidParams, nil)
		return nil
	}

	db := data.Instance()
	merchant, err := common.FindMerchant(db, req.MerchantId)
	if err != nil {
		return err
	}

	sign := common.NewSign(merchant.PrivateKey)

	if !sign.Check(c, req) {
		return errorx.ErrInvalidSign
	}

	o, err := common.FindOrderById(c, req.MerchantTradeNo)
	if err != nil {
		return err
	}

	var payAt string
	if !o.PayAt.IsZero() {
		payAt = cast.ToString(o.PayAt.Unix())
	}

	status := o.Status
	if status == model.OrderStatusRefundSuccessful || status == model.OrderStatusRefundFailed {
		status = model.OrderStatusUnpaid
	}

	resp := types.QueryOrderResp{
		MerchantId:      req.MerchantId,
		MerchantTradeNo: req.MerchantTradeNo,
		Amount:          util.ToDecimal(o.Amount),
		ActualAmount:    util.ToDecimal(o.ReceivedAmount),
		TradeNo:         o.OrderId,
		Status:          int32(status),
		PayAt:           payAt,
	}

	resp.Sign = sign.Generate(contextx.NewContextFromEcho(c), resp)
	response.ResponseSuccess(c, resp)

	return nil
}

// 查询余额
func QueryBalance(c echo.Context) error {
	req := types.QueryOrderReq{}
	err := c.Bind(&req)
	if err != nil {
		response.ResponseError(c, http.StatusBadRequest, errorx.ErrTxtQueryBalanceInvalidParams, nil)
		return nil
	}

	db := data.Instance()
	merchant, err := common.FindMerchant(db, req.MerchantId)
	if err != nil {
		return err
	}

	sign := common.NewSign(merchant.PrivateKey)

	if !sign.Check(c, req) {
		return errorx.ErrInvalidSign
	}

	resp := types.QueryBalanceResp{
		MerchantId: req.MerchantId,
		Balance:    util.ToDecimal(merchant.Balance),
	}

	resp.Sign = sign.Generate(contextx.NewContextFromEcho(c), resp)
	response.ResponseSuccess(c, resp)

	return nil
}
