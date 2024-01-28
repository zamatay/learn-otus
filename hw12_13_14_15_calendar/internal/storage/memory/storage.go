package memorystorage

import (
	"sync"
	"time"

	"github.com/zamatay/learn-otus/hw12_13_14_15_calendar/internal/domain"
)

const initialSize = 100

type Storage struct {
	mu      sync.RWMutex
	storage map[int64]domain.Event
}

func (s Storage) Close() error {
	return nil
}

func (s Storage) AddEvent(event domain.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.storage[event.ID] = event
	return nil
}

func (s Storage) EditEvent(_ int64, event domain.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.storage[event.ID] = event
	return nil
}

func (s Storage) RemoveEvent(id int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.storage, id)
	return nil
}

func (s Storage) List(beginDate time.Time, endDate time.Time) []domain.Event {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]domain.Event, 0, initialSize)
	for _, event := range s.storage {
		if (event.Date.After(beginDate) && event.Date.Before(endDate)) ||
			(event.Date.Add(event.DateInterval).After(beginDate) && event.Date.Add(event.DateInterval).Before(endDate)) ||
			(event.Date.Equal(beginDate)) {
			result = append(result, event)
		}
	}
	return result
}

func (s Storage) GetEvent(id int64) (domain.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.storage[id], nil
}

func New() *Storage {
	return &Storage{
		storage: make(map[int64]domain.Event, initialSize),
	}
}
