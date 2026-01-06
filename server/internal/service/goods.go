package service

import (
	v1 "apollo/server/api/v1"
	"apollo/server/internal/model"
	"apollo/server/internal/repository"
	"apollo/server/pkg/util"
	"fmt"
	"github.com/samber/lo"

	"github.com/labstack/echo/v4"
)

var (
	Goods = new(GoodsService)
)

type GoodsService struct {
}

func (s *GoodsService) Create(c echo.Context, req *v1.GoodsCreateReq) (*v1.GoodsCreateResp, error) {
	goods := &model.Goods{
		PartnerId:  req.PartnerId,
		SkuId:      req.SkuId,
		BrandId:    req.BrandId,
		Amount:     util.ToDecimal(req.Price),
		RealAmount: util.ToDecimal(req.RealPrice),
		ShopName:   req.ShopName,
		Status:     model.GoodsStatus(req.Status),
	}

	err := repository.Goods.Create(c, goods, true)
	if err != nil {
		return nil, err
	}

	return &v1.GoodsCreateResp{}, nil
}

func (s *GoodsService) Update(c echo.Context, req *v1.GoodsUpdateReq) (*v1.GoodsUpdateResp, error) {
	goods := model.Goods{
		PartnerId:  req.PartnerId,
		SkuId:      req.SkuId,
		BrandId:    req.BrandId,
		Amount:     util.ToDecimal(req.Price),
		RealAmount: util.ToDecimal(req.RealPrice),
		ShopName:   req.ShopName,
		Status:     model.GoodsStatus(req.Status),
	}

	err := repository.Goods.Update(c, req.Id, goods)
	if err != nil {
		return nil, err
	}

	return &v1.GoodsUpdateResp{}, nil
}

func (s *GoodsService) Delete(c echo.Context, req *v1.GoodsDeleteReq) (*v1.GoodsDeleteResp, error) {
	err := repository.Goods.Delete(c, req.Id)
	if err != nil {
		return nil, err
	}

	return &v1.GoodsDeleteResp{}, nil
}

func (s *GoodsService) List(c echo.Context, req *v1.ListGoodsReq) (*v1.ListGoodsResp, error) {
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
			return &v1.ListGoodsResp{}, nil
		}
	}

	goodsList, total, err := repository.Goods.List(c, req.Offset(), req.Limit(), partnerIds, req.SkuId)
	if err != nil {
		return nil, err
	}

	list := make([]*v1.Goods, 0, len(goodsList))
	for _, g := range goodsList {
		list = append(list, &v1.Goods{
			Id: g.ID,
			GoodsInfo: v1.GoodsInfo{
				PartnerId: g.PartnerId,
				SkuId:     g.SkuId,
				BrandId:   g.BrandId,
				Price:     util.ToDecimal(g.Amount),
				RealPrice: util.ToDecimal(g.RealAmount),
				ShopName:  g.ShopName,
				Status:    int(g.Status),
			},
			CreateAt: g.CreatedAt.Unix(),
		})
	}

	return &v1.ListGoodsResp{
		ListTableData: v1.ListTableData[v1.Goods]{
			List:  list,
			Total: total,
		},
	}, nil
}
