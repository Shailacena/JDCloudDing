package service

import (
	v1 "apollo/server/api/v1"
	"apollo/server/internal/model"
	"apollo/server/internal/repository"
	"apollo/server/internal/service/adminx"
	"apollo/server/pkg/headerx"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

var (
	Admin = new(AdminService)
)

type AdminService struct {
}

func (s *AdminService) Register(c echo.Context, req *v1.AdminRegisterReq) (*v1.AdminRegisterResp, error) {
	header := headerx.GetDataFromHeader(c)
	adminId := header.AdminId

	role := model.SysUserRole(req.Role)

	creator, err := repository.Admin.GetById(c, adminId)
	if err != nil {
		return nil, err
	}

	var admin adminx.IAdminGenerator
	switch role {
	case model.SuperAdminRole:
		admin = adminx.NewSuperAdmin()
	case model.NormalAdminRole:
		admin = adminx.NewAdmin(creator)
	case model.ClonedAdminRole:
		admin = adminx.NewClonedAdmin(creator)
	case model.AgencyAdminRole:
		admin = adminx.NewAgency(creator)
	}

	if admin == nil {
		return nil, fmt.Errorf("非法类型")
	}

	err = admin.CheckCreator(creator)
	if err != nil {
		return nil, err
	}

	u := admin.Gen(req)

	newUser, err := repository.Admin.Register(c, &u)
	if err != nil {
		return nil, err
	}

	return &v1.AdminRegisterResp{
		Username: req.Username,
		Nickname: req.Nickname,
		Password: newUser.Password,
	}, nil
}

func (s *AdminService) Register11(c echo.Context, req *v1.AdminRegisterReq) (*v1.AdminRegisterResp, error) {
	// Register11 is used to create the first super admin without authentication
	admin := adminx.NewSuperAdmin()
	u := admin.Gen(req)

	newUser, err := repository.Admin.Register(c, &u)
	if err != nil {
		return nil, err
	}

	return &v1.AdminRegisterResp{
		Username: req.Username,
		Nickname: req.Nickname,
		Password: newUser.Password,
	}, nil
}

func (s *AdminService) Login(c echo.Context, req *v1.AdminLoginReq) (*v1.AdminLoginResp, error) {
	user, err := repository.Admin.Login(c, req.Username, req.Password, req.VerifiCode)
	if err != nil {
		return nil, err
	}

	logs := make([]*model.OperationLog, 0)
	logs = append(logs, &model.OperationLog{
		IP:        c.RealIP(),
		Operation: c.Request().RequestURI,
		Operator:  user.ID,
	})

	err = repository.OperationLog.Create(c, logs)
	if err != nil {
		c.Logger().Warn("OperationLog.Create error", err)
	}

	return &v1.AdminLoginResp{
		Id:       user.ID,
		Token:    user.Token,
		Nickname: user.Nickname,
		Role:     int(user.Role),
	}, nil
}

func (s *AdminService) Logout(c echo.Context, req *v1.AdminLogoutReq, token string) (*v1.AdminLogoutResp, error) {
	err := repository.Admin.Logout(c, token)
	if err != nil {
		return nil, err
	}

	return &v1.AdminLogoutResp{}, nil
}

func (s *AdminService) List(c echo.Context, req *v1.ListAdminReq) (*v1.ListAdminResp, error) {
	header := headerx.GetDataFromHeader(c)
	r := header.Role
	adminId := header.AdminId

	role := model.SysUserRole(r)

	users, err := repository.Admin.List(c, adminId, role)
	if err != nil {
		return nil, err
	}

	list := make([]*v1.Admin, 0, len(users))
	for _, u := range users {
		list = append(list, &v1.Admin{
			Id:       u.ID,
			Username: u.Username,
			Nickname: u.Nickname,
			Remark:   u.Remark,
			Enable:   int(u.Enable),
			Role:     int(u.Role),
			UrlKey:   u.UrlKey,
			ParentId: u.ParentId,
		})
	}

	return &v1.ListAdminResp{
		ListTableData: v1.ListTableData[v1.Admin]{
			List:  list,
			Total: int64(len(list)),
		},
	}, nil
}

func (s *AdminService) SetPassword(c echo.Context, req *v1.AdminSetPasswordReq, token string) (*v1.AdminSetPasswordResp, error) {
	if len(req.NewPassword) < 6 {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "密码应大于6位")
	}

	u, err := repository.Admin.SetPassword(c, token, req.OldPassword, req.NewPassword)
	if err != nil {
		return nil, err
	}

	logs := make([]*model.OperationLog, 0)
	logs = append(logs, &model.OperationLog{
		IP:        c.RealIP(),
		Operation: c.Request().RequestURI,
		Operator:  u.ID,
	})

	err = repository.OperationLog.Create(c, logs)
	if err != nil {
		c.Logger().Warn("OperationLog.Create error", err)
	}

	return &v1.AdminSetPasswordResp{}, nil
}

func (s *AdminService) ResetPassword(c echo.Context, req *v1.AdminResetPasswordReq) (*v1.AdminResetPasswordResp, error) {
	user, err := repository.Admin.ResetPassword(c, req.Username)
	if err != nil {
		return nil, err
	}

	return &v1.AdminResetPasswordResp{
		Password: user.Password,
	}, nil
}

func (s *AdminService) Delete(c echo.Context, req *v1.AdminDeleteReq) (*v1.AdminDeleteResp, error) {
	_, err := repository.Admin.Delete(c, req.Username)
	if err != nil {
		return nil, err
	}

	return &v1.AdminDeleteResp{}, nil
}

func (s *AdminService) Update(c echo.Context, req *v1.AdminUpdateReq) (*v1.AdminUpdateResp, error) {
	_, err := repository.Admin.Update(c, req.Username, req.Nickname, req.Remark)
	if err != nil {
		return nil, err
	}

	return &v1.AdminUpdateResp{}, nil
}

func (s *AdminService) Enable(c echo.Context, req *v1.AdminEnableReq) (*v1.AdminEnableResp, error) {
	user, err := repository.Admin.Enable(c, req.Username, req.Enable)
	if err != nil {
		return nil, err
	}

	return &v1.AdminEnableResp{
		Enable: int(user.Enable),
	}, nil
}

func (s *AdminService) ResetVerifiCode(c echo.Context, req *v1.AdminResetVerifiCodeReq) (*v1.AdminResetPasswordResp, error) {
	user, err := repository.Admin.ResetVerifiCode(c, req.Id)
	if err != nil {
		return nil, err
	}

	return &v1.AdminResetPasswordResp{
		Password: user.Password,
	}, nil
}

func (s *AdminService) FindParentIds(c echo.Context) ([]uint, error) {
	header := headerx.GetDataFromHeader(c)
	r := header.Role
	adminId := header.AdminId

	role := model.SysUserRole(r)

	ids, err := repository.Admin.FindAdminIds(c, adminId, role)
	if err != nil {
		return nil, err
	}

	return ids, nil
}

func (s *AdminService) FindAdminIds(c echo.Context, adminId uint, role model.SysUserRole) ([]uint, error) {
	ids, err := repository.Admin.FindAdminIds(c, adminId, role)
	if err != nil {
		return nil, err
	}

	return ids, nil
}

func (s *AdminService) GetMasterIncome(c echo.Context, req *v1.GetMasterIncomeReq) (*v1.GetMasterIncomeResp, error) {
	masterId := req.MasterId

	// 获取主账号总收入
	totalIncome, err := repository.Admin.GetMasterIncome(c, masterId)
	if err != nil {
		return nil, err
	}

	return &v1.GetMasterIncomeResp{
		TotalIncome: totalIncome,
	}, nil
}
