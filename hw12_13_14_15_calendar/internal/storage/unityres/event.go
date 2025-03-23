package unityres

import (
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
