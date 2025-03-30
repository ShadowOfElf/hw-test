package memorystorage

import (
	"testing"
	"time"

	"github.com/ShadowOfElf/hw_test/hw12_13_14_15_calendar/internal/storage/unityres"
	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	memStorage := New()

	t.Run("testAdd", func(t *testing.T) {
		newEvent := unityres.Event{
			ID:                 "1",
			Title:              "title",
			Date:               time.Now(),
			Description:        "desc",
			Duration:           10 * time.Second,
			UserID:             1,
			NotificationMinute: 10 * time.Second,
		}
		newEvent2 := unityres.Event{
			ID:                 "2",
			Title:              "title",
			Date:               time.Now(),
			Description:        "desc",
			Duration:           10 * time.Second,
			UserID:             1,
			NotificationMinute: 10 * time.Second,
		}
		err := memStorage.AddEvent(newEvent)
		require.NoError(t, err)
		err = memStorage.AddEvent(newEvent2)
		require.NoError(t, err)
	})

	t.Run("EditEvent", func(t *testing.T) {
		editingEvent := unityres.Event{
			ID:                 "1",
			Title:              "new title",
			Date:               time.Now(),
			Description:        "desc 2",
			Duration:           10 * time.Second,
			UserID:             2,
			NotificationMinute: 10 * time.Second,
		}
		err := memStorage.EditEvent("1", editingEvent)
		require.NoError(t, err)
	})

	t.Run("list events", func(t *testing.T) {
		events, err := memStorage.ListEventByMonth(
			time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(),
				0, 0, 0, 0, time.UTC),
		)
		require.NoError(t, err)
		require.Len(t, events, 2)
		require.Equal(t, "1", events[0].ID)
		require.Equal(t, "new title", events[0].Title)
	})

	t.Run("delete event", func(t *testing.T) {
		err := memStorage.DeleteEvent("2")
		require.NoError(t, err)

		events, err := memStorage.ListEventByMonth(
			time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(),
				0, 0, 0, 0, time.UTC),
		)
		require.NoError(t, err)
		require.Len(t, events, 1)
	})
}
