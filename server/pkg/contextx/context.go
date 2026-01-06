package contextx

import "github.com/labstack/echo/v4"

type Context struct {
	logger echo.Logger
}

func NewContext(logger echo.Logger) Context {
	return Context{
		logger: logger,
	}
}

func NewContextFromEcho(ctx echo.Context) Context {
	return Context{
		logger: ctx.Logger(),
	}
}

func (c *Context) Logger() echo.Logger {
	return c.logger
}
