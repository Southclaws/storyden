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

var errNotAuthorised = fault.Wrap(fault.New("not authorised"), ftag.With(ftag.PermissionDenied))

type Graph interface {
	// Move moves a node from either orphan state or belonging to one node
	// to another node essentially setting its parent slug to some/new value.
	Move(ctx context.Context, child datagraph.NodeSlug, parent datagraph.NodeSlug) (*datagraph.Node, error)

	// Sever orphans a node by removing it from its parent to the root level.
	Sever(ctx context.Context, child datagraph.NodeSlug, parent datagraph.NodeSlug) (*datagraph.Node, error)
}

type service struct {
	cr node.Repository
}

func New(cr node.Repository) Graph {
	return &service{cr: cr}
}

func (s *service) Move(ctx context.Context, child datagraph.NodeSlug, parent datagraph.NodeSlug) (*datagraph.Node, error) {
	accountID, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	cclus, err := s.cr.Get(ctx, child)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	pclus, err := s.cr.Get(ctx, parent)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if !cclus.Owner.Admin {
		if cclus.Owner.ID != accountID && pclus.Owner.ID != accountID {
			return nil, fault.Wrap(errNotAuthorised, fctx.With(ctx))
		}
	}

	pclus, err = s.cr.Update(ctx, pclus.ID, node.WithChildNodeAdd(xid.ID(cclus.ID)))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return pclus, nil
}

func (s *service) Sever(ctx context.Context, child datagraph.NodeSlug, parent datagraph.NodeSlug) (*datagraph.Node, error) {
	accountID, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	cclus, err := s.cr.Get(ctx, child)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	pclus, err := s.cr.Get(ctx, parent)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if !cclus.Owner.Admin {
		if cclus.Owner.ID != accountID && pclus.Owner.ID != accountID {
			return nil, fault.Wrap(errNotAuthorised, fctx.With(ctx))
		}
	}

	pclus, err = s.cr.Update(ctx, pclus.ID, node.WithChildNodeRemove(xid.ID(cclus.ID)))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return pclus, nil
}
