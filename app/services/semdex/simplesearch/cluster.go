package simplesearch

import (
	"context"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"

	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/ent/cluster"
)

type clusterSearcher struct {
	ec *ent.Client
}

func (s *clusterSearcher) Search(ctx context.Context, query string) (datagraph.NodeReferenceList, error) {
	cq := s.ec.Cluster.Query().Where(
		cluster.Or(
			cluster.NameContainsFold(query),
			cluster.DescriptionContainsFold(query),
			cluster.ContentContainsFold(query),
		),
	)

	rs, err := cq.All(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	clusters, err := dt.MapErr(rs, datagraph.ClusterFromModel)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	results := dt.Map(clusters, indexableToResult)

	return results, nil
}
