package pubsub

import (
	"context"

	"github.com/Southclaws/opt"
	"github.com/rs/xid"
)

type Message[T any] struct {
	ID      string
	Payload T
	Ack     func() bool
	Nack    func() bool
	ActorID opt.Optional[xid.ID]
}

type Topic[T any] interface {
	Subscriber[T]
	Publisher[T]
}

type Subscriber[T any] interface {
	Subscribe(ctx context.Context) (<-chan *Message[T], error)
}

type Publisher[T any] interface {
	Publish(ctx context.Context, messages ...T) error
	PublishAndForget(ctx context.Context, messages ...T)
}
