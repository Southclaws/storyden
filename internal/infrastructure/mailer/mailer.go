package mailer

import (
	"context"

	"go.uber.org/fx"
	"go.uber.org/zap"

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

func newMailer(l *zap.Logger, cfg config.Config) (Sender, error) {
	switch cfg.EmailProvider {
	case "none":
		l.Info("initialising with no mailer")
		return nil, nil

	case "sendgrid":
		l.Info("initialising sendgrid mailer")
		return newSendgridMailer(l, cfg)

	case "mock":
		l.Info("initialising mock mailer")
		return &Mock{}, nil

	default:
		return nil, fault.Newf("unknown email provider: '%s'", cfg.EmailProvider)
	}
}
