package model

import (
	"gorm.io/gorm"
	"time"
)

type Notify struct {
	gorm.Model

	BizId        string        `json:"index;biz_id;comment:业务id"`
	BizType      NotifyBizType `json:"index;biz_type;comment:业务类型"`
	ExpiredAt    time.Time     `gorm:"index;default:null;comment:过期时间"`
	NotifyAt     time.Time     `gorm:"default:null;comment:通知时间"`
	NotifyStatus NotifyStatus  `gorm:"index;default:1;comment:通知状态"`
}

func (Notify) TableName() string {
	return "notify"
}
