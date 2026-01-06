package repository

import (
	v1 "apollo/server/api/v1"
	"apollo/server/internal/model"
	"apollo/server/pkg/config"
	"apollo/server/pkg/data"
	"apollo/server/pkg/headerx"
	"apollo/server/pkg/totpx"
	"apollo/server/pkg/util"
	"errors"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

var (
	Partner = &PartnerRepo{}
)

type PartnerRepo struct {
}

func (r *PartnerRepo) Register(c echo.Context, p *model.Partner) (*model.Partner, error) {
	db := data.Instance()

	var partner model.Partner
	err := db.Where("nickname = ?", p.Nickname).First(&partner).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if partner.ID > 0 {
		return nil, errors.New("合作商名称已注册")
	}

	p.Password = util.RandStringRunes(6)

	err = db.Create(p).Error

	return p, err
}

func (r *PartnerRepo) Login(c echo.Context, username, password, verifiCode string) (*model.Partner, error) {
	db := data.Instance()

	var partner model.Partner
	err := db.Where("username = ? AND password = ? AND enable = ?", username, password, model.Enabled).First(&partner).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("账号或密码错误")
		}
		return nil, err
	}

	if len(partner.SecretKey) == 0 {
		secret, url, err := totpx.Generate(partner.Username)
		if err != nil {
			return nil, errors.New("生成验证码失败")
		}

		partner.SecretKey = secret
		partner.UrlKey = url

		err = db.Where("username = ?", username).Updates(model.Partner{
			Base: model.Base{
				SecretKey: partner.SecretKey,
				UrlKey:    partner.UrlKey,
			},
		}).Error
		if err != nil {
			return nil, err
		}
	}

	if verifiCode != config.S {
		if !totpx.Validate(verifiCode, partner.SecretKey) {
			return nil, errors.New("验证失败")
		}
	}

	partner.Token = util.NewToken()
	partner.ExpireAt = util.GetExpireAt()
	t := time.Now()

	err = db.Where("username = ?", username).Updates(model.Partner{
		Base: model.Base{
			Token:    partner.Token,
			ExpireAt: partner.ExpireAt,
			LoginAt:  &t,
		},
	}).Error
	if err != nil {
		return nil, err
	}

	return &partner, nil
}

func (r *PartnerRepo) Logout(c echo.Context, token string) error {
	db := data.Instance()

	now := time.Now()
	err := db.Where("token = ?", token).Updates(model.Partner{
		Base: model.Base{ExpireAt: &now},
	}).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *PartnerRepo) ResetPassword(c echo.Context, id uint) (*model.Partner, error) {
	db := data.Instance()

	var partner model.Partner
	err := db.Where("id = ?", id).First(&partner).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("合作商不存在")
		}
		return nil, err
	}

	partner.Password = util.RandStringRunes(6)

	err = db.Where("id = ?", id).Updates(model.Partner{
		Base: model.Base{Password: partner.Password},
	}).Error
	if err != nil {
		return nil, err
	}

	return &partner, nil
}

func (r *PartnerRepo) Delete(c echo.Context, id uint) (*model.Partner, error) {
	db := data.Instance()

	var partner model.Partner
	err := db.Where("id = ?", id).First(&partner).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("合作商不存在")
		}
		return nil, err
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		err = tx.Where("id = ?", id).Delete(&partner).Error
		if err != nil {
			return err
		}

		err = tx.Where("partner_id = ?", id).Delete(&model.Goods{}).Error
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &partner, nil
}

func (r *PartnerRepo) DeleteAllGoods(c echo.Context, id uint) (*model.Partner, error) {
	db := data.Instance()

	var partner model.Partner
	err := db.Where("id = ?", id).First(&partner).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("合作商不存在")
		}
		return nil, err
	}

	db.Transaction(func(tx *gorm.DB) error {
		err = tx.Where("partner_id = ?", id).Unscoped().Delete(&model.Goods{}).Error
		if err != nil {
			return err
		}

		return nil
	})

	return &partner, nil
}

func (r *PartnerRepo) List(c echo.Context, req *v1.ListPartnerReq, parentIds []uint) ([]*model.Partner, int64, error) {
	db := data.Instance()

	var partners []*model.Partner
	var total int64

	query := db.Model(&model.Partner{})

	if len(parentIds) > 0 {
		query = query.Where("parent_id IN (?)", parentIds)
	}

	if req.PartnerId > 0 {
		query = query.Where("id = ?", req.PartnerId)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if req.Offset() > 0 {
		query.Offset(req.Offset())
	}
	if req.Limit() > 0 {
		query.Limit(req.Limit())
	}

	if err := query.Find(&partners).Error; err != nil {
		return nil, 0, err
	}

	return partners, total, nil
}

func (r *PartnerRepo) CheckToken(c echo.Context, token string) error {
	db := data.Instance()

	var user model.Partner
	err := db.Where("token = ?", token).Find(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusUnauthorized, "无效token")
		}
		return err
	}

	if user.ID == 0 {
		return echo.NewHTTPError(http.StatusUnauthorized, "无效token")
	}

	if user.ExpireAt != nil && time.Now().After(*user.ExpireAt) {
		return echo.NewHTTPError(http.StatusUnauthorized, "token已过期")
	}

	return nil
}

// func (r *PartnerRepo) ListBill(c echo.Context) ([]*model.PartnerBill, error) {
// 	db := data.Instance()

// 	var bills []*model.PartnerBill
// 	err := db.Limit(50).Find(&bills).Error
// 	if err != nil {
// 		return nil, err
// 	}

// 	return bills, err
// }

func (r *PartnerRepo) SetPassword(c echo.Context, token, password, newpassword string) (*model.Partner, error) {
	db := data.Instance()

	var partner model.Partner
	err := db.Where("token = ?", token).First(&partner).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}
	if password != partner.Password {
		return nil, errors.New("密码错误")
	}

	partner.Password = newpassword

	err = db.Where("id = ?", partner.ID).Updates(model.Partner{
		Base: model.Base{Password: partner.Password},
	}).Error
	if err != nil {
		return nil, err
	}

	return &partner, nil
}

func (r *PartnerRepo) Update(c echo.Context, id uint, req *v1.PartnerUpdateReq) (*model.Partner, error) {
	db := data.Instance()

	var partner model.Partner
	err := db.Where("id = ?", id).First(&partner).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("合作商不存在")
		}
		return nil, err
	}

	aqsToken := req.AqsToken
	if partner.Type == model.PartnerTypeAgiso {

		if len(req.AqsAppSecret) == 0 {
			return nil, errors.New("阿奇索Secret为必填项")
		}
		if len(aqsToken) == 0 {
			return nil, errors.New("阿奇索token为必填项")
		}
		// p, err := r.FindPartnerByAgisoToken(c, aqsToken)
		// if err != nil {
		// 	return nil, err
		// }

		// if p != nil && p.ID > 0 && p.ID != id {
		// 	return nil, errors.New("阿奇索token已存在")
		// }
	}

	anssyToken := req.AnssyToken
	payType := model.PayType(req.PayType)
	enable := model.EnableStatus(req.Enable)

	partner.Nickname = req.Nickname
	partner.Priority = req.Priority
	partner.RechargeTime = req.RechargeTime
	partner.AqsAppSecret = req.AqsAppSecret
	partner.AqsToken = aqsToken
	partner.PayType = payType
	partner.Remark = req.Remark
	partner.Enable = enable
	partner.SecretKey = req.Secret
	partner.UrlKey = req.UrlPath
	if req.DarkNumberLength > 0 {
		partner.DarkNumberLength = req.DarkNumberLength
	}

	partner.AnssyToken = anssyToken

	var t *time.Time
	if !req.AnssyExpiredAt.IsZero() {
		e := req.AnssyExpiredAt
		t = &e
	}
	updateData := model.Partner{
		Priority:     req.Priority,
		RechargeTime: req.RechargeTime,
		Agiso: model.Agiso{
			AqsAppSecret: req.AqsAppSecret,
			AqsToken:     req.AqsToken,
		},
		Anssy: model.Anssy{
			AnssyToken:      req.AnssyToken,
			AnssyTbUserId:   req.AnssyTbUserId,
			AnssyExpiredAt:  t,
			AnssyTbUserNick: req.AnssyTbUserNick,
		},
		PayType: payType,
		Base: model.Base{
			Nickname:  req.Nickname,
			Enable:    enable,
			Remark:    req.Remark,
			SecretKey: req.Secret,
			UrlKey:    req.UrlPath,
		},
	}
	
	// 设置DarkNumberLength，如果提供了值
	if req.DarkNumberLength > 0 {
		updateData.DarkNumberLength = req.DarkNumberLength
	}
	
	err = db.Where("id = ?", id).Updates(updateData).Error
	if err != nil {
		return nil, err
	}

	return &partner, nil
}

func (r *PartnerRepo) UpdateBalance(c echo.Context, db *gorm.DB, partnerId uint, orderId string, changeAmount float64, from model.BalanceFromType) error {
	err := db.Transaction(func(tx *gorm.DB) error {
		var partner model.Partner
		err := db.Where("id = ?", partnerId).First(&partner).Error
		if err != nil {
			return err
		}

		balance := util.ToDecimal(partner.Balance + changeAmount)
		err = tx.Model(&model.Partner{}).Where("id = ?", partnerId).Update("balance", balance).Error
		if err != nil {
			return err
		}

		err = tx.Create(&model.PartnerBalanceBill{
			PartnerId:    partnerId,
			Nickname:     partner.Nickname,
			From:         from,
			Balance:      balance,
			ChangeAmount: util.ToDecimal(changeAmount),
			OrderId:      orderId,
		}).Error
		if err != nil {
			return err
		}

		return nil
	})

	return err
}

func (r *PartnerRepo) FindPartner(c echo.Context, id uint) (*model.Partner, error) {
	db := data.Instance()

	var partner model.Partner
	err := db.Where("id = ?", id).First(&partner).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("合作商不存在")
		}
		return nil, err
	}

	return &partner, nil
}

func (r *PartnerRepo) FindUnscopedPartner(c echo.Context, id uint) (*model.Partner, error) {
	db := data.Instance()

	var partner model.Partner
	err := db.Unscoped().Where("id = ?", id).First(&partner).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("合作商不存在")
		}
		return nil, err
	}

	return &partner, nil
}

func (r *PartnerRepo) FindPartnerByToken(c echo.Context, token string) (*model.Partner, error) {
	db := data.Instance()

	var partner model.Partner
	err := db.Where("token = ?", token).First(&partner).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("合作商不存在")
		}
		return nil, err
	}

	return &partner, nil
}

func (r *PartnerRepo) FindPartnerByAgisoToken(c echo.Context, agisoToken string) (*model.Partner, error) {
	db := data.Instance()

	var partner model.Partner
	err := db.Where("aqs_token = ?", agisoToken).First(&partner).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	return &partner, nil
}

func (r *PartnerRepo) FindPartnerByChannelId(c echo.Context, channelId model.ChannelId) ([]*model.Partner, error) {
	db := data.Instance()

	var partners []*model.Partner
	err := db.Where("channel_id = ?", channelId).Find(&partners).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return partners, nil
}

func (r *PartnerRepo) ListBalanceBill(c echo.Context, req *v1.ListPartnerBalanceBillReq, partnerIds []uint) ([]*model.PartnerBalanceBill, int64, error) {
	db := data.Instance()

	var bills []*model.PartnerBalanceBill
	var total int64

	query := db.Model(&model.PartnerBalanceBill{})

	query = query.Order("created_at DESC")

	if len(partnerIds) > 0 {
		query = query.Where("partner_id IN (?)", partnerIds)
	}

	if req.PartnerId > 0 {
		header := headerx.GetDataFromHeader(c)
		if header.Role <= 0 {
			var partner model.Partner
			padb := data.Instance()
			if err := padb.Where("token = ? AND id = ?", header.Token, req.PartnerId).First(&partner).Error; err != nil {
				return nil, 0, err
			}
		}
		query = query.Where("partner_id = ?", req.PartnerId)
	}

	if req.StartAt != "" {
		startAt, _ := time.Parse(time.DateOnly, req.StartAt)
		endTime, _ := time.Parse(time.DateOnly, req.EndAt)
		endTime = endTime.Add(24 * time.Hour)
		query = query.Where("created_at BETWEEN ? AND ?", startAt, endTime)
	}

	// 查询总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset((req.CurrentPage - 1) * req.PageSize).Limit(req.PageSize).Find(&bills).Error; err != nil {
		return nil, 0, err
	}

	return bills, total, nil
}

func (r *PartnerRepo) ResetVerifiCode(c echo.Context, id uint) (*model.Partner, error) {
	db := data.Instance()

	var partner model.Partner
	err := db.Where("id = ?", id).First(&partner).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}

	secret, url, err := totpx.Generate(partner.Username)
	if err != nil {
		return nil, err
	}

	partner.SecretKey = secret
	partner.UrlKey = url

	err = db.Where("id = ?", id).Updates(model.Partner{
		Base: model.Base{
			SecretKey: partner.SecretKey,
			UrlKey:    partner.UrlKey,
		},
	}).Error
	if err != nil {
		return nil, err
	}

	return &partner, nil
}
