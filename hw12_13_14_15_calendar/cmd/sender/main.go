package main

import (
	"github.com/zamatay/learn-otus/hw12_13_14_15_calendar/configs"
	"github.com/zamatay/learn-otus/hw12_13_14_15_calendar/internal/app"
	"github.com/zamatay/learn-otus/hw12_13_14_15_calendar/internal/logger"
	service "github.com/zamatay/learn-otus/hw12_13_14_15_calendar/internal/sender"
	"github.com/zamatay/learn-otus/hw12_13_14_15_calendar/pkg/background"
	"github.com/zamatay/learn-otus/hw12_13_14_15_calendar/pkg/rabbit"
	"os"
)

func main() {

	ctx, cancel := app.InitShutDowner()
	defer cancel()

	//создали брокер сообщений
	cfg := configs.Configs()
	newRabbit, err := rabbit.NewRabbit(&cfg.Broker)
	if err != nil {
		logger.Logger().Error("Ошибка при подключении к БД", "Error", err.Error())
		os.Exit(1)
	}
	defer newRabbit.Close()
	// создали сервис
	srv := service.NewService(newRabbit)

	err = background.NewScheduling(configs.Configs().Schedule.Interval).Run(ctx, srv.Run)
	if err != nil {
		logger.Logger().Error("Ошибка при запуске джоба", "Error", err.Error())
	}

	quit := app.Shutdown(ctx, app.Calendar.Closers())

	<-quit

}
