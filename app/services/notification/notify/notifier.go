package notify

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/opt"
	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/notification"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/message"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
)

type Notifier struct {
	bus *pubsub.Bus
}

func New(bus *pubsub.Bus) *Notifier {
	return &Notifier{bus: bus}
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
