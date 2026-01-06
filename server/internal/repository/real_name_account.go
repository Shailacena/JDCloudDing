package repository

import (
	v1 "apollo/server/api/v1"
	"apollo/server/internal/model"
	"apollo/server/pkg/data"

	"github.com/labstack/echo/v4"
)

var (
	RealNameAccount = &RealNameAccountRepo{}
)

type RealNameAccountRepo struct {
}

func (r *RealNameAccountRepo) Create(c echo.Context, list []*model.RealNameAccount) error {
	db := data.Instance()

	err := db.Create(list).Error

	return err
}

func (r *RealNameAccountRepo) List(c echo.Context, req *v1.ListRealNameAccountReq, parentIds []uint) ([]*model.RealNameAccount, int64, error) {
	db := data.Instance()

	var accounts []*model.RealNameAccount
	var total int64

	if len(parentIds) > 0 {
		db = db.Where("parent_id IN (?)", parentIds)
	}

	if err := db.Model(model.RealNameAccount{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if req.Offset() > 0 {
		db.Offset(req.Offset())
	}
	if req.Limit() > 0 {
		db.Limit(req.Limit())
	}

	err := db.Where("enable = ?", model.Enabled).Order("created_at desc").Find(&accounts).Error
	if err != nil {
		return nil, total, err
	}

	return accounts, total, err
}
