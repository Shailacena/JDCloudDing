package partnerx

import (
	v1 "apollo/server/api/v1"
	"apollo/server/internal/model"
	"apollo/server/pkg/util"
	"errors"

	"github.com/labstack/echo/v4"
)

type Agiso struct {
	Type model.PartnerType
	req  *v1.PartnerRegisterReq
}

func NewAgiso(typ model.PartnerType, req *v1.PartnerRegisterReq) *Agiso {
	return &Agiso{
		Type: typ,
		req:  req,
	}
}

func (a Agiso) CheckParams(c echo.Context) error {
	if len(a.req.AqsAppSecret) == 0 {
		return errors.New("阿奇索Secret为必填项")
	}

	aqsToken := a.req.AqsToken
	if len(aqsToken) == 0 {
		return errors.New("阿奇索token为必填项")
	}

	// p, err := repository.Partner.FindPartnerByAgisoToken(c, aqsToken)
	// if err != nil {
	// 	return err
	// }

	// if p != nil && p.ID > 0 {
	// 	return errors.New("阿奇索token已存在")
	// }

	return nil
}

func (a Agiso) GenPartner() model.Partner {
	aqsToken := a.req.AqsToken

	p := model.Partner{
		Base: model.Base{
			Nickname: a.req.Nickname,
			Remark:   a.req.Remark,
		},
		Type:       a.Type,
		Balance:    0,
		Priority:   a.req.Priority,
		Level:      a.req.Level,
		PayType:    a.req.PayType,
		ChannelId:  a.req.ChannelId,
		PrivateKey: util.NewPrivateKey(),
		Agiso: model.Agiso{
			AqsAppSecret: a.req.AqsAppSecret,
			AqsToken:     aqsToken,
		},
	}

	return p
}
