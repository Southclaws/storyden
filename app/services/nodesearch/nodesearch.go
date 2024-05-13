package nodesearch

import (
	"context"

	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/internal/ent"
)

type Search interface {
	Search(ctx context.Context, query string, root datagraph.NodeSlug) ([]*datagraph.Node, error)
}

type service struct {
	ec ent.Client
}

func New(ec ent.Client) Search {
	return &service{ec: ec}
}

func (s *service) Search(ctx context.Context, q string, root datagraph.NodeSlug) ([]*datagraph.Node, error) {
	// query := s.ec.Node.Query()
	// query.
	return nil, nil
}
