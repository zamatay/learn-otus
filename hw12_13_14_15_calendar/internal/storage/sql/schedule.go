package sqlstorage

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/zamatay/learn-otus/hw12_13_14_15_calendar/internal/domain"
	"github.com/zamatay/learn-otus/hw12_13_14_15_calendar/internal/logger"
	"time"
)

const sqlGetEventByDate = `SELECT id, title, date, date_interval, description, user_id FROM calendar where date > $1`
const sqlSetEventByDate = `insert into calendar(id, title, date, date_interval, description, user_id)
							values (1, 'Тестовый заголовок', now(), 100, 'Тестовое описание', 0)`

type Scheduler interface {
	GetNew(ctx context.Context) ([]domain.Event, error)
	InsertData(ctx context.Context)
}

type Schedule struct {
	dt      time.Time
	db      *sqlx.DB
	prepare *sqlx.Stmt
}

func (s Schedule) InsertData(ctx context.Context) {
	_, err := s.db.ExecContext(ctx, sqlSetEventByDate)
	if err != nil {
		logger.Logger().Error("Не удалось вставить данные для теста")
	}
}
func (s Schedule) GetNew(ctx context.Context) ([]domain.Event, error) {
	var event []domain.Event //:= make([]domain.Event, 0, 0)
	if err := s.prepare.SelectContext(ctx, &event, s.dt); err != nil {
		logger.Logger().Error("Ошибка при получении данных в шедулере", "Error", err.Error())
		return nil, err
	}
	s.dt = time.Now()
	return event, nil
}

func NewSchedule(db *sqlx.DB) (*Schedule, error) {
	prepare, err := db.Preparex(sqlGetEventByDate)
	if err != nil {
		logger.Logger().Error("Ошибка при подготвоке запроса", "Error", err.Error())
		return nil, err
	}
	if _, err := db.Exec(sqlSetEventByDate); err != nil {
		logger.Logger().Error("Не удалось вставить данные для теста")
	}
	return &Schedule{db: db, prepare: prepare}, nil
}
