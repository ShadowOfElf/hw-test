package internalhttp

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ShadowOfElf/hw_test/hw12_13_14_15_calendar/internal/logger"
)

type Logger struct {
	handler http.Handler
	logger  logger.LogInterface
}

func NewHandler(handlerToWrap http.Handler, logI logger.LogInterface) *Logger {
	return &Logger{handler: handlerToWrap, logger: logI}
}

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func (l *Logger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	rw := &responseWriter{ResponseWriter: w}

	l.handler.ServeHTTP(rw, r)

	l.logger.Info(
		fmt.Sprintf(
			"%s [%s] %s %s %s %d %d %s",
			r.RemoteAddr,
			time.Now().Format("02/Jan/2006:15:04:05 -0700"),
			r.Method,
			r.URL.Path,
			r.Proto,
			rw.status,
			time.Since(start).Milliseconds(),
			r.UserAgent(),
		),
	)
}
