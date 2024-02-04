package main

import (
	"github.com/brianvoe/gofakeit/v6"
	"github.com/zamatay/learn-otus/hw12_13_14_15_calendar/configs"
	"github.com/zamatay/learn-otus/hw12_13_14_15_calendar/internal/app"
	"github.com/zamatay/learn-otus/hw12_13_14_15_calendar/internal/logger"
	"github.com/zamatay/learn-otus/hw12_13_14_15_calendar/pkg/pgsql"
	"os"
	"time"
)

type Event struct {
	Title       string
	Date        time.Time
	Description string
	Interval    time.Duration
	UserID      uint8
}

func main() {
	ctx, cancel := app.InitShutDowner()
	defer cancel()

	//создали брокер сообщений
	cfg := configs.Configs()
	//подключились к базе
	connect, dbCloser, err := pgsql.Connect(ctx, &cfg.DB)
	if err != nil {
		logger.Logger().Error("Ошибка при подключении к БД", "Error", err.Error())
		os.Exit(1)
	}
	defer dbCloser()

	go func() {
		var f Event
		for {
			f = Event{gofakeit.Name(), time.Now(), gofakeit.ProductDescription(), time.Duration(gofakeit.Minute()), gofakeit.Uint8()}
			result, err := connect.ExecContext(ctx,
				`insert into calendar(title, date, description, date_interval, user_id)
					values ($1, $2, $3, $4, $5)`, f.Title, f.Date, f.Description, f.Interval, f.UserID)
			if err != nil {
				logger.Logger().Error("Ошибка вставки данных", "Error", err.Error())
			}
			logger.Logger().Info("Данные вставлены", "Info", result)
			time.Sleep(10 * time.Second)
		}
	}()

	quit := app.Shutdown(ctx, nil)

	<-quit

}
