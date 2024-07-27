package mailer

import (
	"context"
	"fmt"
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
		return newSendgridMailer(l)

	case "":
		return &Mock{}, nil

	default:
		panic(fmt.Sprintf("unknown email provider: %s", cfg.EmailProvider))
	}
}
