// Package refhydrate provides a Semdexer implementation which wraps an instance
// of a RefSemdexer which will provide references for read-path methods instead
// of fully hydrated Storyden objects (Post, Node, etc.) The Semdexer provided
// by this package hydrates those references by looking them up in the database.
package refhydrate

import (
	"context"
	"sort"
	"sync"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/samber/lo"

	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/library/node_querier"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/post/reply"
	"github.com/Southclaws/storyden/app/resources/post/thread"
)

type Hydrator struct {
	threads     thread.Repository
	replies     reply.Repository
	nodeQuerier *node_querier.Querier
}

func New(
	threads thread.Repository,
	replies reply.Repository,
	nodeQuerier *node_querier.Querier,
) *Hydrator {
	return &Hydrator{
		threads:     threads,
		replies:     replies,
		nodeQuerier: nodeQuerier,
	}
}

type withRelevance struct {
	datagraph.Item
	r float64
}

type sortedByRelevance []withRelevance

func (a sortedByRelevance) Len() int           { return len(a) }
func (a sortedByRelevance) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a sortedByRelevance) Less(i, j int) bool { return a[i].r > a[j].r }

func (h *Hydrator) Hydrate(ctx context.Context, refs ...*datagraph.Ref) (datagraph.ItemList, error) {
	parts := lo.GroupBy(refs, func(r *datagraph.Ref) datagraph.Kind { return r.Kind })

	// TODO: Use "GetMany" funcs so this is optimised at DB level.

	results := make(chan []withRelevance, len(refs))
	errChan := make(chan error)

	wg := sync.WaitGroup{}

	for k, v := range parts {
		wg.Add(1)
		go func() {
			var items []withRelevance
			var err error

			switch k {
			case datagraph.KindPost:
				// TODO: Repo for generic post types.
				items, err = dt.MapErr(v, func(r *datagraph.Ref) (withRelevance, error) {
					i, err := h.replies.Get(ctx, post.ID(r.ID))
					return withRelevance{Item: i, r: r.Relevance}, err
				})

			case datagraph.KindThread:
				items, err = dt.MapErr(v, func(r *datagraph.Ref) (withRelevance, error) {
					i, err := h.threads.Get(ctx, post.ID(r.ID), nil)
					return withRelevance{Item: i, r: r.Relevance}, err
				})

			case datagraph.KindReply:
				items, err = dt.MapErr(v, func(r *datagraph.Ref) (withRelevance, error) {
					i, err := h.replies.Get(ctx, post.ID(r.ID))
					return withRelevance{Item: i, r: r.Relevance}, err
				})

			case datagraph.KindNode:
				items, err = dt.MapErr(v, func(r *datagraph.Ref) (withRelevance, error) {
					i, err := h.nodeQuerier.Probe(ctx, library.NodeID(r.ID))
					return withRelevance{Item: i, r: r.Relevance}, err
				})

			case datagraph.KindCollection:
				// TODO

			case datagraph.KindProfile:
				// TODO

			case datagraph.KindEvent:
				// TODO
			}

			if err != nil {
				errChan <- err
			}

			results <- items

			wg.Done()
		}()
	}

	go func() {
		wg.Wait()

		close(results)
		close(errChan)
	}()

	if waitErr := <-errChan; waitErr != nil {
		return nil, fault.Wrap(waitErr, fctx.With(ctx))
	}

	var hydrated sortedByRelevance
	for items := range results {
		hydrated = append(hydrated, items...)
	}

	sort.Sort(hydrated)

	sorted := dt.Map(hydrated, func(i withRelevance) datagraph.Item {
		return i.Item
	})

	return sorted, nil
}
