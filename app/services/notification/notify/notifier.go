package notify

import (
	"context"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/notification"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/mq"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
)

type Notifier struct {
	q pubsub.Topic[mq.Notification]
}

func New(q pubsub.Topic[mq.Notification]) *Notifier {
	return &Notifier{q: q}
}

func (n *Notifier) Send(ctx context.Context, targetID account.AccountID, event notification.Event, item *datagraph.Ref) {
	n.q.PublishAndForget(ctx, mq.Notification{
		Event:    event,
		Item:     item,
		TargetID: targetID,
	})
}
