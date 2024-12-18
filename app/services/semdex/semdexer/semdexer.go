package semdexer

import (
	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"go.uber.org/fx"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/storyden/app/resources/datagraph/hydrate"
	"github.com/Southclaws/storyden/app/services/semdex"
	"github.com/Southclaws/storyden/app/services/semdex/semdexer/chromem_semdexer"
	"github.com/Southclaws/storyden/app/services/semdex/semdexer/weaviate_semdexer"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/infrastructure/ai"
	weaviate_infra "github.com/Southclaws/storyden/internal/infrastructure/weaviate"
)

func newSemdexer(
	cfg config.Config,
	wc *weaviate.Client,

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

	default:
		return &semdex.Disabled{}, nil
	}
}

func Build() fx.Option {
	return fx.Provide(
		fx.Annotate(
			newSemdexer,
			fx.As(new(semdex.Semdexer)),
			fx.As(new(semdex.Querier)),
			fx.As(new(semdex.Mutator)),
			fx.As(new(semdex.Recommender)),
			fx.As(new(semdex.Searcher)),
		),
	)
}
