package mentioner

import (
	"context"

	"go.uber.org/zap"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/mq"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
)

type Mentioner struct {
	l *zap.Logger
	q pubsub.Topic[mq.Mention]
}

func New(l *zap.Logger, q pubsub.Topic[mq.Mention]) *Mentioner {
	return &Mentioner{l: l, q: q}
}

func (n *Mentioner) Send(ctx context.Context, source datagraph.Ref, items ...*datagraph.Ref) {
	sender, err := session.GetAccountID(ctx)
	if err != nil {
		n.l.Warn("cannot send notification without source session", zap.Error(err))
	}

	for _, i := range items {
		if i.Kind == datagraph.KindProfile && sender == account.AccountID(i.ID) {
			// Skip self-mentions
			continue
		}

		err := n.q.Publish(ctx, mq.Mention{
			Source: source,
			Item:   *i,
		})
		if err != nil {
			n.l.Error("failed to publish mention message", zap.Error(err))
		}
	}
}
