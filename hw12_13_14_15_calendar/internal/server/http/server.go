package internalhttp

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/ShadowOfElf/hw_test/hw12_13_14_15_calendar/configs"
	"github.com/ShadowOfElf/hw_test/hw12_13_14_15_calendar/internal/app"
	"github.com/pkg/errors"
)

const defaultInterval = 2 * time.Second

type Server struct {
	Conf        configs.HTTPConf
	application *app.App
	server      *http.Server
}

func NewServer(appl *app.App, conf configs.HTTPConf) *Server {
	return &Server{
		Conf:        conf,
		application: appl,
	}
}

func (s *Server) Start() error {
	h := NewService(s.application.Logger, s.application)
	mux := http.NewServeMux()
	mux.HandleFunc("/", h.Hello)
	mux.HandleFunc("/add", h.AddEvent)
	mux.HandleFunc("/edit/{id}", h.UpdateEvent)
	mux.HandleFunc("/delete/{id}", h.DeleteEvent)
	mux.HandleFunc("/day", h.ListEventByDay)
	mux.HandleFunc("/weak", h.ListEventByWeak)
	mux.HandleFunc("/month", h.ListEventByMonth)

	logMiddleware := NewHandler(mux, s.application.Logger)

	server := &http.Server{
		Addr:              s.Conf.Addr,
		Handler:           logMiddleware,
		ReadHeaderTimeout: 5 * time.Second,
	}
	s.server = server

	s.application.Logger.Warn(fmt.Sprintf("Server start on address: %s", s.Conf.Addr))
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
