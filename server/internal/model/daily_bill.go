package model

import (
	"gorm.io/gorm"
)

type DailyBill struct {
	gorm.Model
	TotalMoney   int `gorm:"comment:成功总金额"`
	WxFee        int `gorm:"default:0;comment:微信缴费"`
	WxManualFee  int `gorm:"default:0;comment:微信手动缴费"`
	AliFee       int `gorm:"default:0;comment:支付宝缴费"`
	AliManualFee int `gorm:"default:0;comment:支付宝手动缴费"`
}

func (DailyBill) TableName() string {
	return "daily_bill"
}
