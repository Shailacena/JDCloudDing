package middleware

import (
	"apollo/server/pkg/response"
	"net/http"

	"github.com/labstack/echo/v4"
)

func HandleErrorMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := next(c)

			if err != nil {
				c.Logger().Error(err)

				if he, ok := err.(*echo.HTTPError); ok {
					return response.ResponseError(c, he.Code, he.Error(), nil)
				}

				return response.ResponseError(c, http.StatusInternalServerError, err.Error(), nil)
			}

			return nil
		}
	}
}
