package response

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

const (
	SUCCESS = 0

	success = "success"

	ERROR = 1
)

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func ResponseSuccess(c echo.Context, data any) error {
	if data == nil {
		data = map[string]any{}
	}

	resp := Response{Code: SUCCESS, Message: success, Data: data}
	return c.JSON(http.StatusOK, resp)
}

func ResponseError(c echo.Context, httpCode int, message string, data any) error {
	if data == nil {
		data = map[string]any{}
	}

	resp := Response{Code: ERROR, Message: message, Data: data}
	return c.JSON(httpCode, resp)
}
