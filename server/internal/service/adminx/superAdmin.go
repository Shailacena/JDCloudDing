package adminx

import (
	v1 "apollo/server/api/v1"
	"apollo/server/internal/model"
	"fmt"
)

type SuperAdmin struct {
}

func NewSuperAdmin() *SuperAdmin {
	return &SuperAdmin{}
}

func (a *SuperAdmin) CheckCreator(creator *model.SysUser) error {
	if creator == nil || creator.Role != model.SuperAdminRole {
		return fmt.Errorf("权限不足，请切换为超级管理员")
	}

	return nil
}

func (a *SuperAdmin) Gen(req *v1.AdminRegisterReq) model.SysUser {

	return model.SysUser{
		Role: model.SysUserRole(req.Role),
		Base: model.Base{
			Username: req.Username,
			Nickname: req.Nickname,
			Remark:   req.Remark,
		},
	}
}
