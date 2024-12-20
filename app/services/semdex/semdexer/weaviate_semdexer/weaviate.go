package weaviate_semdexer

import (
	"github.com/weaviate/weaviate-go-client/v4/weaviate"

	"github.com/Southclaws/storyden/app/resources/datagraph/hydrate"
	"github.com/Southclaws/storyden/internal/infrastructure/ai"
	weaviate_infra "github.com/Southclaws/storyden/internal/infrastructure/weaviate"
)

type weaviateSemdexer struct {
	wc       *weaviate.Client
	cn       weaviate_infra.WeaviateClassName
	ai       ai.Prompter
	hydrator *hydrate.Hydrator
}

func New(
	wc *weaviate.Client,
	cn weaviate_infra.WeaviateClassName,
	ai ai.Prompter,
	hydrator *hydrate.Hydrator,
) *weaviateSemdexer {
	return &weaviateSemdexer{
		wc:       wc,
		cn:       cn,
		ai:       ai,
		hydrator: hydrator,
	}
}
