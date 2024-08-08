package simplesearch

import (
	"context"

	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/internal/ent"
)

type nodeSearcher struct {
	ec *ent.Client
}

func (s *nodeSearcher) Search(ctx context.Context, query string) (datagraph.ItemList, error) {
	
	// TODO: Same as postSearcher
	return nil, nil
}
