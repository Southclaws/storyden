package result_hydrator

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"

	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/datagraph/link"
	"github.com/Southclaws/storyden/app/resources/datagraph/node"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/post/reply"
	"github.com/Southclaws/storyden/app/resources/thread"
)

type Hydrator struct {
	tr thread.Repository
	rr reply.Repository
	nr node.Repository
	lr link.Repository
}

func New(
	tr thread.Repository,
	rr reply.Repository,
	nr node.Repository,
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
		c, err := h.nr.GetByID(ctx, datagraph.NodeID(sr.ID))
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
