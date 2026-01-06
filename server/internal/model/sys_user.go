package model

import (
	"gorm.io/gorm"
)

type SysUser struct {
	gorm.Model

	Base

	Role SysUserRole `gorm:"index;comment:管理员角色"`
}

func (SysUser) TableName() string {
	return "sys_user"
}
