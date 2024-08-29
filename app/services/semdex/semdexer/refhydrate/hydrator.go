// Package refhydrate provides a Semdexer implementation which wraps an instance
// of a RefSemdexer which will provide references for read-path methods instead
// of fully hydrated Storyden objects (Post, Node, etc.) The Semdexer provided
// by this package hydrates those references by looking them up in the database.
package refhydrate

import (
	"context"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/samber/lo"

	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/post/reply"
	"github.com/Southclaws/storyden/app/resources/post/thread"
)

type Hydrator struct {
	threads thread.Repository
	replies reply.Repository
	library library.Repository
}

func New(
	threads thread.Repository,
	replies reply.Repository,
	library library.Repository,
) *Hydrator {
	return &Hydrator{
		threads: threads,
		replies: replies,
		library: library,
	}
}

func (h *Hydrator) Hydrate(ctx context.Context, refs ...*datagraph.Ref) ([]datagraph.Item, error) {
	parts := lo.GroupBy(refs, func(r *datagraph.Ref) datagraph.Kind { return r.Kind })

	// TODO: Use "GetMany" funcs so this is optimised at DB level.

	posts, err := dt.MapErr(parts[datagraph.KindPost], func(r *datagraph.Ref) (datagraph.Item, error) {
		return h.replies.Get(ctx, post.ID(r.ID))
	})
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	nodes, err := dt.MapErr(parts[datagraph.KindNode], func(r *datagraph.Ref) (datagraph.Item, error) {
		return h.library.GetByID(ctx, library.NodeID(r.ID))
	})
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	items := datagraph.ItemList{}

	items = append(items, posts...)
	items = append(items, nodes...)

	return items, nil
}
