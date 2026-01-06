package headerx

import (
	"apollo/server/pkg/util"

	"github.com/labstack/echo/v4"
	"github.com/spf13/cast"
)

type Header struct {
	Token   string
	Role    uint
	AdminId uint
}

func GetDataFromHeader(c echo.Context) Header {
	var h Header
	h.Token = c.Request().Header.Get(util.TokenCookieKey)
	r := c.Request().Header.Get(util.RoleCookieKey)
	d := c.Request().Header.Get(util.AdminIdCookieKey)

	h.Role = cast.ToUint(r)
	h.AdminId = cast.ToUint(d)

	return h
}
