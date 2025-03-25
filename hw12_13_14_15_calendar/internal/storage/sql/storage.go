package sqlstorage

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ShadowOfElf/hw_test/hw12_13_14_15_calendar/configs"
	"github.com/ShadowOfElf/hw_test/hw12_13_14_15_calendar/internal/storage/unityres"
	_ "github.com/jackc/pgx/stdlib" //nolint:depguard,nolintlint
	"github.com/jmoiron/sqlx"
)

// !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
// ДЛЯ РАБОТЫ SQL в консоли необходимо выполнить миграцию
// migrate -path ./migrations -database "postgres://postgres:postgres@localhost/events_db?sslmode=disable" up
// !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!

type Storage struct {
	db *sqlx.DB
}

func New() unityres.UnityStorageInterface {
	return &Storage{}
}

func (s *Storage) Connect(ctx context.Context, config configs.StorageConf) error {
	var err error
	// example "postgres://myuser:mypass@localhost:5432/mydb?sslmode=verify-full"
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=%s",
		config.User, config.Password, config.Address, config.DBName, config.SslMode,
	)
	s.db, err = sqlx.Open("pgx", dsn)
	if err != nil {
		return err
	}

	err = s.db.PingContext(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) Close() error {
	err := s.db.Close()
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) AddEvent(event unityres.Event) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var count int
	err := s.db.QueryRowContext(
		ctx,
		`SELECT COUNT(*) FROM events WHERE date_trunc('hour', date) = $1`,
		event.Date.Format("2006-01-02 15:00:00"),
	).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return unityres.ErrDateBusy
	}

	ctxAdd, cancelAdd := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelAdd()
	queryArgs := map[string]interface{}{
		"id":          event.ID,
		"title":       event.Title,
		"date":        event.Date,
		"duration":    int64(event.Duration.Seconds()),
		"description": event.Description,
		"userID":      event.UserID,
		"not":         int64(event.NotificationMinute.Seconds()),
	}
	_, err = s.db.NamedExecContext(
		ctxAdd,
		`
			INSERT INTO events (id, title, date, duration, description, user_id, notification_minute)
			VALUES (
				:id, :title, :date, make_interval(secs => :duration), :description, :userID, make_interval(secs => :not)
			)
		`,
		queryArgs,
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) EditEvent(id string, event unityres.Event) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var count int
	err := s.db.QueryRowContext(
		ctx,
		"SELECT count(*) FROM events WHERE date_trunc('hour', date) = $1 and id != $2",
		event.Date.Format("2006-01-02 15:00:00"), id,
	).Scan(&count)
	if err != nil {
		return nil
	}
	if count > 0 {
		return unityres.ErrDateBusy
	}

	ctxUpd, cancelUpd := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelUpd()

	argUpd := map[string]interface{}{
		"id":          event.ID,
		"title":       event.Title,
		"date":        event.Date,
		"duration":    int64(event.Duration.Seconds()),
		"description": event.Description,
		"userID":      event.UserID,
		"not":         int64(event.NotificationMinute.Seconds()),
	}

	_, err = s.db.NamedExecContext(
		ctxUpd,
		`
			UPDATE events
			SET title = :title, date = :date, duration = make_interval(secs => :duration), description = :description, 
			    user_id = :userID, notification_minute = make_interval(secs => :not)
			WHERE id = :id
			`,
		argUpd,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) DeleteEvent(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := s.db.ExecContext(ctx, "DELETE FROM events WHERE id = $1", id)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) ListEventByDate(date time.Time) ([]unityres.Event, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	rows, err := s.db.QueryContext(
		ctx,
		`
		SELECT id, title, date, duration, description, user_id, notification_minute FROM events WHERE DATE(date) = $1
		`,
		date.Format("2006-01-02"),
	)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()

	return rowsToEvents(rows)
}

func (s *Storage) ListEventByWeak(startDate time.Time) ([]unityres.Event, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	endDate := startDate.AddDate(0, 0, 7)
	rows, err := s.db.QueryContext(
		ctx,
		`
		SELECT id, title, date, duration, description, user_id, notification_minute FROM events WHERE date > $1 and date < $2
		`,
		startDate.Format("2006-01-02"),
		endDate.Format("2006-01-02"),
	)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()

	return rowsToEvents(rows)
}

func (s *Storage) ListEventByMonth(startDate time.Time) ([]unityres.Event, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	endDate := startDate.AddDate(0, 1, 0)
	rows, err := s.db.QueryContext(
		ctx,
		`
		SELECT id, title, date, duration, description, user_id, notification_minute FROM events WHERE date > $1 and date < $2
		`,
		startDate.Format("2006-01-02"),
		endDate.Format("2006-01-02"),
	)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()

	return rowsToEvents(rows)
}

func rowsToEvents(rows *sql.Rows) ([]unityres.Event, error) {
	var events []unityres.Event
	for rows.Next() {
		var id, title, description, duration, notif string
		var dateT time.Time
		var userID int

		err := rows.Scan(&id, &title, &dateT, &duration, &description, &userID, &notif)
		if err != nil {
			return nil, err
		}

		dur, err := strToDur(duration)
		if err != nil {
			return nil, err
		}

		not, err := strToDur(notif)
		if err != nil {
			return nil, err
		}

		newEvent := unityres.Event{
			ID:                 id,
			Title:              title,
			Date:               dateT,
			Duration:           dur,
			Description:        description,
			UserID:             userID,
			NotificationMinute: not,
		}
		events = append(events, newEvent)
	}
	return events, nil
}

func strToDur(strDur string) (time.Duration, error) {
	durationParts := strings.Split(strDur, ":") // Разделяем строку по ":"
	if len(durationParts) != 3 {
		return 0, fmt.Errorf("invalid duration format: %s", strDur)
	}

	hours, err := parseInt(durationParts[0])
	if err != nil {
		return 0, err
	}
	minutes, err := parseInt(durationParts[1])
	if err != nil {
		return 0, err
	}
	seconds, err := parseInt(durationParts[2])
	if err != nil {
		return 0, err
	}
	dur := time.Duration(hours)*time.Hour +
		time.Duration(minutes)*time.Minute +
		time.Duration(seconds)*time.Second
	return dur, nil
}

func parseInt(s string) (int, error) {
	value, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	return value, nil
}
