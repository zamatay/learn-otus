package background

import (
	"context"
	"errors"
	"sync/atomic"
	"time"
)

type (
	Runner interface {
		Run(context.Context, Handler) error
	}
	Handler func(ctx context.Context)
)

type scheduling struct {
	timeInterval time.Duration
	isStarted    atomic.Bool
}

func NewScheduling(timeInterval time.Duration) Runner {
	scheduling := &scheduling{
		timeInterval: timeInterval,
	}
	scheduling.isStarted.Store(false)

	return scheduling
}

func (s *scheduling) Run(ctx context.Context, handler Handler) error {
	prevStarted := s.isStarted.Swap(true)
	if prevStarted {
		return errors.New("Сервис уже работает")
	}

	ticker := time.NewTicker(s.timeInterval)
	defer ticker.Stop()

	for {
		handler(ctx)

		select {
		case <-ctx.Done():
			// logger.Log.Debug("Работа завершена.")
			return nil
		case <-ticker.C:
			continue
		}
	}
}
