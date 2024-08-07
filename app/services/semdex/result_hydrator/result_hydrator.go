package result_hydrator

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"

	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/link"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/post/reply"
	"github.com/Southclaws/storyden/app/resources/post/thread"
)

type Hydrator struct {
	tr thread.Repository
	rr reply.Repository
	nr library.Repository
	lr link.Repository
}

func New(
	tr thread.Repository,
	rr reply.Repository,
	nr library.Repository,
	lr link.Repository,
) *Hydrator {
	return &Hydrator{tr, rr, nr, lr}
}

func (h *Hydrator) Hydrate(ctx context.Context, sr *datagraph.NodeReference) (*datagraph.NodeReference, error) {
	switch sr.Kind {
	case datagraph.KindPost:
		r, err := h.rr.Get(ctx, post.ID(sr.ID))
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}

		sr.Name = r.RootThreadTitle
		sr.Description = r.Content.Short()
		sr.Slug = r.RootThreadMark

		return sr, nil

	case datagraph.KindNode:
		c, err := h.nr.GetByID(ctx, library.NodeID(sr.ID))
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}

		sr.Name = c.Name
		sr.Description = c.Content.OrZero().Short()

		return sr, nil

	default:
		return sr, nil
	}
}
