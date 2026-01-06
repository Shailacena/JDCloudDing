package repository

import (
	v1 "apollo/server/api/v1"
	"apollo/server/internal/model"
	"apollo/server/pkg/data"
	"github.com/labstack/echo/v4"
)

var (
	OperationLog = &OperationLogRepo{}
)

type OperationLogRepo struct {
}

func (r *OperationLogRepo) Create(c echo.Context, logs []*model.OperationLog) error {
	db := data.Instance()

	err := db.Create(&logs).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *OperationLogRepo) List(c echo.Context, pagination v1.Pagination, operatorIds []uint) ([]*model.OperationLog, int64, error) {
	db := data.Instance()

	var logs []*model.OperationLog
	var total int64

	db = db.Model(&model.OperationLog{})

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if len(operatorIds) > 0 {
		db = db.Where("operator IN (?)", operatorIds)
	}
	if pagination.Offset() > 0 {
		db.Offset(pagination.Offset())
	}
	if pagination.Limit() > 0 {
		db.Limit(pagination.Limit())
	}

	err := db.Order("created_at desc").Find(&logs).Error
	if err != nil {
		return nil, 0, err
	}

	return logs, total, err
}
