package sqlstorage

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ShadowOfElf/hw_test/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	db *sql.DB
}

func New() storage.UnityStorageInterface {
	return &Storage{}
}

func (s *Storage) Connect(ctx context.Context, config interface{}) error {
	var err error
	// example "postgres://myuser:mypass@localhost:5432/mydb?sslmode=verify-full"
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s")
	s.db, err = sql.Open("postgres", dsn)
	if err != nil {
		return err
	}

	err = s.db.Ping()
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) Close(ctx context.Context) error {
	// TODO
	return nil
}
