package node_visibility

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/datagraph/node"
	"github.com/Southclaws/storyden/app/resources/datagraph/node_children"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/services/authentication/session"
)

var errNotAuthorised = fault.Wrap(fault.New("not authorised"), ftag.With(ftag.PermissionDenied))

type Controller struct {
	ar account.Repository
	nr node.Repository
	nc node_children.Repository
}

func New(
	ar account.Repository,
	nr node.Repository,
	nc node_children.Repository,
) *Controller {
	return &Controller{
		ar: ar,
		nr: nr,
		nc: nc,
	}
}

func (m *Controller) ChangeVisibility(ctx context.Context, slug datagraph.NodeSlug, visibility post.Visibility) (*datagraph.Node, error) {
	accountID, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	acc, err := m.ar.GetByID(ctx, accountID)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	n, err := m.nr.Get(ctx, slug)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if !acc.Admin {
		if n.Owner.ID != accountID {
			return nil, fault.Wrap(errNotAuthorised, fctx.With(ctx))
		}
	}

	n, err = m.nr.Update(ctx, n.ID, node.WithVisibility(visibility))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if visibility == post.VisibilityPublished {
		// TODO: Emit events, send notifications, etc.
	}

	return n, nil
}
