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
	Partner = new(PartnerHandler)
)

type PartnerHandler struct {
}

func (h *PartnerHandler) Register(c echo.Context) error {
	req := new(v1.PartnerRegisterReq)
	err := c.Bind(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	resp, err := service.Partner.Register(c, req)
	if err != nil {
		return err
	}

	return response.ResponseSuccess(c, resp)
}

func (h *PartnerHandler) Login(c echo.Context) error {
	req := new(v1.PartnerLoginReq)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	resp, err := service.Partner.Login(c, req)
	if err != nil {
		return err
	}

	return response.ResponseSuccess(c, resp)
}

func (h *PartnerHandler) Logout(c echo.Context) error {
	req := new(v1.PartnerLogoutReq)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	header := headerx.GetDataFromHeader(c)
	resp, err := service.Partner.Logout(c, req, header.Token)
	if err != nil {
		return err
	}

	return response.ResponseSuccess(c, resp)
}

func (h *PartnerHandler) ResetPassword(c echo.Context) error {
	req := new(v1.PartnerResetPasswordReq)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	resp, err := service.Partner.ResetPassword(c, req)
	if err != nil {
		return err
	}

	return response.ResponseSuccess(c, resp)
}

func (h *PartnerHandler) Delete(c echo.Context) error {
	req := new(v1.PartnerDeleteReq)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	resp, err := service.Partner.Delete(c, req)
	if err != nil {
		return err
	}

	return response.ResponseSuccess(c, resp)
}

func (h *PartnerHandler) List(c echo.Context) error {
	req := new(v1.ListPartnerReq)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	resp, err := service.Partner.List(c, req)
	if err != nil {
		return err
	}

	return response.ResponseSuccess(c, resp)
}

// func (h *PartnerHandler) ListBill(c echo.Context) error {
// 	req := new(v1.ListPartnerBillReq)
// 	if err := c.Bind(req); err != nil {
// 		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
// 	}

// 	if err := c.Validate(req); err != nil {
// 		return err
// 	}

// 	resp, err := service.Partner.ListBill(c, req)
// 	if err != nil {
// 		return err
// 	}

// 	return response.ResponseSuccess(c, resp)
// }

func (h *PartnerHandler) SetPassword(c echo.Context) error {
	req := new(v1.PartnerSetPasswordReq)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	header := headerx.GetDataFromHeader(c)
	resp, err := service.Partner.SetPassword(c, req, header.Token)
	if err != nil {
		return err
	}

	return response.ResponseSuccess(c, resp)
}

func (h *PartnerHandler) Update(c echo.Context) error {
	req := new(v1.PartnerUpdateReq)
	err := c.Bind(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	resp, err := service.Partner.Update(c, req)
	if err != nil {
		return err
	}

	return response.ResponseSuccess(c, resp)
}

func (h *PartnerHandler) UpdateBalance(c echo.Context) error {
	req := new(v1.PartnerUpdateBalanceReq)
	err := c.Bind(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	resp, err := service.Partner.UpdateBalance(c, req)
	if err != nil {
		return err
	}

	return response.ResponseSuccess(c, resp)
}

func (h *PartnerHandler) SyncGoods(c echo.Context) error {
	req := new(v1.PartnerSyncGoodsReq)
	err := c.Bind(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	resp, err := service.Partner.SyncGoods(c, req)
	if err != nil {
		return err
	}

	return response.ResponseSuccess(c, resp)
}

func (h *PartnerHandler) ListBalanceBill(c echo.Context) error {
	req := new(v1.ListPartnerBalanceBillReq)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	resp, err := service.Partner.ListBalanceBill(c, req)
	if err != nil {
		return err
	}

	return response.ResponseSuccess(c, resp)
}

func (h *PartnerHandler) ResetVerifiCode(c echo.Context) error {
	req := new(v1.PartnerResetVerifiCodeReq)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	resp, err := service.Partner.ResetVerifiCode(c, req)
	if err != nil {
		return err
	}

	return response.ResponseSuccess(c, resp)
}
