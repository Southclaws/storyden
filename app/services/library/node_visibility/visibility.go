package node_visibility

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"

	"github.com/Southclaws/storyden/app/resources/account/account_querier"
	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/library/node_children"
	"github.com/Southclaws/storyden/app/resources/library/node_querier"
	"github.com/Southclaws/storyden/app/resources/library/node_writer"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/resources/visibility"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
	"github.com/Southclaws/storyden/lib/plugin/rpc"
)

var errNotAuthorised = fault.Wrap(fault.New("not authorised"), ftag.With(ftag.PermissionDenied))

type Controller struct {
	accountQuery *account_querier.Querier
	nodeQuerier  *node_querier.Querier
	nodeWriter   *node_writer.Writer
	nc           *node_children.Writer
	bus          *pubsub.Bus
}

func New(
	accountQuery *account_querier.Querier,
	nodeQuerier *node_querier.Querier,
	nodeWriter *node_writer.Writer,
	nc *node_children.Writer,
	bus *pubsub.Bus,
) *Controller {
	return &Controller{
		accountQuery: accountQuery,
		nodeQuerier:  nodeQuerier,
		nodeWriter:   nodeWriter,
		nc:           nc,
		bus:          bus,
	}
}

func (m *Controller) ChangeVisibility(ctx context.Context, qk library.QueryKey, vis visibility.Visibility) (*library.Node, error) {
	accountID, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	acc, err := m.accountQuery.GetByID(ctx, accountID)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	n, err := m.nodeQuerier.Get(ctx, qk)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if err := acc.Roles.Permissions().Authorise(ctx, func() error {
		if n.Owner.ID != accountID {
			return fault.Wrap(errNotAuthorised, fctx.With(ctx))
		}
		return nil
	}, rbac.PermissionManageLibrary); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	oldVisibility := n.Visibility

	n, err = m.nodeWriter.Update(ctx, qk, node_writer.WithVisibility(vis))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	// Emit visibility transition events
	// NOTE: If this changes, remove the node_visibility service and consolidate
	if oldVisibility != vis {
		switch vis {
		case visibility.VisibilityPublished:
			m.bus.Publish(ctx, &rpc.EventNodePublished{
				ID:   library.NodeID(n.Mark.ID()),
				Slug: n.GetSlug(),
			})
		case visibility.VisibilityReview:
			m.bus.Publish(ctx, &rpc.EventNodeSubmittedForReview{
				ID:   library.NodeID(n.Mark.ID()),
				Slug: n.GetSlug(),
			})
		case visibility.VisibilityUnlisted, visibility.VisibilityDraft:
			if oldVisibility == visibility.VisibilityPublished {
				m.bus.Publish(ctx, &rpc.EventNodeUnpublished{
					ID:   library.NodeID(n.Mark.ID()),
					Slug: n.GetSlug(),
				})
			}
		}
	}

	return n, nil
}
