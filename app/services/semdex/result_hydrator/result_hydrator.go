package result_hydrator

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"

	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/datagraph/cluster"
	"github.com/Southclaws/storyden/app/resources/datagraph/item"
	"github.com/Southclaws/storyden/app/resources/datagraph/link"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/reply"
	"github.com/Southclaws/storyden/app/resources/thread"
	"github.com/Southclaws/storyden/app/services/semdex"
)

type Hydrator struct {
	tr thread.Repository
	rr reply.Repository
	ir item.Repository
	cr cluster.Repository
	lr link.Repository
}

func New(
	tr thread.Repository,
	rr reply.Repository,
	ir item.Repository,
	cr cluster.Repository,
	lr link.Repository,
) *Hydrator {
	return &Hydrator{tr, rr, ir, cr, lr}
}

func (h *Hydrator) Hydrate(ctx context.Context, sr *semdex.Result) (*semdex.Result, error) {
	switch sr.Type {
	case datagraph.KindThread:
		t, err := h.tr.Get(ctx, post.ID(sr.Id))
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}

		sr.Name = t.Title
		sr.Description = t.Short
		sr.Slug = t.Slug

		return sr, nil

	case datagraph.KindReply:

		r, err := h.rr.Get(ctx, post.ID(sr.Id))
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}

		sr.Name = r.RootThreadTitle
		sr.Description = r.Short
		sr.Slug = r.RootThreadMark

		return sr, nil

	case datagraph.KindCluster:
		c, err := h.cr.GetByID(ctx, datagraph.ClusterID(sr.Id))
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}

		sr.Name = c.Name
		sr.Description = c.Description

		return sr, nil

	case datagraph.KindItem:
		i, err := h.ir.GetByID(ctx, datagraph.ItemID(sr.Id))
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}

		sr.Name = i.Name
		sr.Description = i.Description

		return sr, nil

	case datagraph.KindLink:
		ln, err := h.lr.GetByID(ctx, datagraph.LinkID(sr.Id))
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}

		sr.Name = ln.Title.OrZero()
		sr.Description = ln.Description.OrZero()

		return sr, nil

	default:
		return sr, nil
	}
}
