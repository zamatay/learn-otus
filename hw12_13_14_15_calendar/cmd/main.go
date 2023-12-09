package main

import (
	"context"
	"github.com/zamatay/learn-otus/hw12_13_14_15_calendar/configs"
	"github.com/zamatay/learn-otus/hw12_13_14_15_calendar/internal/app"
	"github.com/zamatay/learn-otus/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/zamatay/learn-otus/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/zamatay/learn-otus/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/zamatay/learn-otus/hw12_13_14_15_calendar/internal/storage/sql"
	"io"
	"os"
	"os/signal"
	"syscall"
)

type Constructor func(...any) app.Storage

type CLoserStorage interface {
	io.Closer
	app.Storage
}

func main() {
	// Инициализируем контекст завершения
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	quit := initQuitChannelOnSignals()

	// Подгружаем конфигурацию
	app.LoadConfig()
	config := configs.NewConfig()

	// Инициализируем хранилище
	app.New(logger.New(config.Log.Level))
	storage := getStorage(ctx, config)
	app.Calendar.Storage = storage

	// Стартуем сервер
	server := internalhttp.NewServer(config.Http)

	app.Calendar.Shutdown(ctx, quit, server, storage)

	app.Calendar.Logger.Info("calendar is running...")

	if err := server.Start(ctx); err != nil {
		app.Calendar.Logger.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}

func initQuitChannelOnSignals() <-chan os.Signal {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	return quit
}

func getStorage(ctx context.Context, cfg *configs.Config) CLoserStorage {
	switch cfg.DB.Driver {
	case "inMemory":
		return memorystorage.New()
	case "postgresql":
		return sqlstorage.New(ctx, &cfg.DB)
	default:
		app.Calendar.Logger.Fatal("Неизвестный драйвер")
	}
	return nil
}
