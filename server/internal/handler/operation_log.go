package handler

import (
	v1 "apollo/server/api/v1"
	"apollo/server/internal/service"
	"apollo/server/pkg/response"
	"github.com/labstack/echo/v4"
	"net/http"
)

var (
	OperationLog = new(OperationLogHandler)
)

type OperationLogHandler struct {
}

func (h *OperationLogHandler) List(c echo.Context) error {
	req := new(v1.ListOperationLogReq)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	resp, err := service.OperationLog.List(c, req)
	if err != nil {
		return err
	}

	return response.ResponseSuccess(c, resp)
}
