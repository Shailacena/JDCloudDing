package service

import (
	v1 "apollo/server/api/v1"
	"apollo/server/internal/model"
	"apollo/server/pkg/data"
	"apollo/server/pkg/util"
	"time"

	"github.com/labstack/echo/v4"
)

var (
	Status = new(StatusService)
)

type StatusService struct{}

// 获取服务器与数据库状态
func (s *StatusService) GetServerStatus(c echo.Context, _ *v1.GetServerStatusReq) (*v1.GetServerStatusResp, error) {
	db := data.Instance()

	// 数据库版本
	var version string
	// 兼容MySQL
	db.Raw("SELECT VERSION() AS version").Scan(&version)

	// 统计订单与实体数量
	var orderTotal, orderToday, successTodayCount int64
	var successTodayAmount float64
	var partnerTotal, merchantTotal int64

	_ = db.Model(&model.Order{}).Count(&orderTotal).Error

	today := time.Now()
	_ = db.Model(&model.Order{}).
		Where("DATE(CONVERT_TZ(created_at, '+00:00', '+08:00')) = DATE(CONVERT_TZ(?, '+00:00', '+08:00'))", today).
		Count(&orderToday).Error

	// 成功订单统计（当天）
	_ = db.Model(&model.Order{}).
		Select("COUNT(*)").
		Where("DATE(CONVERT_TZ(created_at, '+00:00', '+08:00')) = DATE(CONVERT_TZ(?, '+00:00', '+08:00')) AND status IN (?)", today, model.SuccessOrderStatus).
		Count(&successTodayCount).Error

	var totalAmount struct{ Total float64 }
	_ = db.Model(&model.Order{}).
		Select("COALESCE(SUM(received_amount), 0) AS total").
		Where("DATE(CONVERT_TZ(created_at, '+00:00', '+08:00')) = DATE(CONVERT_TZ(?, '+00:00', '+08:00')) AND status IN (?)", today, model.SuccessOrderStatus).
		Scan(&totalAmount).Error
	successTodayAmount = util.ToDecimal(totalAmount.Total)

	_ = db.Model(&model.Partner{}).Count(&partnerTotal).Error
	_ = db.Model(&model.Merchant{}).Count(&merchantTotal).Error

	return &v1.GetServerStatusResp{
		DbVersion:          version,
		OrderTotal:         orderTotal,
		OrderToday:         orderToday,
		SuccessTodayCount:  successTodayCount,
		SuccessTodayAmount: successTodayAmount,
		PartnerTotal:       partnerTotal,
		MerchantTotal:      merchantTotal,
	}, nil
}

// 当日订单5分钟趋势
func (s *StatusService) TodayTrend(c echo.Context, req *v1.GetTodayTrendReq) (*v1.GetTodayTrendResp, error) {
	db := data.Instance()

	if req.PartnerId > 0 {
		db = db.Where("partner_id = ?", req.PartnerId)
	}
	if req.MerchantId > 0 {
		db = db.Where("merchant_id = ?", req.MerchantId)
	}

	// 5分钟分桶：使用东八区时间
	// bucket = FROM_UNIXTIME(FLOOR(UNIX_TIMESTAMP(CONVERT_TZ(created_at,'+00:00','+08:00'))/300)*300)
	// 分组统计：订单数量与金额、成功订单数量与到账金额
	type row struct {
		Bucket        time.Time
		OrderCount    int64
		OrderAmount   float64
		SuccessCount  int64
		SuccessAmount float64
	}

	var rows []row
	err := db.Model(&model.Order{}).
		Select(`
            FROM_UNIXTIME(FLOOR(UNIX_TIMESTAMP(CONVERT_TZ(created_at,'+00:00','+08:00'))/300)*300) AS bucket,
            COUNT(*) AS order_count,
            COALESCE(SUM(amount), 0) AS order_amount,
            SUM(CASE WHEN status IN (?) THEN 1 ELSE 0 END) AS success_count,
            COALESCE(SUM(CASE WHEN status IN (?) THEN received_amount ELSE 0 END), 0) AS success_amount
        `, model.SuccessOrderStatus, model.SuccessOrderStatus).
		Where("DATE(CONVERT_TZ(created_at,'+00:00','+08:00')) = DATE(CONVERT_TZ(?, '+00:00', '+08:00'))", time.Now()).
		Group("bucket").
		Order("bucket").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	points := make([]*v1.TrendPoint, 0, len(rows))
	for _, r := range rows {
		points = append(points, &v1.TrendPoint{
			Time:          r.Bucket.Format("15:04"),
			OrderCount:    r.OrderCount,
			OrderAmount:   util.ToDecimal(r.OrderAmount),
			SuccessCount:  r.SuccessCount,
			SuccessAmount: util.ToDecimal(r.SuccessAmount),
		})
	}

	return &v1.GetTodayTrendResp{Points: points}, nil
}
