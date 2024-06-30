package pubsub

import (
	"context"
)

type Message[T any] struct {
	ID      string
	Payload T
	Ack     func() bool
	Nack    func() bool
}

type Topic[T any] interface {
	Subscriber[T]
	Publisher[T]
}

type Subscriber[T any] interface {
	Subscribe(ctx context.Context) (<-chan *Message[T], error)
}

type Publisher[T any] interface {
	Publish(ctx context.Context, messages ...Message[T]) error
}
