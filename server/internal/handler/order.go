package handler

import (
	v1 "apollo/server/api/v1"
	"apollo/server/internal/service"
	"apollo/server/pkg/response"
	"net/http"

	"github.com/labstack/echo/v4"
)

var (
	Order = new(OrderHandler)
)

type OrderHandler struct {
}

func (h *OrderHandler) List(c echo.Context) error {
	req := new(v1.ListOrderReq)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	resp, err := service.Order.List(c, req)
	if err != nil {
		return err
	}

	return response.ResponseSuccess(c, resp)
}

func (h *OrderHandler) Confirm(c echo.Context) error {
	req := new(v1.ConfirmOrderReq)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	resp, err := service.Order.Confirm(c, req)
	if err != nil {
		return err
	}

	return response.ResponseSuccess(c, resp)
}

func (h *OrderHandler) GetOrderSummary(c echo.Context) error {
	req := new(v1.GetOrderSummaryReq)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	resp, err := service.Order.GetOrderSummary(c, req)
	if err != nil {
		return err
	}

	return response.ResponseSuccess(c, resp)
}

func (h *OrderHandler) Archive(c echo.Context) error {
    req := new(v1.ArchiveOrdersReq)
    if err := c.Bind(req); err != nil {
        return echo.NewHTTPError(http.StatusBadRequest, err.Error())
    }
    if err := c.Validate(req); err != nil {
        return err
    }
    resp, err := service.Order.Archive(c, req)
    if err != nil {
        return err
    }
    return response.ResponseSuccess(c, resp)
}
