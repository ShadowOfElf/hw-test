package memorystorage

import (
	"context"
	"sync"
	"time"

	"github.com/ShadowOfElf/hw_test/hw12_13_14_15_calendar/configs"
	"github.com/ShadowOfElf/hw_test/hw12_13_14_15_calendar/internal/storage/unityres"
)

type EventList map[string]unityres.Event

type StorageMem struct {
	events EventList
	mu     sync.RWMutex //nolint:unused
}

func New() unityres.UnityStorageInterface {
	return &StorageMem{
		events: make(EventList),
	}
}

func (s *StorageMem) ListEventByDate(date time.Time) ([]unityres.Event, error) {
	var result []unityres.Event
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, event := range s.events {
		if date.Year() == event.Date.Year() &&
			date.Month() == event.Date.Month() &&
			date.Day() == event.Date.Day() {
			result = append(result, event)
		}
	}
	return result, nil
}

func (s *StorageMem) ListEventByWeak(startDate time.Time) ([]unityres.Event, error) {
	var result []unityres.Event
	endDate := startDate.AddDate(0, 0, 7)
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, event := range s.events {
		if !event.Date.Before(startDate) && !event.Date.After(endDate) {
			result = append(result, event)
		}
	}
	return result, nil
}

func (s *StorageMem) ListEventByMonth(startDate time.Time) ([]unityres.Event, error) {
	var result []unityres.Event
	endDate := startDate.AddDate(0, 1, 0)
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, event := range s.events {
		if !event.Date.Before(startDate) && !event.Date.After(endDate) {
			result = append(result, event)
		}
	}
	return result, nil
}

func (s *StorageMem) AddEvent(event unityres.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, existEvent := range s.events {
		if event.Date.Equal(existEvent.Date) {
			return unityres.ErrDateBusy
		}
	}
	s.events[event.ID] = event
	return nil // в реализации с БД могут быть ошибки
}

func (s *StorageMem) EditEvent(id string, event unityres.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, existEvent := range s.events {
		if event.Date.Equal(existEvent.Date) {
			return unityres.ErrDateBusy
		}
	}

	s.events[id] = event
	return nil
}

func (s *StorageMem) DeleteEvent(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.events, id)
	return nil
}

func (s *StorageMem) Connect(ctx context.Context, config configs.StorageConf) error {
	_ = ctx
	_ = config
	return nil
}

func (s *StorageMem) Close() error {
	return nil
}
