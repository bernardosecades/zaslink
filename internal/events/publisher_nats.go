package events

import (
	"context"
	"encoding/json"

	"github.com/nats-io/nats.go"

	"github.com/bernardosecades/zaslink/pkg/events"
	"github.com/bernardosecades/zaslink/pkg/masking"
)

type NatsPublisher struct {
	url string
}

func NewNatsPublisher(url string) *NatsPublisher {
	return &NatsPublisher{url: url}
}

func (m *NatsPublisher) Publish(_ context.Context, event events.Event[map[string]string]) error {
	nc, err := nats.Connect(m.url)
	if err != nil {
		return err
	}
	defer nc.Close()

	subject := "notifications.telegram.zaslink_service"

	event.Data["eventName"] = event.Name
	event.Data["id"] = masking.MaskString(event.Data["id"])
	event.Data["privateId"] = masking.MaskString(event.Data["id"])
	message, err := json.MarshalIndent(event.Data, "", "    ")
	if err != nil {
		return err
	}

	// Publish the message
	err = nc.Publish(subject, message)
	if err != nil {
		return err
	}

	return nil
}
