// Package semdex provides an interface for semantic indexing of the datagraph.
package semdex

import (
	"context"
	"net/url"

	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/pagination"
	"github.com/Southclaws/storyden/app/services/search/searcher"
)

type Semdexer interface {
	Mutator
	Querier
}

type Mutator interface {
	Index(ctx context.Context, object datagraph.Item) error
	Delete(ctx context.Context, object xid.ID) error
}

type Querier interface {
	Searcher
	Recommender

	GetMany(ctx context.Context, limit uint, ids ...xid.ID) (datagraph.RefList, error)
}

type Chunk struct {
	ID      xid.ID
	Kind    datagraph.Kind
	URL     url.URL
	Content string
}

type Searcher interface {
	Search(ctx context.Context, q string, p pagination.Parameters, opts searcher.Options) (*pagination.Result[datagraph.Item], error)
	SearchRefs(ctx context.Context, q string, p pagination.Parameters, opts searcher.Options) (*pagination.Result[*datagraph.Ref], error)
	SearchChunks(ctx context.Context, q string, p pagination.Parameters, opts searcher.Options) ([]*Chunk, error)
}

type Asker interface {
	Ask(ctx context.Context, q string) (chan string, chan error)
}

type Recommender interface {
	Recommend(ctx context.Context, object datagraph.Item) (datagraph.ItemList, error)
	RecommendRefs(ctx context.Context, object datagraph.Item) (datagraph.RefList, error)
}
