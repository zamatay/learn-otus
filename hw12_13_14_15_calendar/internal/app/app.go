package app

import (
	"context"
	"errors"
	"flag"
	"github.com/zamatay/learn-otus/hw12_13_14_15_calendar/internal/domain"
	"io"
	"os"
	"time"
)

var (
	ErrDateBusy   = errors.New("Date is Busy")
	ErrEventEmpty = errors.New("Event is Empty")
)

var Calendar *App

type App struct {
	Logger  Logger
	Storage Storage
}

type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Error(msg string, args ...any)
	Warn(msg string, args ...any)
	Fatal(msg string)
}

type Storage interface {
	AddEvent(event domain.Event) error
	EditEvent(id int64, event domain.Event) error
	RemoveEvent(id int64) error
	List(beginDate time.Time, endDate time.Time) []domain.Event
	GetEvent(id int64) (domain.Event, error)
}

func New(logger Logger) *App {
	Calendar = &App{
		Logger: logger,
	}
	return Calendar
}

func LoadConfig() {
	var configFile string
	flag.StringVar(&configFile, "config", "/etc/calendar/config.yaml", "Path to configuration file")
	flag.Parse()

	if flag.Arg(0) == "version" {
		PrintVersion()
		os.Exit(1)
	}

}

func (a App) Shutdown(ctx context.Context, quit <-chan os.Signal, closers ...io.Closer) {
	go func() {
		select {
		case <-ctx.Done():
			a.Logger.Info("api - Start - ctx.Done")
		case s := <-quit:
			a.Logger.Info("app - Start - signal: " + s.String())
		}

		for _, close := range closers {
			close.Close()
		}
	}()
}
