package sms

import (
	"context"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/storyden/internal/config"
)

type Sender interface {
	Send(ctx context.Context, phone string, message string) error
}

func Build() fx.Option {
	return fx.Provide(func(cfg config.Config, l *zap.Logger) (Sender, error) {
		switch cfg.SMSProvider {
		case "none":
			l.Info("initialising with no SMS provider")
			return nil, nil

		case "twilio":
			l.Info("initialising Twilio SMS provider")
			return newTwilio(l, cfg)

		case "mock":
			return newMock(l)

		default:
			return nil, fault.Newf("unknown SMS provider: '%s'", cfg.SMSProvider)
		}
	})
}
