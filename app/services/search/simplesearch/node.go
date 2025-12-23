package simplesearch

import (
	"context"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/library/node_search"
	"github.com/Southclaws/storyden/app/resources/pagination"
	"github.com/Southclaws/storyden/app/resources/tag/tag_ref"
	"github.com/Southclaws/storyden/app/services/search/searcher"
)

type nodeSearcher struct {
	node_search node_search.Search
}

func (s *nodeSearcher) Search(ctx context.Context, query string, p pagination.Parameters, opts searcher.Options) (*pagination.Result[datagraph.Item], error) {
	o := []node_search.Option{
		node_search.WithNameContains(query),
		node_search.WithContentContains(query),
	}

	opts.Authors.Call(func(value []account.AccountID) {
		o = append(o, node_search.WithAuthors(value...))
	})

	opts.Tags.Call(func(value []tag_ref.Name) {
		o = append(o, node_search.WithTags(value...))
	})

	rs, err := s.node_search.Search(ctx, p, o...)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	items := dt.Map(rs.Items, func(r *library.Node) datagraph.Item { return r })

	result := pagination.ConvertPageResult(*rs, items)

	return &result, nil
}

func (s *nodeSearcher) MatchFast(ctx context.Context, q string, limit int, opts searcher.Options) (datagraph.MatchList, error) {
	return nil, searcher.ErrFastMatchesUnavailable
}
