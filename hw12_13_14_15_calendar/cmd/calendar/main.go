package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ShadowOfElf/hw_test/hw12_13_14_15_calendar/configs"
	"github.com/ShadowOfElf/hw_test/hw12_13_14_15_calendar/internal/app"
	"github.com/ShadowOfElf/hw_test/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/ShadowOfElf/hw_test/hw12_13_14_15_calendar/internal/server/http"
	"github.com/ShadowOfElf/hw_test/hw12_13_14_15_calendar/internal/storage"
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

	config := configs.NewConfig(configFile)
	logg := logger.New(config.Logger.Level)
	logg.Info("APP Started")

	ctxStor, cancelStorage := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelStorage()

	Storage := storage.NewStorage(config.StorageDB)
	err := Storage.Connect(ctxStor, config.Storage)
	if err != nil {
		logg.Error(fmt.Sprintf("Error connect: %s", err))
	}
	defer func() {
		err := Storage.Close()
		if err != nil {
			logg.Error(fmt.Sprintf("closing DB error: %s", err))
		}
	}()
	//event := unityres.Event{
	//	ID:                 "3",
	//	Title:              "new title 2",
	//	Date:               time.Now(),
	//	Description:        "desc",
	//	Duration:           time.Duration(10 * time.Second),
	//	UserID:             1,
	//	NotificationMinute: time.Duration(10 * time.Second),
	//}
	//err = Storage.AddEvent(event)

	//err = Storage.DeleteEvent("1")
	//if err != nil {
	//	logg.Error(fmt.Sprintf("error on Add: %s", err))
	//}
	fmt.Println("!!!!!!!!!!!!!!!")
	fmt.Println("!!!!!!!!!!!!!!!")
	fmt.Println(Storage.ListEventByMonth(time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.UTC)))
	fmt.Println("!!!!!!!!!!!!!!!")
	fmt.Println("!!!!!!!!!!!!!!!")
	calendar := app.New(logg, Storage)

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
