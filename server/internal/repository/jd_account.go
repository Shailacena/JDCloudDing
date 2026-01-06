package repository

import (
	v1 "apollo/server/api/v1"
	"apollo/server/internal/model"
	"apollo/server/pkg/data"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

var (
	JDAccount = &JDAccountRepo{}
)

type JDAccountRepo struct {
}

func (r *JDAccountRepo) Create(c echo.Context, list []*model.JDAccount) error {
	db := data.Instance()

	err := db.Create(list).Error

	return err
}

func (r *JDAccountRepo) Enable(c echo.Context, id uint, status model.JDAccountStatus) error {
	db := data.Instance()

	err := db.Where("id = ?", id).Updates(model.JDAccount{Status: status}).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *JDAccountRepo) List(c echo.Context, req *v1.ListJDAccountReq, parentIds []uint) ([]*model.JDAccount, int64, error) {
	db := data.Instance()

	var accounts []*model.JDAccount
	var total int64

	db = filterJDAccount(db, &req.JDAccountSearchParams)
	if len(parentIds) > 0 {
		db = db.Where("parent_id IN (?)", parentIds)
	}

	if err := db.Model(model.JDAccount{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if req.Offset() > 0 {
		db.Offset(req.Offset())
	}
	if req.Limit() > 0 {
		db.Limit(req.Limit())
	}

	err := db.Find(&accounts).Error
	if err != nil {
		return nil, 0, err
	}

	return accounts, total, err
}

func filterJDAccount(db *gorm.DB, req *v1.JDAccountSearchParams) *gorm.DB {
	if req.Id > 0 {
		db = db.Where("id = ?", req.Id)
	}
	if len(req.Account) > 0 {
		db = db.Where("account = ?", req.Account)
	}
	if len(req.Status) > 0 {
		db = db.Where("status IN (?)", req.Status)
	}

	return db
}

func (r *JDAccountRepo) Delete(c echo.Context, req *v1.JDAccountDeleteReq) error {
	db := data.Instance()

	var err error
	if req.IsAll {
		err = db.Where("1 = 1").Delete(&model.JDAccount{}).Error
	} else {
		db = filterJDAccount(db, &req.JDAccountSearchParams)
		err = db.Delete(&model.JDAccount{}).Error
	}

	return err
}

func (r *JDAccountRepo) Reset(c echo.Context, req *v1.JDAccountResetReq) error {
	db := data.Instance()

	db = filterJDAccount(db, &req.JDAccountSearchParams)

	err := db.Model(&model.JDAccount{}).Update("status", model.JDAccountStatusNormal).Error

	return err
}

func (r *JDAccountRepo) UseCount(id uint) error {
	db := data.Instance()

	err := db.Model(&model.JDAccount{}).Where("id = ?", id).Update("use_count", gorm.Expr("use_count + ?", 1)).Error
	if err != nil {
		return err
	}

	return nil
}
