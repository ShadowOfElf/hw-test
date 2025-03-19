package storage

import (
	"context"
	"time"
)

type Event struct {
	ID                 string
	Title              string
	Date               time.Time
	Duration           time.Duration
	Description        string
	UserID             int
	NotificationMinute time.Duration
}

type UnityStorageInterface interface {
	AddEvent(event Event) error
	EditEvent(id string, event Event) error
	DeleteEvent(id string) error
	ListEventByDate(date time.Time) []Event
	ListEventByWeak(startDate time.Time) []Event
	ListEventByMonth(startDate time.Time) []Event
	Connect(ctx context.Context, config interface{}) error
	Close(ctx context.Context) error
}
