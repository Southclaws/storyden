package semdexer

import (
	"context"

	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/datagraph/hydrate"
	"github.com/Southclaws/storyden/app/resources/robot/llm_provider"
	"github.com/Southclaws/storyden/app/services/semdex"
	"github.com/Southclaws/storyden/app/services/semdex/semdexer/chromem_semdexer"
	"github.com/Southclaws/storyden/app/services/semdex/semdexer/pinecone_semdexer"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/infrastructure/vector/pinecone"
)

func newSemdexer(
	ctx context.Context,
	cfg config.Config,
	pc *pinecone.Client,
	modelFactory *llm_provider.Factory,
	hydrator *hydrate.Hydrator,
) (semdex.Semdexer, error) {
	embedder := resolveEmbedder(modelFactory)

	switch cfg.SemdexProvider {
	case "chromem":
		return chromem_semdexer.New(cfg, hydrator, embedder)

	case "pinecone":
		return pinecone_semdexer.New(ctx, cfg, pc, hydrator, embedder)

	default:
		return &semdex.Disabled{}, nil
	}
}

func resolveEmbedder(modelFactory *llm_provider.Factory) semdex.Embedder {
	return func(ctx context.Context, text string) ([]float32, error) {
		embedder, err := modelFactory.GetEmbedder(ctx)
		if err != nil {
			return nil, err
		}

		return embedder(ctx, text)
	}
}

func Build() fx.Option {
	return fx.Options(
		fx.Provide(
			fx.Annotate(
				newSemdexer,
				fx.As(new(semdex.Semdexer)),
				fx.As(new(semdex.Querier)),
				fx.As(new(semdex.Mutator)),
				fx.As(new(semdex.Recommender)),
				fx.As(new(semdex.Searcher)),
			),
		),
	)
}
