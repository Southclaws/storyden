package simplesearch

import (
	"context"
	"sync"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"golang.org/x/sync/errgroup"

	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/datagraph/semdex"
	"github.com/Southclaws/storyden/app/resources/library/node_search"
	"github.com/Southclaws/storyden/app/resources/post/post_search"
)

type ParallelSearcher struct {
	searchers []semdex.Searcher
}

func NewParallelSearcher(
	post_searcher post_search.Repository,
	node_searcher node_search.Search,
) *ParallelSearcher {
	return &ParallelSearcher{
		searchers: []semdex.Searcher{
			&postSearcher{post_searcher},
			&nodeSearcher{node_searcher},
		},
	}
}

func (s *ParallelSearcher) Search(ctx context.Context, query string) (datagraph.ItemList, error) {
	mx := sync.Mutex{}
	results := datagraph.ItemList{}

	eg, ctx := errgroup.WithContext(ctx)

	for _, v := range s.searchers {
		v := v
		eg.Go(func() error {
			r, err := v.Search(ctx, query)
			if err != nil {
				return err
			}

			mx.Lock()
			results = append(results, r...)
			mx.Unlock()

			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return results, nil
}
