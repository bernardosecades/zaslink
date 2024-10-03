package events

import "context"

type Publisher[T any] interface {
	Publish(ctx context.Context, event Event[T]) error
}
