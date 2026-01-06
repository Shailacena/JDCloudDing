package model

import (
	"time"

	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	OrderId         string        `gorm:"size:32;uniqueIndex;comment:订单号"`
	ChannelId       ChannelId     `gorm:"index;comment:渠道id"`
	MerchantId      uint          `gorm:"index;comment:商户号"`
	MerchantOrderId string        `gorm:"size:32;uniqueIndex;comment:商户订单号"`
	PartnerOrderId  string        `gorm:"index;comment:合作商订单号"`
	Amount          float64       `gorm:"default:0;comment:订单金额"`
	ReceivedAmount  float64       `gorm:"default:0;comment:实收金额"`
	PayType         PayType       `gorm:"comment:支付类型"`
	PayAccount      string        `gorm:"index;comment:下单账号"`
	Status          OrderStatus   `gorm:"default:1;comment:订单状态"`
	SkuId           string        `gorm:"index;comment:skuId"`
	PartnerId       uint          `gorm:"index;comment:合作商id"`
	MerchantName    string        `gorm:"comment:商户名称"`
	PartnerName     string        `gorm:"comment:合作商名称"`
	Shop            string        `gorm:"comment:店铺"`
	NotifyStatus    NotifyStatus  `gorm:"default:1;comment:回调状态"`
	NotifyUrl       string        `gorm:"comment:回调地址"`
	PayUrl          string        `gorm:"comment:支付地址"`
	ExtParam        string        `gorm:"comment:回调参数"`
	IP              string        `gorm:"comment:ip"`
	Device          string        `gorm:"comment:设备类型"`
	DarkNumber      string        `gorm:"index;comment:暗号码"`
	Remark          string        `gorm:"comment:备注"`
	PayAt           time.Time     `gorm:"default:null;comment:支付时间"`
	NotifyAt        time.Time     `gorm:"default:null;comment:通知时间"`
	EndLockAt       time.Time     `gorm:"default:null;comment:结束锁定时间"`
	ConfirmStatus   ConfirmStatus `gorm:"default:1;comment:确认收货状态"`
}

func (Order) TableName() string {
	return "order"
}
