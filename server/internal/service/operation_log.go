package service

import (
	v1 "apollo/server/api/v1"
	"apollo/server/internal/model"
	"apollo/server/internal/repository"
	"github.com/labstack/echo/v4"
)

var (
	OperationLog = new(OperationLogService)
)

type OperationLogService struct {
}

func (s *OperationLogService) List(c echo.Context, req *v1.ListOperationLogReq) (*v1.ListOperationLogResp, error) {
	var parentIds []uint
	parentIds, _ = Admin.FindParentIds(c)

	logs, total, err := repository.OperationLog.List(c, req.Pagination, parentIds)
	if err != nil {
		return nil, err
	}

	list := make([]*v1.OperationLog, 0, len(logs))
	for _, l := range logs {
		u, err := repository.Admin.GetById(c, l.Operator)
		if err != nil {
			continue
		}

		list = append(list, &v1.OperationLog{
			Id:           l.ID,
			IP:           l.IP,
			Operation:    l.Operation,
			OperationStr: model.GetOperationStr(l.Operation),
			Operator:     l.Operator,
			OperatorName: u.Nickname,
			CreateAt:     l.CreatedAt.Unix(),
		})
	}

	return &v1.ListOperationLogResp{
		ListTableData: v1.ListTableData[v1.OperationLog]{
			List:  list,
			Total: total,
		},
	}, nil
}
