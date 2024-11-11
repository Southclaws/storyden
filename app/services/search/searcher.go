package search

import (
	"context"

	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/datagraph/semdex"
	"github.com/Southclaws/storyden/app/services/search/simplesearch"
	"github.com/Southclaws/storyden/internal/config"
)

type Searcher interface {
	Search(ctx context.Context, query string) (datagraph.ItemList, error)
}

func NewSearcher(
	cfg config.Config,
	simpleSearcher *simplesearch.ParallelSearcher,
	semdexSearcher semdex.Searcher,
) Searcher {
	if cfg.SemdexEnabled {
		return semdexSearcher
	}

	return simpleSearcher
}
