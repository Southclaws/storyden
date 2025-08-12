package watermill

import (
	"log/slog"

	"github.com/ThreeDotsLabs/watermill-amqp/v3/pkg/amqp"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/google/uuid"

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
		l.Debug("using amqp pubsub")

		// Use PubSub config for fan-out behavior (not Queue config which creates competing consumers)
		// Each subscriber gets a unique queue name for fan-out
		// Use a random suffix to ensure unique queue names per subscriber
		instanceSuffix := uuid.New().String()[:8] // Short unique suffix
		apsc := amqp.NewDurablePubSubConfig(cfg.AmqpURL, amqp.GenerateQueueNameTopicNameWithSuffix(instanceSuffix))

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
