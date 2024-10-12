package events

import "context"

type InMemoryPublisher struct{}

func NewInMemoryPublisher() *InMemoryPublisher {
	return &InMemoryPublisher{}
}

func (p *InMemoryPublisher) Publish(_ context.Context, _ Event[map[string]string]) error {
	return nil
}
