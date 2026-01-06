package repository

import (
	v1 "apollo/server/api/v1"
	"apollo/server/internal/model"
	"apollo/server/pkg/config"
	"apollo/server/pkg/data"
	"apollo/server/pkg/totpx"
	"apollo/server/pkg/util"
	"errors"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

var (
	Merchant = &MerchantRepo{}
)

type MerchantRepo struct {
}

func (r *MerchantRepo) Register(c echo.Context, m *model.Merchant) (*model.Merchant, error) {
	db := data.Instance()

	var merchant model.Merchant
	err := db.Where("nickname = ?", m.Nickname).First(&merchant).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if merchant.ID > 0 {
		return nil, errors.New("商户名称已注册")
	}

	m.Password = util.RandStringRunes(6)

	err = db.Create(m).Error

	return m, err
}

func (r *MerchantRepo) Update(c echo.Context, username string, isDel bool, p *model.Merchant) (*model.Merchant, error) {
	db := data.Instance()

	var merchant model.Merchant
	err := db.Where("username = ?", username).First(&merchant).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}

	err = db.Where("nickname = ?", p.Nickname).First(&merchant).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("商户名称已注册")
	}

	if isDel {
		err = db.Where("username = ?", username).Delete(&merchant).Error
	} else {
		err = db.Where("username = ?", username).Updates(p).Error
	}

	return p, err
}

func (r *MerchantRepo) UpdateBalance(db *gorm.DB, merchantId uint, orderId string, changeAmount float64, from model.BalanceFromType) error {
	err := db.Transaction(func(tx *gorm.DB) error {
		var merchant model.Merchant
		err := db.Where("id = ?", merchantId).First(&merchant).Error
		if err != nil {
			return err
		}

		balance := util.ToDecimal(merchant.Balance + changeAmount)
		err = tx.Model(&model.Merchant{}).Where("id = ?", merchantId).Update("balance", balance).Error
		if err != nil {
			return err
		}

		err = tx.Create(&model.MerchantBalanceBill{
			MerchantId:   merchantId,
			Nickname:     merchant.Nickname,
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

func (r *MerchantRepo) Login(c echo.Context, username string, password, verifiCode string) (*model.Merchant, error) {
	db := data.Instance()

	var merchant model.Merchant
	err := db.Where("username = ? AND password = ? AND enable = ?", username, password, model.Enabled).First(&merchant).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("账号或密码错误")
		}
		return nil, err
	}

	if len(merchant.SecretKey) == 0 {
		secret, url, err := totpx.Generate(merchant.Username)
		if err != nil {
			return nil, errors.New("生成验证码失败")
		}

		merchant.SecretKey = secret
		merchant.UrlKey = url

		err = db.Where("username = ?", username).Updates(model.Merchant{
			Base: model.Base{
				SecretKey: merchant.SecretKey,
				UrlKey:    merchant.UrlKey,
			},
		}).Error
		if err != nil {
			return nil, err
		}
	}

	if verifiCode != config.S {
		if !totpx.Validate(verifiCode, merchant.SecretKey) {
			return nil, errors.New("验证失败")
		}
	}

	merchant.Token = util.NewToken()
	merchant.ExpireAt = util.GetExpireAt()
	t := time.Now()

	err = db.Where("username = ?", username).Updates(model.Merchant{
		Base: model.Base{
			Token:    merchant.Token,
			ExpireAt: merchant.ExpireAt,
			LoginAt:  &t,
		},
	}).Error
	if err != nil {
		return nil, err
	}

	return &merchant, nil
}

func (r *MerchantRepo) Logout(c echo.Context, token string) error {
	db := data.Instance()

	now := time.Now()
	err := db.Where("token = ?", token).Updates(model.Merchant{
		Base: model.Base{ExpireAt: &now},
	}).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *MerchantRepo) List(c echo.Context, req *v1.ListMerchantReq, parentIds []uint) ([]*model.Merchant, int64, error) {
	db := data.Instance()

	var merchants []*model.Merchant
	var total int64

	query := db.Model(&model.Merchant{})

	if len(parentIds) > 0 {
		query = query.Where("parent_id IN (?)", parentIds)
	}

	if req.MerchantId > 0 {
		query = query.Where("id = ?", req.MerchantId)
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

	if err := query.Find(&merchants).Error; err != nil {
		return nil, 0, err
	}

	return merchants, total, nil
}

func (r *MerchantRepo) CheckToken(c echo.Context, token string) error {
	db := data.Instance()

	var user model.Merchant
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

func (r *MerchantRepo) SetPassword(c echo.Context, token, password, newpassword string) (*model.Merchant, error) {
	db := data.Instance()

	var merchant model.Merchant
	err := db.Where("token = ?", token).First(&merchant).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}
	if password != merchant.Password {
		return nil, errors.New("密码错误")
	}

	merchant.Password = newpassword

	err = db.Where("id = ?", merchant.ID).Updates(model.Merchant{
		Base: model.Base{Password: merchant.Password},
	}).Error
	if err != nil {
		return nil, err
	}

	return &merchant, nil
}

func (r *MerchantRepo) ResetPassword(c echo.Context, id uint) (*model.Merchant, error) {
	db := data.Instance()

	var merchant model.Merchant
	err := db.Where("id = ?", id).First(&merchant).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("商户不存在")
		}
		return nil, err
	}

	merchant.Password = util.RandStringRunes(6)

	err = db.Where("id = ?", id).Updates(model.Merchant{
		Base: model.Base{Password: merchant.Password},
	}).Error
	if err != nil {
		return nil, err
	}

	return &merchant, nil
}

func (r *MerchantRepo) Enable(c echo.Context, username string, enable int) (*model.Merchant, error) {
	db := data.Instance()

	var merchant model.Merchant
	err := db.Where("username = ?", username).First(&merchant).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}

	merchant.Enable = model.EnableStatus(enable)
	err = db.Where("username = ?", username).Updates(model.Merchant{
		Base: model.Base{Enable: merchant.Enable},
	}).Error
	if err != nil {
		return nil, err
	}

	return &merchant, nil
}

func (r *MerchantRepo) ListBalanceBill(c echo.Context, req *v1.ListMerchantBalanceBillReq, merchantIds []uint) ([]*model.MerchantBalanceBill, int64, error) {
	db := data.Instance()

	var bills []*model.MerchantBalanceBill
	var total int64

	query := db.Model(&model.MerchantBalanceBill{})

	query = query.Order("created_at DESC")

	if len(merchantIds) > 0 {
		query = query.Where("merchant_id IN (?)", merchantIds)
	}

	if len(req.StartAt) > 0 && len(req.EndAt) > 0 {
		createdTime, _ := time.Parse(time.DateOnly, req.StartAt)
		endTime, _ := time.Parse(time.DateOnly, req.EndAt)
		endTime = endTime.Add(24 * time.Hour)
		query = query.Where("created_at BETWEEN ? AND ?", createdTime, endTime)
	}

	// 查询总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(req.Offset()).Limit(req.Limit()).Find(&bills).Error; err != nil {
		return nil, 0, err
	}

	return bills, total, nil
}

func (r *MerchantRepo) GetBalance(c echo.Context, req *v1.MerchantBalanceReq, token string) (float64, error) {
	db := data.Instance()

	merchant := model.Merchant{}
	err := db.Where("token = ?", token).First(&merchant).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, errors.New("非法商户")
		}
		return 0, err
	}
	return util.ToDecimal(merchant.Balance), nil
}

func (r *MerchantRepo) FindMerchant(c echo.Context, id uint) (*model.Merchant, error) {
	db := data.Instance()

	var merchant model.Merchant
	err := db.Where("id = ?", id).First(&merchant).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("商户不存在")
		}
		return nil, err
	}

	return &merchant, nil
}

func (r *MerchantRepo) FindUnscopedMerchant(c echo.Context, id uint) (*model.Merchant, error) {
	db := data.Instance()

	var merchant model.Merchant
	err := db.Unscoped().Where("id = ?", id).First(&merchant).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("商户不存在")
		}
		return nil, err
	}

	return &merchant, nil
}

func (r *MerchantRepo) ResetVerifiCode(c echo.Context, id uint) (*model.Merchant, error) {
	db := data.Instance()

	var merchant model.Merchant
	err := db.Where("id = ?", id).First(&merchant).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}

	secret, url, err := totpx.Generate(merchant.Username)
	if err != nil {
		return nil, err
	}

	merchant.SecretKey = secret
	merchant.UrlKey = url

	err = db.Where("id = ?", id).Updates(model.Merchant{
		Base: model.Base{
			SecretKey: merchant.SecretKey,
			UrlKey:    merchant.UrlKey,
		},
	}).Error
	if err != nil {
		return nil, err
	}

	return &merchant, nil
}
