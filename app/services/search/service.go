package search

import (
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/services/search/searcher"
	"github.com/Southclaws/storyden/app/services/search/simplesearch"
	"github.com/Southclaws/storyden/app/services/semdex"
	"github.com/Southclaws/storyden/internal/config"
)

func New(
	cfg config.Config,
	simpleSearcher *simplesearch.ParallelSearcher,
	semdexSearcher semdex.Searcher,
) searcher.Searcher {
	switch cfg.SemdexProvider {
	case "chromem", "weaviate", "pinecone":
		return semdexSearcher

	default:
		return simpleSearcher
	}
}

func Build() fx.Option {
	return fx.Options(
		fx.Provide(
			New,
			simplesearch.NewParallelSearcher,
		),
	)
}
