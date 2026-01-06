package repository

import (
	v1 "apollo/server/api/v1"
	"apollo/server/internal/model"
	"apollo/server/pkg/data"
	"apollo/server/pkg/util"
	"github.com/samber/lo"
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

var (
	Order = &OrderRepo{}
)

type OrderRepo struct {
}

func (r *OrderRepo) GetByOrderId(db *gorm.DB, orderId string) (*model.Order, error) {
	var order model.Order
	err := db.Where("order_id = ?", orderId).First(&order).Error
	if err != nil {
		return nil, err
	}

	return &order, err
}

func (r *OrderRepo) GetByMerchantOrderId(c echo.Context, merchantOrderId string) (*model.Order, error) {
	db := data.Instance()

	var order model.Order
	err := db.Where("merchant_order_id = ?", merchantOrderId).First(&order).Error
	if err != nil {
		return nil, err
	}

	return &order, err
}

func (r *OrderRepo) GetBySkuId(c echo.Context, db *gorm.DB, skuId string, offsetTime time.Time) (*model.Order, error) {
	var order model.Order
	err := db.Where("sku_id = ? AND created_at <= ?", skuId, offsetTime).Order("created_at desc").First(&order).Error
	if err != nil {
		return nil, err
	}

	return &order, err
}

func (r *OrderRepo) GetBySkuIdDarkNumber(c echo.Context, db *gorm.DB, skuId string, darkNumber string) (*model.Order, error) {
	var order model.Order
	err := db.Where("sku_id = ? AND dark_number = ?", skuId, darkNumber).Order("created_at desc").First(&order).Error
	if err != nil {
		return nil, err
	}

	return &order, err
}

func (r *OrderRepo) GetByPartnerOrderId(c echo.Context, db *gorm.DB, partnerOrderId string) (*model.Order, error) {
	var order model.Order
	err := db.Where("partner_order_id = ?", partnerOrderId).First(&order).Error
	if err != nil {
		return nil, err
	}

	return &order, err
}

type QueryMerchantAmountResult struct {
	MerchantId           uint
	TotalSuccessAmount   float64
	TodayOrderNum        float64
	TodayOrderAmount     float64
	TodaySuccessAmount   float64
	TodaySuccessOrderNum float64
}

func (r *OrderRepo) QueryResultByMerchant(c echo.Context, db *gorm.DB, merchantIds []uint, day time.Time) ([]QueryMerchantAmountResult, error) {
	var results []QueryMerchantAmountResult
	err := db.
		Model(&model.Order{}).
		// Select(`
		// 	merchant_id,
		// 	COUNT(*) AS today_order_num,
		// 	SUM(amount) AS today_order_amount,
		// 	SUM(CASE WHEN status IN (?) THEN amount ELSE 0 END) AS today_success_amount,
		// 	SUM(CASE WHEN status IN (?) THEN 1 ELSE 0 END) AS today_success_order_num
		// `, model.OrderStatusFinish, model.OrderStatusFinish).
		Select(`
			merchant_id,
			SUM(CASE WHEN status IN (?) THEN amount ELSE 0 END) AS today_success_amount
		`, model.OrderStatusFinish).
		Where("merchant_id IN (?) AND DATE(CONVERT_TZ(created_at, '+00:00', '+08:00')) = DATE(CONVERT_TZ(?, '+00:00', '+08:00'))", merchantIds, day).
		Group("merchant_id").
		Find(&results).Error
	if err != nil {
		return nil, err
	}

	return results, err
}

type QueryMerchantOrderResult struct {
	QueryMerchantAmountResult
	Date time.Time
}

func (r *OrderRepo) QueryMerchantOrder(c echo.Context, db *gorm.DB, merchantIds []uint, startAt, endAt time.Time) ([]QueryMerchantOrderResult, int64, error) {
	if len(merchantIds) > 0 {
		db = db.Where("merchant_id IN (?)", merchantIds)
	}

	if !startAt.IsZero() && !endAt.IsZero() {
		db = db.Where("created_at BETWEEN ? AND ?", startAt, endAt)
	}

	var results []QueryMerchantOrderResult
	var total int64
	db = db.
		Model(&model.Order{}).
		Select(`
		merchant_id,
		DATE(CONVERT_TZ(created_at, '+00:00', '+08:00')) AS date,
		COUNT(*) AS today_order_num,
		SUM(amount) AS today_order_amount,
		SUM(CASE WHEN status IN (?) THEN amount ELSE 0 END) AS today_success_amount,
		SUM(CASE WHEN status IN (?) THEN 1 ELSE 0 END) AS today_success_order_num
	`, model.OrderStatusFinish, model.OrderStatusFinish)

	err := db.Group("merchant_id, date").
		Order("date desc, merchant_id").
		Having("today_order_amount > ?", 0).
		Find(&results).Error
	if err != nil {
		c.Logger().Errorf("QueryMerchantOrder error=%s", err)
		return nil, 0, err
	}

	return results, total, nil
}

type QueryPartnerAmountResult struct {
	PartnerId            uint
	TodayOrderNum        float64
	TodayOrderAmount     float64
	TodaySuccessAmount   float64
	TodaySuccessOrderNum float64
}

func (r *OrderRepo) QueryResultByPartner(c echo.Context, db *gorm.DB, partnerIds []uint, day time.Time) ([]QueryPartnerAmountResult, error) {
	var results []QueryPartnerAmountResult
	err := db.
		Model(&model.Order{}).
		Select(`
			partner_id,
			COUNT(*) AS today_order_num,
			SUM(amount) AS today_order_amount,
			SUM(CASE WHEN status IN (?) THEN received_amount ELSE 0 END) AS today_success_amount,
			SUM(CASE WHEN status IN (?) THEN 1 ELSE 0 END) AS today_success_order_num
		`, model.SuccessOrderStatus, model.SuccessOrderStatus).
		Where("partner_id IN (?) AND DATE(CONVERT_TZ(created_at, '+00:00', '+08:00')) = DATE(CONVERT_TZ(?, '+00:00', '+08:00'))", partnerIds, day).
		Group("partner_id").
		Find(&results).Error
	if err != nil {
		c.Logger().Errorf("QueryResultByPartner error=%s", err)
	}

	return results, nil
}

type QueryPartnerOrderResult struct {
	QueryPartnerAmountResult
	Date time.Time
}

func (r *OrderRepo) QueryPartnerOrder(c echo.Context, db *gorm.DB, partnerIds []uint, startAt, endAt time.Time) ([]QueryPartnerOrderResult, int64, error) {
	if len(partnerIds) > 0 {
		db = db.Where("partner_id IN (?)", partnerIds)
	}

	if !startAt.IsZero() && !endAt.IsZero() {
		db = db.Where("created_at BETWEEN ? AND ?", startAt, endAt)
	}

	var results []QueryPartnerOrderResult
	var total int64
	db = db.
		Model(&model.Order{}).
		Select(`
		partner_id,
		DATE(CONVERT_TZ(created_at, '+00:00', '+08:00')) AS date,
		COUNT(*) AS today_order_num,
		SUM(amount) AS today_order_amount,
		SUM(CASE WHEN status IN (?) THEN received_amount ELSE 0 END) AS today_success_amount,
		SUM(CASE WHEN status IN (?) THEN 1 ELSE 0 END) AS today_success_order_num
	`, model.OrderStatusFinish, model.OrderStatusFinish)

	err := db.Group("partner_id, date").
		Order("date desc, partner_id").
		Having("today_order_amount > ?", 0).
		Find(&results).Error
	if err != nil {
		c.Logger().Errorf("QueryPartnerOrder error=%s", err)
		return nil, 0, err
	}

	return results, total, nil
}

type OrderStats struct {
	Last1HourTotal       int64
	Last1HourSuccess     int64
	Last30MinutesTotal   int64
	Last30MinutesSuccess int64
}

func (r *OrderRepo) QueryLastOrderResult(c echo.Context, db *gorm.DB, partnerIds []uint) map[uint]OrderStats {
	resultMap := make(map[uint]OrderStats)

	for _, id := range partnerIds {
		hour1 := time.Now().Add(-1 * time.Hour)
		var last1HourTotal, last1HourSuccess int64
		err := db.Model(&model.Order{}).
			Where("partner_id = ? AND created_at >= ? ", id, hour1).
			Count(&last1HourTotal).Error
		if err != nil {
			c.Logger().Errorf("last1HourTotal error=%s", err)
		}
		err = db.Model(&model.Order{}).
			Where("partner_id = ? AND created_at >= ? AND status IN (?)", id, hour1, model.SuccessOrderStatus).
			Count(&last1HourSuccess).Error
		if err != nil {
			c.Logger().Errorf("last1HourSuccess error=%s", err)
		}

		minute30 := time.Now().Add(-30 * time.Minute)
		var last30MinutesTotal, last30MinutesSuccess int64
		err = db.Model(&model.Order{}).
			Where("partner_id = ? AND created_at >= ?", id, minute30).
			Count(&last30MinutesTotal).Error
		if err != nil {
			c.Logger().Errorf("last30MinutesTotal error=%s", err)
		}
		err = db.Model(&model.Order{}).
			Where("partner_id = ? AND created_at >= ? AND status IN (?)", id, minute30, model.SuccessOrderStatus).
			Count(&last30MinutesSuccess).Error
		if err != nil {
			c.Logger().Errorf("last30MinutesSuccess error=%s", err)
		}

		resultMap[id] = OrderStats{
			Last1HourTotal:       last1HourTotal,
			Last1HourSuccess:     last1HourSuccess,
			Last30MinutesTotal:   last30MinutesTotal,
			Last30MinutesSuccess: last30MinutesSuccess,
		}
	}

	return resultMap
}

func (r *OrderRepo) List(c echo.Context, req *v1.ListOrderReq, parentIds []uint) ([]*model.Order, int64, error) {
	db := data.Instance()
	var total int64
	var orders []*model.Order

	db = db.Model(&model.Order{})

	if req.PartnerId > 0 {
		db = db.Where("partner_id = ?", req.PartnerId)
	}
	if req.MerchantId > 0 {
		db = db.Where("merchant_id = ?", req.MerchantId)
	}
	if req.PartnerId == 0 && req.MerchantId == 0 {
		merchants, _, err := Merchant.List(c, &v1.ListMerchantReq{}, parentIds)
		if err != nil {
			return nil, 0, err
		}

		merchantIds := lo.Map(merchants, func(item *model.Merchant, _ int) uint {
			return item.ID
		})

		partners, _, err := Partner.List(c, &v1.ListPartnerReq{}, parentIds)
		if err != nil {
			return nil, 0, err
		}

		partnerIds := lo.Map(partners, func(item *model.Partner, _ int) uint {
			return item.ID
		})

		db = db.Where("partner_id IN (?) OR merchant_id IN (?)", partnerIds, merchantIds)
	}

	if len(req.OrderId) > 0 {
		db = db.Where("order_id = ?", req.OrderId)
	}
	if len(req.PartnerOrderId) > 0 {
		db = db.Where("partner_order_id = ?", req.PartnerOrderId)
	}
	if len(req.MerchantOrderId) > 0 {
		db = db.Where("merchant_order_id = ?", req.MerchantOrderId)
	}
	if req.PayType > 0 {
		db = db.Where("pay_type = ?", req.PayType)
	}
	if req.OrderStatus > 0 {
		db = db.Where("status = ?", req.OrderStatus)
	}
	if req.NotifyStatus > 0 {
		db = db.Where("notify_status = ?", req.NotifyStatus)
	}
	if len(req.StartAt) > 0 && len(req.EndAt) > 0 {
		// 获取 PRC 时区
		prcLocation, _ := time.LoadLocation("Asia/Shanghai")
		
		// 尝试解析完整的时间格式 "2006-01-02 15:04:05"
		startAt, err := time.ParseInLocation("2006-01-02 15:04:05", req.StartAt, prcLocation)
		if err != nil {
			// 如果失败，尝试解析日期格式 "2006-01-02"
			startAt, _ = time.ParseInLocation(time.DateOnly, req.StartAt, prcLocation)
		}
		
		endTime, err := time.ParseInLocation("2006-01-02 15:04:05", req.EndAt, prcLocation)
		if err != nil {
			// 如果失败，尝试解析日期格式 "2006-01-02" 并添加24小时
			endTime, _ = time.ParseInLocation(time.DateOnly, req.EndAt, prcLocation)
			endTime = endTime.Add(24 * time.Hour)
		}
		
		db = db.Where("created_at BETWEEN ? AND ?", startAt, endTime)
	}

	// 查询总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := db.Offset(req.Offset()).Limit(req.Limit()).Order("created_at desc").Find(&orders).Error; err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

func (r *OrderRepo) Create(c echo.Context, order model.Order) error {
	db := data.Instance()

	err := db.Create(&order).Error
	if err != nil {
		return err
	}

	return err
}

func (r *OrderRepo) Update(db *gorm.DB, orderId string, order model.Order) error {
	err := db.Where("order_id = ?", orderId).Updates(order).Error
	if err != nil {
		return err
	}

	return err
}

type StatisticsReq struct {
	PartnerIds  []uint
	MerchantIds []uint
	StartAt     time.Time
	EndAt       time.Time
}

func (r *OrderRepo) Statistics(c echo.Context, req StatisticsReq) ([]*v1.BaseDailyBill, error) {
	db := data.Instance()
	var bills []*v1.BaseDailyBill
	startAt := req.StartAt
	endAt := req.EndAt

	if len(req.PartnerIds) > 0 {
		db = db.Where("partner_id IN (?)", req.PartnerIds)
	}
	if len(req.MerchantIds) > 0 {
		db = db.Where("merchant_id IN (?)", req.MerchantIds)
	}
	db = db.Where("created_at BETWEEN ? AND ?", startAt, endAt)

	err := db.
		Model(&model.Order{}).
		Select(`
		DATE_FORMAT(DATE(CONVERT_TZ(created_at, '+00:00', '+08:00')), '%Y-%m-%d') AS date,
		SUM(amount) AS total_order_amount,
		COUNT(*) AS total_order_num,
		SUM(CASE WHEN status IN (?) THEN received_amount ELSE 0 END) AS total_success_amount,
		SUM(CASE WHEN status IN (?) THEN 1 ELSE 0 END) AS total_success_order_num
	`, model.SuccessOrderStatus, model.SuccessOrderStatus).
		Group("date").
		Order("date desc").
		Find(&bills).Error
	if err != nil {
		return nil, err
	}
	return bills, nil
}

type OrderArchiveSummary struct {
    Total float64
}

func (r *OrderRepo) ArchiveByAdmin(c echo.Context, adminId uint) (*model.OrderArchive, error) {
    db := data.Instance()

    var admin model.SysUser
    if err := db.Where("id = ?", adminId).First(&admin).Error; err != nil {
        return nil, err
    }

    if admin.Role != model.NormalAdminRole || admin.MasterId != 0 {
        return nil, echo.NewHTTPError(403, "仅主账号可归档")
    }

    ids, err := Admin.FindAdminIds(c, admin.ID, admin.Role)
    if err != nil {
        return nil, err
    }

    query := "`order`.partner_id IN (SELECT id FROM partner WHERE parent_id IN ?) OR `order`.merchant_id IN (SELECT id FROM merchant WHERE parent_id IN ?)"

    var totalOrderCount int64
    if err := db.Model(&model.Order{}).Where(query, ids, ids).Count(&totalOrderCount).Error; err != nil {
        return nil, err
    }

    var summary OrderSummary
    if err := db.Model(&model.Order{}).Select("SUM(received_amount) as total_amount").Where(query, ids, ids).Where("status IN (?)", model.SuccessOrderStatus).Find(&summary).Error; err != nil {
        return nil, err
    }

    record := model.OrderArchive{
        MasterId:    adminId,
        ArchiveDate: time.Now(),
        TotalAmount: util.ToDecimal(summary.TotalAmount),
        OrderCount:  totalOrderCount,
    }

    if err := db.Transaction(func(tx *gorm.DB) error {
        if err := tx.Create(&record).Error; err != nil {
            return err
        }
        if err := tx.Where(query, ids, ids).Delete(&model.Order{}).Error; err != nil {
            return err
        }
        return nil
    }); err != nil {
        return nil, err
    }

    return &record, nil
}

func (r *OrderRepo) SumArchivedAmountByAdminIds(c echo.Context, adminIds []uint) (float64, error) {
    db := data.Instance()
    var total float64
    if len(adminIds) == 0 {
        return 0, nil
    }
    if err := db.Model(&model.OrderArchive{}).Select("COALESCE(SUM(total_amount),0)").Where("master_id IN (?)", adminIds).Scan(&total).Error; err != nil {
        return 0, err
    }
    return util.ToDecimal(total), nil
}

func (r *OrderRepo) SumAllArchivedAmount(c echo.Context) (float64, error) {
    db := data.Instance()
    var total float64
    if err := db.Model(&model.OrderArchive{}).Select("COALESCE(SUM(total_amount),0)").Scan(&total).Error; err != nil {
        return 0, err
    }
    return util.ToDecimal(total), nil
}

type OrderSummary struct {
	TotalAmount float64
}

func (r *OrderRepo) GetOrderSummary(c echo.Context, partnerIds, merchantIds []uint) (OrderSummary, error) {
	db := data.Instance()

	if len(partnerIds) > 0 {
		db = db.Where("partner_id IN (?)", partnerIds)
	}
	if len(merchantIds) > 0 {
		db = db.Where("merchant_id IN (?)", merchantIds)
	}

	var summary OrderSummary
	err := db.Model(&model.Order{}).Select("SUM(received_amount) as total_amount").Where("status IN (?)", model.SuccessOrderStatus).Find(&summary).Error
	if err != nil {
		return summary, err
	}

	return summary, err
}
