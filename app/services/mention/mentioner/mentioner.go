package mentioner

import (
	"context"

	"go.uber.org/zap"

	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/mq"
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
	for _, i := range items {
		err := n.q.Publish(ctx, mq.Mention{
			Source: source,
			Item:   *i,
		})
		if err != nil {
			n.l.Error("failed to publish mention message", zap.Error(err))
		}
	}
}
