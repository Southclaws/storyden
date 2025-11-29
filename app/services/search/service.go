package search

import (
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/services/search/bleve_search"
	"github.com/Southclaws/storyden/app/services/search/redis_search"
	"github.com/Southclaws/storyden/app/services/search/searcher"
	"github.com/Southclaws/storyden/app/services/search/simplesearch"
	"github.com/Southclaws/storyden/internal/config"
)

func newSearcher(
	cfg config.Config,
	simpleSearcher *simplesearch.ParallelSearcher,
	bleveSearcher *bleve_search.BleveSearcher,
	redisSearcher *redis_search.RedisSearcher,
) searcher.Searcher {
	switch cfg.SearchProvider {
	case "bleve":
		return bleveSearcher

	case "redis":
		return redisSearcher

	case "database":
		fallthrough
	default:
		return simpleSearcher
	}
}

func newIndexer(
	cfg config.Config,
	bleveSearcher *bleve_search.BleveSearcher,
	redisSearcher *redis_search.RedisSearcher,
) searcher.Indexer {
	switch cfg.SearchProvider {
	case "bleve":
		return bleveSearcher

	case "redis":
		return redisSearcher

	default:
		return nil
	}
}

func Build() fx.Option {
	return fx.Options(
		fx.Provide(
			newSearcher,
			newIndexer,
			simplesearch.NewParallelSearcher,
		),
	)
}
