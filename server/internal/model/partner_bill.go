package model

import (
	"gorm.io/gorm"
)

type BillType int

const (
	BillTypePartner = iota + 1 // 合作商
)

type PartnerBill struct {
	gorm.Model
	PartnerId   uint     `gorm:"index;default:0;comment:合作商编号"`
	Type        BillType `gorm:"index;comment:用户类型"`
	ChangeMoney int      `gorm:"default:0;comment:变更金额"`
	Money       int      `gorm:"default:0;comment:当前金额"`
	Remark      string   `gorm:"comment:备注"`
}

func (PartnerBill) TableName() string {
	return "partner_bill"
}
