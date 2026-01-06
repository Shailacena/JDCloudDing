package service

import (
	v1 "apollo/server/api/v1"
	"apollo/server/internal/model"
	"apollo/server/internal/repository"
	"apollo/server/pkg/headerx"

	"github.com/labstack/echo/v4"
)

var (
	JDAccount = new(JDAccountService)
)

type JDAccountService struct {
}

func (s *JDAccountService) Create(c echo.Context, req *v1.JDAccountCreateReq) (*v1.JDAccountCreateResp, error) {
	header := headerx.GetDataFromHeader(c)
	adminId := header.AdminId

	creator, err := repository.Admin.GetById(c, adminId)
	if err != nil {
		return nil, err
	}

	list := make([]*model.JDAccount, 0, len(req.AccountList))
	for _, a := range req.AccountList {
		if len(a.Account) == 0 || len(a.WsKey) == 0 {
			continue
		}

		list = append(list, &model.JDAccount{
			Account:  a.Account,
			WsKey:    a.WsKey,
			Remark:   req.Remark,
			ParentId: adminId,
			MasterId: creator.MasterId,
		})
	}
	err = repository.JDAccount.Create(c, list)
	if err != nil {
		return nil, err
	}

	return &v1.JDAccountCreateResp{}, nil
}

func (s *JDAccountService) Enable(c echo.Context, req *v1.JDAccountEnableReq) (*v1.JDAccountEnableResp, error) {
	err := repository.JDAccount.Enable(c, req.Id, model.JDAccountStatus(req.Status))
	if err != nil {
		return nil, err
	}

	return &v1.JDAccountEnableResp{}, nil
}

func (s *JDAccountService) List(c echo.Context, req *v1.ListJDAccountReq) (*v1.ListJDAccountResp, error) {
	parentIds, _ := Admin.FindParentIds(c)

	accounts, total, err := repository.JDAccount.List(c, req, parentIds)
	if err != nil {
		return nil, err
	}

	list := make([]*v1.JDAccount, 0, len(accounts))
	for _, a := range accounts {
		list = append(list, &v1.JDAccount{
			Id:       a.ID,
			Account:  a.Account,
			Status:   int(a.Status),
			Remark:   a.Remark,
			CreateAt: a.CreatedAt.Unix(),
			UpdateAt: a.UpdatedAt.Unix(),
		})
	}

	return &v1.ListJDAccountResp{
		ListTableData: v1.ListTableData[v1.JDAccount]{
			List:  list,
			Total: total,
		},
	}, nil
}

func (s *JDAccountService) Delete(c echo.Context, req *v1.JDAccountDeleteReq) (*v1.JDAccountDeleteResp, error) {
	err := repository.JDAccount.Delete(c, req)
	if err != nil {
		return nil, err
	}

	return &v1.JDAccountDeleteResp{}, nil
}

func (s *JDAccountService) Reset(c echo.Context, req *v1.JDAccountResetReq) (*v1.JDAccountResetResp, error) {
	err := repository.JDAccount.Reset(c, req)
	if err != nil {
		return nil, err
	}

	return &v1.JDAccountResetResp{}, nil
}
