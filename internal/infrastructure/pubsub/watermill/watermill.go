package watermill

import (
	"log/slog"

	"github.com/ThreeDotsLabs/watermill-amqp/v3/pkg/amqp"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"

	"github.com/Southclaws/storyden/internal/config"
)

func NewWatermillQueue(cfg config.Config, l *slog.Logger) (message.Subscriber, message.Publisher, error) {
	logger := &logAdapter{l}

	switch cfg.QueueType {
	default:
		l.Debug("using channel queue")

		pubsub := gochannel.NewGoChannel(
			gochannel.Config{},
			logger,
		)

		return pubsub, pubsub, nil

	case "amqp":
		l.Debug("using amqp queue")

		aqc := amqp.NewDurableQueueConfig(cfg.AmqpURL)

		publisher, err := amqp.NewPublisher(aqc, logger)
		if err != nil {
			return nil, nil, err
		}

		subscriber, err := amqp.NewSubscriber(aqc, logger)
		if err != nil {
			return nil, nil, err
		}

		return subscriber, publisher, nil
	}
}
