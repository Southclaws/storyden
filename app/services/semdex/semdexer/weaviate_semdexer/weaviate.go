package weaviate_semdexer

import (
	"github.com/weaviate/weaviate-go-client/v5/weaviate"

	"github.com/Southclaws/storyden/app/resources/datagraph/hydrate"
	weaviate_infra "github.com/Southclaws/storyden/internal/infrastructure/weaviate"
)

type weaviateSemdexer struct {
	wc       *weaviate.Client
	cn       weaviate_infra.WeaviateClassName
	hydrator *hydrate.Hydrator
}

func New(
	wc *weaviate.Client,
	cn weaviate_infra.WeaviateClassName,
	hydrator *hydrate.Hydrator,
) *weaviateSemdexer {
	return &weaviateSemdexer{
		wc:       wc,
		cn:       cn,
		hydrator: hydrator,
	}
}
