package item_tree

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/datagraph/cluster"
	"github.com/Southclaws/storyden/app/resources/datagraph/item"
	"github.com/Southclaws/storyden/app/services/authentication/session"
)

var errNotAuthorised = fault.Wrap(fault.New("not authorised"), ftag.With(ftag.PermissionDenied))

type Graph interface {
	// Link adds an item to a cluster.
	Link(ctx context.Context, item datagraph.ItemSlug, cluster datagraph.ClusterSlug) (*datagraph.Item, error)

	// Sever removes an item from a cluster if it was a member.
	Sever(ctx context.Context, item datagraph.ItemSlug, cluster datagraph.ClusterSlug) (*datagraph.Item, error)
}

type service struct {
	cr cluster.Repository
	ir item.Repository
}

func New(cr cluster.Repository, ir item.Repository) Graph {
	return &service{cr: cr, ir: ir}
}

func (s *service) Link(ctx context.Context, is datagraph.ItemSlug, cs datagraph.ClusterSlug) (*datagraph.Item, error) {
	accountID, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	clus, err := s.cr.Get(ctx, cs)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	itm, err := s.ir.Get(ctx, is)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if !clus.Owner.Admin {
		if clus.Owner.ID != accountID && itm.Owner.ID != accountID {
			return nil, fault.Wrap(errNotAuthorised, fctx.With(ctx))
		}
	}

	itm, err = s.ir.Update(ctx, itm.ID, item.WithParentClusterAdd(xid.ID(clus.ID)))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return itm, nil
}

func (s *service) Sever(ctx context.Context, is datagraph.ItemSlug, cs datagraph.ClusterSlug) (*datagraph.Item, error) {
	accountID, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	clus, err := s.cr.Get(ctx, cs)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	itm, err := s.ir.Get(ctx, is)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if !clus.Owner.Admin {
		if clus.Owner.ID != accountID && itm.Owner.ID != accountID {
			return nil, fault.Wrap(errNotAuthorised, fctx.With(ctx))
		}
	}

	itm, err = s.ir.Update(ctx, itm.ID, item.WithParentClusterRemove(xid.ID(clus.ID)))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return itm, nil
}
