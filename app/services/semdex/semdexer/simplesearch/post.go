package simplesearch

import (
	"context"

	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/internal/ent"
)

type postSearcher struct {
	ec *ent.Client
}

func (s *postSearcher) Search(ctx context.Context, query string) (datagraph.ItemList, error) {
	// TODO: Implement searcher that returns resources that implement Item

	return nil, nil
}
