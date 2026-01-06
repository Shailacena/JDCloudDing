package service

import (
	v1 "apollo/server/api/v1"
	"apollo/server/internal/model"
	"apollo/server/internal/repository"
	"apollo/server/pkg/data"
	"apollo/server/pkg/timex"
	"fmt"
	"github.com/golang-module/carbon"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
)

var (
	Statistics = new(StatisticsService)
)

type StatisticsService struct {
}

func (s *StatisticsService) List(c echo.Context, req *v1.ListDailyBillReq) (*v1.ListDailyBillResp, error) {
	var partnerIds, merchantIds []uint
	if req.PartnerId > 0 {
		partnerIds = append(partnerIds, req.PartnerId)
	}
	if req.MerchantId > 0 {
		merchantIds = append(merchantIds, req.MerchantId)
	}
	if req.PartnerId == 0 && req.MerchantId == 0 {
		parentIds, _ := Admin.FindParentIds(c)
		merchants, _, err := repository.Merchant.List(c, &v1.ListMerchantReq{}, parentIds)
		if err != nil {
			return nil, err
		}

		merchantIds = lo.Map(merchants, func(item *model.Merchant, _ int) uint {
			return item.ID
		})

		partners, _, err := repository.Partner.List(c, &v1.ListPartnerReq{}, parentIds)
		if err != nil {
			return nil, err
		}

		partnerIds = lo.Map(partners, func(item *model.Partner, _ int) uint {
			return item.ID
		})

		fmt.Println(parentIds, partnerIds, merchantIds)
		if len(partnerIds) == 0 && len(merchantIds) == 0 {
			return &v1.ListDailyBillResp{}, nil
		}
	}

	days := 14
	now := timex.GetPRCNowTime()
	results, err := repository.Order.Statistics(c, repository.StatisticsReq{
		PartnerIds:  partnerIds,
		MerchantIds: merchantIds,
		StartAt:     now.AddDays(-days).StartOfDay().Carbon2Time(),
		EndAt:       now.Carbon2Time(),
	})
	if err != nil {
		return nil, err
	}

	resultMap := lo.SliceToMap(results, func(item *v1.BaseDailyBill) (string, *v1.BaseDailyBill) {
		return item.Date, item
	})

	list := make([]*v1.BaseDailyBill, 0, days)
	for i := 0; i < days; i++ {
		prcNow := timex.GetPRCNowTime().Carbon2Time()
		date := prcNow.AddDate(0, 0, -i)
		dateStr := date.Format(time.DateOnly)

		d, ok := resultMap[dateStr]
		if ok {
			list = append(list, d)
			continue
		}

		list = append(list, &v1.BaseDailyBill{
			Date: dateStr,
		})

	}

	return &v1.ListDailyBillResp{
		List: list,
	}, nil
}

func (s *StatisticsService) ListByPartner(c echo.Context, req *v1.ListDailyBillByPartnerReq) (*v1.ListDailyBillByPartnerResp, error) {
	var partnerIds []uint
	if req.PartnerId > 0 {
		partnerIds = append(partnerIds, req.PartnerId)
	} else {
		parentIds, _ := Admin.FindParentIds(c)

		partners, _, err := repository.Partner.List(c, &v1.ListPartnerReq{}, parentIds)
		if err != nil {
			return nil, err
		}

		partnerIds = lo.Map(partners, func(item *model.Partner, _ int) uint {
			return item.ID
		})

		fmt.Println(parentIds, partnerIds)
		if len(partnerIds) == 0 {
			return &v1.ListDailyBillByPartnerResp{}, nil
		}
	}

	endAt := timex.GetPRCNowTime()
	startAt := endAt.StartOfDay()
	if len(req.StartAt) > 0 {
		startAt = carbon.ParseByLayout(req.StartAt, time.DateOnly, carbon.PRC)
	}
	if len(req.EndAt) > 0 {
		endAt = carbon.ParseByLayout(req.EndAt, time.DateOnly, carbon.PRC)
		endAt = endAt.AddDay()
	}

	db := data.Instance()
	result, total, err := repository.Order.QueryPartnerOrder(c, db, partnerIds, startAt.Carbon2Time(), endAt.Carbon2Time())
	if err != nil {
		return nil, err
	}

	list := make([]*v1.DailyBill, 0, len(result))
	for _, r := range result {
		if r.PartnerId <= 0 {
			total--
			continue
		}

		p, err := repository.Partner.FindUnscopedPartner(c, r.PartnerId)
		if err != nil {
			continue
		}

		list = append(list, &v1.DailyBill{
			Id:       r.PartnerId,
			Nickname: p.Nickname,
			Balance:  p.Balance,
			Time:     r.Date.Unix(),
			BaseDailyBill: v1.BaseDailyBill{
				TotalOrderNum:        r.TodayOrderNum,
				TotalOrderAmount:     r.TodayOrderAmount,
				TotalSuccessAmount:   r.TodaySuccessAmount,
				TotalSuccessOrderNum: r.TodaySuccessOrderNum,
			},
		})
	}

	return &v1.ListDailyBillByPartnerResp{
		ListTableData: v1.ListTableData[v1.DailyBill]{
			List:  list,
			Total: int64(len(list)),
		},
	}, nil
}

func (s *StatisticsService) ListByMerchant(c echo.Context, req *v1.ListDailyBillByMerchantReq) (*v1.ListDailyBillByMerchantResp, error) {
	var merchantIds []uint
	if req.MerchantId > 0 {
		merchantIds = append(merchantIds, req.MerchantId)
	} else {
		parentIds, _ := Admin.FindParentIds(c)

		merchants, _, err := repository.Merchant.List(c, &v1.ListMerchantReq{}, parentIds)
		if err != nil {
			return nil, err
		}

		merchantIds = lo.Map(merchants, func(item *model.Merchant, _ int) uint {
			return item.ID
		})

		if len(merchantIds) == 0 {
			return &v1.ListDailyBillByMerchantResp{}, nil
		}
		fmt.Println(merchantIds)
	}

	endAt := timex.GetPRCNowTime()
	startAt := endAt.StartOfDay()
	if len(req.StartAt) > 0 {
		startAt = carbon.ParseByLayout(req.StartAt, time.DateOnly, carbon.PRC)
	}
	if len(req.EndAt) > 0 {
		endAt = carbon.ParseByLayout(req.EndAt, time.DateOnly, carbon.PRC)
		endAt = endAt.AddDay()
	}

	db := data.Instance()
	result, total, err := repository.Order.QueryMerchantOrder(c, db, merchantIds, startAt.Carbon2Time(), endAt.Carbon2Time())
	if err != nil {
		return nil, err
	}

	list := make([]*v1.DailyBill, 0, len(result))
	for _, r := range result {
		if r.MerchantId <= 0 {
			total--
			continue
		}

		m, err := repository.Merchant.FindUnscopedMerchant(c, r.MerchantId)
		if err != nil {
			continue
		}

		list = append(list, &v1.DailyBill{
			Id:       r.MerchantId,
			Nickname: m.Nickname,
			Balance:  m.Balance,
			Time:     r.Date.Unix(),
			BaseDailyBill: v1.BaseDailyBill{
				TotalOrderNum:        r.TodayOrderNum,
				TotalOrderAmount:     r.TodayOrderAmount,
				TotalSuccessAmount:   r.TodaySuccessAmount,
				TotalSuccessOrderNum: r.TodaySuccessOrderNum,
			},
		})
	}

	return &v1.ListDailyBillByMerchantResp{
		ListTableData: v1.ListTableData[v1.DailyBill]{
			List:  list,
			Total: int64(len(list)),
		},
	}, nil
}
