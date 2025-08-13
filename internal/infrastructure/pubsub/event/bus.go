package event

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-amqp/v3/pkg/amqp"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub/queuename"
)

type Bus struct {
	logger           *slog.Logger
	cfg              config.Config
	pub              message.Publisher
	sub              message.Subscriber
	router           *message.Router
	eventBus         *cqrs.EventBus
	commandBus       *cqrs.CommandBus
	eventProcessor   *cqrs.EventProcessor
	commandProcessor *cqrs.CommandProcessor

	mu            sync.RWMutex
	subscriptions map[string]*Subscription
}

func New(
	lc fx.Lifecycle,
	l *slog.Logger,
	ctx context.Context,
	cfg config.Config,
	pub message.Publisher,
	sub message.Subscriber,
	eventTypes ...any,
) (*Bus, error) {
	logger := watermill.NewSlogLogger(l.With("component", "watermill"))

	router, err := message.NewRouter(message.RouterConfig{
		CloseTimeout: time.Second * 30,
	}, logger)
	if err != nil {
		return nil, fault.Wrap(err)
	}

	router.AddMiddleware(middleware.Recoverer)

	marshaler := cqrs.JSONMarshaler{
		GenerateName: func(v interface{}) string {
			return queuename.FromValue(v)
		},
	}

	eventBus, err := cqrs.NewEventBusWithConfig(pub, cqrs.EventBusConfig{
		GeneratePublishTopic: func(params cqrs.GenerateEventPublishTopicParams) (string, error) {
			return params.EventName, nil
		},
		Marshaler: marshaler,
	})
	if err != nil {
		return nil, fault.Wrap(err)
	}

	commandBus, err := cqrs.NewCommandBusWithConfig(pub, cqrs.CommandBusConfig{
		GeneratePublishTopic: func(params cqrs.CommandBusGeneratePublishTopicParams) (string, error) {
			return params.CommandName, nil
		},
		Marshaler: marshaler,
	})
	if err != nil {
		return nil, fault.Wrap(err)
	}

	eventProcessor, err := cqrs.NewEventProcessorWithConfig(router, cqrs.EventProcessorConfig{
		GenerateSubscribeTopic: func(params cqrs.EventProcessorGenerateSubscribeTopicParams) (string, error) {
			return params.EventName, nil
		},
		SubscriberConstructor: func(params cqrs.EventProcessorSubscriberConstructorParams) (message.Subscriber, error) {
			// NOTE: When we're using AMQP as the broker, because the fanout
			// logic requires separate subscribers per queue, we need to do
			// some additional setup to ensure that the subscriber is unique
			// to each consumer. This is how we get events properly fanned out
			// to subscribers of the same event. Internally this creates a new
			// subscriber for each event+service key and AMQP handles delivery.
			if cfg.QueueType == "amqp" {
				apsc := amqp.NewDurablePubSubConfig(cfg.AmqpURL, amqp.GenerateQueueNameTopicNameWithSuffix(params.HandlerName))
				subscriber, err := amqp.NewSubscriber(apsc, logger)
				if err != nil {
					return nil, err
				}
				return subscriber, nil
			}

			return sub, nil
		},
		Marshaler: marshaler,
	})
	if err != nil {
		return nil, fault.Wrap(err)
	}

	commandProcessor, err := cqrs.NewCommandProcessorWithConfig(router, cqrs.CommandProcessorConfig{
		GenerateSubscribeTopic: func(params cqrs.CommandProcessorGenerateSubscribeTopicParams) (string, error) {
			return params.CommandName, nil
		},
		SubscriberConstructor: func(params cqrs.CommandProcessorSubscriberConstructorParams) (message.Subscriber, error) {
			return sub, nil
		},
		Marshaler: marshaler,
	})
	if err != nil {
		return nil, fault.Wrap(err)
	}

	lc.Append(fx.StartHook(func() {
		go func() {
			for {
				select {
				case <-ctx.Done():
					return

				default:
					err := router.Run(ctx)
					if err != nil {
						l.Error("message router stopped unexpectedly",
							slog.String("error", err.Error()),
						)
					}

					l.Warn("restarting message router in 5 seconds")

					time.Sleep(5 * time.Second)
				}
			}
		}()
	}))

	lc.Append(fx.StopHook(func(ctx context.Context) error {
		if err := router.Close(); err != nil {
			return err
		}

		l.Info("message router stopped successfully")
		return nil
	}))

	return &Bus{
		logger:           l,
		cfg:              cfg,
		pub:              pub,
		sub:              sub,
		router:           router,
		eventBus:         eventBus,
		commandBus:       commandBus,
		eventProcessor:   eventProcessor,
		commandProcessor: commandProcessor,
		subscriptions:    make(map[string]*Subscription),
	}, nil
}

func (b *Bus) Publish(ctx context.Context, events ...any) error {
	var errs []error

	for _, event := range events {
		if err := b.eventBus.Publish(ctx, event); err != nil {
			b.logger.Error("failed to publish event",
				slog.String("event_type", fmt.Sprintf("%T", event)),
				slog.String("error", err.Error()),
			)
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}

func (b *Bus) SendCommand(ctx context.Context, command any) error {
	if err := b.commandBus.Send(ctx, command); err != nil {
		b.logger.Error("failed to send command",
			slog.String("command_type", fmt.Sprintf("%T", command)),
			slog.String("error", err.Error()),
		)
		return fault.Wrap(err, fctx.With(ctx))
	}
	return nil
}

type Subscription struct {
	bus            *Bus
	handlerID      string
	topic          string
	closed         bool
	mu             sync.Mutex
	messageHandler *message.Handler
}

func (s *Subscription) Close() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return
	}

	s.messageHandler.Stop()
	s.closed = true

	s.bus.mu.Lock()
	delete(s.bus.subscriptions, s.handlerID)
	s.bus.mu.Unlock()
}

type (
	HandlerFunc[T any]        func(ctx context.Context, event *T) error
	CommandHandlerFunc[T any] func(ctx context.Context, command *T) error
)

func Subscribe[T any](ctx context.Context, bus *Bus, handlerName string, handler HandlerFunc[T]) (*Subscription, error) {
	var zero T
	topic := queuename.FromValue(zero)
	handlerID := fmt.Sprintf("%s_%s", topic, handlerName)

	bus.mu.Lock()
	defer bus.mu.Unlock()
	if _, exists := bus.subscriptions[handlerID]; exists {
		return nil, fmt.Errorf("subscription already exists: %s", handlerID)
	}

	cqrsHandler := cqrs.NewEventHandler(handlerName, func(ctx context.Context, event *T) error {
		return handler(ctx, event)
	})

	messageHandler, err := bus.eventProcessor.AddHandler(cqrsHandler)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	err = bus.router.RunHandlers(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	sub := &Subscription{
		bus:            bus,
		handlerID:      handlerID,
		topic:          topic,
		messageHandler: messageHandler,
	}

	bus.subscriptions[handlerID] = sub

	return sub, nil
}

func SubscribeCommand[T any](ctx context.Context, bus *Bus, handlerName string, handler CommandHandlerFunc[T]) (*Subscription, error) {
	var zero T
	topic := queuename.FromValue(zero)
	handlerID := fmt.Sprintf("%s_%s", topic, handlerName)

	bus.mu.Lock()
	defer bus.mu.Unlock()
	if _, exists := bus.subscriptions[handlerID]; exists {
		return nil, fmt.Errorf("subscription already exists: %s", handlerID)
	}

	cqrsHandler := cqrs.NewCommandHandler(handlerName, func(ctx context.Context, command *T) error {
		return handler(ctx, command)
	})

	if err := bus.commandProcessor.AddHandlers(cqrsHandler); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	sub := &Subscription{
		bus:       bus,
		handlerID: handlerID,
		topic:     topic,
	}

	bus.subscriptions[handlerID] = sub

	return sub, nil
}
