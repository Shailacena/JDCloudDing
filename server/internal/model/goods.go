package model

import (
	"gorm.io/gorm"
)

type Goods struct {
	gorm.Model
	PartnerId uint `gorm:"index;default:0;comment:合作商编号"`
	// 天猫、京东使用skuId、淘宝使用skuId、numId
	SkuId      string      `gorm:"uniqueIndex;size:30;comment:skuId"`
	NumId      string      `gorm:"index;comment:numId"`
	BrandId    string      `gorm:"comment:brandId"`
	Amount     float64     `gorm:"default:0;comment:商品金额"`
	RealAmount float64     `gorm:"default:0;comment:商品实际金额"`
	ShopName   string      `gorm:"comment:店铺名称"`
	Weight     int         `gorm:"default:0;comment:权重"`
	Status     GoodsStatus `gorm:"index;default:1;comment:状态"`
}

func (Goods) TableName() string {
	return "goods"
}
