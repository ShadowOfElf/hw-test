package queue

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ShadowOfElf/hw_test/hw12_13_14_15_calendar/internal/storage/unityres"
)

func (s *ServiceMQ) GetNotification(ctx context.Context) (chan unityres.Notification, error) {
	messages, err := s.ChanMQ.Consume(
		s.Queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}
	notification := make(chan unityres.Notification)

	go func() {
		defer close(notification)
		for {
			select {
			case <-ctx.Done():
				return
			case message := <-messages:
				body := message.Body
				var msg unityres.Notification
				err := json.Unmarshal(body, &msg)
				if err != nil {
					s.log.Error(fmt.Sprintf("failed parse message Rabbit: %s", err))
					return
				}
				notification <- msg
			}
		}
	}()
	return notification, err
}
