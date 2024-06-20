package nodetree

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/datagraph/node"
	"github.com/Southclaws/storyden/app/services/authentication/session"
)

var errNotAuthorised = fault.New("not authorised", ftag.With(ftag.PermissionDenied))

var ErrIdenticalParentChild = fault.New("cannot relate a node to itself", ftag.With(ftag.InvalidArgument))

type Graph interface {
	// Move moves a node from either orphan state or belonging to one node
	// to another node essentially setting its parent slug to some/new value.
	Move(ctx context.Context, child datagraph.NodeSlug, parent datagraph.NodeSlug) (*datagraph.Node, error)

	// Sever orphans a node by removing it from its parent to the root level.
	Sever(ctx context.Context, child datagraph.NodeSlug, parent datagraph.NodeSlug) (*datagraph.Node, error)
}

type service struct {
	nr node.Repository
}

func New(nr node.Repository) Graph {
	return &service{nr: nr}
}

func (s *service) Move(ctx context.Context, child datagraph.NodeSlug, parent datagraph.NodeSlug) (*datagraph.Node, error) {
	if child == parent {
		return nil, fault.Wrap(ErrIdenticalParentChild, fctx.With(ctx))
	}

	accountID, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	cnode, err := s.nr.Get(ctx, child)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	pnode, err := s.nr.Get(ctx, parent)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if !cnode.Owner.Admin {
		if cnode.Owner.ID != accountID && pnode.Owner.ID != accountID {
			return nil, fault.Wrap(errNotAuthorised, fctx.With(ctx))
		}
	}

	// If the target parent is actually a child of the target child, sever this
	// connection before adding the target child to the target parent.
	if parentParent, ok := pnode.Parent.Get(); ok {
		if parentParent.ID == cnode.ID {
			cnode, err = s.nr.Update(ctx, cnode.ID, node.WithChildNodeRemove(xid.ID(pnode.ID)))
			if err != nil {
				return nil, fault.Wrap(err, fctx.With(ctx))
			}
		}
	}

	pnode, err = s.nr.Update(ctx, pnode.ID, node.WithChildNodeAdd(xid.ID(cnode.ID)))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return pnode, nil
}

func (s *service) Sever(ctx context.Context, child datagraph.NodeSlug, parent datagraph.NodeSlug) (*datagraph.Node, error) {
	if child == parent {
		return nil, fault.Wrap(ErrIdenticalParentChild, fctx.With(ctx))
	}

	accountID, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	n, err := s.nr.Get(ctx, child)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	pclus, err := s.nr.Get(ctx, parent)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if !n.Owner.Admin {
		if n.Owner.ID != accountID && pclus.Owner.ID != accountID {
			return nil, fault.Wrap(errNotAuthorised, fctx.With(ctx))
		}
	}

	pclus, err = s.nr.Update(ctx, pclus.ID, node.WithChildNodeRemove(xid.ID(n.ID)))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return pclus, nil
}
