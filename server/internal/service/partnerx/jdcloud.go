package partnerx

import (
	v1 "apollo/server/api/v1"
	"apollo/server/internal/model"
	"apollo/server/pkg/util"
	"errors"

	"github.com/labstack/echo/v4"
)

type JDCloud struct {
	Type model.PartnerType
	req  *v1.PartnerRegisterReq
}

func NewJDCloud(typ model.PartnerType, req *v1.PartnerRegisterReq) *JDCloud {
	return &JDCloud{
		Type: typ,
		req:  req,
	}
}

func (j JDCloud) CheckParams(c echo.Context) error {
	if j.req.ChannelId != model.ChannelJDCard {
		return errors.New("京东云鼎类型必须选择京东云鼎通道")
	}
	return nil
}

func (j JDCloud) GenPartner() model.Partner {
	p := model.Partner{
		Base: model.Base{
			Nickname: j.req.Nickname,
			Remark:   j.req.Remark,
		},
		Type:       j.Type,
		Balance:    0,
		Priority:   j.req.Priority,
		Level:      j.req.Level,
		PayType:    j.req.PayType,
		ChannelId:  j.req.ChannelId,
		PrivateKey: util.NewPrivateKey(),
	}

	return p
}
