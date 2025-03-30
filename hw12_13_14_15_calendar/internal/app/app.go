package app

import (
	"context"
)

type App struct { // TODO
}

type Logger interface { // TODO
}

type Storage interface { // TODO
}

func New(logger Logger, storage Storage) *App {
	_ = logger
	_ = storage
	return &App{}
}

func (a *App) CreateEvent(ctx context.Context, id, title string) error {
	_ = ctx
	_ = id
	_ = title
	return nil
	// return a.storage.CreateEvent(storage.Event{ID: id, Title: title})
}
