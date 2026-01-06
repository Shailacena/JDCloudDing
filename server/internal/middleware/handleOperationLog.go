package middleware

import (
	"apollo/server/internal/model"
	"apollo/server/internal/repository"
	"apollo/server/pkg/headerx"
	"github.com/labstack/echo/v4"
)

func HandleOperationLogMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := next(c)
			if err == nil {
				h := headerx.GetDataFromHeader(c)

				if h.AdminId <= 0 || h.Role <= 0 || len(c.Request().RequestURI) == 0 {
					return err
				}

				logs := make([]*model.OperationLog, 0)
				logs = append(logs, &model.OperationLog{
					IP:        c.RealIP(),
					Operation: c.Request().RequestURI,
					Operator:  h.AdminId,
				})

				err1 := repository.OperationLog.Create(c, logs)
				if err1 != nil {
					c.Logger().Warn("OperationLog.Create error", err1)
				}
			}

			return err
		}
	}
}
