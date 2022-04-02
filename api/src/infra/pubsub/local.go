package pubsub

import (
	"errors"

	"github.com/cskr/pubsub"
	"go.uber.org/zap"
)

var _ Bus = &Embedded{}

type Embedded struct {
	ps *pubsub.PubSub
}

func NewEmbedded() *Embedded {
	return &Embedded{pubsub.New(0)}
}

func (b *Embedded) Declare(topic string) Topic {
	return Topic(topic)
}

func (b *Embedded) Publish(topic Topic, message []byte) error {
	b.ps.Pub(message, string(topic))
	return nil
}

func (b *Embedded) Subscribe(topic Topic, handler func([]byte) (bool, error)) error {
	ch := b.ps.Sub(string(topic))

	for m := range ch {
		ack, err := handler(m.([]byte))
		if err != nil {
			zap.L().Error("mock pubsub handler failed", zap.Error(err))
		}

		if !ack {
			zap.L().Error("nack not implemented for embedded pubsub")
		}
	}

	return errors.New("pubsub subscriber hung up")
}
