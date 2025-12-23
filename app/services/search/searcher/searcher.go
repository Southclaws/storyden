package searcher

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/pagination"
	"github.com/Southclaws/storyden/app/resources/post/category"
	"github.com/Southclaws/storyden/app/resources/tag/tag_ref"
)

type Options struct {
	Kinds      opt.Optional[[]datagraph.Kind]
	Authors    opt.Optional[[]account.AccountID]
	Categories opt.Optional[[]category.CategoryID]
	Tags       opt.Optional[[]tag_ref.Name]
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
