package unityres

import (
	"context"
	"errors"
	"time"

	"github.com/ShadowOfElf/hw_test/hw12_13_14_15_calendar/configs"
)

type UnityStorageInterface interface {
	AddEvent(event Event) error
	EditEvent(id string, event Event) error
	DeleteEvent(id string) error
	ListEventByDate(date time.Time) ([]Event, error)
	ListEventByWeak(startDate time.Time) ([]Event, error)
	ListEventByMonth(startDate time.Time) ([]Event, error)
	Connect(ctx context.Context, config configs.StorageConf) error
	Close() error
}

var ErrDateBusy = errors.New("event already exists for this date")
