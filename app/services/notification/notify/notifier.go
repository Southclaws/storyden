package notify

import (
	"context"

	"github.com/Southclaws/opt"
	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/notification"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/mq"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub/event"
)

type Notifier struct {
	bus *event.Bus
}

func New(bus *event.Bus) *Notifier {
	return &Notifier{bus: bus}
}

func (n *Notifier) Send(ctx context.Context, targetID account.AccountID, sourceID opt.Optional[account.AccountID], event notification.Event, item *datagraph.Ref) {
	n.bus.SendCommand(ctx, &mq.CommandSendNotification{
		Event:    event,
		Item:     item,
		TargetID: targetID,
		SourceID: sourceID,
	})
}
