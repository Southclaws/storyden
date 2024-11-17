package mailer

import (
	"context"
	"net/mail"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/Southclaws/storyden/internal/config"
)

type Sender interface {
	Send(
		ctx context.Context,
		address mail.Address,
		name string,
		subject string,
		html string,
		plain string,
	) error
}

func Build() fx.Option {
	return fx.Options(
		fx.Provide(newMailer),
	)
}

func newMailer(l *zap.Logger, cfg config.Config) (Sender, error) {
	switch cfg.EmailProvider {
	case "sendgrid":
		l.Info("initialising sendgrid mailer")
		return newSendgridMailer(l)

	case "mock":
		l.Info("initialising mock mailer")
		return &Mock{}, nil

	default:
		l.Info("initialising with no mailer")
		return nil, nil
	}
}
