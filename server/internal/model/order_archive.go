package model

import (
    "time"
    "gorm.io/gorm"
)

type OrderArchive struct {
    gorm.Model
    MasterId    uint      `gorm:"index;comment:归档所属主账号ID"`
    ArchiveDate time.Time `gorm:"index;comment:归档日期"`
    TotalAmount float64   `gorm:"default:0;comment:归档金额"`
    OrderCount  int64     `gorm:"default:0;comment:归档订单数"`
}

func (OrderArchive) TableName() string {
    return "order_archive"
}