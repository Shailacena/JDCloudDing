package repository

import (
	"apollo/server/internal/model"
	"apollo/server/pkg/data"
	"apollo/server/pkg/headerx"
	"errors"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

var (
	Goods = &GoodsRepo{}
)

type GoodsRepo struct {
}

func (r *GoodsRepo) Create(c echo.Context, goods *model.Goods, isCheck bool) error {
	db := data.Instance()

	if isCheck {
		err := db.Where("sku_id = ?", goods.SkuId).First(&goods).Error
		if err == nil {
			return errors.New("商品的skuId已存在")
		}
	}

	err := db.Create(goods).Error

	return err
}

func (r *GoodsRepo) Update(c echo.Context, id uint, newGoods model.Goods) error {
	db := data.Instance()

	var goods model.Goods
	err := db.Where("id = ?", id).First(&goods).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("商品不存在")
		}
		return err
	}

	var g model.Goods
	err = db.Where("sku_id = ?", newGoods.SkuId).First(&g).Error
	if err == nil && id != g.ID {
		return errors.New("商品的skuId已存在")
	}

	err = db.Where("id = ?", id).Updates(newGoods).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *GoodsRepo) Delete(c echo.Context, id uint) error {
	db := data.Instance()

	goods := model.Goods{}
	return db.Where("id = ?", id).Delete(&goods).Error
}

func (r *GoodsRepo) List(c echo.Context, offset, limit int, partnerIds []uint, skuId string) ([]*model.Goods, int64, error) {
	header := headerx.GetDataFromHeader(c)

	db := data.Instance()

	var goodsList []*model.Goods
	var total int64

	if header.Role > 0 { // 管理员
		query := db.Model(&model.Goods{})
		if len(partnerIds) > 0 {
			query = query.Where("partner_id IN (?)", partnerIds)
		}

		if len(skuId) > 0 {
			query = query.Where("sku_id = ?", skuId)
		}

		// 查询总数
		if err := query.Count(&total).Error; err != nil {
			return nil, 0, err
		}

		if err := query.Offset(offset).Limit(limit).Order("created_at desc").Find(&goodsList).Error; err != nil {
			return nil, 0, err
		}
		return goodsList, total, nil

	} else {
		var partner model.Partner

		if err := db.Where("token = ?", header.Token).First(&partner).Error; err != nil {
			return nil, 0, err
		}

		query := db.Model(&model.Goods{})

		query = query.Where("partner_id = ?", partner.ID)

		if len(skuId) > 0 {
			query = query.Where("sku_id = ?", skuId)
		}

		// 先查询总数
		if err := query.Count(&total).Error; err != nil {
			return nil, 0, err
		}

		if err := query.Offset(offset).Limit(limit).Order("created_at desc").Find(&goodsList).Error; err != nil {
			return nil, 0, err
		}

		return goodsList, total, nil
	}
}

func (r *GoodsRepo) ResetWeight(c echo.Context, db *gorm.DB, partnerIds []uint) error {
	err := db.Model(&model.Goods{}).
		Where("partner_id IN (?)", partnerIds).
		Update("weight", 0).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *GoodsRepo) GetBySkuId(c echo.Context, db *gorm.DB, skuId string) (*model.Goods, error) {
	var goods model.Goods
	err := db.Where("sku_id = ?", skuId).Unscoped().First(&goods).Error
	if err != nil {
		return nil, err
	}

	return &goods, err
}
