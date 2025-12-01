package search_indexer

import (
	"go.uber.org/fx"
)

func Build() fx.Option {
	return fx.Options(
		fx.Provide(newIndexer),
	)
}
