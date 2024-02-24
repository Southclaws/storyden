package weaviate

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate/entities/models"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/Southclaws/storyden/app/services/semdex"
	"github.com/Southclaws/storyden/app/services/semdex/result_hydrator"
	"github.com/Southclaws/storyden/app/services/semdex/simplesearch"
)

// NOT PROD READY: Just using local transformers for now.

type weaviateSemdexer struct {
	wc *weaviate.Client
	mc models.Class
}

func newWeaviateSemdexer(lc fx.Lifecycle, wc *weaviate.Client) semdex.Semdexer {
	class := models.Class{
		Class:      "ContentText2vecTransformers",
		Vectorizer: "text2vec-transformers",
		ModuleConfig: map[string]interface{}{
			// "text2vec-openai":   map[string]interface{}{},
			// "generative-openai": map[string]interface{}{},
			"text2vec-transformers": map[string]interface{}{},
		},
	}

	lc.Append(fx.StartHook(func(ctx context.Context) error {
		r, err := wc.Schema().
			ClassExistenceChecker().
			WithClassName(class.Class).
			Do(ctx)
		if err != nil {
			return fault.Wrap(err)
		}

		if !r {
			err := wc.Schema().
				ClassCreator().
				WithClass(&class).
				Do(ctx)
			if err != nil {
				return fault.Wrap(err)
			}
		}

		return nil
	}))

	return &weaviateSemdexer{wc, class}
}

func Build() fx.Option {
	return fx.Provide(
		result_hydrator.New,
		simplesearch.NewParallelSearcher,
		fx.Annotate(
			func(lc fx.Lifecycle, l *zap.Logger, wc *weaviate.Client, rh *result_hydrator.Hydrator, ss *simplesearch.ParallelSearcher) semdex.Semdexer {
				if wc == nil {
					return &semdex.OnlySearcher{ss}
				}

				return &withHydration{l, newWeaviateSemdexer(lc, wc), rh}
			},
			fx.As(new(semdex.Semdexer)),
			fx.As(new(semdex.Indexer)),
			fx.As(new(semdex.Searcher)),
			fx.As(new(semdex.Recommender)),
		),
	)
}
