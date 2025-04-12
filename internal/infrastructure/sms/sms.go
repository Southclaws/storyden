package sms

import (
	"context"
	"log/slog"

	"go.uber.org/fx"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/storyden/internal/config"
)

type Sender interface {
	Send(ctx context.Context, phone string, message string) error
}

func Build() fx.Option {
	return fx.Provide(func(cfg config.Config, l *slog.Logger) (Sender, error) {
		switch cfg.SMSProvider {
		case "none":
			return nil, nil

		case "twilio":
			return newTwilio(l, cfg)

		case "mock":
			return newMock(l)

		default:
			return nil, fault.Newf("unknown SMS provider: '%s'", cfg.SMSProvider)
		}
	})
}
