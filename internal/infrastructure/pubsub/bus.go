package pubsub

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"reflect"
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
	subscriptions map[subscriptionKey]*Subscription
}

type subscriptionKey string

func newBus(
	lc fx.Lifecycle,
	l *slog.Logger,
	ctx context.Context,
	cfg config.Config,
	pub message.Publisher,
	sub message.Subscriber,
) (*Bus, error) {
	logger := watermill.NewSlogLogger(l.With("component", "watermill"))

	router, err := message.NewRouter(message.RouterConfig{
		CloseTimeout: time.Second * 30,
	}, nil)
	if err != nil {
		return nil, fault.Wrap(err)
	}

	router.AddMiddleware(middleware.Recoverer)
	router.AddMiddleware(newSessionContextMiddleware(l))
	router.AddMiddleware(newChaosDelayMiddleware(cfg.DevChaosSlowModeQueue, l))

	poisonQueue, err := middleware.PoisonQueue(pub, "poison_queue")
	if err != nil {
		return nil, fault.Wrap(err)
	}
	router.AddMiddleware(poisonQueue)

	retryMiddleware := middleware.Retry{
		MaxRetries:      cfg.QueueMaxRetries,
		InitialInterval: cfg.QueueRetryInitialInterval,
		MaxInterval:     cfg.QueueRetryMaxInterval,
		Multiplier:      2.0,
		OnRetryHook: func(retryNum int, delay time.Duration) {
			log := fmt.Sprintf("a message consumer returned an error: retrying %d/%d after %s/%s",
				retryNum, cfg.QueueMaxRetries, delay, cfg.QueueRetryMaxInterval,
			)
			l.Error(log,
				slog.Int("retry_num", retryNum),
				slog.String("delay", delay.String()),
			)
		},
	}
	router.AddMiddleware(retryMiddleware.Middleware)

	marshaler := cqrs.JSONMarshaler{
		GenerateName: func(v interface{}) string {
			return topicFromValue(v)
		},
	}

	// Wrap publisher with session context middleware
	contextPub := publisherContextMiddleware(pub)

	eventBus, err := cqrs.NewEventBusWithConfig(contextPub, cqrs.EventBusConfig{
		GeneratePublishTopic: func(params cqrs.GenerateEventPublishTopicParams) (string, error) {
			return params.EventName, nil
		},
		Marshaler: marshaler,
	})
	if err != nil {
		return nil, fault.Wrap(err)
	}

	commandBus, err := cqrs.NewCommandBusWithConfig(contextPub, cqrs.CommandBusConfig{
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
			// NOTE 2: This is durable and needs to be optionally durable based
			// on parameters passed to Subscribe. But that's really awkward...
			if cfg.QueueType == "amqp" {
				apsc := amqp.NewDurablePubSubConfig(cfg.AmqpURL, func(topic string) string {
					return topic + "." + params.HandlerName
				})
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

	router.AddNoPublisherHandler("poison_queue_logger", "poison_queue", sub, func(msg *message.Message) error {
		l.Error("poisoned message received after all retries failed",
			slog.String("message_id", msg.UUID),
			slog.String("message_type", msg.Metadata.Get("name")),
			slog.String("reason", msg.Metadata.Get("reason_poisoned")),
			slog.String("original_topic", msg.Metadata.Get("topic_poisoned")),
			slog.String("handler", msg.Metadata.Get("handler_poisoned")),
			slog.String("payload", string(msg.Payload)),
		)
		return nil
	})

	lc.Append(fx.StartHook(func() {
		go func() {
			err := router.Run(ctx)
			if err != nil {
				l.Error("message router stopped unexpectedly",
					slog.String("error", err.Error()),
				)
				os.Exit(0x12)
			}
		}()

		<-router.Running()
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
		subscriptions:    make(map[subscriptionKey]*Subscription),
	}, nil
}

// Publish publishes an event and does not provide error handling semantics.
// Most simple events can use this where a failure to publish isn't critical.
func (b *Bus) Publish(ctx context.Context, event any) {
	if err := b.eventBus.Publish(ctx, event); err != nil {
		b.logger.Error("failed to publish event",
			slog.String("event_type", fmt.Sprintf("%T", event)),
			slog.String("error", err.Error()),
		)
	}
}

func (b *Bus) PublishMany(ctx context.Context, events ...any) {
	for _, e := range events {
		b.Publish(ctx, e)
	}
}

// MustPublish is for when publishing is a critical requirement and errors must
// prevent further procedures, for things like sending emails, etc.
func (b *Bus) MustPublish(ctx context.Context, event any) error {
	if err := b.eventBus.Publish(ctx, event); err != nil {
		b.logger.Error("failed to publish event",
			slog.String("event_type", fmt.Sprintf("%T", event)),
			slog.String("error", err.Error()),
		)

		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}

func (b *Bus) MustPublishMany(ctx context.Context, events ...any) error {
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
	subkey         subscriptionKey
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

	if s.messageHandler != nil {
		s.messageHandler.Stop()
	}

	s.closed = true

	s.bus.mu.Lock()
	delete(s.bus.subscriptions, s.subkey)
	s.bus.mu.Unlock()
}

type (
	HandlerFunc[T any]        func(ctx context.Context, event *T) error
	CommandHandlerFunc[T any] func(ctx context.Context, command *T) error
)

func Subscribe[T any](ctx context.Context, bus *Bus, handlerName string, handler HandlerFunc[T]) (*Subscription, error) {
	topic := topicFromT[T]()
	subkey := subscriptionKey(handlerName)

	bus.mu.Lock()
	defer bus.mu.Unlock()
	if _, exists := bus.subscriptions[subkey]; exists {
		return nil, fmt.Errorf("subscription already exists: %s", subkey)
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
		subkey:         subkey,
		topic:          topic,
		messageHandler: messageHandler,
	}

	bus.subscriptions[subkey] = sub

	return sub, nil
}

func SubscribeCommand[T any](ctx context.Context, bus *Bus, handlerName string, handler CommandHandlerFunc[T]) (*Subscription, error) {
	var zero T
	topic := topicFromValue(zero)
	subkey := subscriptionKey(handlerName)

	bus.mu.Lock()
	defer bus.mu.Unlock()
	if _, exists := bus.subscriptions[subkey]; exists {
		return nil, fmt.Errorf("subscription already exists: %s", subkey)
	}

	cqrsHandler := cqrs.NewCommandHandler(handlerName, func(ctx context.Context, command *T) error {
		return handler(ctx, command)
	})

	if err := bus.commandProcessor.AddHandlers(cqrsHandler); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if err := bus.router.RunHandlers(ctx); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	sub := &Subscription{
		bus:    bus,
		subkey: subkey,
		topic:  topic,
	}

	bus.subscriptions[subkey] = sub

	return sub, nil
}

func topicFromValue(zero any) string {
	t := reflect.TypeOf(zero)

	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	to := t.String()

	return to
}

func topicFromT[T any]() string {
	var zero T
	t := reflect.TypeOf(zero)

	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	to := t.String()

	return to
}
