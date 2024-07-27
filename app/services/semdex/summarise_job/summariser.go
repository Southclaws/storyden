package summarise_job

import (
	"context"

	"go.uber.org/zap"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/mq"
	"github.com/Southclaws/storyden/app/services/semdex"
	"github.com/Southclaws/storyden/internal/pubsub"
)

type summariseConsumer struct {
	l *zap.Logger

	nodeRepo library.Repository

	qnode pubsub.Topic[mq.SummariseNode]

	summariser semdex.Summariser
}

func newSummariseConsumer(
	l *zap.Logger,

	nodeRepo library.Repository,

	qnode pubsub.Topic[mq.SummariseNode],

	summariser semdex.Summariser,
) *summariseConsumer {
	return &summariseConsumer{
		l: l,

		nodeRepo: nodeRepo,

		qnode:      qnode,
		summariser: summariser,
	}
}

func (i *summariseConsumer) summariseNode(ctx context.Context, id library.NodeID) error {
	summary, err := i.summariser.Summarise(ctx, &library.Node{ID: id})
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	_, err = i.nodeRepo.Update(ctx, id, library.WithDescription(summary))
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}
