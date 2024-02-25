package app

import (
	"context"
	"errors"
	memorystorage "github.com/zamatay/learn-otus/hw12_13_14_15_calendar/internal/storage/memory"
	"io"
	"os/signal"
	"syscall"
	"time"

	"github.com/zamatay/learn-otus/hw12_13_14_15_calendar/configs"
	"github.com/zamatay/learn-otus/hw12_13_14_15_calendar/internal/domain"
	"github.com/zamatay/learn-otus/hw12_13_14_15_calendar/internal/logger"
	sqlstorage "github.com/zamatay/learn-otus/hw12_13_14_15_calendar/internal/storage/sql"
)

var (
	ErrDateBusy   = errors.New("Date is Busy")
	ErrEventEmpty = errors.New("Event is Empty")
)

var (
	Calendar *App
	log      Logger
)

type App struct {
	Logger  Logger
	Storage Storage
	closers []io.Closer
	Config  *configs.Config
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

type CLoserStorage interface {
	io.Closer
	Storage
}

func New() *App {
	Calendar = &App{
		closers: make([]io.Closer, 0, 2),
		Config:  configs.NewConfig(),
	}
	return Calendar
}

func (a *App) AddClosers(closer io.Closer) {
	a.closers = append(a.closers, closer)
}

func (a *App) Init(ctx context.Context) {
	//a.LoadConfig()
	storage := getStorage(ctx, a.Config)
	a.AddClosers(storage)
	a.Storage = storage
	log = logger.New(a.Config.Log.Level)
	a.Logger = log
}

func (a *App) Closers() []io.Closer {
	return a.closers
}

func getStorage(ctx context.Context, cfg *configs.Config) CLoserStorage {
	switch cfg.DB.Driver {
	case "postgres":
		storage, _ := sqlstorage.New(ctx, &cfg.DB)
		return storage
	default:
		return memorystorage.New()
	}
	return nil
}

func Shutdown(ctx context.Context, closers []io.Closer) chan struct{} {
	quit := make(chan struct{})
	go func() {
		select {
		case <-ctx.Done():
			log.Info("api - Start - ctx.Done")
		}

		for _, close := range closers {
			if err := close.Close(); err != nil {
				log.Error("error close", err)
			}
		}
		quit <- struct{}{}
		close(quit)
	}()
	return quit
}

func InitShutDowner() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	return ctx, cancel
}
