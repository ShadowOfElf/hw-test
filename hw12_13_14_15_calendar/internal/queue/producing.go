package queue

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ShadowOfElf/hw_test/hw12_13_14_15_calendar/internal/storage/unityres"
	amqp "github.com/rabbitmq/amqp091-go"
)

func (s *ServiceMQ) SendNotification(ctx context.Context, notification unityres.Notification) error {
	body, err := json.Marshal(notification)
	if err != nil {
		return err
	}

	err = s.ChanMQ.PublishWithContext(
		ctx,
		"",
		s.Queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		return err
	}
	s.log.Debug(fmt.Sprintf("Send message %v", notification))
	return nil
}
