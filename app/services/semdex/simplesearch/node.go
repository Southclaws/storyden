package simplesearch

import (
	"context"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"

	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/ent/node"
)

type nodeSearcher struct {
	ec *ent.Client
}

func (s *nodeSearcher) Search(ctx context.Context, query string) (datagraph.NodeReferenceList, error) {
	cq := s.ec.Node.Query().Where(
		node.Or(
			node.NameContainsFold(query),
			node.DescriptionContainsFold(query),
			node.ContentContainsFold(query),
		),
	)

	rs, err := cq.All(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	nodes, err := dt.MapErr(rs, datagraph.NodeFromModel)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	results := dt.Map(nodes, indexableToResult)

	return results, nil
}
