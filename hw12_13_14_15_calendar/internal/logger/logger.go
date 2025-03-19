package logger

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

type LogLevel string

const (
	DebugLevel LogLevel = "DEBUG"
	InfoLevel  LogLevel = "INFO"
	WarnLevel  LogLevel = "WARN"
	ErrorLevel LogLevel = "ERROR"
)

type LogInterface interface {
	Info(msg string)
	Debug(msg string)
	Error(msg string)
	Warn(msg string)
}

type Logger struct {
	mu       sync.Mutex
	logLevel LogLevel
	output   *log.Logger
}

func New(level LogLevel) LogInterface {
	// тут есть возможность добавить из конфига файл для складывания лога вместо stdout
	return &Logger{
		logLevel: level,
		output:   log.New(os.Stdout, "", 0),
	}
}

func (l *Logger) Info(msg string) {
	l.writeLog(msg, InfoLevel)
}

func (l *Logger) Error(msg string) {
	l.writeLog(msg, ErrorLevel)
}

func (l *Logger) Debug(msg string) {
	l.writeLog(msg, DebugLevel)
}

func (l *Logger) Warn(msg string) {
	l.writeLog(msg, WarnLevel)
}

func (l *Logger) writeLog(msg string, level LogLevel) {
	if l.shouldShow(level) {
		msgFormat := fmt.Sprintf("%s [%s] %s", time.Now().Format("2006-01-02 15:04:05"), level, msg)
		l.output.Println(msgFormat)
	}
}

func (l *Logger) shouldShow(level LogLevel) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	switch l.logLevel {
	case DebugLevel:
		return true
	case InfoLevel:
		return level == InfoLevel || level == WarnLevel || level == ErrorLevel
	case WarnLevel:
		return level == WarnLevel || level == ErrorLevel
	case ErrorLevel:
		return level == ErrorLevel
	default:
		return false
	}
}
