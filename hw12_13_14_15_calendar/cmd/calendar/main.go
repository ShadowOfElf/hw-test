package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ShadowOfElf/hw_test/hw12_13_14_15_calendar/internal/app"
	"github.com/ShadowOfElf/hw_test/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/ShadowOfElf/hw_test/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/ShadowOfElf/hw_test/hw12_13_14_15_calendar/internal/storage/memory"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	config := NewConfig()
	logg := logger.New(config.Logger.Level)
	logg.Info("APP Started")

	memStorage := memorystorage.New()
	//storage = sqlstorage.New()
	//event := storage.Event{
	//	ID:                 "1",
	//	Title:              "sdfsdf",
	//	Date:               time.Now(),
	//	Description:        "desc",
	//	Duration:           time.Duration(10),
	//	UserID:             1,
	//	NotificationMinute: time.Duration(10),
	//}
	//_ = memStorage.AddEvent(event)
	//fmt.Println(memStorage.ListEventByMonth(time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, time.UTC)))
	calendar := app.New(logg, memStorage)

	server := internalhttp.NewServer(logg, calendar)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")

	if err := server.Start(ctx); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}
