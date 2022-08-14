package magiclink

import (
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/backend/pkg/resources/account"
)

type Magiclink interface {
	Send(email string) (*account.Account, error)
	Callback(token []byte) (*account.Account, error)
}

func Build() fx.Option {
	return fx.Provide(
		// TODO: Switch between regular email and Magic.link based.
		NewEmail,
	)
}
