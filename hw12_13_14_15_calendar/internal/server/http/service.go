package internalhttp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/ShadowOfElf/hw_test/hw12_13_14_15_calendar/internal/app"
	"github.com/ShadowOfElf/hw_test/hw12_13_14_15_calendar/internal/logger"
	"github.com/ShadowOfElf/hw_test/hw12_13_14_15_calendar/internal/storage/unityres"
)

type Response struct {
	Data  interface{} `json:"data"`
	Error struct {
		Message string `json:"message,omitempty"`
	} `json:"error,omitempty"`
}

type EventsResponse struct {
	Events []unityres.Event `json:"events,omitempty"`
}

type Service struct {
	mu          sync.RWMutex
	application *app.App
	Stats       map[uint32]uint32
	Interval    time.Duration
	logger      logger.LogInterface
}

func NewService(log logger.LogInterface, application *app.App) *Service {
	return &Service{
		application: application,
		Stats:       make(map[uint32]uint32),
		Interval:    defaultInterval,
		logger:      log,
	}
}

func (s *Service) Hello(w http.ResponseWriter, r *http.Request) {
	_ = r
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err := json.NewEncoder(w).Encode("Hello world!")
	if err != nil {
		s.logger.Error(fmt.Sprintf("resp marshal error: %s", err))
	}
}

func (s *Service) AddEvent(w http.ResponseWriter, r *http.Request) {
	resp := &Response{}
	event, ok := getEventFromForm(w, r, resp, s.logger)
	if !ok {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	err := s.application.CreateEvent(event)
	if err != nil {
		resp.Error.Message = fmt.Sprintf("failed add to DB: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		WriteResponse(w, resp, s.logger)
		return
	}

	resp.Data = map[string]interface{}{
		"message": "Event added successfully",
		"event":   event,
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	WriteResponse(w, resp, s.logger)
}

func (s *Service) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	resp := &Response{}

	event, ok := getEventFromForm(w, r, resp, s.logger)
	if !ok {
		return
	}
	eventID := r.PathValue("id")
	if eventID == "" {
		resp.Error.Message = "failed get event id"
		w.WriteHeader(http.StatusBadRequest)
		WriteResponse(w, resp, s.logger)
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	err := s.application.EditEvent(eventID, event)
	if err != nil {
		resp.Error.Message = fmt.Sprintf("failed add to DB: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		WriteResponse(w, resp, s.logger)
		return
	}

	resp.Data = map[string]interface{}{
		"message": "Event update successfully",
		"event":   event,
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	WriteResponse(w, resp, s.logger)
}

func (s *Service) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	resp := &Response{}

	if r.Method != http.MethodDelete {
		resp.Error.Message = fmt.Sprintf("method %s not not supported on uri %s", r.Method, r.URL.Path)
		w.WriteHeader(http.StatusMethodNotAllowed)
		WriteResponse(w, resp, s.logger)
		return
	}

	eventID := r.PathValue("id")
	if eventID == "" {
		resp.Error.Message = "failed get event id"
		w.WriteHeader(http.StatusBadRequest)
		WriteResponse(w, resp, s.logger)
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	err := s.application.DeleteEvent(eventID)
	if err != nil {
		resp.Error.Message = fmt.Sprintf("failed delete from DB: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		WriteResponse(w, resp, s.logger)
		return
	}

	resp.Data = map[string]interface{}{
		"message": "Event delete successfully",
		"eventID": eventID,
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	WriteResponse(w, resp, s.logger)
}

func (s *Service) ListEventByDay(w http.ResponseWriter, r *http.Request) {
	resp := &Response{}
	date, ok := getDateFromReq(w, r, resp, s.logger)
	if !ok {
		return
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	events, err := s.application.ListEventByDay(date)
	if err != nil {
		s.logger.Error(fmt.Sprintf("error when receiving events: %s", err))
	}

	resp.Data = &EventsResponse{Events: events}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	WriteResponse(w, resp, s.logger)
}

func (s *Service) ListEventByWeak(w http.ResponseWriter, r *http.Request) {
	resp := &Response{}
	startDate, ok := getDateFromReq(w, r, resp, s.logger)
	if !ok {
		return
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	events, err := s.application.ListEventByWeak(startDate)
	if err != nil {
		s.logger.Error(fmt.Sprintf("error when receiving events: %s", err))
	}

	resp.Data = &EventsResponse{Events: events}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	WriteResponse(w, resp, s.logger)
}

func (s *Service) ListEventByMonth(w http.ResponseWriter, r *http.Request) {
	resp := &Response{}
	startDate, ok := getDateFromReq(w, r, resp, s.logger)
	if !ok {
		return
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	events, err := s.application.ListEventByMonth(startDate)
	if err != nil {
		s.logger.Error(fmt.Sprintf("error when receiving events: %s", err))
	}

	resp.Data = &EventsResponse{Events: events}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	WriteResponse(w, resp, s.logger)
}

func WriteResponse(w http.ResponseWriter, resp *Response, logg logger.LogInterface) {
	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		logg.Error(fmt.Sprintf("response marshal error: %s", err))
	}
}

func getDateFromReq(
	w http.ResponseWriter, r *http.Request, resp *Response, logg logger.LogInterface,
) (time.Time, bool) {
	if r.Method != http.MethodGet {
		resp.Error.Message = fmt.Sprintf("method %s not not supported on uri %s", r.Method, r.URL.Path)
		w.WriteHeader(http.StatusMethodNotAllowed)
		WriteResponse(w, resp, logg)
		return time.Time{}, false
	}
	args := r.URL.Query()

	dateHTTP := args.Get("date")
	date, err := time.Parse("2006-01-02", dateHTTP)
	if err != nil {
		resp.Error.Message = fmt.Sprintf("error in date parse: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		WriteResponse(w, resp, logg)
		return time.Time{}, false
	}
	return date, true
}

func getEventFromForm(
	w http.ResponseWriter, r *http.Request, resp *Response, logg logger.LogInterface,
) (unityres.Event, bool) {
	if r.Method != http.MethodPost {
		resp.Error.Message = fmt.Sprintf("method %s not not supported on uri %s", r.Method, r.URL.Path)
		w.WriteHeader(http.StatusMethodNotAllowed)
		WriteResponse(w, resp, logg)
		return unityres.Event{}, false
	}

	var event unityres.Event

	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		resp.Error.Message = fmt.Sprintf("failed parse form data: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		WriteResponse(w, resp, logg)
		return unityres.Event{}, false
	}

	return event, true
}
