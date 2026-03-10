package model

import (
	"time"

	"gorm.io/gorm"
)

type CardType string

const (
	CardTypeReal    CardType = "real"
	CardTypeVirtual CardType = "virtual"
)

type PriceCard struct {
	gorm.Model
	CardNo     string    `gorm:"size:50;index;comment:卡号"`
	Password   string    `gorm:"size:50;comment:密码"`
	CardGroup  string    `gorm:"size:50;index;comment:卡组"`
	Amount     float64   `gorm:"type:decimal(10,2);comment:面额"`
	CardType   CardType  `gorm:"size:10;index;comment:卡密类型:real真实,virtual虚拟"`
	BatchNo    string    `gorm:"size:20;index;comment:批次号(YYYYMMDD)"`
	UsedStatus bool      `gorm:"default:false;comment:使用状态:false未使用,true已使用"`
	OrderId    string    `gorm:"size:50;index;comment:订单ID"`
	PartnerId  uint      `gorm:"index;comment:合作商ID"`
	UseIP      string    `gorm:"comment:使用IP"`
	UsedAt     *time.Time `gorm:"comment:使用时间"`
}

func (PriceCard) TableName() string {
	return "price_card"
}
