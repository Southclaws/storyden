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
	if cfg.EmailProvider != "" && len(cfg.JWTSecret) == 0 {
		return nil, fault.New("JWT secret must be provided when enabling email features, set JWT_SECRET in the environment")
	}

	switch cfg.EmailProvider {
	case "":
		return nil, nil

	case "sendgrid":
		return newSendgridMailer(logger, cfg)

	case "mock":
		return &Mock{}, nil

	default:
		return nil, fault.Newf("unknown email provider: '%s'", cfg.EmailProvider)
	}
}
