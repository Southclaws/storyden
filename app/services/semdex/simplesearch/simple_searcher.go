package simplesearch

import (
	"context"
	"sync"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"golang.org/x/sync/errgroup"

	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/services/semdex"
	"github.com/Southclaws/storyden/internal/ent"
)

type ParallelSearcher struct {
	searchers []semdex.Searcher
}

func NewParallelSearcher(ec *ent.Client) *ParallelSearcher {
	return &ParallelSearcher{
		searchers: []semdex.Searcher{
			&postSearcher{ec},
			&nodeSearcher{ec},
		},
	}
}

func (s *ParallelSearcher) Search(ctx context.Context, query string) (datagraph.NodeReferenceList, error) {
	mx := sync.Mutex{}
	results := []*datagraph.NodeReference{}

	eg, ctx := errgroup.WithContext(ctx)

	for _, v := range s.searchers {
		v := v
		eg.Go(func() error {
			r, err := v.Search(ctx, query)
			if err != nil {
				return err // nolint:wrapcheck
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

func indexableToResult[T datagraph.Indexable](v T) *datagraph.NodeReference {
	return &datagraph.NodeReference{
		ID:   v.GetID(),
		Kind: v.GetKind(),
		Name: v.GetName(),
		Slug: v.GetSlug(),
	}
}
