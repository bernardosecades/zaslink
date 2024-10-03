package events

type Event[T any] struct {
	Name string
	Data T
}
