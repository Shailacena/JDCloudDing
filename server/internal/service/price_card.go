package service

import (
	v1 "apollo/server/api/v1"
	"apollo/server/internal/model"
	"apollo/server/internal/repository"
	"time"

	"github.com/labstack/echo/v4"
)

type PriceCardService struct{}

var PriceCardServiceInst = &PriceCardService{}

func (s *PriceCardService) CreateCards(c echo.Context, req *v1.CardCreateReq) (*v1.CardCreateResp, error) {
	cards := make([]model.PriceCard, 0, len(req.Cards))
	batchNo := time.Now().Format("20060102")

	for _, card := range req.Cards {
		cards = append(cards, model.PriceCard{
			CardNo:      card.CardNo,
			Password:    card.Password,
			CardGroup:   card.CardGroup,
			Amount:      card.Amount,
			CardType:    model.CardTypeReal,
			BatchNo:     batchNo,
			CardStatus: model.CardStatusPending,
		})
	}

	err := repository.PriceCard.Create(c, cards)
	if err != nil {
		return nil, err
	}

	return &v1.CardCreateResp{
		Count: len(cards),
	}, nil
}

func (s *PriceCardService) GenerateVirtualCards(c echo.Context, req *v1.VirtualCardGenerateReq) (*v1.VirtualCardGenerateResp, error) {
	batchNo := time.Now().Format("20060102")

	cards, err := repository.PriceCard.GenerateVirtualCards(c, req.Prefix, req.CardNoLen, req.PasswordLen, req.CardGroup, req.Amount, req.Count, batchNo)
	if err != nil {
		return nil, err
	}

	return &v1.VirtualCardGenerateResp{
		Count: len(cards),
	}, nil
}

func (s *PriceCardService) List(c echo.Context, req *v1.ListCardReq) (*v1.ListCardResp, error) {
	page := req.CurrentPage
	pageSize := req.PageSize

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	cards, total, err := repository.PriceCard.List(c, req.CardNo, req.CardGroup, req.BatchNo, req.StartTime, req.EndTime, req.CardType, page, pageSize)
	if err != nil {
		return nil, err
	}

	result := make([]*v1.PriceCard, 0, len(cards))
	for _, card := range cards {
		var usedAt int64
		if card.UsedAt != nil {
			usedAt = card.UsedAt.Unix()
		}
		result = append(result, &v1.PriceCard{
			Id:         card.ID,
			CardNo:     card.CardNo,
			Password:   card.Password,
			CardGroup:  card.CardGroup,
			Amount:     card.Amount,
			CardType:   string(card.CardType),
			BatchNo:    card.BatchNo,
			CardStatus: string(card.CardStatus),
			OrderId:    card.OrderId,
			Remark:     card.Remark,
			UsedAt:     usedAt,
			CreateAt:   card.CreatedAt.Unix(),
		})
	}

	return &v1.ListCardResp{
		ListTableData: v1.ListTableData[v1.PriceCard]{
			List:  result,
			Total: total,
		},
	}, nil
}

func (s *PriceCardService) Delete(c echo.Context, req *v1.DeleteCardReq) (*v1.DeleteCardResp, error) {
	count, err := repository.PriceCard.Delete(c, req.Ids)
	if err != nil {
		return nil, err
	}

	return &v1.DeleteCardResp{
		Count: int(count),
	}, nil
}

func (s *PriceCardService) DeleteByCondition(c echo.Context, req *v1.ListCardReq) (*v1.DeleteCardResp, error) {
	count, err := repository.PriceCard.DeleteByCondition(req.CardNo, req.CardGroup, req.BatchNo, req.StartTime, req.EndTime, req.CardType)
	if err != nil {
		return nil, err
	}

	return &v1.DeleteCardResp{
		Count: int(count),
	}, nil
}
