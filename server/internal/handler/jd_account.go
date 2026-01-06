package handler

import (
	v1 "apollo/server/api/v1"
	"apollo/server/internal/service"
	"apollo/server/pkg/response"
	"net/http"

	"github.com/labstack/echo/v4"
)

var (
	JDAccount = new(JDAccountHandler)
)

type JDAccountHandler struct {
}

func (h *JDAccountHandler) Create(c echo.Context) error {
	req := new(v1.JDAccountCreateReq)
	err := c.Bind(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	resp, err := service.JDAccount.Create(c, req)
	if err != nil {
		return err
	}

	return response.ResponseSuccess(c, resp)
}

func (h *JDAccountHandler) Enable(c echo.Context) error {
	req := new(v1.JDAccountEnableReq)
	err := c.Bind(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	resp, err := service.JDAccount.Enable(c, req)
	if err != nil {
		return err
	}

	return response.ResponseSuccess(c, resp)
}

func (h *JDAccountHandler) List(c echo.Context) error {
	req := new(v1.ListJDAccountReq)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	resp, err := service.JDAccount.List(c, req)
	if err != nil {
		return err
	}

	return response.ResponseSuccess(c, resp)
}

func (h *JDAccountHandler) Delete(c echo.Context) error {
	req := new(v1.JDAccountDeleteReq)
	err := c.Bind(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	resp, err := service.JDAccount.Delete(c, req)
	if err != nil {
		return err
	}

	return response.ResponseSuccess(c, resp)
}

func (h *JDAccountHandler) Reset(c echo.Context) error {
	req := new(v1.JDAccountResetReq)
	err := c.Bind(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	resp, err := service.JDAccount.Reset(c, req)
	if err != nil {
		return err
	}

	return response.ResponseSuccess(c, resp)
}
