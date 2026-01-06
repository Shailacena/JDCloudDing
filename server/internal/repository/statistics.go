package repository

import (
	"apollo/server/internal/model"
	"apollo/server/pkg/data"

	"github.com/labstack/echo/v4"
)

var (
	Statistics = &StatisticsRepo{}
)

type StatisticsRepo struct {
}

func (r *StatisticsRepo) ListDailyBill(c echo.Context) ([]*model.DailyBill, error) {
	db := data.Instance()

	var bills []*model.DailyBill
	err := db.Limit(20).Find(&bills).Error
	if err != nil {
		return nil, err
	}

	return bills, err
}
