package clustersearch

import (
	"context"

	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/internal/ent"
)

type Search interface {
	Search(ctx context.Context, query string, root datagraph.ClusterSlug) ([]*datagraph.Cluster, error)
}

type service struct {
	ec ent.Client
}

func New(ec ent.Client) Search {
	return &service{ec: ec}
}

func (s *service) Search(ctx context.Context, q string, root datagraph.ClusterSlug) ([]*datagraph.Cluster, error) {
	// query := s.ec.Cluster.Query()
	// query.
	return nil, nil
}
