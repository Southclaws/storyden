package node_versioning

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/notification"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/library/node_querier"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/app/services/notification/notify"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
	"github.com/Southclaws/storyden/lib/plugin/rpc"
)

func (s *Service) subscribeNotifications(
	ctx context.Context,
	lc fx.Lifecycle,
	nodeQuerier *node_querier.Querier,
	notifier *notify.Notifier,
) {
	lc.Append(fx.StartHook(func(hctx context.Context) error {
		if _, err := pubsub.Subscribe(ctx, s.bus, "node_versioning.notify.draft_created", func(ctx context.Context, evt *rpc.EventNodeVersionDraftCreated) error {
			return notifyDraftCreated(ctx, notifier, nodeQuerier, evt)
		}); err != nil {
			return err
		}

		if _, err := pubsub.Subscribe(ctx, s.bus, "node_versioning.notify.draft_applied", func(ctx context.Context, evt *rpc.EventNodeVersionDraftApplied) error {
			return notifyDraftApplied(ctx, notifier, evt)
		}); err != nil {
			return err
		}

		if _, err := pubsub.Subscribe(ctx, s.bus, "node_versioning.notify.draft_deleted", func(ctx context.Context, evt *rpc.EventNodeVersionDraftDeleted) error {
			return notifyDraftDeleted(ctx, notifier, evt)
		}); err != nil {
			return err
		}

		return nil
	}))
}

func notifyDraftCreated(
	ctx context.Context,
	notifier *notify.Notifier,
	nodeQuerier *node_querier.Querier,
	evt *rpc.EventNodeVersionDraftCreated,
) error {
	n, err := nodeQuerier.Probe(ctx, evt.NodeID)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	if n.Owner.ID == evt.AuthorID {
		return nil
	}

	return sendNodeVersionNotification(ctx, notifier, n.Owner.ID, evt.AuthorID, notification.EventNodeVersionCreated, evt.NodeID)
}

func notifyDraftApplied(
	ctx context.Context,
	notifier *notify.Notifier,
	evt *rpc.EventNodeVersionDraftApplied,
) error {
	actorID, err := session.GetAccountID(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	if evt.AuthorID == actorID {
		return nil
	}

	return sendNodeVersionNotification(ctx, notifier, evt.AuthorID, actorID, notification.EventNodeVersionApplied, evt.NodeID)
}

func notifyDraftDeleted(
	ctx context.Context,
	notifier *notify.Notifier,
	evt *rpc.EventNodeVersionDraftDeleted,
) error {
	actorID, err := session.GetAccountID(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	if evt.AuthorID == actorID {
		return nil
	}

	return sendNodeVersionNotification(ctx, notifier, evt.AuthorID, actorID, notification.EventNodeVersionDeleted, evt.NodeID)
}

func sendNodeVersionNotification(
	ctx context.Context,
	notifier *notify.Notifier,
	targetID account.AccountID,
	sourceID account.AccountID,
	event notification.Event,
	nodeID library.NodeID,
) error {
	if err := notifier.Send(ctx,
		targetID,
		opt.New(sourceID),
		event,
		&datagraph.Ref{
			ID:   xid.ID(nodeID),
			Kind: datagraph.KindNode,
		},
	); err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}
