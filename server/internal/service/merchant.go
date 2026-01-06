package service

import (
	v1 "apollo/server/api/v1"
	"apollo/server/internal/model"
	"apollo/server/internal/repository"
	"apollo/server/pkg/data"
	"apollo/server/pkg/headerx"
	"apollo/server/pkg/timex"
	"apollo/server/pkg/totpx"
	"apollo/server/pkg/util"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
	"net/http"
)

var (
	Merchant = new(MerchantService)
)

type MerchantService struct {
}

func (s *MerchantService) Register(c echo.Context, req *v1.MerchantRegisterReq) (*v1.MerchantRegisterResp, error) {
	header := headerx.GetDataFromHeader(c)
	adminId := header.AdminId

	creator, err := repository.Admin.GetById(c, adminId)
	if err != nil {
		return nil, err
	}

	m := model.Merchant{
		Base: model.Base{
			Nickname: req.Nickname,
			Remark:   req.Remark,
			ParentId: adminId,
			MasterId: creator.MasterId,
		},
		PrivateKey: util.NewPrivateKey(),
	}
	merchant, err := repository.Merchant.Register(c, &m)
	if err != nil {
		return nil, err
	}

	secret, url, err := totpx.Generate(merchant.Username)
	if err != nil {
		return nil, err
	}

	merchant.SecretKey = secret
	merchant.UrlKey = url

	_, err = repository.Merchant.Update(c, merchant.Username, false, merchant)
	if err != nil {
		return nil, err
	}

	return &v1.MerchantRegisterResp{
		Nickname: merchant.Nickname,
		Password: merchant.Password,
	}, nil
}

func (s *MerchantService) Login(c echo.Context, req *v1.MerchantLoginReq) (*v1.MerchantLoginResp, error) {
	merchant, err := repository.Merchant.Login(c, req.Username, req.Password, req.VerifiCode)
	if err != nil {
		return nil, err
	}

	return &v1.MerchantLoginResp{
		Id:       merchant.ID,
		Token:    merchant.Token,
		Nickname: merchant.Nickname,
	}, nil
}

func (s *MerchantService) Logout(c echo.Context, req *v1.MerchantLogoutReq, token string) (*v1.MerchantLogoutResp, error) {
	err := repository.Merchant.Logout(c, token)
	if err != nil {
		return nil, err
	}

	return &v1.MerchantLogoutResp{}, nil
}

func (s *MerchantService) Update(c echo.Context, req *v1.MerchantUpdateReq) (*v1.MerchantUpdateResp, error) {
	m := model.Merchant{
		Base: model.Base{
			Nickname: req.Nickname,
			Remark:   req.Remark,
		},
	}
	_, err := repository.Merchant.Update(c, req.Username, req.IsDel, &m)
	if err != nil {
		return nil, err
	}

	return &v1.MerchantUpdateResp{}, nil
}

func (s *MerchantService) UpdateBalance(c echo.Context, req *v1.MerchantUpdateBalanceReq) (*v1.MerchantUpdateBalanceResp, error) {
	if req.ChangeAmount == 0 {
		return &v1.MerchantUpdateBalanceResp{}, nil
	}

	from := model.BalanceFromTypeSystemAdd
	if req.ChangeAmount < 0 {
		from = model.BalanceFromTypeSystemDeduct
	}

	err := repository.Admin.CheckPassword(c, req.AdminId, req.Password)
	if err != nil {
		return nil, err
	}

	db := data.Instance()
	err = repository.Merchant.UpdateBalance(db, req.MerchantId, "", req.ChangeAmount, from)
	if err != nil {
		return nil, err
	}

	return &v1.MerchantUpdateBalanceResp{}, nil
}

func (s *MerchantService) List(c echo.Context, req *v1.ListMerchantReq) (*v1.ListMerchantResp, error) {
	var parentIds []uint
	if req.MerchantId == 0 {
		parentIds, _ = Admin.FindParentIds(c)
	}

	merchants, total, err := repository.Merchant.List(c, req, parentIds)
	if err != nil {
		return nil, err
	}

	ids := lo.Map(merchants, func(item *model.Merchant, _ int) uint {
		return item.ID
	})

	dataMap := make(map[uint]repository.QueryMerchantAmountResult)

	if !req.IgnoreStatistics {
		results, err := repository.Order.QueryResultByMerchant(c, data.Instance(), ids, timex.GetPRCNowTime().Carbon2Time())
		if err != nil {
			return nil, err
		}

		dataMap = lo.SliceToMap(results, func(item repository.QueryMerchantAmountResult) (uint, repository.QueryMerchantAmountResult) {
			return item.MerchantId, item
		})
	}

	list := make([]*v1.Merchant, 0, len(merchants))
	for _, m := range merchants {
		var todayAmount, totalAmount float64
		item, ok := dataMap[m.ID]
		if ok {
			todayAmount = item.TodaySuccessAmount
			// totalAmount = item.TotalSuccessAmount
		}

		list = append(list, &v1.Merchant{
			Id:          m.ID,
			Username:    m.Username,
			Nickname:    m.Nickname,
			PrivateKey:  m.PrivateKey,
			Enable:      int(m.Enable),
			Remark:      m.Remark,
			Balance:     util.ToDecimal(m.Balance),
			UrlKey:      m.UrlKey,
			CreateAt:    m.CreatedAt.Unix(),
			TotalAmount: totalAmount,
			TodayAmount: todayAmount,
			ParentId:    m.ParentId,
		})
	}

	return &v1.ListMerchantResp{
		ListTableData: v1.ListTableData[v1.Merchant]{
			List:  list,
			Total: total,
		},
	}, nil
}

func (s *MerchantService) SetPassword(c echo.Context, req *v1.MerchantSetPasswordReq, token string) (*v1.MerchantSetPasswordResp, error) {
	if len(req.NewPassword) < 6 {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "密码应大于6位")
	}

	_, err := repository.Merchant.SetPassword(c, token, req.OldPassword, req.NewPassword)
	if err != nil {
		return nil, err
	}

	return &v1.MerchantSetPasswordResp{}, nil
}

func (s *MerchantService) ResetPassword(c echo.Context, req *v1.MerchantResetPasswordReq) (*v1.MerchantResetPasswordResp, error) {
	user, err := repository.Merchant.ResetPassword(c, req.Id)
	if err != nil {
		return nil, err
	}

	return &v1.MerchantResetPasswordResp{
		Password: user.Password,
	}, nil
}

func (s *MerchantService) Enable(c echo.Context, req *v1.MerchantEnableReq) (*v1.MerchantEnableResp, error) {
	user, err := repository.Merchant.Enable(c, req.Username, req.Enable)
	if err != nil {
		return nil, err
	}

	return &v1.MerchantEnableResp{
		Enable: int(user.Enable),
	}, nil
}

func (s *MerchantService) ListBalanceBill(c echo.Context, req *v1.ListMerchantBalanceBillReq) (*v1.ListMerchantBalanceBillResp, error) {
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
			return &v1.ListMerchantBalanceBillResp{}, nil
		}
		fmt.Println(merchantIds)
	}

	bills, total, err := repository.Merchant.ListBalanceBill(c, req, merchantIds)
	if err != nil {
		return nil, err
	}

	list := make([]*v1.MerchantBalanceBill, 0, len(bills))
	for _, b := range bills {
		list = append(list, &v1.MerchantBalanceBill{
			Id:           b.ID,
			MerchantId:   b.MerchantId,
			Nickname:     b.Nickname,
			OrderId:      b.OrderId,
			From:         int(b.From),
			Balance:      util.ToDecimal(b.Balance),
			ChangeAmount: util.ToDecimal(b.ChangeAmount),
			CreateAt:     b.CreatedAt.Unix(),
		})
	}

	return &v1.ListMerchantBalanceBillResp{
		ListTableData: v1.ListTableData[v1.MerchantBalanceBill]{
			List:  list,
			Total: total,
		},
	}, nil
}

func (s *MerchantService) GetBalance(c echo.Context, req *v1.MerchantBalanceReq, token string) (*v1.MerchantBalanceResp, error) {
	balance, err := repository.Merchant.GetBalance(c, req, token)
	if err != nil {
		return nil, err
	}

	return &v1.MerchantBalanceResp{
		Balance: balance,
	}, nil
}

func (s *MerchantService) ResetVerifiCode(c echo.Context, req *v1.MerchantResetVerifiCodeReq) (*v1.MerchantResetPasswordResp, error) {
	user, err := repository.Merchant.ResetVerifiCode(c, req.Id)
	if err != nil {
		return nil, err
	}

	return &v1.MerchantResetPasswordResp{
		Password: user.Password,
	}, nil
}
