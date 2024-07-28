package node_visibility

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"

	"github.com/Southclaws/storyden/app/resources/account/account_querier"
	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/library/node_children"
	"github.com/Southclaws/storyden/app/resources/visibility"
	"github.com/Southclaws/storyden/app/services/authentication/session"
)

var errNotAuthorised = fault.Wrap(fault.New("not authorised"), ftag.With(ftag.PermissionDenied))

type Controller struct {
	accountQuery account_querier.Querier
	nr           library.Repository
	nc           node_children.Repository
}

func New(
	accountQuery account_querier.Querier,
	nr library.Repository,
	nc node_children.Repository,
) *Controller {
	return &Controller{
		accountQuery: accountQuery,
		nr:           nr,
		nc:           nc,
	}
}

func (m *Controller) ChangeVisibility(ctx context.Context, slug library.NodeSlug, vis visibility.Visibility) (*library.Node, error) {
	accountID, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	acc, err := m.accountQuery.GetByID(ctx, accountID)
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

	n, err = m.nr.Update(ctx, n.ID, library.WithVisibility(vis))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if vis == visibility.VisibilityPublished {
		// TODO: Emit events, send notifications, etc.
	}

	return n, nil
}
