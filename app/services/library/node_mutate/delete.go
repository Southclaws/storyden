package node_mutate

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/opt"

	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/services/library/nodetree"
	"github.com/Southclaws/storyden/lib/plugin/rpc"
)

type DeleteOptions struct {
	NewParent opt.Optional[library.QueryKey]
}

func (s *Manager) Delete(ctx context.Context, qk library.QueryKey, d DeleteOptions) (*library.Node, error) {
	n, err := s.nodeQuerier.Get(ctx, qk)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if err := authoriseDelete(ctx, n); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	destination, err := opt.MapErr(d.NewParent, func(target library.QueryKey) (*library.Node, error) {
		parent, err := s.nodeQuerier.Get(ctx, target)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}

		if err := nodetree.AuthoriseOwner(ctx, parent); err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}

		return s.nc.Move(ctx, qk, target)
	})
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if err := s.cache.InvalidateBeforeWrite(ctx); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	err = s.nodeWriter.Delete(ctx, qk)
	s.cache.InvalidateAfterWrite(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	s.bus.Publish(ctx, &rpc.EventNodeDeleted{
		ID:   library.NodeID(n.GetID()),
		Slug: n.GetSlug(),
	})

	return destination.OrZero(), nil
}

func authoriseDelete(ctx context.Context, node *library.Node) error {
	if err := nodetree.AuthoriseOwner(ctx, node); err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	for _, child := range node.Nodes {
		if err := nodetree.AuthoriseOwner(ctx, child); err != nil {
			return fault.Wrap(err, fctx.With(ctx))
		}
	}

	return nil
}
