package queue

import (
	"context"
	"encoding/json"
	"reflect"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/opt"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/rs/xid"
	"go.uber.org/zap"

	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
)

const actorIDMetadataKey = "actor_id"

type QueueFactory struct {
	log *zap.Logger
	pub message.Publisher
	sub message.Subscriber
}

func New[T any](q *QueueFactory) pubsub.Topic[T] {
	topic := typename[T]()

	logger := q.log.With(zap.String("topic", topic))

	logger.Debug("registered new queue")

	return &watermillQueue[T]{
		logger,
		topic,
		q.pub,
		q.sub,
	}
}

type watermillQueue[T any] struct {
	log   *zap.Logger
	topic string
	pub   message.Publisher
	sub   message.Subscriber
}

func (q *watermillQueue[T]) Subscribe(ctx context.Context) (<-chan *pubsub.Message[T], error) {
	ch, subscribeErr := q.sub.Subscribe(ctx, q.topic)
	if subscribeErr != nil {
		return nil, fault.Wrap(subscribeErr, fctx.With(ctx))
	}

	recv := make(chan *pubsub.Message[T], 100)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return

			case msg := <-ch:
				if msg == nil {
					q.log.Warn("nil message received by subscriber")
					continue
				}

				var payload T
				if err := json.Unmarshal(msg.Payload, &payload); err != nil {
					q.log.Error("failed to decode message payload",
						zap.Error(err))

					msg.Nack()

					continue
				}

				actorID, err := getActorID(msg)
				if err != nil {
					q.log.Error("failed to get actor ID from message metadata",
						zap.Error(err))
				}

				recv <- &pubsub.Message[T]{
					ID:      msg.UUID,
					Payload: payload,
					Ack:     msg.Ack,
					Nack:    msg.Nack,
					ActorID: actorID,
				}
			}
		}
	}()

	return recv, nil
}

func (q *watermillQueue[T]) Publish(ctx context.Context, payloads ...T) error {
	// If the publish was acted by a session account, store in the payload.
	actorID := session.GetOptAccountID(ctx)

	messages, err := dt.MapErr(payloads, func(p T) (*message.Message, error) {
		payload, err := json.Marshal(p)
		if err != nil {
			return nil, err
		}

		msg := message.NewMessage(watermill.NewUUID(), payload)

		msg.Metadata.Set(actorIDMetadataKey, actorID.String())

		return msg, nil
	})
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	err = q.pub.Publish(q.topic, messages...)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}

// so we don't need to be manually specifying topic names, derive from name of T
func typename[T any]() string {
	var zero [0]T
	to := reflect.TypeOf(zero).Elem()
	return to.String()
}

func getActorID(msg *message.Message) (opt.Optional[xid.ID], error) {
	raw := msg.Metadata.Get(actorIDMetadataKey)
	if raw == "" {
		return opt.NewEmpty[xid.ID](), nil
	}

	actorID, err := xid.FromString(raw)
	if err != nil {
		return nil, fault.Wrap(err, fmsg.With("malformed actor ID in message metadata"))
	}

	return opt.New[xid.ID](actorID), nil
}
