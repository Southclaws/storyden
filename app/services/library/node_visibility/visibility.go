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
	"github.com/Southclaws/storyden/app/resources/mq"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/resources/visibility"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
)

var errNotAuthorised = fault.Wrap(fault.New("not authorised"), ftag.With(ftag.PermissionDenied))

type Controller struct {
	accountQuery *account_querier.Querier
	nodeQuerier  *node_querier.Querier
	nodeWriter   *node_writer.Writer
	nc           node_children.Repository
	indexQueue   pubsub.Topic[mq.IndexNode]
}

func New(
	accountQuery *account_querier.Querier,
	nodeQuerier *node_querier.Querier,
	nodeWriter *node_writer.Writer,
	nc node_children.Repository,
	indexQueue pubsub.Topic[mq.IndexNode],
) *Controller {
	return &Controller{
		accountQuery: accountQuery,
		nodeQuerier:  nodeQuerier,
		nodeWriter:   nodeWriter,
		nc:           nc,
		indexQueue:   indexQueue,
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

	n, err = m.nodeWriter.Update(ctx, qk, node_writer.WithVisibility(vis))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if vis == visibility.VisibilityPublished {
		if err := m.indexQueue.Publish(ctx, mq.IndexNode{ID: library.NodeID(n.Mark.ID())}); err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}
	}

	return n, nil
}
