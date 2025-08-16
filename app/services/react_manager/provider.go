package react_manager

import (
	"go.uber.org/fx"
)

func Build() fx.Option {
	return fx.Options(
		fx.Provide(New),
	)
}
