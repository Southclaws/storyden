package cluster_visibility

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/datagraph/cluster"
	"github.com/Southclaws/storyden/app/resources/datagraph/cluster_children"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/services/authentication/session"
)

var errNotAuthorised = fault.Wrap(fault.New("not authorised"), ftag.With(ftag.PermissionDenied))

type Controller struct {
	ar account.Repository
	cr cluster.Repository
	cc cluster_children.Repository
}

func New(
	ar account.Repository,
	cr cluster.Repository,
	cc cluster_children.Repository,
) *Controller {
	return &Controller{
		ar: ar,
		cr: cr,
		cc: cc,
	}
}

func (m *Controller) ChangeVisibility(ctx context.Context, slug datagraph.ClusterSlug, visibility post.Visibility) (*datagraph.Cluster, error) {
	accountID, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	acc, err := m.ar.GetByID(ctx, accountID)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	clus, err := m.cr.Get(ctx, slug)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if !acc.Admin {
		if clus.Owner.ID != accountID {
			return nil, fault.Wrap(errNotAuthorised, fctx.With(ctx))
		}
	}

	clus, err = m.cr.Update(ctx, clus.ID, cluster.WithVisibility(visibility))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if visibility == post.VisibilityPublished {
		// TODO: Emit events, send notifications, etc.
	}

	return clus, nil
}
