package app

import (
	"time"

	"github.com/ShadowOfElf/hw_test/hw12_13_14_15_calendar/internal/logger"
	"github.com/ShadowOfElf/hw_test/hw12_13_14_15_calendar/internal/storage/unityres"
)

type App struct {
	Logger  logger.LogInterface
	Storage unityres.UnityStorageInterface
}

func New(logg logger.LogInterface, storage unityres.UnityStorageInterface) *App {
	return &App{
		Logger:  logg,
		Storage: storage,
	}
}

func (a *App) CreateEvent(event unityres.Event) error {
	return a.Storage.AddEvent(event)
}

func (a *App) EditEvent(id string, event unityres.Event) error {
	return a.Storage.EditEvent(id, event)
}

func (a *App) DeleteEvent(id string) error {
	return a.Storage.DeleteEvent(id)
}

func (a *App) ListEventByDay(date time.Time) ([]unityres.Event, error) {
	return a.Storage.ListEventByDate(date)
}

func (a *App) ListEventByWeak(startDate time.Time) ([]unityres.Event, error) {
	return a.Storage.ListEventByWeak(startDate)
}

func (a *App) ListEventByMonth(startDate time.Time) ([]unityres.Event, error) {
	return a.Storage.ListEventByMonth(startDate)
}
