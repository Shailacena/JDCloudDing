package adminx

import (
	v1 "apollo/server/api/v1"
	"apollo/server/internal/model"
)

type IAdminGenerator interface {
	CheckCreator(creator *model.SysUser) error
	Gen(req *v1.AdminRegisterReq) model.SysUser
}

type IAdminList struct {
}
