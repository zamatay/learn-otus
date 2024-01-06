package memorystorage

import (
	"github.com/zamatay/learn-otus/hw12_13_14_15_calendar/internal/domain"
	"reflect"
	"testing"
	"time"
)

func TestStorage_AddEvent(t *testing.T) {
	type args struct {
		event domain.Event
	}
	time := time.Now()
	tests := []struct {
		name string
		args args
	}{
		{
			name: "addEvent",
			args: args{event: domain.Event{ID: 1, Title: "New", Date: time, DateInterval: 10, Description: "description"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New()
			if err := s.AddEvent(tt.args.event); (err != nil) && reflect.DeepEqual(s.storage[1], tt.args.event) {
				t.Errorf("AddEvent() error = %v, wantErr %v", err, s.storage)
			}
		})
	}
}

func TestStorage_EditEvent(t *testing.T) {
	type args struct {
		event domain.Event
	}
	time := time.Now()
	tests := []struct {
		name string
		args args
	}{
		{
			name: "addEvent&EditEvent",
			args: args{event: domain.Event{ID: 1, Title: "New", Date: time, DateInterval: 10, Description: "description"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New()
			if err := s.AddEvent(tt.args.event); err != nil {
				t.Errorf("AddEvent() error = %v, wantErr %v", err, s.storage)
			}
			event := domain.Event{ID: 1, Title: "Edit", Date: time, DateInterval: 10, Description: "description"}
			if err := s.EditEvent(1, event); (err != nil) && reflect.DeepEqual(event, s.storage[1]) {
				t.Errorf("AddEvent() error = %v, wantErr %v", err, s.storage)
			}
		})
	}
}

func TestStorage_RemoveEvent(t *testing.T) {
	time := time.Now()
	tests := []struct {
		name string
	}{
		{
			name: "addEvent&RemoveEvent",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New()
			event1 := domain.Event{ID: 1, Title: "New1", Date: time, DateInterval: 10, Description: "description"}
			event2 := domain.Event{ID: 2, Title: "New2", Date: time, DateInterval: 10, Description: "description"}
			if err := s.AddEvent(event1); err != nil {
				t.Errorf("AddEvent() error = %v, wantErr %v", err, s.storage)
			}
			if err := s.AddEvent(event2); err != nil {
				t.Errorf("AddEvent() error = %v, wantErr %v", err, s.storage)
			}
			if err := s.RemoveEvent(1); (err != nil) && reflect.DeepEqual(event2, s.storage[2]) && len(s.storage) == 1 {
				t.Errorf("AddEvent() error = %v, wantErr %v", err, s.storage)
			}
		})
	}
}
