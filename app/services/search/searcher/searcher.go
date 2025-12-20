package searcher

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"

	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/pagination"
)

type Options struct {
	Kinds opt.Optional[[]datagraph.Kind]
}

var ErrFastMatchesUnavailable = fault.New("datagraph matches are not enabled", ftag.With(ftag.InvalidArgument))

type Searcher interface {
	Search(ctx context.Context, q string, p pagination.Parameters, opts Options) (*pagination.Result[datagraph.Item], error)
	MatchFast(ctx context.Context, q string, limit int, opts Options) (datagraph.MatchList, error)
}

type Indexer interface {
	Index(ctx context.Context, item datagraph.Item) error
	Deindex(ctx context.Context, ir datagraph.ItemRef) error
}
