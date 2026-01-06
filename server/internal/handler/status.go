package handler

import (
	v1 "apollo/server/api/v1"
	"apollo/server/internal/service"
	"apollo/server/pkg/response"
	"net/http"

	"github.com/labstack/echo/v4"
)

var (
	Status = new(StatusHandler)
)

type StatusHandler struct{}

// GET /web_api/admin/server/status
func (h *StatusHandler) GetServerStatus(c echo.Context) error {
	req := new(v1.GetServerStatusReq)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	resp, err := service.Status.GetServerStatus(c, req)
	if err != nil {
		return err
	}
	return response.ResponseSuccess(c, resp)
}

// GET /web_api/admin/order/trend/today
func (h *StatusHandler) TodayTrend(c echo.Context) error {
	req := new(v1.GetTodayTrendReq)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	resp, err := service.Status.TodayTrend(c, req)
	if err != nil {
		return err
	}
	return response.ResponseSuccess(c, resp)
}
