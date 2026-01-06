package service

import (
    v1 "apollo/server/api/v1"
    "apollo/server/internal/model"
    "apollo/server/internal/repository"
    "apollo/server/pkg/contextx"
    "apollo/server/pkg/data"
    "apollo/server/pkg/headerx"
    "apollo/server/pkg/util"
    "errors"
    "fmt"
    "time"

    "apollo/server/cmd/payment/api/common"

    "github.com/labstack/echo/v4"
    "gorm.io/gorm"
)

var (
	Order = new(OrderService)
)

type OrderService struct {
}

func (s *OrderService) List(c echo.Context, req *v1.ListOrderReq) (*v1.ListOrderResp, error) {
	parentIds, _ := Admin.FindParentIds(c)

	orderList, total, err := repository.Order.List(c, req, parentIds)
	if err != nil {
		return nil, err
	}

	list := make([]*v1.Order, 0, len(orderList))
	for _, o := range orderList {
		var payAt int64
		if !o.PayAt.IsZero() {
			payAt = o.PayAt.Unix()
		}

		list = append(list, &v1.Order{
			OrderId:         o.OrderId,
			MerchantId:      o.MerchantId,
			PartnerId:       o.PartnerId,
			MerchantOrderId: o.MerchantOrderId,
			PartnerOrderId:  o.PartnerOrderId,
			Amount:          util.ToDecimal(o.Amount),
			ReceivedAmount:  util.ToDecimal(o.ReceivedAmount),
			PayType:         int(o.PayType),
			PayAccount:      o.PayAccount,
			Status:          uint(o.Status),
			SkuId:           o.SkuId,
			Shop:            o.Shop,
			NotifyStatus:    uint(o.NotifyStatus),
			IP:              o.IP,
			Device:          o.Device,
			Remark:          o.Remark,
			CreateAt:        o.CreatedAt.Unix(),
			PayAt:           payAt,
			ChannelId:       string(o.ChannelId),
			MerchantName:    o.MerchantName,
			PartnerName:     o.PartnerName,
		})
	}

	return &v1.ListOrderResp{
		ListTableData: v1.ListTableData[v1.Order]{
			List:  list,
			Total: total,
		},
	}, nil
}

func (s *OrderService) Confirm(c echo.Context, req *v1.ConfirmOrderReq) (*v1.ConfirmOrderResp, error) {
	db := data.Instance()
	o, err := repository.Order.GetByOrderId(db, req.OrderId)
	if err != nil {
		return nil, err
	}

	if o.Status == model.OrderStatusRefundSuccessful {
		return nil, errors.New("订单已退款")
	}

	if o.Status == model.OrderStatusFinish && o.NotifyStatus == model.NotifyDone {
		return nil, errors.New("订单已完成并且已经通知商户")
	}

	if (o.Status == model.OrderStatusPaid || o.Status == model.OrderStatusFinish) && o.NotifyStatus == model.NotNotify {
		// 只通知
		if o.ReceivedAmount == 0 {
			o.ReceivedAmount = util.ToDecimal(o.Amount)
		}
		err = common.NotifyMerchant(contextx.NewContextFromEcho(c), db, o)
		if err != nil {
			return nil, err
		}
	} else if (o.Status == model.OrderStatusUnpaid || o.Status == model.OrderStatusRefundFailed) && o.NotifyStatus == model.NotNotify {
		// 结算并通知
		err = db.Transaction(func(tx *gorm.DB) error {
			o1, err := repository.Order.GetByOrderId(tx, req.OrderId)
			if err != nil {
				return err
			}

			if o.ReceivedAmount == 0 {
				o1.ReceivedAmount = util.ToDecimal(o1.Amount)
			}
			o = o1

			err = repository.Order.Update(tx, o1.OrderId, model.Order{
				Status:         model.OrderStatusPaid,
				ReceivedAmount: util.ToDecimal(o1.ReceivedAmount),
				ConfirmStatus:  model.ConfirmStatusManual,
				PayAt:          time.Now(),
			})
			if err != nil {
				return err
			}

			err = common.UpdatePartnerBalance(c, tx, o1)
			if err != nil {
				return err
			}

			return nil
		})
		if err != nil {
			c.Logger().Errorf("Confirm Transaction error=%s", err.Error())
			return nil, err
		}

		err = common.NotifyMerchant(contextx.NewContextFromEcho(c), db, o)
		if err != nil {
			return nil, fmt.Errorf("confirm NotifyMerchant error=%s", err.Error())
		}
	}

	return &v1.ConfirmOrderResp{}, nil
}

func (s *OrderService) GetOrderSummary(c echo.Context, _ *v1.GetOrderSummaryReq) (*v1.GetOrderSummaryesp, error) {
    parentIds, _ := Admin.FindParentIds(c)
    db := data.Instance().Unscoped()
    if len(parentIds) > 0 {
        query := "order.partner_id IN (SELECT id FROM partner WHERE parent_id IN ?) OR order.merchant_id IN (SELECT id FROM merchant WHERE parent_id IN ?)"
        db = db.Where(query, parentIds, parentIds)
    }
    var summary repository.OrderSummary
    err := db.Model(&model.Order{}).Select("SUM(received_amount) as total_amount").Where("status IN (?)", model.SuccessOrderStatus).Find(&summary).Error
    if err != nil {
        return nil, err
    }
    var archivedTotal float64
    if len(parentIds) == 0 {
        archivedTotal, err = repository.Order.SumAllArchivedAmount(c)
    } else {
        archivedTotal, err = repository.Order.SumArchivedAmountByAdminIds(c, parentIds)
    }
    if err != nil {
        return nil, err
    }
    return &v1.GetOrderSummaryesp{
        TotalAmount: util.ToDecimal(summary.TotalAmount + archivedTotal),
    }, nil
}

func (s *OrderService) Archive(c echo.Context, req *v1.ArchiveOrdersReq) (*v1.ArchiveOrdersResp, error) {
    header := headerx.GetDataFromHeader(c)
    if model.SysUserRole(header.Role) != model.SuperAdminRole {
        return nil, errors.New("仅超级管理员可归档")
    }

    target, err := repository.Admin.GetById(c, req.AdminId)
    if err != nil {
        return nil, err
    }
    if target.Role != model.NormalAdminRole || target.MasterId != 0 {
        return nil, errors.New("仅支持归档主账号订单")
    }

    record, err := repository.Order.ArchiveByAdmin(c, req.AdminId)
    if err != nil {
        return nil, err
    }
    return &v1.ArchiveOrdersResp{
        ArchiveDate: record.ArchiveDate.Format("2006-01-02"),
        TotalAmount: record.TotalAmount,
        OrderCount:  record.OrderCount,
    }, nil
}
