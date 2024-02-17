package pgsql

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/zamatay/learn-otus/hw12_13_14_15_calendar/configs"
	"github.com/zamatay/learn-otus/hw12_13_14_15_calendar/internal/logger"
)

func Connect(ctx context.Context, cfg *configs.DBConfig) (*sqlx.DB, func(), error) {
	connectionString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable application_name=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DB, cfg.AppName)
	conn, err := sqlx.Open(cfg.Driver, connectionString)
	if err != nil {
		return nil, nil, err
	}
	if err := conn.PingContext(ctx); err != nil {
		logger.Logger().Error("Ошибка при инициализации БД", "error", err.Error())
		return conn, nil, err
	}

	closer := func() {
		err := conn.Close()
		if err != nil {
			logger.Logger().Error("Ошибка при закрытии подключения к БД", "error", err.Error())
		}
	}

	return conn, closer, nil
}
