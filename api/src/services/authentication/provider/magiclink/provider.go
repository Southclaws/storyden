package magiclink

import (
	"github.com/Southclaws/storyden/api/src/resources/user"
	"go.uber.org/fx"
)

type Magiclink interface {
	Send(email string) (*user.User, error)
	Callback(token []byte) (*user.User, error)
}

func Build() fx.Option {
	return fx.Provide(
		// TODO: Switch between regular email and Magic.link based.
		NewEmail,
	)
}
