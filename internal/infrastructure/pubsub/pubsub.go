package pubsub

import (
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/internal/config"
)

type Topic string

type Bus interface {
	Declare(topic string) Topic
	Publish(topic Topic, message []byte) error
	Subscribe(topic Topic, handler func([]byte) (bool, error)) error
}

func Build() fx.Option {
	return fx.Provide(func(cfg config.Config) (Bus, error) {
		if cfg.Production {
			return NewRabbit(cfg)
		} else {
			return NewEmbedded(), nil
		}
	})
}
