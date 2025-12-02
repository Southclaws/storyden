// Package hydrate provides a generic datagraph item lookup conversion.
package hydrate

import (
	"context"
	"fmt"
	"sort"

	"github.com/Southclaws/dt"
	"github.com/samber/lo"
	"golang.org/x/sync/errgroup"

	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/library/node_querier"
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

	results := make(chan withRelevance, len(refs))

	eg, ctx := errgroup.WithContext(ctx)

	for k, v := range parts {
		eg.Go(func() error {
			relevanceMap := make(map[string]float64)
			for _, r := range v {
				relevanceMap[r.ID.String()] = r.Relevance
			}

			switch k {
			case datagraph.KindPost:
				ids := dt.Map(v, func(r *datagraph.Ref) post.ID { return post.ID(r.ID) })
				items, err := h.replies.GetMany(ctx, ids...)
				if err != nil {
					return fmt.Errorf("failed to get posts for kind=%s ids=%v: %w", k, ids, err)
				}
				for _, item := range items {
					results <- withRelevance{item, relevanceMap[item.ID.String()]}
				}

			case datagraph.KindThread:
				ids := dt.Map(v, func(r *datagraph.Ref) post.ID { return post.ID(r.ID) })
				items, err := h.threads.GetMany(ctx, ids, nil)
				if err != nil {
					return fmt.Errorf("failed to get threads for kind=%s ids=%v: %w", k, ids, err)
				}
				for _, item := range items {
					results <- withRelevance{item, relevanceMap[item.ID.String()]}
				}

			case datagraph.KindReply:
				ids := dt.Map(v, func(r *datagraph.Ref) post.ID { return post.ID(r.ID) })
				items, err := h.replies.GetMany(ctx, ids...)
				if err != nil {
					return fmt.Errorf("failed to get replies for kind=%s ids=%v: %w", k, ids, err)
				}
				for _, item := range items {
					results <- withRelevance{item, relevanceMap[item.ID.String()]}
				}

			case datagraph.KindNode:
				ids := dt.Map(v, func(r *datagraph.Ref) library.NodeID { return library.NodeID(r.ID) })
				items, err := h.nodeQuerier.ProbeMany(ctx, ids...)
				if err != nil {
					return fmt.Errorf("failed to get nodes for kind=%s ids=%v: %w", k, ids, err)
				}
				for _, item := range items {
					results <- withRelevance{item, relevanceMap[item.GetID().String()]}
				}

			case datagraph.KindCollection:
				// TODO

			case datagraph.KindProfile:
				// TODO

			case datagraph.KindEvent:
				// TODO
			}

			return nil
		})
	}

	go func() {
		eg.Wait()
		close(results)
	}()

	var hydrated sortedByRelevance
	for items := range results {
		hydrated = append(hydrated, items)
	}

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	sort.Sort(hydrated)

	sorted := dt.Map(hydrated, func(i withRelevance) datagraph.Item {
		return i.Item
	})

	return sorted, nil
}
