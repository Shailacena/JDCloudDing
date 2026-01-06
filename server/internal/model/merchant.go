package model

import (
	"fmt"

	"gorm.io/gorm"
)

type Merchant struct {
	gorm.Model

	Base

	PrivateKey string  `gorm:"comment:秘钥"`
	Balance    float64 `gorm:"default:0;comment:余额"`
}

func (*Merchant) TableName() string {
	return "merchant"
}

func (m *Merchant) AfterCreate(tx *gorm.DB) error {
	m.Username = fmt.Sprintf("%d", m.ID)
	tx.Save(m)
	return nil
}
