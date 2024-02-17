package sender

import (
	"context"
	"github.com/zamatay/learn-otus/hw12_13_14_15_calendar/internal/logger"
	"github.com/zamatay/learn-otus/hw12_13_14_15_calendar/pkg/rabbit"
	"time"
)

type Service struct {
	interval time.Duration
	broker   rabbit.ConsumerRabbited
}

func NewService(broker *rabbit.Rabbit) *Service {
	return &Service{broker: broker, interval: time.Second * 10}
}

func (srv *Service) Run(ctx context.Context) {
	for {
		if err := srv.broker.GetMessage(ctx); err != nil {
			logger.Logger().Error("Ошибка при чтении сообщения из очереди", "Error", err.Error())
		}
		time.Sleep(srv.interval)
	}
}
