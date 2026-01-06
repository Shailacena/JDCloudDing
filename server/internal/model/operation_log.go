package model

import (
	"gorm.io/gorm"
)

type OperationLog struct {
	gorm.Model

	IP        string `gorm:"comment:ip"`
	Operation string `gorm:"comment:操作项"`
	Operator  uint   `gorm:"comment:操作人"`
}

func (OperationLog) TableName() string {
	return "operation_log"
}
