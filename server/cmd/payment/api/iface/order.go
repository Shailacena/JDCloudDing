package iface

import (
	"apollo/server/cmd/payment/api/common"
	"apollo/server/cmd/payment/pkg/types"
	"apollo/server/internal/model"

	"github.com/labstack/echo/v4"
)

type IOrder interface {
	FindGoods(c echo.Context, channelId string, merchantId int32, amount float64) (*common.FindGoodsResult, error)
	GenOrder(c echo.Context, req types.CreateOrderReq, goods *common.FindGoodsResult, orderId string) (model.Order, error)
	GenPayUrl(c echo.Context, baseUrl string, o model.Order, numId string) string
	GetMerchant(c echo.Context) *model.Merchant
}

type IHander interface {
	Handle(c echo.Context) error
}
