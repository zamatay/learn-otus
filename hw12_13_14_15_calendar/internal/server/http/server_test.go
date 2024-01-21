package internalhttp

import (
	"context"
	v1 "github.com/zamatay/learn-otus/hw12_13_14_15_calendar/api/gen"
	"github.com/zamatay/learn-otus/hw12_13_14_15_calendar/configs"
	"github.com/zamatay/learn-otus/hw12_13_14_15_calendar/internal/app"
	"reflect"
	"testing"
	"time"
)

var (
	server *Server
	ctx    context.Context = context.Background()
)

func init() {
	app.New()
	app.Calendar.Init(ctx)
	server = NewServer(ctx, configs.Grpc{Port: 44044}, configs.HTTP{Port: "8080"})
}

func TestServer_AddEvent(t *testing.T) {
	ti := time.Now()
	tests := []struct {
		name  string
		event v1.EventRequest
		want  v1.OkResponse
	}{
		{
			name: "AllMethod",
			event: v1.EventRequest{
				ID:           0,
				Title:        "Тест",
				Description:  "Тестовое событие",
				Date:         uint64(ti.Unix()),
				DateInterval: 100,
				UserID:       0,
			},
			want: v1.OkResponse{IsOk: true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := server.AddEvent(ctx, &tt.event)
			if err != nil {
				t.Errorf("AddEvent() error = %v", err)
				return
			}
			if got.IsOk != tt.want.IsOk {
				t.Errorf("AddEvent() got = %v, want %v", got, tt.want)
			}

			request := &v1.IdRequest{Id: 0}
			event, err := server.GetEvent(ctx, request)
			if err != nil {
				t.Errorf("GetEvent()")
			}
			if !reflect.DeepEqual(*event, tt.event) {
				t.Errorf("GetEvent() got = %v, want %v", event, tt.event)
			}

			tt.event.Title = "Новое тестовое событие"
			_, err = server.EditEvent(ctx, &tt.event)
			if err != nil {
				t.Errorf("EditEvent()")
			}
			event, err = server.GetEvent(ctx, request)
			if err != nil {
				t.Errorf("GetEvent()")
			}
			if !reflect.DeepEqual(*event, tt.event) {
				t.Errorf("GetEvent() got = %v, want %v", event, tt.event)
			}
			dr := v1.DateRequest{DateFrom: uint64(ti.Unix()), DateTo: uint64(time.Now().Unix())}

			list, err := server.List(ctx, &dr)
			if err != nil {
				t.Errorf("List()")
			}
			if len(list.Data) != 1 {
				t.Errorf("List() got = %d, want %d", len(list.Data), 1)
			}

			_, err = server.RemoveEvent(ctx, request)
			if err != nil {
				t.Errorf("GetEvent()")
			}
			list, err = server.List(ctx, &dr)
			if err != nil {
				t.Errorf("List()")
			}
			if len(list.Data) > 0 {
				t.Errorf("List() got = %d, want %d", len(list.Data), 0)
			}
		})
	}
}
