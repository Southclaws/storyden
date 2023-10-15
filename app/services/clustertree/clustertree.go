package clustertree

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/cluster"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/services/authentication"
)

var errNotAuthorised = fault.Wrap(fault.New("not authorised"), ftag.With(ftag.PermissionDenied))

type Graph interface {
	// Move moves a cluster from either orphan state or belonging to one cluster
	// to another cluster essentially setting its parent slug to some/new value.
	Move(ctx context.Context, child datagraph.ClusterSlug, parent datagraph.ClusterSlug) (*datagraph.Cluster, error)

	// Sever orphans a cluster by removing it from its parent to the root level.
	Sever(ctx context.Context, child datagraph.ClusterSlug, parent datagraph.ClusterSlug) (*datagraph.Cluster, error)
}

type service struct {
	cr cluster.Repository
}

func New(cr cluster.Repository) Graph {
	return &service{cr: cr}
}

func (s *service) Move(ctx context.Context, child datagraph.ClusterSlug, parent datagraph.ClusterSlug) (*datagraph.Cluster, error) {
	accountID, err := authentication.GetAccountID(ctx)
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

	pclus, err = s.cr.Update(ctx, pclus.ID, cluster.WithChildClusterAdd(xid.ID(cclus.ID)))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return pclus, nil
}

func (s *service) Sever(ctx context.Context, child datagraph.ClusterSlug, parent datagraph.ClusterSlug) (*datagraph.Cluster, error) {
	accountID, err := authentication.GetAccountID(ctx)
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

	pclus, err = s.cr.Update(ctx, pclus.ID, cluster.WithChildClusterRemove(xid.ID(cclus.ID)))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return pclus, nil
}
