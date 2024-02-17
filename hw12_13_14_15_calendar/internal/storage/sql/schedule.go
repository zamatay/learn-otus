package sqlstorage

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/zamatay/learn-otus/hw12_13_14_15_calendar/internal/domain"
	"github.com/zamatay/learn-otus/hw12_13_14_15_calendar/internal/logger"
	"time"
)

const sqlGetEventByDate = `SELECT id, title, date, date_interval, description, user_id FROM calendar where date > $1`

type Scheduler interface {
	GetNew(ctx context.Context, dt time.Time) ([]domain.Event, error)
}

type Schedule struct {
	db      *sqlx.DB
	prepare *sqlx.Stmt
}

func (s Schedule) GetNew(ctx context.Context, dt time.Time) ([]domain.Event, error) {
	var event []domain.Event //:= make([]domain.Event, 0, 0)
	if err := s.prepare.SelectContext(ctx, &event, dt); err != nil {
		logger.Logger().Error("Ошибка при получении данных в шедулере", "Error", err.Error())
		return nil, err
	}
	return event, nil
}

func NewSchedule(db *sqlx.DB) (*Schedule, error) {
	prepare, err := db.Preparex(sqlGetEventByDate)
	if err != nil {
		logger.Logger().Error("Ошибка при подготвоке запроса", "Error", err.Error())
		return nil, err
	}
	return &Schedule{db: db, prepare: prepare}, nil
}
