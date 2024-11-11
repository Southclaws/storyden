package simplesearch

import (
	"context"
	"sync"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/samber/lo"
	"golang.org/x/sync/errgroup"

	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/library/node_search"
	"github.com/Southclaws/storyden/app/resources/pagination"
	"github.com/Southclaws/storyden/app/resources/post/post_search"
	"github.com/Southclaws/storyden/app/services/search/searcher"
)

type Basic interface {
	Search(ctx context.Context, query string) (datagraph.ItemList, error)
}

type ParallelSearcher struct {
	searchers map[datagraph.Kind]Basic
}

func NewParallelSearcher(
	post_searcher post_search.Repository,
	node_searcher node_search.Search,
) *ParallelSearcher {
	return &ParallelSearcher{
		searchers: map[datagraph.Kind]Basic{
			datagraph.KindPost: &postSearcher{post_searcher},
			datagraph.KindNode: &nodeSearcher{node_searcher},
		},
	}
}

func (s *ParallelSearcher) Search(ctx context.Context, q string, p pagination.Parameters, opts searcher.Options) (*pagination.Result[datagraph.Item], error) {
	mx := sync.Mutex{}
	results := datagraph.ItemList{}

	eg, ctx := errgroup.WithContext(ctx)

	var searchers []Basic

	if kinds, ok := opts.Kinds.Get(); ok {
		for _, k := range kinds {
			if s, ok := s.searchers[k]; ok {
				searchers = append(searchers, s)
			}
		}
	} else {
		searchers = lo.Values(s.searchers)
	}

	for _, v := range searchers {
		v := v
		eg.Go(func() error {
			r, err := v.Search(ctx, q)
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

	// TODO: Fix
	total := len(results)

	result := pagination.NewPageResult(p, total, results)

	return &result, nil
}
