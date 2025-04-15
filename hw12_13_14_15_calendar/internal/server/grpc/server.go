package internalgrpc

import (
	"context"
	"errors"
	"fmt"
	"net"
	"sync"

	"github.com/ShadowOfElf/hw_test/hw12_13_14_15_calendar/configs"
	"github.com/ShadowOfElf/hw_test/hw12_13_14_15_calendar/internal/app"
	"github.com/ShadowOfElf/hw_test/hw12_13_14_15_calendar/internal/storage/unityres"
	pb "github.com/ShadowOfElf/hw_test/hw12_13_14_15_calendar/pkg"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var ErrNoContent = errors.New("no events in response")

type ServiceGRPC struct {
	pb.UnimplementedCalendarServer
	mu  sync.RWMutex
	app *app.App
}

func NewServiceGRPC(app *app.App) *ServiceGRPC {
	return &ServiceGRPC{
		app: app,
	}
}

func (s *ServiceGRPC) CreateEventProto(ctx context.Context, req *pb.CreateEventRequest) (*pb.EventResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	_ = ctx

	err := s.app.CreateEvent(messageToEvent(req.GetEvent()))
	if err != nil {
		return nil, err
	}
	return &pb.EventResponse{Success: true}, nil
}

func (s *ServiceGRPC) EditEventProto(ctx context.Context, req *pb.EditEventRequest) (*pb.EventResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	_ = ctx

	err := s.app.EditEvent(req.GetId(), messageToEvent(req.GetEvent()))
	if err != nil {
		return nil, err
	}
	return &pb.EventResponse{Success: true}, nil
}

func (s *ServiceGRPC) DeleteEventProto(ctx context.Context, req *pb.DeleteEventRequest) (*pb.EventResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	_ = ctx

	err := s.app.DeleteEvent(req.GetId())
	if err != nil {
		return nil, err
	}

	return &pb.EventResponse{Success: true}, nil
}

func (s *ServiceGRPC) ListEventByDateProto(
	ctx context.Context, req *pb.ListEventByDateRequest,
) (*pb.ListEventResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_ = ctx

	date := req.GetData().AsTime()

	events, err := s.app.ListEventByDay(date)
	if err != nil {
		return nil, err
	}

	return eventToMessage(events), nil
}

func (s *ServiceGRPC) ListEventByWeakProto(
	ctx context.Context, req *pb.ListEventByWeakRequest,
) (*pb.ListEventResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_ = ctx

	startData := req.GetStartData().AsTime()

	events, err := s.app.ListEventByWeak(startData)
	if err != nil {
		return nil, err
	}

	return eventToMessage(events), nil
}

func (s *ServiceGRPC) ListEventByMonthProto(
	ctx context.Context, req *pb.ListEventByMonthRequest,
) (*pb.ListEventResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_ = ctx

	startData := req.GetStartData().AsTime()

	events, err := s.app.ListEventByMonth(startData)
	if err != nil {
		return nil, err
	}

	return eventToMessage(events), nil
}

func messageToEvent(event *pb.EventProto) unityres.Event {
	return unityres.Event{
		ID:                 event.GetId(),
		Title:              event.GetTitle(),
		Date:               event.GetDate().AsTime(),
		Duration:           event.GetDuration().AsDuration(),
		Description:        event.GetDescription(),
		UserID:             int(event.GetUserId()),
		NotificationMinute: event.GetNotificationMinute().AsDuration(),
	}
}

func eventToMessage(events []unityres.Event) *pb.ListEventResponse {
	resEvents := make([]*pb.EventProto, 0, len(events))

	for _, event := range events {
		resEvents = append(resEvents, &pb.EventProto{
			Id:                 event.ID,
			Title:              event.Title,
			Date:               timestamppb.New(event.Date),
			Duration:           durationpb.New(event.Duration),
			Description:        event.Description,
			UserId:             int64(event.UserID),
			NotificationMinute: durationpb.New(event.NotificationMinute),
		})
	}
	return &pb.ListEventResponse{Events: resEvents}
}

type ServerGRPC struct {
	Conf        configs.GRPCConf
	application *app.App
	service     *ServiceGRPC
	server      *grpc.Server
}

func NewGRPCServer(
	app *app.App, conf configs.GRPCConf,
) *ServerGRPC {
	server := grpc.NewServer(grpc.UnaryInterceptor(UnaryServerLogRequestInterceptor(app.Logger)))
	return &ServerGRPC{
		Conf:        conf,
		application: app,
		service:     NewServiceGRPC(app),
		server:      server,
	}
}

func (s *ServerGRPC) Start() error {
	lsn, err := net.Listen("tcp", s.Conf.Addr)
	if err != nil {
		return err
	}

	pb.RegisterCalendarServer(s.server, s.service)
	reflection.Register(s.server)
	s.application.Logger.Warn(fmt.Sprintf("GRPC server is started on: %s", s.Conf.Addr))
	go func() {
		if err := s.server.Serve(lsn); err != nil {
			s.application.Logger.Error(fmt.Sprintf("GRPC server start error: %s", err))
		}
	}()

	return nil
}

func (s *ServerGRPC) Stop() error {
	s.application.Logger.Warn("GRPC server has been stopped")
	s.server.Stop()
	return nil
}
