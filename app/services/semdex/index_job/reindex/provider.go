package reindex

import (
	"go.uber.org/fx"
)

func Build() fx.Option {
	return fx.Options(
		fx.Provide(newReindexer),
		fx.Invoke(runReindexer),
	)
}
