package model

import (
	"time"
)

type Base struct {
	Username string `gorm:"size:20;uniqueIndex;not null;comment:登录账号"`
	Password string `gorm:"not null;comment:登录密码"`
	Nickname string `gorm:"comment:昵称"`

	Token    string       `gorm:"index;comment:登录token"`
	ExpireAt *time.Time   `gorm:"comment:token有效期"`
	Enable   EnableStatus `gorm:"default:1;comment:用户是否被冻结 1正常 2冻结"`
	Remark   string       `gorm:"comment:备注"`

	SecretKey string `gorm:"comment:密钥"`
	UrlKey    string `gorm:"comment:二维码链接"`

	LoginAt *time.Time `gorm:"comment:登录时间"`

	ParentId uint `gorm:"index;default:0;comment:父id"`
	MasterId uint `gorm:"index;default:0;comment:分身主id"`
}

var SuccessOrderStatus = []OrderStatus{OrderStatusPaid, OrderStatusFinish}

func IsJDShop(channelId ChannelId) bool {
	switch channelId {
	case ChannelJDPay:
		return true
	default:
		return false
	}
}

type Agiso struct {
	AqsAppSecret string `gorm:"index;comment:阿奇索密钥"`
	AqsToken     string `gorm:"index;comment:阿奇索token"`
}

type Anssy struct {
	AnssyAppSecret  string     `gorm:"index;comment:安式密钥"` // 本地程序生成
	AnssyToken      string     `gorm:"index;comment:安式token"`
	AnssyExpiredAt  *time.Time `gorm:"comment:token有效期"`
	AnssyTbUserId   string     `gorm:"comment:淘宝账户id"`
	AnssyTbUserNick string     `gorm:"comment:店铺昵称"`
}
