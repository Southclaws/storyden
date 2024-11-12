package simplesearch

import (
	"context"
	"sort"
	"sync"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/opt"
	"github.com/samber/lo"
	"golang.org/x/sync/errgroup"

	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/library/node_search"
	"github.com/Southclaws/storyden/app/resources/pagination"
	"github.com/Southclaws/storyden/app/resources/post/post_search"
	"github.com/Southclaws/storyden/app/services/search/searcher"
)

type ParallelSearcher struct {
	searchers map[datagraph.Kind]searcher.SingleKindSearcher
}

func NewParallelSearcher(
	post_searcher post_search.Repository,
	node_searcher node_search.Search,
) *ParallelSearcher {
	return &ParallelSearcher{
		searchers: map[datagraph.Kind]searcher.SingleKindSearcher{
			datagraph.KindThread: &postSearcher{post_searcher},
			datagraph.KindNode:   &nodeSearcher{node_searcher},
		},
	}
}

func (s *ParallelSearcher) Search(ctx context.Context, q string, p pagination.Parameters, opts searcher.Options) (*pagination.Result[datagraph.Item], error) {
	mx := sync.Mutex{}
	results := []*pagination.Result[datagraph.Item]{}

	eg, ctx := errgroup.WithContext(ctx)

	var searchers []searcher.SingleKindSearcher

	if kinds, ok := opts.Kinds.Get(); ok {
		for _, k := range kinds {
			if s, ok := s.searchers[k]; ok {
				searchers = append(searchers, s)
			}
		}
	} else {
		searchers = lo.Values(s.searchers)
	}

	if len(searchers) == 0 {
		results := pagination.NewPageResult(p, 0, []datagraph.Item{})
		return &results, nil
	}

	// Earch searcher receives a smaller page size
	subsearchPageSize := uint(p.Size() / len(searchers))
	subsearchParams := pagination.NewPageParams(uint(p.PageOneIndexed()), subsearchPageSize)

	for _, v := range searchers {
		v := v
		eg.Go(func() error {
			r, err := v.Search(ctx, q, subsearchParams)
			if err != nil {
				return err
			}

			mx.Lock()
			results = append(results, r)
			mx.Unlock()

			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	result := dt.Reduce(results, func(acc *pagination.Result[datagraph.Item], prev *pagination.Result[datagraph.Item]) *pagination.Result[datagraph.Item] {
		nextPage := acc.NextPage
		if !nextPage.Ok() {
			nextPage = prev.NextPage
		}

		totalPages := max(acc.TotalPages, prev.TotalPages)

		return &pagination.Result[datagraph.Item]{
			Size:        p.Size(),
			Results:     acc.Results + prev.Results,
			TotalPages:  totalPages,
			CurrentPage: p.PageOneIndexed(),
			NextPage:    nextPage,
			Items:       append(acc.Items, prev.Items...),
		}
	}, &pagination.Result[datagraph.Item]{
		Size:        0,
		Results:     0,
		TotalPages:  0,
		CurrentPage: 0,
		NextPage:    opt.NewEmpty[int](),
		Items:       []datagraph.Item{},
	})

	sort.Sort(datagraph.ByCreatedDesc(result.Items))

	return result, nil
}
