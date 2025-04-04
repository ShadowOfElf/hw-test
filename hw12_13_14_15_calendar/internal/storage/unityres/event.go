package unityres

import (
	"time"
)

type Event struct {
	ID                 string        `json:"id"`
	Title              string        `json:"title"`
	Date               time.Time     `json:"date"`
	Duration           time.Duration `json:"duration"`
	Description        string        `json:"description"`
	UserID             int           `json:"userid"`
	NotificationMinute time.Duration `json:"notificationMinute"`
}
