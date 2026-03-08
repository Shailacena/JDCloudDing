package repository

import (
	"apollo/server/internal/model"
	"apollo/server/pkg/data"
	"apollo/server/pkg/rand"
	"apollo/server/pkg/util"
	"errors"
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type PriceCardRepo struct{}

var PriceCard = &PriceCardRepo{}

func (r *PriceCardRepo) Create(c echo.Context, cards []model.PriceCard) error {
	db := data.Instance()
	return db.Create(&cards).Error
}

func (r *PriceCardRepo) GenerateVirtualCards(c echo.Context, prefix string, cardNoLen, passwordLen int, cardGroup string, amount float64, count int, batchNo string) ([]model.PriceCard, error) {
	db := data.Instance()

	cards := make([]model.PriceCard, count)
	for i := 0; i < count; i++ {
		cardNo := generateCardNo(prefix, cardNoLen)
		password := generatePasswordUpper(passwordLen)
		cards[i] = model.PriceCard{
			CardNo:     cardNo,
			Password:   password,
			CardGroup:  cardGroup,
			Amount:     amount,
			CardType:  model.CardTypeVirtual,
			BatchNo:   batchNo,
			UsedStatus: false,
		}
	}

	err := db.Create(&cards).Error
	return cards, err
}

func (r *PriceCardRepo) List(c echo.Context, cardNo, cardGroup, batchNo, startTime, endTime, cardType string, page, pageSize int) ([]model.PriceCard, int64, error) {
	db := data.Instance()

	query := db.Model(&model.PriceCard{})

	if cardNo != "" {
		query = query.Where("card_no LIKE ?", "%"+cardNo+"%")
	}
	if cardGroup != "" {
		query = query.Where("card_group = ?", cardGroup)
	}
	if batchNo != "" {
		query = query.Where("batch_no = ?", batchNo)
	}
	if startTime != "" {
		t, err := time.Parse("2006-01-02", startTime)
		if err == nil {
			query = query.Where("created_at >= ?", t)
		}
	}
	if endTime != "" {
		t, err := time.Parse("2006-01-02", endTime)
		if err == nil {
			query = query.Where("created_at <= ?", t.Add(24*time.Hour))
		}
	}
	if cardType != "" {
		query = query.Where("card_type = ?", cardType)
	}

	var total int64
	query.Count(&total)

	var cards []model.PriceCard
	err := query.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&cards).Error
	return cards, total, err
}

func (r *PriceCardRepo) Delete(c echo.Context, ids []uint) (int64, error) {
	db := data.Instance()
	result := db.Where("id IN ?", ids).Delete(&model.PriceCard{})
	return result.RowsAffected, result.Error
}

func (r *PriceCardRepo) DeleteByCondition(cardNo, cardGroup, batchNo, startTime, endTime, cardType string) (int64, error) {
	db := data.Instance()
	query := db.Model(&model.PriceCard{})

	if cardNo != "" {
		query = query.Where("card_no LIKE ?", "%"+cardNo+"%")
	}
	if cardGroup != "" {
		query = query.Where("card_group = ?", cardGroup)
	}
	if batchNo != "" {
		query = query.Where("batch_no = ?", batchNo)
	}
	if startTime != "" {
		t, err := time.Parse("2006-01-02", startTime)
		if err == nil {
			query = query.Where("created_at >= ?", t)
		}
	}
	if endTime != "" {
		t, err := time.Parse("2006-01-02", endTime)
		if err == nil {
			query = query.Where("created_at <= ?", t.Add(24*time.Hour))
		}
	}
	if cardType != "" {
		query = query.Where("card_type = ?", cardType)
	}

	result := query.Delete(&model.PriceCard{})
	return result.RowsAffected, result.Error
}

func (r *PriceCardRepo) GetByCardNo(c echo.Context, cardNo string) (*model.PriceCard, error) {
	db := data.Instance()
	var card model.PriceCard
	err := db.Where("card_no = ?", cardNo).First(&card).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &card, nil
}

func (r *PriceCardRepo) UseCard(c echo.Context, cardNo, orderId string) error {
	db := data.Instance()
	now := time.Now()
	return db.Model(&model.PriceCard{}).Where("card_no = ?", cardNo).Updates(map[string]interface{}{
		"used_status": true,
		"order_id":    orderId,
		"used_at":     now,
	}).Error
}

func generateCardNo(prefix string, length int) string {
	digitCount := length - len(prefix)
	if digitCount < 0 {
		digitCount = 0
	}
	result := make([]byte, length)
	copy(result, []byte(prefix))
	for i := len(prefix); i < length; i++ {
		result[i] = byte('0' + rand.Random.Intn(10))
	}
	return string(result)
}

func generatePassword(length int) string {
	return util.RandStringRunes(length)
}

func generatePasswordUpper(length int) string {
	result := []byte(util.RandStringRunes(length))
	for i := 0; i < len(result); i++ {
		c := result[i]
		if c >= 'a' && c <= 'z' {
			result[i] = c - 'a' + 'A'
		}
	}
	return string(result)
}
