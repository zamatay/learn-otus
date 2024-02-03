package main

import (
	"github.com/zamatay/learn-otus/hw12_13_14_15_calendar/internal/app"
	internalhttp "github.com/zamatay/learn-otus/hw12_13_14_15_calendar/internal/server/http"
	"os"
)

func main() {
	//Инициализируем контекст завершения
	ctx, cancel := app.InitShutDowner()

	//Инициализируем хранилище
	app.New()
	app.Calendar.Init(ctx)

	//Стартуем сервер
	server := internalhttp.NewServer(ctx, app.Calendar.Config.Grpc, app.Calendar.Config.HTTP)
	app.Calendar.AddClosers(server)

	//Shutdown
	quit := app.Shutdown(ctx, app.Calendar.Closers())

	if err := server.Start(ctx); err != nil {
		app.Calendar.Logger.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
	<-quit
}
