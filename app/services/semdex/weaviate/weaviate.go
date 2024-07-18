package weaviate

import (
	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/Southclaws/storyden/app/services/semdex"
	"github.com/Southclaws/storyden/app/services/semdex/result_hydrator"
	"github.com/Southclaws/storyden/app/services/semdex/simplesearch"
	weaviate_client "github.com/Southclaws/storyden/internal/weaviate"
)

type weaviateSemdexer struct {
	wc *weaviate.Client
	cn weaviate_client.WeaviateClassName
}

func newWeaviateSemdexer(wc *weaviate.Client, cn weaviate_client.WeaviateClassName) semdex.Semdexer {
	return &weaviateSemdexer{wc, cn}
}

func newSemdexer(
	l *zap.Logger,
	wc *weaviate.Client,
	cn weaviate_client.WeaviateClassName,
	rh *result_hydrator.Hydrator,
	ss *simplesearch.ParallelSearcher,
) semdex.Semdexer {
	if wc == nil {
		return &semdex.OnlySearcher{Searcher: ss}
	}

	return &withHydration{l, newWeaviateSemdexer(wc, cn), rh}
}

func Build() fx.Option {
	return fx.Provide(
		result_hydrator.New,
		simplesearch.NewParallelSearcher,
		fx.Annotate(
			newSemdexer,
			fx.As(new(semdex.Semdexer)),
			fx.As(new(semdex.Indexer)),
			fx.As(new(semdex.Searcher)),
			fx.As(new(semdex.Recommender)),
			fx.As(new(semdex.Retriever)),
			fx.As(new(semdex.RelevanceScorer)),
		),
	)
}
