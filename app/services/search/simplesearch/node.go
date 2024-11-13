package simplesearch

import (
	"context"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"

	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/library/node_search"
	"github.com/Southclaws/storyden/app/resources/pagination"
)

type nodeSearcher struct {
	node_search node_search.Search
}

func (s *nodeSearcher) Search(ctx context.Context, query string, p pagination.Parameters) (*pagination.Result[datagraph.Item], error) {
	rs, err := s.node_search.Search(ctx, p, node_search.WithNameContains(query), node_search.WithContentContains(query))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	items := dt.Map(rs.Items, func(r *library.Node) datagraph.Item { return r })

	result := pagination.ConvertPageResult(*rs, items)

	return &result, nil
}
