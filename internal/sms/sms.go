package sms

import (
	"context"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/Southclaws/storyden/internal/config"
)

type Sender interface {
	Send(ctx context.Context, phone string, message string) error
}

func Build() fx.Option {
	return fx.Provide(func(cfg config.Config, l *zap.Logger) (Sender, error) {
		if !cfg.Production { // CHANGE
			return newTwilio(l)
		}

		return newMock(l)
	})
}
