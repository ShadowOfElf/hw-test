package internalgrpc

import (
	"context"
	"testing"
	"time"

	"github.com/ShadowOfElf/hw_test/hw12_13_14_15_calendar/internal/app"
	"github.com/ShadowOfElf/hw_test/hw12_13_14_15_calendar/internal/logger"
	"github.com/ShadowOfElf/hw_test/hw12_13_14_15_calendar/internal/storage"
	pb "github.com/ShadowOfElf/hw_test/hw12_13_14_15_calendar/pkg"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestServicesGRPC(t *testing.T) {
	logg := logger.New(logger.DebugLevel)
	store := storage.NewStorage(false)
	application := app.New(logg, store)
	service := NewServiceGRPC(application)
	date := time.Date(2025, time.April, 29, 0, 0, 0, 0, time.UTC)

	t.Run("grpc create event", func(t *testing.T) {
		req := &pb.CreateEventRequest{
			Event: &pb.EventProto{
				Id:                 "grpc1",
				Title:              "Test Event",
				Date:               timestamppb.New(date),
				Duration:           durationpb.New(60 * time.Minute),
				Description:        "Test Description",
				UserId:             1,
				NotificationMinute: durationpb.New(60 * time.Minute),
			},
		}

		resp, err := service.CreateEventProto(context.Background(), req)
		require.NoError(t, err)
		require.True(t, resp.GetSuccess())
	})

	t.Run("grpc list by day", func(t *testing.T) {
		req := &pb.ListEventByDateRequest{
			Data: timestamppb.New(date),
		}

		resp, err := service.ListEventByDateProto(context.Background(), req)
		require.NoError(t, err)
		require.Len(t, resp.GetEvents(), 1)
		require.Equal(t, resp.GetEvents()[0].Title, "Test Event")
	})

	t.Run("grpc delete event", func(t *testing.T) {
		req := &pb.DeleteEventRequest{Id: "grpc1"}

		resp, err := service.DeleteEventProto(context.Background(), req)
		require.NoError(t, err)
		require.True(t, resp.GetSuccess())
	})

	t.Run("grpc list by day after delete", func(t *testing.T) {
		req := &pb.ListEventByDateRequest{
			Data: timestamppb.New(date),
		}

		_, err := service.ListEventByDateProto(context.Background(), req)
		require.ErrorIs(t, err, ErrNoContent)
	})
}
