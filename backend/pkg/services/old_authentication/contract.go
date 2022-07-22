package authentication

import (
	"go.uber.org/fx"
)

func Build() fx.Option {
	return fx.Provide(NewCookieAuth)
}
