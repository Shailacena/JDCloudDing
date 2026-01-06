package partnerx

import (
	v1 "apollo/server/api/v1"
	"apollo/server/internal/model"
	"apollo/server/pkg/util"
	"errors"
	"github.com/labstack/echo/v4"
)

type Anssy struct {
	Type model.PartnerType
	req  *v1.PartnerRegisterReq
}

func NewAnssy(typ model.PartnerType, req *v1.PartnerRegisterReq) *Anssy {
	return &Anssy{
		Type: typ,
		req:  req,
	}
}

func (a Anssy) CheckParams(c echo.Context) error {
	if model.IsJDShop(a.req.ChannelId) {
		return errors.New("不支持京东通道")
	}

	return nil
}

func (a Anssy) GenPartner() model.Partner {
	p := model.Partner{
		Base: model.Base{
			Nickname: a.req.Nickname,
			Remark:   a.req.Remark,
		},
		Type:       a.Type,
		PayType:    a.req.PayType,
		ChannelId:  a.req.ChannelId,
		Balance:    0,
		Priority:   a.req.Priority,
		Level:      a.req.Level,
		PrivateKey: util.NewPrivateKey(),
		Anssy: model.Anssy{
			AnssyAppSecret: "519dd33a110440f1b6e0cbda500f8797",
		},
	}

	return p
}
