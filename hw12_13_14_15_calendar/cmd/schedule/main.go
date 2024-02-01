package main

import (
	"github.com/zamatay/learn-otus/hw12_13_14_15_calendar/configs"
	"github.com/zamatay/learn-otus/hw12_13_14_15_calendar/internal/app"
	"github.com/zamatay/learn-otus/hw12_13_14_15_calendar/internal/logger"
	service "github.com/zamatay/learn-otus/hw12_13_14_15_calendar/internal/schedule"
	sqlstorage "github.com/zamatay/learn-otus/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/zamatay/learn-otus/hw12_13_14_15_calendar/pkg/background"
	"github.com/zamatay/learn-otus/hw12_13_14_15_calendar/pkg/pgsql"
	"github.com/zamatay/learn-otus/hw12_13_14_15_calendar/pkg/rabbit"
	"os"
)

func main() {
	ctx, cancel := app.InitShutDowner()
	defer cancel()

	//создали брокер сообщений
	cfg := configs.Configs()
	newRabbit, err := rabbit.NewRabbit(cfg.Broker.Login, cfg.Broker.Password, cfg.Broker.Url)
	if err != nil {
		logger.Logger().Error("Ошибка при подключении к БД", "Error", err.Error())
		os.Exit(1)
	}
	defer newRabbit.Close()

	//подключились к базе
	connect, dbCloser, err := pgsql.Connect(ctx, &configs.Configs().DB)
	if err != nil {
		logger.Logger().Error("Ошибка при подключении к БД", "Error", err.Error())
		os.Exit(1)
	}
	defer dbCloser()

	//создали репу
	schedule, err := sqlstorage.NewSchedule(connect)
	if err != nil {
		logger.Logger().Error("Ошибка при создании репы", "Error", err.Error())
		os.Exit(1)
	}

	// создали сервис
	srv := service.NewService(schedule, newRabbit)

	err = background.NewScheduling(configs.Configs().Schedule.Interval).Run(ctx, srv.Run)
	if err != nil {
		logger.Logger().Error("Ошибка при запуске джоба", "Error", err.Error())
	}

	logger.Logger().Info("Сервис Scheduler стартовал")
}
