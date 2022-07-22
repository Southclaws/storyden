package mailer

import (
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/backend/internal/config"
)

type Mailer interface {
	Mail(toname, toaddr, subj, rich, text string) error
}

func Build() fx.Option {
	return fx.Provide(func(cfg config.Config) (Mailer, error) {
		if cfg.Production {
			return NewSendGrid()
		} else {
			return &Mock{}, nil
		}
	})
}
