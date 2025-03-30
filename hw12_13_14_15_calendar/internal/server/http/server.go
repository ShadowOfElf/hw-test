package internalhttp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/ShadowOfElf/hw_test/hw12_13_14_15_calendar/configs"
	"github.com/ShadowOfElf/hw_test/hw12_13_14_15_calendar/internal/app"
	"github.com/ShadowOfElf/hw_test/hw12_13_14_15_calendar/internal/logger"
	"github.com/pkg/errors"
)

const defaultInterval = 2 * time.Second

func NewServer(log logger.LogInterface, appl *app.App, conf configs.HTTPConf) *Server {
	return &Server{
		Conf:        conf,
		logger:      log,
		application: appl,
	}
}

func (s *Server) Start() error {
	h := NewService(s.logger)
	mux := http.NewServeMux()
	mux.HandleFunc("/", h.Hello)
	logMiddleware := NewHandler(mux, s.logger)

	server := &http.Server{
		Addr:              s.Conf.Addr,
		Handler:           logMiddleware,
		ReadHeaderTimeout: 5 * time.Second,
	}
	s.server = server

	s.logger.Warn(fmt.Sprintf("Server start on address: %s", s.Conf.Addr))
	err := server.ListenAndServe()
	if err != nil && !errors.Is(http.ErrServerClosed, err) {
		return err
	}

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	err := s.server.Shutdown(ctx)
	if err != nil {
		return err
	}
	return nil
}

type Service struct {
	sync.RWMutex
	Stats    map[uint32]uint32
	Interval time.Duration
	logger   logger.LogInterface
}

func NewService(log logger.LogInterface) *Service {
	return &Service{
		Stats:    make(map[uint32]uint32),
		Interval: defaultInterval,
		logger:   log,
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
