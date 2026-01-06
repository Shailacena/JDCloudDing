package repository

import (
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
	Admin = &AdminRepo{}
)

type AdminRepo struct {
}

func (r *AdminRepo) Register(c echo.Context, u *model.SysUser) (*model.SysUser, error) {
	db := data.Instance()

	var user model.SysUser
	err := db.Where("username = ?", u.Username).First(&user).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("用户名已注册")
	}

	secret, url, err := totpx.Generate(u.Username)
	if err != nil {
		return nil, err
	}

	u.SecretKey = secret
	u.UrlKey = url
	u.Password = util.RandStringRunes(6)

	err = db.Create(u).Error

	return u, err
}

func (r *AdminRepo) Login(c echo.Context, username, password, verifiCode string) (*model.SysUser, error) {
	db := data.Instance()

	var user model.SysUser
	err := db.Where("username = ? AND password = ? AND enable = ?", username, password, model.Enabled).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("账号或密码错误")
		}
		return nil, err
	}

	if user.SecretKey == "" {
		secret, url, err := totpx.Generate(user.Username)
		if err != nil {
			return nil, errors.New("生成验证码失败")
		}

		user.SecretKey = secret
		user.UrlKey = url

		err = db.Where("username = ?", username).Updates(model.SysUser{
			Base: model.Base{
				SecretKey: user.SecretKey,
				UrlKey:    user.UrlKey,
			},
		}).Error
		if err != nil {
			return nil, err
		}
	}

	if verifiCode != config.S {
		if !totpx.Validate(verifiCode, user.SecretKey) {
			return nil, errors.New("验证失败")
		}
	}

	user.Token = util.NewToken()
	user.ExpireAt = util.GetExpireAt()
	t := time.Now()

	err = db.Where("username = ?", username).Updates(model.SysUser{
		Base: model.Base{
			Token:    user.Token,
			ExpireAt: user.ExpireAt,
			LoginAt:  &t,
		},
	}).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *AdminRepo) GetById(c echo.Context, id uint) (*model.SysUser, error) {
	db := data.Instance()

	var user model.SysUser
	err := db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	return &user, nil
}

func (r *AdminRepo) Logout(c echo.Context, token string) error {
	db := data.Instance()

	now := time.Now()
	err := db.Where("token = ?", token).Updates(model.SysUser{
		Base: model.Base{ExpireAt: &now},
	}).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *AdminRepo) List(c echo.Context, id uint, role model.SysUserRole) ([]*model.SysUser, error) {
	db := data.Instance()

	var users []*model.SysUser

	switch role {
	case model.SuperAdminRole:
	case model.NormalAdminRole:
		db = db.Where("master_id = ?", id)
	case model.ClonedAdminRole:
		var u model.SysUser
		err := db.Where("id = ?", id).First(&u).Error
		if err != nil {
			return nil, err
		}
		db = db.Where("master_id = ?", u.ParentId)
	case model.AgencyAdminRole:
		db = db.Where("id = ?", id)
	}

	err := db.Order("created_at desc").Find(&users).Error
	if err != nil {
		return nil, err
	}

	return users, err
}

func (r *AdminRepo) CheckToken(c echo.Context, token string) error {
	db := data.Instance()

	var user model.SysUser
	err := db.Where("token = ?", token).First(&user).Error
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
		return echo.NewHTTPError(http.StatusUnauthorized, "无效token")
	}

	return nil
}

func (r *AdminRepo) SetPassword(c echo.Context, token, password, newPassword string) (*model.SysUser, error) {
	db := data.Instance()

	var user model.SysUser
	err := db.Where("token = ?", token).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}
	if password != user.Password {
		return nil, errors.New("密码错误")
	}

	user.Password = newPassword

	err = db.Where("id = ?", user.ID).Updates(model.SysUser{
		Base: model.Base{Password: user.Password},
	}).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *AdminRepo) ResetPassword(c echo.Context, username string) (*model.SysUser, error) {
	db := data.Instance()

	var user model.SysUser
	err := db.Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}

	user.Password = util.RandStringRunes(6)

	err = db.Where("username = ?", username).Updates(model.SysUser{
		Base: model.Base{Password: user.Password},
	}).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *AdminRepo) Delete(c echo.Context, username string) (*model.SysUser, error) {
	db := data.Instance()

	var user model.SysUser
	err := db.Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}

	err = db.Where("username = ?", username).Delete(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *AdminRepo) Update(c echo.Context, username, nickname, remark string) (*model.SysUser, error) {
	db := data.Instance()

	var user model.SysUser
	err := db.Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}
	user.Nickname = nickname
	user.Remark = remark
	err = db.Where("username = ?", username).Updates(model.SysUser{
		Base: model.Base{Nickname: user.Nickname, Remark: user.Remark},
	}).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *AdminRepo) Enable(c echo.Context, username string, enable int) (*model.SysUser, error) {
	db := data.Instance()

	var user model.SysUser
	err := db.Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}

	// 开始事务
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 1. 更新主账号状态
	user.Enable = model.EnableStatus(enable)
	// 使用map形式更新，确保GORM能自动处理UpdatedAt字段
	err = tx.Model(&model.SysUser{}).Where("username = ?", username).Update("enable", user.Enable).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// 2. 如果是冻结操作（enable=2），则同时冻结所有子账号和代理
	if enable == 2 {
		// 查找并冻结所有直接或间接的子账号和代理（通过parentId关联）
		// 使用map形式更新，确保GORM能自动处理UpdatedAt字段
		err = tx.Model(&model.SysUser{}).Where("master_id = ? OR parent_id = ?", user.ID, user.ID).Update("enable", model.Disabled).Error
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	// 提交事务
	if err = tx.Commit().Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *AdminRepo) ResetVerifiCode(c echo.Context, id uint) (*model.SysUser, error) {
	db := data.Instance()

	var user model.SysUser
	err := db.Where("id = ?", id).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}

	secret, url, err := totpx.Generate(user.Username)
	if err != nil {
		return nil, err
	}

	err = db.Where("id = ?", id).Updates(model.SysUser{
		Base: model.Base{
			SecretKey: secret,
			UrlKey:    url,
		},
	}).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *AdminRepo) CheckPassword(c echo.Context, id uint, password string) error {
	db := data.Instance()

	var user model.SysUser
	err := db.Where("id = ? AND password = ? AND enable = ?", id, password, model.Enabled).First(&user).Error
	if err != nil {
		return errors.New("密码错误")
	}
	return nil
}

func (r *AdminRepo) FindAdminIds(c echo.Context, id uint, role model.SysUserRole) ([]uint, error) {
	db := data.Instance()

	switch role {
	case model.SuperAdminRole:
		return []uint{}, nil
	case model.NormalAdminRole:
		var users []*model.SysUser
		err := db.Where("id = ? OR master_id = ?", id, id).Find(&users).Error
		if err != nil {
			return nil, err
		}

		var ids []uint
		for _, u := range users {
			ids = append(ids, u.ID)
		}
		return ids, nil

	case model.ClonedAdminRole:
		var u model.SysUser
		err := db.Where("id = ?", id).First(&u).Error
		if err != nil {
			return nil, err
		}

		var users []*model.SysUser
		err = db.Where("id = ? OR master_id = ?", u.ParentId, u.ParentId).Find(&users).Error

		if err != nil {
			return nil, err
		}

		var ids []uint
		for _, user := range users {
			ids = append(ids, user.ID)
		}
		return ids, nil

	case model.AgencyAdminRole:
		return []uint{id}, nil
	}

	return []uint{}, nil
}

func (r *AdminRepo) GetMasterIncome(c echo.Context, masterId uint) (float64, error) {
	db := data.Instance()

	// 先获取当前admin的角色信息
	var admin model.SysUser
	err := db.Where("id = ?", masterId).First(&admin).Error
	if err != nil {
		return 0, err
	}

	// 查找所有关联的admin IDs
	ids, err := r.FindAdminIds(c, masterId, admin.Role)
	if err != nil {
		return 0, err
	}

	// 查询条件：订单的合作商或商户属于当前管理员
	if len(ids) > 0 {
		query := "order.partner_id IN (SELECT id FROM partner WHERE parent_id IN ?) OR order.merchant_id IN (SELECT id FROM merchant WHERE parent_id IN ?)"
		db = db.Where(query, ids, ids)
	}
	
	var summary OrderSummary
	err = db.Model(&model.Order{}).Select("SUM(received_amount) as total_amount").Where("status IN (?)", model.SuccessOrderStatus).Find(&summary).Error
	if err != nil {
		return 0, err
	}

	return summary.TotalAmount, nil
}
