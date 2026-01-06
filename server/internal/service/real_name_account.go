package service

import (
	v1 "apollo/server/api/v1"
	"apollo/server/internal/model"
	"apollo/server/internal/repository"
	"apollo/server/pkg/headerx"

	"github.com/labstack/echo/v4"
)

var (
	RealNameAccount = new(RealNameAccountService)
)

type RealNameAccountService struct {
}

func (s *RealNameAccountService) Create(c echo.Context, req *v1.RealNameAccountCreateReq) (*v1.RealNameAccountCreateResp, error) {
	header := headerx.GetDataFromHeader(c)
	adminId := header.AdminId
	creator, err := repository.Admin.GetById(c, adminId)
	if err != nil {
		return nil, err
	}

	list := make([]*model.RealNameAccount, 0, len(req.AccountList))
	for _, a := range req.AccountList {
		if a == nil || len(a.IdNumber) == 0 || len(a.Name) == 0 {
			continue
		}

		list = append(list, &model.RealNameAccount{
			IdNumber: a.IdNumber,
			Name:     a.Name,
			Mobile:   a.Mobile,
			Address:  a.Address,
			Remark:   req.Remark,
			ParentId: adminId,
			MasterId: creator.MasterId,
		})
	}
	err = repository.RealNameAccount.Create(c, list)
	if err != nil {
		return nil, err
	}

	return &v1.RealNameAccountCreateResp{}, nil
}

func (s *RealNameAccountService) List(c echo.Context, req *v1.ListRealNameAccountReq) (*v1.ListRealNameAccountResp, error) {
	parentIds, _ := Admin.FindParentIds(c)

	accounts, total, err := repository.RealNameAccount.List(c, req, parentIds)
	if err != nil {
		return nil, err
	}

	list := make([]*v1.RealNameAccount, 0, len(accounts))
	for _, a := range accounts {
		list = append(list, &v1.RealNameAccount{
			BaseRealNameAccount: v1.BaseRealNameAccount{
				IdNumber: a.IdNumber,
				Name:     a.Name,
				Mobile:   a.Mobile,
				Address:  a.Address,
			},
			RealNameCount: a.RealNameCount,
			Enable:        int(a.Enable),
			Remark:        a.Remark,
		})
	}

	return &v1.ListRealNameAccountResp{
		ListTableData: v1.ListTableData[v1.RealNameAccount]{
			List:  list,
			Total: total,
		},
	}, nil
}
