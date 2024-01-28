package internalhttp

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	v1 "github.com/zamatay/learn-otus/hw12_13_14_15_calendar/api/gen"
	"github.com/zamatay/learn-otus/hw12_13_14_15_calendar/configs"
	"github.com/zamatay/learn-otus/hw12_13_14_15_calendar/internal/app"
	"github.com/zamatay/learn-otus/hw12_13_14_15_calendar/internal/domain"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
	"net/http"
	"time"
)

type GetHost interface {
	GetHost() string
}

type Server struct {
	v1.UnimplementedEventsServer
	grpcCfg configs.Grpc
	httpCfg configs.HTTP
	server  *grpc.Server
}

func (server *Server) GetHost() string {
	return fmt.Sprintf(":%d", server.grpcCfg.Port)
}

func (server *Server) AddEvent(ctx context.Context, request *v1.EventRequest) (*v1.OkResponse, error) {
	app.Calendar.Storage.AddEvent(*domain.NewEvent(request.ID, request.UserID, request.Title, request.Description, request.DateInterval, request.Date))
	return &v1.OkResponse{IsOk: true}, nil
}

func (server *Server) EditEvent(ctx context.Context, request *v1.EventRequest) (*v1.OkResponse, error) {
	event := *domain.NewEvent(request.ID, request.UserID, request.Title, request.Description, request.DateInterval, request.Date)
	app.Calendar.Storage.EditEvent(event.ID, event)
	return &v1.OkResponse{}, nil
}

func (server *Server) RemoveEvent(ctx context.Context, request *v1.IdRequest) (*v1.OkResponse, error) {
	app.Calendar.Storage.RemoveEvent(request.Id)
	return &v1.OkResponse{}, nil
}

func (server *Server) List(ctx context.Context, request *v1.DateRequest) (*v1.EventDataSet, error) {
	events := app.Calendar.Storage.List(time.Unix(int64(request.DateFrom), 0), time.Unix(int64(request.DateTo), 0))
	eds := v1.EventDataSet{
		Data: make([]*v1.EventRequest, 0, len(events)),
	}
	for _, event := range events {
		eds.Data = append(eds.Data, getEvent(event))
	}
	return &eds, nil
}

func getEvent(event domain.Event) *v1.EventRequest {
	return &v1.EventRequest{ID: event.ID, Title: event.Title, Description: event.Description,
		Date: uint64(event.Date.Unix()), DateInterval: int32(event.DateInterval), UserID: int32(event.UserID)}
}

func (server *Server) GetEvent(ctx context.Context, request *v1.IdRequest) (er *v1.EventRequest, err error) {
	if event, err := app.Calendar.Storage.GetEvent(request.Id); err != nil {
		app.Calendar.Logger.Error("Ошибка при получении GetEvent", err)
	} else {
		er = getEvent(event)
	}
	return er, nil
}

func (server *Server) runRest(ctx context.Context) {
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := v1.RegisterEventsHandlerFromEndpoint(ctx, mux, server.GetHost(), opts)
	if err != nil {
		app.Calendar.Logger.Error("Ошибка при регистрации ендпоинта", err)
	}

	httpSrv := &http.Server{Addr: fmt.Sprintf(":%s", server.httpCfg.Port)}
	httpSrv.Handler = logMiddleware(mux)

	if err := httpSrv.ListenAndServe(); err != nil {
		app.Calendar.Logger.Error("Ошибка при запуске http сервера", err)
	}
}

func NewServer(ctx context.Context, grpcCfg configs.Grpc, httpCfg configs.HTTP) *Server {
	s := &Server{
		grpcCfg: grpcCfg,
		httpCfg: httpCfg,
		server:  grpc.NewServer(),
	}
	v1.RegisterEventsServer(s.server, s)
	return s
}

func (s *Server) getListener() net.Listener {
	lis, err := net.Listen("tcp", s.GetHost())
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	return lis
}

func (s *Server) Start(ctx context.Context) error {
	app.Calendar.Logger.Info("Запускаем календарь...")

	go s.runRest(ctx)

	if err := s.server.Serve(s.getListener()); err != nil {
		app.Calendar.Logger.Error("Не смогли запустить http сервер" + err.Error())
		return err
	}

	return nil
}

func (s *Server) Close() error {
	return nil
}
