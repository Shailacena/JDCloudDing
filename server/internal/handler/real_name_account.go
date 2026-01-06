package handler

import (
	v1 "apollo/server/api/v1"
	"apollo/server/internal/service"
	"apollo/server/pkg/response"
	"net/http"

	"github.com/labstack/echo/v4"
)

var (
	RealNameAccount = new(RealNameAccountHandler)
)

type RealNameAccountHandler struct {
}

func (h *RealNameAccountHandler) Create(c echo.Context) error {
	req := new(v1.RealNameAccountCreateReq)
	err := c.Bind(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	resp, err := service.RealNameAccount.Create(c, req)
	if err != nil {
		return err
	}

	return response.ResponseSuccess(c, resp)
}

func (h *RealNameAccountHandler) List(c echo.Context) error {
	req := new(v1.ListRealNameAccountReq)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	resp, err := service.RealNameAccount.List(c, req)
	if err != nil {
		return err
	}

	return response.ResponseSuccess(c, resp)
}
