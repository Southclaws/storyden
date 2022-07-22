package magiclink

import (
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/backend/pkg/resources/user"
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
