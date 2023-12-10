package app

import (
	"context"
	"errors"
	"flag"
	"github.com/zamatay/learn-otus/hw12_13_14_15_calendar/configs"
	"github.com/zamatay/learn-otus/hw12_13_14_15_calendar/internal/domain"
	"github.com/zamatay/learn-otus/hw12_13_14_15_calendar/internal/logger"
	memorystorage "github.com/zamatay/learn-otus/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/zamatay/learn-otus/hw12_13_14_15_calendar/internal/storage/sql"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"
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

func (a App) LoadConfig() {
	var configFile string
	flag.StringVar(&configFile, "config", "/etc/calendar/config.yaml", "Path to configuration file")
	flag.Parse()

	if flag.Arg(0) == "version" {
		PrintVersion()
		os.Exit(1)
	}
}

func (a App) Init(ctx context.Context) {
	a.LoadConfig()
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
	case "inMemory":
		return memorystorage.New()
	case "postgresql":
		return sqlstorage.New(ctx, &cfg.DB)
	default:
		Calendar.Logger.Fatal("Неизвестный драйвер")
	}
	return nil
}

func Shutdown(ctx context.Context, quit <-chan os.Signal, closers []io.Closer) {
	go func() {
		select {
		case <-ctx.Done():
			log.Info("api - Start - ctx.Done")
		case s := <-quit:
			log.Info("app - Start - signal: " + s.String())
		}

		for _, close := range closers {
			close.Close()
		}
	}()
}

func InitShutDowner() (context.Context, chan os.Signal, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	return ctx, quit, cancel
}
