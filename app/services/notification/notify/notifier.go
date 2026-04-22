package notify

import (
	"context"
	"log/slog"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/opt"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/notification"
	"github.com/Southclaws/storyden/app/resources/account/notification/notify_writer"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/message"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
)

type Notifier struct {
	bus          *pubsub.Bus
	notifyWriter *notify_writer.Writer
}

func New(
	ctx context.Context,
	lc fx.Lifecycle,
	logger *slog.Logger,
	bus *pubsub.Bus,
	notifyWriter *notify_writer.Writer,
) *Notifier {
	n := &Notifier{
		bus:          bus,
		notifyWriter: notifyWriter,
	}

	lc.Append(fx.StartHook(func(hctx context.Context) error {
		_, err := pubsub.SubscribeCommand(ctx, bus, "notify_job.send_notification", func(ctx context.Context, cmd *message.CommandSendNotification) error {
			if err := n.notify(ctx, cmd.TargetID, cmd.SourceID, cmd.Event, cmd.Item); err != nil {
				logger.Error("failed to notify", slog.String("error", err.Error()))
				return err
			}
			return nil
		})

		return err
	}))

	return n
}

func (n *Notifier) Send(ctx context.Context, targetID account.AccountID, sourceID opt.Optional[account.AccountID], event notification.Event, item *datagraph.Ref) error {
	if err := n.bus.SendCommand(ctx, &message.CommandSendNotification{
		Event:    event,
		Item:     item,
		TargetID: targetID,
		SourceID: sourceID,
	}); err != nil {
		return fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to publish notification command"))
	}
	return nil
}

func (s *Notifier) notify(ctx context.Context,
	targetID account.AccountID,
	sourceID opt.Optional[account.AccountID],
	event notification.Event,
	item *datagraph.Ref,
) error {
	itemref := opt.Map(opt.NewPtr(item), func(i datagraph.Ref) datagraph.ItemRef {
		return &i
	})

	_, err := s.notifyWriter.Notification(ctx, targetID, event, itemref, sourceID)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}
