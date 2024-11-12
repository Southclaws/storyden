package searcher

import (
	"context"

	"github.com/Southclaws/opt"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/pagination"
)

type Options struct {
	Kinds opt.Optional[[]datagraph.Kind]
}

type Searcher interface {
	Search(ctx context.Context, q string, p pagination.Parameters, opts Options) (*pagination.Result[datagraph.Item], error)
}

type SingleKindSearcher interface {
	Search(ctx context.Context, query string, p pagination.Parameters) (*pagination.Result[datagraph.Item], error)
}
