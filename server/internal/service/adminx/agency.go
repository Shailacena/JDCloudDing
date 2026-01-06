package adminx

import (
	v1 "apollo/server/api/v1"
	"apollo/server/internal/model"
	"fmt"
)

type Agency struct {
	creator *model.SysUser
}

func NewAgency(creator *model.SysUser) *Agency {
	return &Agency{
		creator: creator,
	}
}

func (a *Agency) CheckCreator(creator *model.SysUser) error {
	if creator == nil || (creator.Role != model.ClonedAdminRole) {
		return fmt.Errorf("权限不足，请切换为管理员分身")
	}

	return nil
}

func (a *Agency) Gen(req *v1.AdminRegisterReq) model.SysUser {
	masterId := a.creator.MasterId
	if masterId == 0 {
		masterId = a.creator.ID
	}

	return model.SysUser{
		Role: model.SysUserRole(req.Role),
		Base: model.Base{
			Username: req.Username,
			Nickname: req.Nickname,
			Remark:   req.Remark,
			ParentId: a.creator.ID,
			MasterId: masterId,
		},
	}
}
