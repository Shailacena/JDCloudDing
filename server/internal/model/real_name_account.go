package model

import (
	"gorm.io/gorm"
)

type RealNameAccount struct {
	gorm.Model
	IdNumber      string       `gorm:"uniqueIndex;size:20;comment:身份证号码"`
	Name          string       `gorm:"comment:姓名"`
	Mobile        string       `gorm:"comment:手机号"`
	Address       string       `gorm:"comment:地址"`
	RealNameCount int64        `gorm:"default:0;comment:实名次数"`
	Enable        EnableStatus `gorm:"default:1;comment:用户是否被冻结 1正常 2冻结"`
	Remark        string       `gorm:"comment:备注"`
	ParentId      uint         `gorm:"index;default:0;comment:父id"`
	MasterId      uint         `gorm:"index;default:0;comment:分身主id"`
}

func (RealNameAccount) TableName() string {
	return "real_name_account"
}
