// Package hydrate provides a generic datagraph item lookup conversion.
package hydrate

import (
	"context"
	"sort"
	"sync"

	"github.com/Southclaws/dt"
	"github.com/samber/lo"

	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/library/node_querier"
	"github.com/Southclaws/storyden/app/resources/pagination"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/post/reply"
	"github.com/Southclaws/storyden/app/resources/post/thread_querier"
	"github.com/Southclaws/storyden/internal/infrastructure/instrumentation/spanner"
)

type Hydrator struct {
	ins         spanner.Instrumentation
	threads     *thread_querier.Querier
	replies     reply.Repository
	nodeQuerier *node_querier.Querier
}

func New(
	ins spanner.Builder,
	threads *thread_querier.Querier,
	replies reply.Repository,
	nodeQuerier *node_querier.Querier,
) *Hydrator {
	return &Hydrator{
		ins:         ins.Build(),
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
	ctx, span := h.ins.Instrument(ctx)
	defer span.End()

	parts := lo.GroupBy(refs, func(r *datagraph.Ref) datagraph.Kind { return r.Kind })

	// TODO: Use "GetMany" funcs so this is optimised at DB level.

	results := make(chan withRelevance, len(refs))

	wg := sync.WaitGroup{}

	for k, v := range parts {
		wg.Add(1)
		go func() {
			switch k {
			case datagraph.KindPost:
				// TODO: Repo for generic post types.
				for _, r := range v {
					i, err := h.replies.Get(ctx, post.ID(r.ID))
					if err == nil {
						results <- withRelevance{i, r.Relevance}
					}
				}

			case datagraph.KindThread:
				for _, r := range v {
					i, err := h.threads.Get(ctx, post.ID(r.ID), pagination.Parameters{}, nil)
					if err == nil {
						results <- withRelevance{i, r.Relevance}
					}
				}

			case datagraph.KindReply:
				for _, r := range v {
					i, err := h.replies.Get(ctx, post.ID(r.ID))
					if err == nil {
						results <- withRelevance{i, r.Relevance}
					}
				}

			case datagraph.KindNode:
				for _, r := range v {
					i, err := h.nodeQuerier.Probe(ctx, library.NodeID(r.ID))
					if err == nil {
						results <- withRelevance{i, r.Relevance}
					}
				}

			case datagraph.KindCollection:
				// TODO

			case datagraph.KindProfile:
				// TODO

			case datagraph.KindEvent:
				// TODO
			}

			wg.Done()
		}()
	}

	go func() {
		wg.Wait()

		close(results)
	}()

	var hydrated sortedByRelevance
	for items := range results {
		hydrated = append(hydrated, items)
	}

	sort.Sort(hydrated)

	sorted := dt.Map(hydrated, func(i withRelevance) datagraph.Item {
		return i.Item
	})

	return sorted, nil
}
