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
	searchers map[datagraph.Kind]searcher.Searcher
}

func NewParallelSearcher(
	post_searcher post_search.Repository,
	node_searcher node_search.Search,
) *ParallelSearcher {
	// NOTE: We use the same postSearcher instance for Post, Thread, and Reply.
	// Post is an abstract type over thread and reply, searcher returns both.
	ps := &postSearcher{post_searcher}
	return &ParallelSearcher{
		searchers: map[datagraph.Kind]searcher.Searcher{
			datagraph.KindPost:   ps, // Same instance
			datagraph.KindThread: ps, // Same instance
			datagraph.KindReply:  ps, // Same instance
			datagraph.KindNode:   &nodeSearcher{node_searcher},
		},
	}
}

func (s *ParallelSearcher) Search(ctx context.Context, q string, p pagination.Parameters, opts searcher.Options) (*pagination.Result[datagraph.Item], error) {
	mx := sync.Mutex{}
	results := []*pagination.Result[datagraph.Item]{}

	eg, ctx := errgroup.WithContext(ctx)

	var searchers []searcher.Searcher

	if kinds, ok := opts.Kinds.Get(); ok {
		// TRICK: Deduplicate searchers by pointer identity. Post/Thread/Reply
		// all use the same postSearcher instance, this map will prevent running
		// the same searcher multiple times when searching multiple post kinds.
		seenSearchers := make(map[searcher.Searcher]bool)
		for _, k := range kinds {
			if searcher, ok := s.searchers[k]; ok {
				if !seenSearchers[searcher] {
					searchers = append(searchers, searcher)
					seenSearchers[searcher] = true
				}
			}
		}
	} else {
		searchers = lo.Uniq(lo.Values(s.searchers))
	}

	if len(searchers) == 0 {
		results := pagination.NewPageResult(p, 0, []datagraph.Item{})
		return &results, nil
	}

	// Earch searcher receives a smaller page size
	subsearchPageSize := uint(p.Size() / len(searchers))
	subsearchParams := pagination.NewPageParams(uint(p.PageOneIndexed()), subsearchPageSize)

	for _, v := range searchers {
		eg.Go(func() error {
			r, err := v.Search(ctx, q, subsearchParams, opts)
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

func (s *ParallelSearcher) MatchFast(ctx context.Context, q string, limit int, opts searcher.Options) (datagraph.MatchList, error) {
	return nil, searcher.ErrFastMatchesUnavailable
}
