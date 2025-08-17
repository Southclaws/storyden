package pubsub

import (
	"context"
	"log/slog"

	"go.uber.org/fx"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-amqp/v3/pkg/amqp"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"

	"github.com/Southclaws/storyden/internal/config"
)

func Build() fx.Option {
	return fx.Options(
		fx.Provide(func(
			lc fx.Lifecycle,
			ctx context.Context,
			cfg config.Config,
			l *slog.Logger,
		) (*Bus, error) {
			sub, pub, err := newWatermillPubsub(cfg, l)
			if err != nil {
				return nil, err
			}

			bus, err := newBus(lc, l, ctx, cfg, pub, sub)
			if err != nil {
				return nil, err
			}

			return bus, nil
		}),
	)
}

func newWatermillPubsub(cfg config.Config, l *slog.Logger) (message.Subscriber, message.Publisher, error) {
	logger := watermill.NewSlogLogger(l)

	switch cfg.QueueType {
	default:
		l.Debug("using channel queue")

		pubsub := gochannel.NewGoChannel(
			gochannel.Config{},
			logger,
		)

		return pubsub, pubsub, nil

	case "amqp":
		l.Debug("using amqp pubsub")

		apsc := amqp.NewDurablePubSubConfig(cfg.AmqpURL, amqp.GenerateExchangeNameTopicName)

		publisher, err := amqp.NewPublisher(apsc, logger)
		if err != nil {
			return nil, nil, err
		}

		subscriber, err := amqp.NewSubscriber(apsc, logger)
		if err != nil {
			return nil, nil, err
		}

		return subscriber, publisher, nil
	}
}
