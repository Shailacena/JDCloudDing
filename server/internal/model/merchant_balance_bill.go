package model

import (
	"gorm.io/gorm"
)

type MerchantBalanceBill struct {
	gorm.Model

	MerchantId   uint            `gorm:"index;comment:合作商id"`
	OrderId      string          `gorm:"comment:订单id"`
	Nickname     string          `gorm:"comment:名称"`
	From         BalanceFromType `gorm:"index;comment:操作类型"`
	Balance      float64         `gorm:"comment:账户余额"`
	ChangeAmount float64         `gorm:"comment:交易金额"`
}

func (MerchantBalanceBill) TableName() string {
	return "merchant_balance_bill"
}
