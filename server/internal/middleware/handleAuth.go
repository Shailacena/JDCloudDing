package middleware

import (
	"apollo/server/pkg/headerx"
	"github.com/labstack/echo/v4"
)

type TokenChecker interface {
	CheckToken(echo.Context, string) error
}

func GenAuthHandler(checker TokenChecker) func() echo.MiddlewareFunc {
	return func() echo.MiddlewareFunc {
		return func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				h := headerx.GetDataFromHeader(c)

				// 校验token
				err := checker.CheckToken(c, h.Token)
				if err != nil {
					return err
				}

				return next(c)
			}
		}
	}
}
