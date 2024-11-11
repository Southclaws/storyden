package simplesearch

import (
	"context"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/library/node_search"
)

type nodeSearcher struct {
	node_search node_search.Search
}

func (s *nodeSearcher) Search(ctx context.Context, query string) (datagraph.ItemList, error) {
	rs, err := s.node_search.Search(ctx, node_search.WithNameContains(query))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	items := dt.Map(rs, func(r *library.Node) datagraph.Item { return r })

	return items, nil
}
