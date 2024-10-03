package events

import (
	"context"

	"github.com/bernardosecades/sharesecret/pkg/events"
)

type DummyPublisher struct {
}

func NewDummyPublisher() *DummyPublisher {
	return &DummyPublisher{}
}

func (m *DummyPublisher) Publish(_ context.Context, _ events.Event[map[string]string]) error {
	return nil
}
