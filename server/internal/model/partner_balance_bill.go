package model

import (
	"gorm.io/gorm"
)

type PartnerBalanceBill struct {
	gorm.Model

	PartnerId    uint            `gorm:"index;comment:合作商id"`
	OrderId      string          `gorm:"comment:订单id"`
	Nickname     string          `gorm:"comment:名称"`
	From         BalanceFromType `gorm:"index;comment:来源类型"`
	Balance      float64         `gorm:"comment:账户余额"`
	ChangeAmount float64         `gorm:"comment:交易金额"`
}

func (PartnerBalanceBill) TableName() string {
	return "partner_balance_bill"
}
