package internalhttp

import (
	"context"
	"fmt"
	"github.com/zamatay/learn-otus/hw12_13_14_15_calendar/configs"
	"github.com/zamatay/learn-otus/hw12_13_14_15_calendar/internal/app"

	"net/http"
)

type Server struct {
	cfg     configs.Http
	httpSrv *http.Server
}

func NewServer(cfg configs.Http) *Server {
	httpSrv := &http.Server{Addr: fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)}
	m := http.NewServeMux()

	m.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
		w.WriteHeader(http.StatusOK)
	})

	httpSrv.Handler = loggingMiddleware(m)

	return &Server{
		cfg:     cfg,
		httpSrv: httpSrv,
	}
}

func (s *Server) Start(ctx context.Context) error {
	err := s.httpSrv.ListenAndServe()
	if err != nil {
		app.Calendar.Logger.Fatal("Не смогли запустить http сервер" + err.Error())
	}
	app.Calendar.Logger.Info("http сервер запущен на порту", "port: ", s.cfg.Port)

	<-ctx.Done()
	return nil
}

func (s *Server) Close() error {

	return nil
}
