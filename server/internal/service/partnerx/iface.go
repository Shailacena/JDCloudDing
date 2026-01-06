package partnerx

import (
	"apollo/server/internal/model"
	"github.com/labstack/echo/v4"
)

type IPartnerGenerator interface {
	CheckParams(ctx echo.Context) error
	GenPartner() model.Partner
}
