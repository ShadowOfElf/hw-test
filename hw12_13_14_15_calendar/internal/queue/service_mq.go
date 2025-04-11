package queue

import (
	"context"
	"fmt"

	"github.com/ShadowOfElf/hw_test/hw12_13_14_15_calendar/configs"
	"github.com/ShadowOfElf/hw_test/hw12_13_14_15_calendar/internal/logger"
	"github.com/ShadowOfElf/hw_test/hw12_13_14_15_calendar/internal/storage/unityres"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQInterface interface {
	Start() error
	Stop() error
	SendNotification(ctx context.Context, notification unityres.Notification) error
	GetNotification(ctx context.Context) (chan unityres.Notification, error)
}

type ServiceMQ struct {
	log    logger.LogInterface
	Conn   *amqp.Connection
	Config configs.RabbitMQConf
	ChanMQ *amqp.Channel
	Queue  *amqp.Queue
}

func NewServiceMQ(log logger.LogInterface, conf configs.RabbitMQConf) RabbitMQInterface {
	return &ServiceMQ{
		log:    log,
		Config: conf,
	}
}

func (s *ServiceMQ) Start() error {
	conn, err := amqp.Dial(
		fmt.Sprintf(
			"amqp://%s:%s@%s:%s/event_vhost",
			s.Config.User,
			s.Config.Password,
			s.Config.Host,
			s.Config.Port,
		),
	)
	if err != nil {
		return err
	}
	s.Conn = conn
	s.log.Info(fmt.Sprintf("connected to RabbitMQ at %s, %s", s.Config.Host, s.Config.Port))

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	s.ChanMQ = ch

	q, err := ch.QueueDeclare(
		"notif_queue",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	s.Queue = &q

	return nil
}

func (s *ServiceMQ) Stop() error {
	s.log.Info("close connection RabbitMQ")
	err := s.ChanMQ.Close()
	if err != nil {
		return err
	}
	err = s.Conn.Close()
	if err != nil {
		return err
	}
	return nil
}
