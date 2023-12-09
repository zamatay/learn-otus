package sqlstorage

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/zamatay/learn-otus/hw12_13_14_15_calendar/internal/domain"
	"github.com/zamatay/learn-otus/hw12_13_14_15_calendar/internal/logger"
	"time"
)

const (
	sqlAddEvent = `
	INSERT INTO calendar(title, date, date_interval, description, user_id)
	VALUES (@title, @date, @date_interval, @description, @user_id)`
	sqlUpdateEvent = `
	UPDATE calendar SET title = @title, date = @date, date_interval = @date_interval, 
	                    description = @description, user_id= @user_id
	WHERE id = @id`
	sqlRemoveEvent = `DELETE FROM calendar where id = @id`
	sqlListEvent   = `SELECT id, title, date, date_interval, description, user_id FROM calendar where date between @date1 and @date2`
	sqlGetEvent    = `SELECT id, title, date, date_interval, description, user_id FROM calendar where id=@id`
)

type PrepareItem struct {
	stmt *sqlx.Stmt
	sql  string
}

type Prepare struct {
	AddEvent    PrepareItem
	EditEvent   PrepareItem
	RemoveEvent PrepareItem
	ListEvent   PrepareItem
	GetEvent    PrepareItem
}

type Storage struct {
	connect    *sqlx.DB
	ctx        context.Context
	prepareSql Prepare
}

func New(ctx context.Context, cfg *main.DBConfig) *Storage {
	connectionString := fmt.Sprintf(
		cfg.Host, cfg.Port, cfg.User, cfg.Password)
	conn, err := sqlx.Open(cfg.Driver, connectionString)
	if err != nil {
		return nil
	}
	pSql := Prepare{
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
	}
}

func (s *Storage) prepare(item *PrepareItem) error {
	if item.stmt == nil {
		var err error
		item.stmt, err = s.connect.Preparex(item.sql)
		if err != nil {
			logger.GetLog().Error("Ошибка при подготовки запроса", err)
			return err
		}
	}
	return nil
}

func (s *Storage) AddEvent(event domain.Event) error {
	s.prepare(&s.prepareSql.AddEvent)
	s.prepareSql.AddEvent.stmt.ExecContext(s.ctx,
		sql.Named("title", event.Title),
		sql.Named("date", event.Date),
		sql.Named("date_interval", event.DateInterval),
		sql.Named("description", event.Description),
		sql.Named("user_id", event.UserId),
	)
	return nil
}

func (s *Storage) EditEvent(id int64, event domain.Event) error {
	s.prepare(&s.prepareSql.EditEvent)
	s.prepareSql.EditEvent.stmt.ExecContext(s.ctx,
		sql.Named("title", event.Title),
		sql.Named("date", event.Date),
		sql.Named("date_interval", event.DateInterval),
		sql.Named("description", event.Description),
		sql.Named("user_id", event.UserId),
		sql.Named("id", event.Id),
	)
	return nil
}

func (s *Storage) RemoveEvent(id int64) error {
	s.prepare(&s.prepareSql.RemoveEvent)
	s.prepareSql.EditEvent.stmt.ExecContext(s.ctx,
		sql.Named("id", id),
	)
	return nil
}

func (s *Storage) List(beginDate time.Time, endDate time.Time) []domain.Event {
	s.prepare(&s.prepareSql.ListEvent)
	list := make([]domain.Event, 0, 0)
	s.prepareSql.EditEvent.stmt.GetContext(s.ctx, &list,
		sql.Named("date1", beginDate),
		sql.Named("date2", endDate),
	)
	return list
}

func (s *Storage) GetEvent(id int64) (domain.Event, error) {
	s.prepare(&s.prepareSql.ListEvent)
	value := domain.Event{}
	s.prepareSql.EditEvent.stmt.GetContext(s.ctx, &value,
		sql.Named("id", id),
	)
	return value, nil
}

func (s *Storage) Close(ctx context.Context) error {
	s.connect.Close()
	return nil
}
