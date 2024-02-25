package sqlstorage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/zamatay/learn-otus/hw12_13_14_15_calendar/configs"
	"github.com/zamatay/learn-otus/hw12_13_14_15_calendar/internal/domain"
	"github.com/zamatay/learn-otus/hw12_13_14_15_calendar/internal/logger"
)

const (
	sqlAddEvent = `
	INSERT INTO calendar(title, date, date_interval, description, user_id)
	VALUES ($1, $2, $3, $4, $5)`
	sqlUpdateEvent = `
	UPDATE calendar SET title = $1, date = $2, date_interval = $3, 
	                    description = $4, user_id= $5
	WHERE id = $6`
	sqlRemoveEvent = `DELETE FROM calendar where id = $1`
	sqlListEvent   = `SELECT id, title, date, date_interval, description, user_id FROM calendar where date between $1 and $2`
	sqlGetEvent    = `SELECT id, title, date, date_interval, description, user_id FROM calendar where id=$1`
)

type PrepareItem struct {
	stmt *sqlx.Stmt
	sql  string
}

type prepare struct {
	AddEvent    PrepareItem
	EditEvent   PrepareItem
	RemoveEvent PrepareItem
	ListEvent   PrepareItem
	GetEvent    PrepareItem
}

type Storage struct {
	connect    *sqlx.DB
	ctx        context.Context
	prepareSql prepare
}

func New(ctx context.Context, cfg *configs.DBConfig) (*Storage, error) {
	connectionString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable application_name=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DB, cfg.AppName)
	conn, err := sqlx.Open(cfg.Driver, connectionString)
	if err != nil {
		return nil, err
	}
	if err := conn.PingContext(ctx); err != nil {
		logger.Logger().Error("Ошибка при инициализации БД", "error", err.Error())
		return nil, err
	}
	pSql := prepare{
		AddEvent:    PrepareItem{stmt: nil, sql: sqlAddEvent},
		EditEvent:   PrepareItem{stmt: nil, sql: sqlUpdateEvent},
		RemoveEvent: PrepareItem{stmt: nil, sql: sqlRemoveEvent},
		ListEvent:   PrepareItem{stmt: nil, sql: sqlListEvent},
		GetEvent:    PrepareItem{stmt: nil, sql: sqlGetEvent},
	}
	return &Storage{
		connect:    conn,
		ctx:        ctx,
		prepareSql: pSql,
	}, nil
}

func (s Storage) prepare(item *PrepareItem) error {
	if item.stmt == nil {
		var err error
		item.stmt, err = s.connect.Preparex(item.sql)
		if err != nil {
			logger.Logger().Error("Ошибка при подготовки запроса", err)
			return err
		}
	}
	return nil
}

func (s Storage) AddEvent(event domain.Event) error {
	if err := s.prepare(&s.prepareSql.AddEvent); err != nil {
		return err
	}
	s.prepareSql.AddEvent.stmt.ExecContext(s.ctx,
		sql.Named("title", event.Title),
		sql.Named("date", event.Date),
		sql.Named("date_interval", event.DateInterval),
		sql.Named("description", event.Description),
		sql.Named("user_id", event.UserID),
	)
	return nil
}

func (s Storage) EditEvent(_ int64, event domain.Event) error {
	if err := s.prepare(&s.prepareSql.EditEvent); err != nil {
		return err
	}
	if _, err := s.prepareSql.EditEvent.stmt.ExecContext(s.ctx, event.Title, event.Date, event.DateInterval, event.Description,
		event.UserID, event.ID); err != nil {
		return err
	}

	return nil
}

func (s Storage) RemoveEvent(id int64) error {
	if err := s.prepare(&s.prepareSql.RemoveEvent); err != nil {
		return err
	}
	if _, err := s.prepareSql.RemoveEvent.stmt.ExecContext(s.ctx, sql.Named("id", id)); err != nil {
		return err
	}
	return nil
}

func (s Storage) List(beginDate time.Time, endDate time.Time) []domain.Event {
	if err := s.prepare(&s.prepareSql.ListEvent); err != nil {
		return nil
	}
	list := make([]domain.Event, 0, 0)
	rows, err := s.prepareSql.ListEvent.stmt.QueryContext(s.ctx,
		sql.Named("date1", beginDate),
		sql.Named("date2", endDate),
	)
	if err != nil {
		return nil
	}
	value := domain.Event{}
	for rows.Next() {
		err := rows.Scan(&value.ID, &value.Title, &value.Date, &value.DateInterval, &value.Description, &value.UserID)
		if err != nil {
			continue
		}
		list = append(list, value)
	}
	return list
}

func (s Storage) GetEvent(id int64) (domain.Event, error) {
	if err := s.prepare(&s.prepareSql.GetEvent); err != nil {
		return domain.Event{}, err
	}
	value := domain.Event{}
	if err := s.prepareSql.GetEvent.stmt.GetContext(s.ctx, &value, id); err != nil {
		return value, err
	}
	return value, nil
}

func (s Storage) Close() error {
	s.connect.Close()
	return nil
}
