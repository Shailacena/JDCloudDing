package model

import (
	"fmt"

	"gorm.io/gorm"
)

type Partner struct {
	gorm.Model
	Base

	Type           PartnerType `gorm:"index;default:1;comment:合作商类型"`
	PayType        PayType     `gorm:"index;comment:支付类型"`
	ChannelId      ChannelId   `gorm:"index;comment:渠道"`
	Balance        float64     `gorm:"default:0;comment:余额"`
	Priority       int         `gorm:"default:0;comment:优先级"`
	SuperiorAgent  string      `gorm:"comment:上级代理"`
	Level          int         `gorm:"default1;comment:等级"`
	StockAmount    int64       `gorm:"default:0;comment:剩余库存金额"`
	RechargeTime   int64       `gorm:"default:0;comment:充值时间"`
	PrivateKey     string      `gorm:"comment:私钥"`
	DarkNumberLength int       `gorm:"default:11;comment:DarkNumber位数"`

	Agiso
	Anssy
}

func (*Partner) TableName() string {
	return "partner"
}

func (p *Partner) AfterCreate(tx *gorm.DB) error {
	p.Username = fmt.Sprintf("%d", p.ID)
	tx.Save(p)
	return nil
}
