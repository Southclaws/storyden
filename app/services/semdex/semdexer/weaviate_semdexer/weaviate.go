package weaviate_semdexer

import (
	"github.com/weaviate/weaviate-go-client/v4/weaviate"

	"github.com/Southclaws/storyden/app/services/semdex/semdexer/refhydrate"
	weaviate_infra "github.com/Southclaws/storyden/internal/infrastructure/weaviate"
)

// weaviateRefIndex implements what looks slightly like the Semdexer interface
// but all of its methods return references, rather than fully hydrated objects.
// This is because hydration is somewhat costly and not always what you need.
// It also separates the responsibility of hydrating content from the resource
// layer from the vector database. If you need to operate on lower level refs,
// this is what is best to use because it doesn't make costly database calls.
type weaviateRefIndex struct {
	wc *weaviate.Client
	cn weaviate_infra.WeaviateClassName
}

func New(
	wc *weaviate.Client,
	cn weaviate_infra.WeaviateClassName,
	rh *refhydrate.Hydrator,
) *refhydrate.HydratedSemdexer {
	ws := &weaviateRefIndex{wc, cn}

	return &refhydrate.HydratedSemdexer{
		RefSemdex: ws,
		Hydrator:  rh,
	}
}
