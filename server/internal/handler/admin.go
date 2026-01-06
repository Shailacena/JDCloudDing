package handler

import (
	v1 "apollo/server/api/v1"
	"apollo/server/internal/model"
	"apollo/server/internal/service"
	"apollo/server/pkg/config"
	"apollo/server/pkg/headerx"
	"apollo/server/pkg/response"
	"net/http"

	"github.com/labstack/echo/v4"
)

var (
	Admin = new(AdminHandler)
)

type AdminHandler struct {
}

func (h *AdminHandler) Register(c echo.Context) error {
	req := new(v1.AdminRegisterReq)
	err := c.Bind(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	resp, err := service.Admin.Register(c, req)
	if err != nil {
		return err
	}

	return response.ResponseSuccess(c, resp)
}

func (h *AdminHandler) Register11(c echo.Context) error {
	req := new(v1.AdminRegister11Req)
	err := c.Bind(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	if req.S != config.S {
		return echo.NewHTTPError(http.StatusNotFound)
	}

	req.Role = int(model.SuperAdminRole)
	resp, err := service.Admin.Register11(c, &req.AdminRegisterReq)
	if err != nil {
		return err
	}

	return response.ResponseSuccess(c, resp)
}

func (h *AdminHandler) Login(c echo.Context) error {
	req := new(v1.AdminLoginReq)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	resp, err := service.Admin.Login(c, req)
	if err != nil {
		return err
	}

	return response.ResponseSuccess(c, resp)
}

func (h *AdminHandler) Logout(c echo.Context) error {
	req := new(v1.AdminLogoutReq)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	header := headerx.GetDataFromHeader(c)
	resp, err := service.Admin.Logout(c, req, header.Token)
	if err != nil {
		return err
	}

	return response.ResponseSuccess(c, resp)
}

func (h *AdminHandler) List(c echo.Context) error {
	req := new(v1.ListAdminReq)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	resp, err := service.Admin.List(c, req)
	if err != nil {
		return err
	}

	return response.ResponseSuccess(c, resp)
}

func (h *AdminHandler) SetPassword(c echo.Context) error {
	req := new(v1.AdminSetPasswordReq)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	header := headerx.GetDataFromHeader(c)
	resp, err := service.Admin.SetPassword(c, req, header.Token)
	if err != nil {
		return err
	}

	return response.ResponseSuccess(c, resp)
}

func (h *AdminHandler) ResetPassword(c echo.Context) error {
	req := new(v1.AdminResetPasswordReq)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	resp, err := service.Admin.ResetPassword(c, req)
	if err != nil {
		return err
	}

	return response.ResponseSuccess(c, resp)
}

func (h *AdminHandler) Delete(c echo.Context) error {
	req := new(v1.AdminDeleteReq)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	resp, err := service.Admin.Delete(c, req)
	if err != nil {
		return err
	}

	return response.ResponseSuccess(c, resp)
}

func (h *AdminHandler) Update(c echo.Context) error {
	req := new(v1.AdminUpdateReq)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	resp, err := service.Admin.Update(c, req)
	if err != nil {
		return err
	}

	return response.ResponseSuccess(c, resp)
}

func (h *AdminHandler) Enable(c echo.Context) error {
	req := new(v1.AdminEnableReq)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	resp, err := service.Admin.Enable(c, req)
	if err != nil {
		return err
	}

	return response.ResponseSuccess(c, resp)
}

func (h *AdminHandler) ResetVerifiCode(c echo.Context) error {
	req := new(v1.AdminResetVerifiCodeReq)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	resp, err := service.Admin.ResetVerifiCode(c, req)
	if err != nil {
		return err
	}

	return response.ResponseSuccess(c, resp)
}

func (h *AdminHandler) GetMasterIncome(c echo.Context) error {
	req := new(v1.GetMasterIncomeReq)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	resp, err := service.Admin.GetMasterIncome(c, req)
	if err != nil {
		return err
	}

	return response.ResponseSuccess(c, resp)
}
