package model

import (
	"gorm.io/gorm"
)

type JDAccount struct {
	gorm.Model
	Account  string          `gorm:"index;comment:ck中pin字段"`
	WsKey    string          `gorm:"index;comment:ck中wskey字段"`
	Status   JDAccountStatus `gorm:"index;default:1"`
	Remark   string          `gorm:"comment:备注"`
	Weight   int             `gorm:"default:0;comment:权重"`
	UseCount int             `gorm:"default:0;comment:成功使用次数"`
	ParentId uint            `gorm:"index;default:0;comment:父id"`
	MasterId uint            `gorm:"index;default:0;comment:分身主id"`
}

func (JDAccount) TableName() string {
	return "jd_account"
}
