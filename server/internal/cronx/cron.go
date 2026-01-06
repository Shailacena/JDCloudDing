package cronx

import (
	"apollo/server/cmd/payment/api/common"
	"apollo/server/internal/model"
	"apollo/server/internal/repository"
	"apollo/server/pkg/contextx"
	"apollo/server/pkg/data"
	"github.com/labstack/echo/v4"
	"github.com/robfig/cron/v3"
	"github.com/samber/lo"
	"time"
)

func Init(logger echo.Logger) {
	logger.Info("start corn")

	ctx := contextx.NewContext(logger)
	handleFailedNotify(ctx)
}

func handleFailedNotify(ctx contextx.Context) {
	c := cron.New()
	logger := ctx.Logger()

	_, err := c.AddFunc("@every 1m", func() {
		list, err := repository.Notify.List(ctx, time.Now())
		if err != nil {
			logger.Errorf("repository.Notify.List error=%s", err)
			return
		}

		ids := lo.Map(list, func(item *model.Notify, index int) uint {
			return item.ID
		})

		if len(ids) == 0 {
			return
		}

		logger.Info("handleFailedNotify ids=", ids)

		db := data.Instance()
		successful := make(map[string]bool)
		for _, n := range list {
			if successful[n.BizId] {
				continue
			}

			switch n.BizType {
			case model.NotifyBizTypeOrder:
				order, err := repository.Order.GetByOrderId(db, n.BizId)
				if err != nil {
					logger.Errorf("Order.GetByOrderId order=%s, error=%s", n.BizId, err)
					continue
				}

				err = common.NotifyMerchant(ctx, db, order)
				if err != nil {
					logger.Errorf("common.NotifyMerchant order=%s, error=%s", n.BizId, err)
					continue
				}

				err = repository.Notify.Update(ctx, n.ID, time.Now(), model.NotifyDone)
				if err != nil {
					logger.Errorf("Notify.Update order=%s, error=%s", n.BizId, err)
					continue
				}
				successful[n.BizId] = true
			default:
				continue
			}
		}
	})
	if err != nil {
		logger.Errorf("handleFailedNotify error=%s", err)
		return
	}

	c.Start()
}
