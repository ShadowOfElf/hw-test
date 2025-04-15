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
	grpcclient "github.com/ShadowOfElf/hw_test/hw12_13_14_15_calendar/internal/client/grpc_client"
	"github.com/ShadowOfElf/hw_test/hw12_13_14_15_calendar/internal/logger"
	"github.com/ShadowOfElf/hw_test/hw12_13_14_15_calendar/internal/queue"
	internalgrpc "github.com/ShadowOfElf/hw_test/hw12_13_14_15_calendar/internal/server/grpc"
	internalhttp "github.com/ShadowOfElf/hw_test/hw12_13_14_15_calendar/internal/server/http"
	"github.com/ShadowOfElf/hw_test/hw12_13_14_15_calendar/internal/storage"
	"github.com/ShadowOfElf/hw_test/hw12_13_14_15_calendar/internal/storage/unityres"
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
		logg.Error(fmt.Sprintf("Error connect DB: %s", err))
	}
	defer func() {
		err := Storage.Close()
		if err != nil {
			logg.Error(fmt.Sprintf("closing DB error: %s", err))
		}
	}()

	var runApp func(logg logger.LogInterface, Storage unityres.UnityStorageInterface, config configs.Config)
	switch config.AppType {
	case configs.MainApp:
		runApp = mainApp
	case configs.SchedulerApp:
		runApp = SchedulerApp
	case configs.SenderApp:
		runApp = SenderApp
	}
	runApp(logg, Storage, config)
}

func mainApp(logg logger.LogInterface, storage unityres.UnityStorageInterface, config configs.Config) {
	calendar := app.New(logg, storage)

	server := internalhttp.NewServer(calendar, config.HTTP)
	grpc := internalgrpc.NewGRPCServer(calendar, config.GRPC)

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

		if err := grpc.Stop(); err != nil {
			logg.Error("failed to stop grpc server: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")

	if err := grpc.Start(); err != nil {
		logg.Error("failed to start grpc server: " + err.Error())
	}

	if err := server.Start(); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}

func SchedulerApp(logg logger.LogInterface, _ unityres.UnityStorageInterface, config configs.Config) {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	ticker := time.NewTicker(config.UpdateTime)

	go func() {
		<-ctx.Done()
		ticker.Stop()
	}()

	rabbit := queue.NewServiceMQ(logg, config.Rabbit)
	err := rabbit.Start()
	if err != nil {
		logg.Error(fmt.Sprintf("failed connect RabbitMQ: %s", err))
		return
	}
	defer func() {
		_ = rabbit.Stop()
	}()
	logg.Warn("scheduler is running...")

	grpcService := grpcclient.NewGRPClient(logg, config.GRPC)
	client, err := grpcService.Start()
	if err != nil {
		logg.Error(fmt.Sprintf("failed create grpc client: %s", err))
		return
	}

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			resp, err := client.ListEventByDateProto(ctx, grpcService.GetReqByDay(time.Now()))
			if err != nil {
				logg.Error(fmt.Sprintf("failed get events from DB: %s", err))
				continue
			}

			for _, event := range resp.Events {
				notificationTime := event.Date.AsTime().Add(-event.NotificationMinute.AsDuration())
				if notificationTime.Before(time.Now()) && event.Date.AsTime().After(time.Now()) {
					notif := unityres.Notification{
						EventID:    event.Id,
						EventTitle: event.Title,
						EventDate:  event.Date.AsTime(),
						UserID:     int(event.UserId),
					}
					err = rabbit.SendNotification(context.Background(), notif)
					if err != nil {
						logg.Error(fmt.Sprintf("failed send notif: %s", err))
					}
				}
			}
		}
	}
}

func SenderApp(logg logger.LogInterface, _ unityres.UnityStorageInterface, config configs.Config) {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	rabbit := queue.NewServiceMQ(logg, config.Rabbit)
	err := rabbit.Start()
	if err != nil {
		logg.Error(fmt.Sprintf("failed connect RabbitMQ: %s", err))
		return
	}
	defer func() {
		_ = rabbit.Stop()
	}()
	logg.Warn("sender is running...")

	messages, err := rabbit.GetNotification(ctx)
	if err != nil {
		logg.Error(fmt.Sprintf("failed get message from RabbitMQ: %s", err))
		return
	}
	for msg := range messages {
		logg.Info(
			fmt.Sprintf(
				"Notification for user %v: \n event_id: %s event: %s start: %s ",
				msg.UserID, msg.EventID, msg.EventTitle, msg.EventDate,
			),
		)
	}
}
