package notify

import (
	"context"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/notification"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/mq"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
	"go.uber.org/zap"
)

type Notifier struct {
	l *zap.Logger
	q pubsub.Topic[mq.Notification]
}

func New(l *zap.Logger, q pubsub.Topic[mq.Notification]) *Notifier {
	return &Notifier{l: l, q: q}
}

func (n *Notifier) Send(ctx context.Context, targetID account.AccountID, event notification.Event, item *datagraph.Ref) {
	err := n.q.Publish(ctx, mq.Notification{
		Event:    event,
		Item:     item,
		TargetID: targetID,
	})
	if err != nil {
		n.l.Error("failed to publish notification message", zap.Error(err))
	}
}
