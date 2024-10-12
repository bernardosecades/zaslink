package events

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/nats-io/nats.go"

	"github.com/bernardosecades/zaslink/pkg/events"
)

type NatsPublisher struct {
}

func NewNatsPublisher() *NatsPublisher {
	return &NatsPublisher{}
}

func (m *NatsPublisher) Publish(_ context.Context, event events.Event[map[string]string]) error {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer nc.Close()

	subject := "notifications.telegram.zaslink_service"

	event.Data["eventName"] = event.Name
	message, err := json.MarshalIndent(event.Data, "", "    ")
	if err != nil {
		fmt.Println(err)
		return err
	}

	// Publish the message
	err = nc.Publish(subject, message)
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println("Published message to ", string(message))

	return nil
}
