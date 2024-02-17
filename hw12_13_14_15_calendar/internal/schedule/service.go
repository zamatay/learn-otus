package schedule

import (
	"context"
	"fmt"
	"github.com/zamatay/learn-otus/hw12_13_14_15_calendar/internal/logger"
	sqlstorage "github.com/zamatay/learn-otus/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/zamatay/learn-otus/hw12_13_14_15_calendar/pkg/rabbit"
	"time"
)

type Service struct {
	dt     time.Time
	repo   sqlstorage.Scheduler
	broker rabbit.ProduceRabbited
}

func NewService(repo *sqlstorage.Schedule, broker *rabbit.Rabbit) *Service {
	return &Service{repo: repo, broker: broker}
}

func (srv *Service) Run(ctx context.Context) {
	events, err := srv.repo.GetNew(ctx, srv.dt)
	srv.dt = time.Now()
	if err != nil {
		logger.Logger().Error("Щшибка получения данных с БД", "error", err.Error())
		return
	}
	if len(events) > 0 {
		for _, event := range events {
			srv.broker.SendMessage(ctx, fmt.Sprintf("%s\n%s", event.Title, event.Description))
		}
	}
}
