package app

import (
	"sync"

	"github.com/labstack/echo/v4"
)

var (
	e    *echo.Echo
	once sync.Once
)

func Engine() *echo.Echo {
	once.Do(func() {
		e = echo.New()
	})

	return e
}
