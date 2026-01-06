package repository

import (
	"apollo/server/internal/model"
	"apollo/server/pkg/contextx"
	"apollo/server/pkg/data"
	"apollo/server/pkg/timex"
	"time"
)

var (
	Notify = &NotifyRepo{}
)

type NotifyRepo struct {
}

func (r *NotifyRepo) Create(c contextx.Context, n *model.Notify) (*model.Notify, error) {
	db := data.Instance()

	err := db.Create(n).Error

	return n, err
}

func (r *NotifyRepo) List(c contextx.Context, expiredAt time.Time) ([]*model.Notify, error) {
	db := data.Instance()

	endAt := expiredAt
	startAt := expiredAt.Add(-time.Minute)
	var list []*model.Notify
	err := db.Where("notify_status = ? AND expired_at BETWEEN ? AND ?", model.NotNotify, startAt, endAt).Find(&list).Error

	return list, err
}

func (r *NotifyRepo) Update(c contextx.Context, id uint, notifyAt time.Time, notifyStatus model.NotifyStatus) error {
	db := data.Instance()

	err := db.Model(&model.Notify{}).Where("id = ?", id).Updates(map[string]any{
		"notify_at":     notifyAt,
		"notify_status": notifyStatus,
	}).Error

	return err
}

func (r *NotifyRepo) Delete(c contextx.Context, bizId string) error {
	db := data.Instance()

	err := db.Delete(&model.Notify{}, "biz_id = ?", bizId).Error

	return err
}

func (r *NotifyRepo) CountOfToday(c contextx.Context, bizId string, bizType model.NotifyBizType) (int64, error) {
	db := data.Instance()

	endAt := timex.GetPRCNowTime()
	startAt := endAt.StartOfDay()
	var total int64
	err := db.Model(&model.Notify{}).Where("biz_id = ? AND biz_type = ? AND created_at BETWEEN ? AND ?", bizId, bizType, startAt, endAt).Count(&total).Error

	return total, err
}
