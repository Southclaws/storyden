package semdexer

import (
	"context"

	"github.com/weaviate/weaviate-go-client/v5/weaviate"
	"go.uber.org/fx"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/storyden/app/resources/datagraph/hydrate"
	"github.com/Southclaws/storyden/app/services/semdex"
	"github.com/Southclaws/storyden/app/services/semdex/asker"
	"github.com/Southclaws/storyden/app/services/semdex/semdexer/chromem_semdexer"
	"github.com/Southclaws/storyden/app/services/semdex/semdexer/pinecone_semdexer"
	"github.com/Southclaws/storyden/app/services/semdex/semdexer/weaviate_semdexer"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/infrastructure/ai"
	"github.com/Southclaws/storyden/internal/infrastructure/vector/pinecone"
	weaviate_infra "github.com/Southclaws/storyden/internal/infrastructure/weaviate"
)

func newSemdexer(
	ctx context.Context,
	cfg config.Config,
	wc *weaviate.Client,
	pc *pinecone.Client,

	weaviateClassName weaviate_infra.WeaviateClassName,
	hydrator *hydrate.Hydrator,
	prompter ai.Prompter,
) (semdex.Semdexer, error) {
	if cfg.SemdexProvider != "" && cfg.LanguageModelProvider == "" {
		return nil, fault.New("semdex requires a language model provider to be enabled")
	}

	switch cfg.SemdexProvider {
	case "chromem":
		return chromem_semdexer.New(cfg, hydrator, prompter)

	case "weaviate":
		return weaviate_semdexer.New(wc, weaviateClassName, hydrator), nil

	case "pinecone":
		return pinecone_semdexer.New(ctx, cfg, pc, hydrator, prompter)

	default:
		return &semdex.Disabled{}, nil
	}
}

func Build() fx.Option {
	return fx.Options(
		fx.Provide(
			asker.New,
		),
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
