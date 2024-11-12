package search

import (
	"github.com/Southclaws/storyden/app/services/search/searcher"
	"github.com/Southclaws/storyden/app/services/search/simplesearch"
	"github.com/Southclaws/storyden/app/services/semdex"
	"github.com/Southclaws/storyden/internal/config"
)

func NewSearcher(
	cfg config.Config,
	simpleSearcher *simplesearch.ParallelSearcher,
	semdexSearcher semdex.Searcher,
) searcher.Searcher {
	if cfg.SemdexEnabled {
		return semdexSearcher
	}

	return simpleSearcher
}
