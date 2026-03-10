package handler

import (
	v1 "apollo/server/api/v1"
	"apollo/server/internal/service"
	"apollo/server/pkg/response"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

var (
	PriceCard = new(PriceCardHandler)
)

type PriceCardHandler struct {
}

func (h *PriceCardHandler) Create(c echo.Context) error {
	req := new(v1.CardCreateReq)
	err := c.Bind(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	resp, err := service.PriceCardServiceInst.CreateCards(c, req)
	if err != nil {
		return err
	}

	return response.ResponseSuccess(c, resp)
}

func (h *PriceCardHandler) GenerateVirtual(c echo.Context) error {
	req := new(v1.VirtualCardGenerateReq)
	err := c.Bind(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	resp, err := service.PriceCardServiceInst.GenerateVirtualCards(c, req)
	if err != nil {
		return err
	}

	return response.ResponseSuccess(c, resp)
}

func (h *PriceCardHandler) List(c echo.Context) error {
	req := new(v1.ListCardReq)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	if req.CurrentPage == 0 {
		req.CurrentPage, _ = strconv.Atoi(c.QueryParam("currentPage"))
	}
	if req.PageSize == 0 {
		req.PageSize, _ = strconv.Atoi(c.QueryParam("pageSize"))
	}

	resp, err := service.PriceCardServiceInst.List(c, req)
	if err != nil {
		return err
	}

	return response.ResponseSuccess(c, resp)
}

func (h *PriceCardHandler) Delete(c echo.Context) error {
	req := new(v1.DeleteCardReq)
	err := c.Bind(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	resp, err := service.PriceCardServiceInst.Delete(c, req)
	if err != nil {
		return err
	}

	return response.ResponseSuccess(c, resp)
}

func (h *PriceCardHandler) DeleteByCondition(c echo.Context) error {
	req := new(v1.ListCardReq)
	err := c.Bind(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	resp, err := service.PriceCardServiceInst.DeleteByCondition(c, req)
	if err != nil {
		return err
	}

	return response.ResponseSuccess(c, resp)
}
