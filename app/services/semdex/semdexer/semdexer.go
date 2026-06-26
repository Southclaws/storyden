package semdexer

import (
	"context"

	"go.uber.org/fx"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/storyden/app/resources/datagraph/hydrate"
	"github.com/Southclaws/storyden/app/services/semdex"
	"github.com/Southclaws/storyden/app/services/semdex/semdexer/chromem_semdexer"
	"github.com/Southclaws/storyden/app/services/semdex/semdexer/pinecone_semdexer"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/infrastructure/ai"
	"github.com/Southclaws/storyden/internal/infrastructure/vector/pinecone"
)

func newSemdexer(
	ctx context.Context,
	cfg config.Config,
	pc *pinecone.Client,

	hydrator *hydrate.Hydrator,
	prompter ai.Prompter,
) (semdex.Semdexer, error) {
	if cfg.SemdexProvider != "" && cfg.LanguageModelProvider == "" {
		return nil, fault.New("semdex requires a language model provider to be enabled")
	}

	switch cfg.SemdexProvider {
	case "chromem":
		return chromem_semdexer.New(cfg, hydrator, prompter)

	case "pinecone":
		return pinecone_semdexer.New(ctx, cfg, pc, hydrator, prompter)

	default:
		return &semdex.Disabled{}, nil
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
