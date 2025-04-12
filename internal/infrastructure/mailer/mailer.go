package mailer

import (
	"context"
	"log/slog"

	"go.uber.org/fx"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/storyden/internal/config"
)

type Sender interface {
	Send(ctx context.Context, msg Message) error
}

func Build() fx.Option {
	return fx.Options(
		fx.Provide(newMailer),
	)
}

func newMailer(logger *slog.Logger, cfg config.Config) (Sender, error) {
	switch cfg.EmailProvider {
	case "none":
		return nil, nil

	case "sendgrid":
		return newSendgridMailer(logger, cfg)

	case "mock":
		return &Mock{}, nil

	default:
		return nil, fault.Newf("unknown email provider: '%s'", cfg.EmailProvider)
	}
}
