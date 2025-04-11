package grpcclient

import (
	"time"

	"github.com/ShadowOfElf/hw_test/hw12_13_14_15_calendar/configs"
	"github.com/ShadowOfElf/hw_test/hw12_13_14_15_calendar/internal/logger"
	pb "github.com/ShadowOfElf/hw_test/hw12_13_14_15_calendar/pkg"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ServiceGRPCClient struct {
	log  logger.LogInterface
	conf configs.GRPCConf
}

func NewGRPClient(log logger.LogInterface, conf configs.GRPCConf) *ServiceGRPCClient {
	return &ServiceGRPCClient{
		log:  log,
		conf: conf,
	}
}

func (s *ServiceGRPCClient) Start() (pb.CalendarClient, error) {
	conn, err := grpc.NewClient(s.conf.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	client := pb.NewCalendarClient(conn)
	return client, nil
}

func (s *ServiceGRPCClient) GetReqByDay(timeDay time.Time) *pb.ListEventByDateRequest {
	day := time.Date(timeDay.Year(), timeDay.Month(), timeDay.Day(), 0, 0, 0, 0, time.UTC)
	req := &pb.ListEventByDateRequest{Data: timestamppb.New(day)}
	return req
}
