package handler

import (
	v1 "apollo/server/api/v1"
	"apollo/server/internal/service"
	"apollo/server/pkg/headerx"
	"apollo/server/pkg/response"
	"net/http"

	"github.com/labstack/echo/v4"
)

var (
	Merchant = new(MerchantHandler)
)

type MerchantHandler struct {
}

func (h *MerchantHandler) Register(c echo.Context) error {
	req := new(v1.MerchantRegisterReq)
	err := c.Bind(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	resp, err := service.Merchant.Register(c, req)
	if err != nil {
		return err
	}

	return response.ResponseSuccess(c, resp)
}

func (h *MerchantHandler) Update(c echo.Context) error {
	req := new(v1.MerchantUpdateReq)
	err := c.Bind(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	resp, err := service.Merchant.Update(c, req)
	if err != nil {
		return err
	}

	return response.ResponseSuccess(c, resp)
}

func (h *MerchantHandler) UpdateBalance(c echo.Context) error {
	req := new(v1.MerchantUpdateBalanceReq)
	err := c.Bind(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	resp, err := service.Merchant.UpdateBalance(c, req)
	if err != nil {
		return err
	}

	return response.ResponseSuccess(c, resp)
}

func (h *MerchantHandler) Login(c echo.Context) error {
	req := new(v1.MerchantLoginReq)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	resp, err := service.Merchant.Login(c, req)
	if err != nil {
		return err
	}

	return response.ResponseSuccess(c, resp)
}

func (h *MerchantHandler) Logout(c echo.Context) error {
	req := new(v1.MerchantLogoutReq)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	header := headerx.GetDataFromHeader(c)
	resp, err := service.Merchant.Logout(c, req, header.Token)
	if err != nil {
		return err
	}

	return response.ResponseSuccess(c, resp)
}

func (h *MerchantHandler) List(c echo.Context) error {
	req := new(v1.ListMerchantReq)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	resp, err := service.Merchant.List(c, req)
	if err != nil {
		return err
	}

	return response.ResponseSuccess(c, resp)
}

func (h *MerchantHandler) SetPassword(c echo.Context) error {
	req := new(v1.MerchantSetPasswordReq)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	header := headerx.GetDataFromHeader(c)
	resp, err := service.Merchant.SetPassword(c, req, header.Token)
	if err != nil {
		return err
	}

	return response.ResponseSuccess(c, resp)
}

func (h *MerchantHandler) ResetPassword(c echo.Context) error {
	req := new(v1.MerchantResetPasswordReq)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	resp, err := service.Merchant.ResetPassword(c, req)
	if err != nil {
		return err
	}

	return response.ResponseSuccess(c, resp)
}

func (h *MerchantHandler) Enable(c echo.Context) error {
	req := new(v1.MerchantEnableReq)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	resp, err := service.Merchant.Enable(c, req)
	if err != nil {
		return err
	}

	return response.ResponseSuccess(c, resp)
}

func (h *MerchantHandler) ListBalanceBill(c echo.Context) error {
	req := new(v1.ListMerchantBalanceBillReq)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	resp, err := service.Merchant.ListBalanceBill(c, req)
	if err != nil {
		return err
	}

	return response.ResponseSuccess(c, resp)
}

func (h *MerchantHandler) GetBalance(c echo.Context) error {
	req := new(v1.MerchantBalanceReq)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	header := headerx.GetDataFromHeader(c)
	resp, err := service.Merchant.GetBalance(c, req, header.Token)
	if err != nil {
		return err
	}

	return response.ResponseSuccess(c, resp)
}

func (h *MerchantHandler) ResetVerifiCode(c echo.Context) error {
	req := new(v1.MerchantResetVerifiCodeReq)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	resp, err := service.Merchant.ResetVerifiCode(c, req)
	if err != nil {
		return err
	}

	return response.ResponseSuccess(c, resp)
}
