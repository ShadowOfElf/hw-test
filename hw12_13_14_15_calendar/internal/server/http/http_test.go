package internalhttp

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ShadowOfElf/hw_test/hw12_13_14_15_calendar/internal/app"
	"github.com/ShadowOfElf/hw_test/hw12_13_14_15_calendar/internal/logger"
	"github.com/ShadowOfElf/hw_test/hw12_13_14_15_calendar/internal/storage"
	"github.com/ShadowOfElf/hw_test/hw12_13_14_15_calendar/internal/storage/unityres"
	"github.com/stretchr/testify/require"
)

func TestHTTPHandlers(t *testing.T) { //nolint
	logg := logger.New(logger.DebugLevel)
	store := storage.NewStorage(false)
	application := app.New(logg, store)

	h := NewService(logg, application)

	mux := http.NewServeMux()
	mux.HandleFunc("/", h.Hello)
	mux.HandleFunc("/add", h.AddEvent)
	mux.HandleFunc("/edit/{id}", h.UpdateEvent)
	mux.HandleFunc("/delete/{id}", h.DeleteEvent)
	mux.HandleFunc("/day", h.ListEventByDay)
	mux.HandleFunc("/weak", h.ListEventByWeak)
	mux.HandleFunc("/month", h.ListEventByMonth)

	server := httptest.NewServer(mux)
	defer server.Close()

	baseURL := server.URL

	t.Run("base_test", func(t *testing.T) {
		req, err := http.NewRequestWithContext(context.Background(), "GET", baseURL+"/", nil)
		require.NoError(t, err)
		client := &http.Client{}
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer func() {
			_ = resp.Body.Close()
		}()

		require.Equal(t, http.StatusOK, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		require.Equal(t, "\"Hello world!\"\n", string(body))
	})

	t.Run("Add event", func(t *testing.T) {
		event := map[string]interface{}{
			"id":                 "123",
			"title":              "Team Meeting",
			"date":               "2025-04-29T10:00:00Z",
			"duration":           3600000000000,
			"description":        "Weekly team meeting",
			"userid":             1,
			"notificationMinute": 60000000000,
		}

		jsonData, err := json.Marshal(event)
		require.NoError(t, err)

		req, err := http.NewRequestWithContext(context.Background(), "POST", baseURL+"/add", bytes.NewBuffer(jsonData))
		require.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer func() {
			_ = resp.Body.Close()
		}()

		require.Equal(t, http.StatusCreated, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var response map[string]map[string]interface{}

		err = json.Unmarshal(body, &response)
		require.NoError(t, err)
		require.Equal(t, "Event added successfully", response["data"]["message"])
	})

	t.Run("edit event", func(t *testing.T) {
		event := map[string]interface{}{
			"id":                 "123",
			"title":              "Edit Meeting",
			"date":               "2025-04-29T10:00:00Z",
			"duration":           3600000000000,
			"description":        "Weekly team meeting",
			"userid":             1,
			"notificationMinute": 60000000000,
		}

		jsonData, err := json.Marshal(event)
		require.NoError(t, err)

		req, err := http.NewRequestWithContext(context.Background(), "POST", baseURL+"/edit/123", bytes.NewBuffer(jsonData))
		require.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer func() { _ = resp.Body.Close() }()

		require.Equal(t, http.StatusOK, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var response map[string]map[string]interface{}
		err = json.Unmarshal(body, &response)
		require.NoError(t, err)

		require.Equal(t, "Event update successfully", response["data"]["message"])
	})

	t.Run("list event by day", func(t *testing.T) {
		req, err := http.NewRequestWithContext(context.Background(), "GET", baseURL+"/day?date=2025-04-29", nil)
		require.NoError(t, err)

		client := http.Client{}
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer func() { _ = resp.Body.Close() }()

		require.Equal(t, http.StatusOK, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var response map[string]map[string][]unityres.Event
		err = json.Unmarshal(body, &response)
		require.NoError(t, err)
		events := response["data"]["events"]

		require.Len(t, events, 1)
	})

	t.Run("list event by day for empty day", func(t *testing.T) {
		req, err := http.NewRequestWithContext(context.Background(), "GET", baseURL+"/day?date=2025-04-30", nil)
		require.NoError(t, err)

		client := http.Client{}
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer func() { _ = resp.Body.Close() }()

		require.Equal(t, http.StatusOK, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var response map[string]map[string][]unityres.Event
		err = json.Unmarshal(body, &response)
		require.NoError(t, err)
		events := response["data"]["events"]

		require.Len(t, events, 0)
	})

	t.Run("delete event", func(t *testing.T) {
		eventID := "123"
		req, err := http.NewRequestWithContext(context.Background(), "DELETE", baseURL+"/delete/"+eventID, nil)
		require.NoError(t, err)

		client := http.Client{}
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer func() { _ = resp.Body.Close() }()

		require.Equal(t, http.StatusOK, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		// можно было завести и отдельную структуру, но для примера сделал так
		var response map[string]map[string]interface{}
		err = json.Unmarshal(body, &response)
		require.NoError(t, err)

		require.Equal(t, "Event delete successfully", response["data"]["message"])
	})

	t.Run("list event by day after delete", func(t *testing.T) {
		req, err := http.NewRequestWithContext(context.Background(), "GET", baseURL+"/day?date=2025-04-29", nil)
		require.NoError(t, err)

		client := http.Client{}
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer func() { _ = resp.Body.Close() }()

		require.Equal(t, http.StatusOK, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var response map[string]map[string][]unityres.Event
		err = json.Unmarshal(body, &response)
		require.NoError(t, err)
		events := response["data"]["events"]

		require.Len(t, events, 0)
	})
}
