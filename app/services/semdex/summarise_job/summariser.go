package summarise_job

import (
	"context"

	"go.uber.org/zap"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/storyden/app/resources/datagraph/semdex"
	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/library/node_writer"
	"github.com/Southclaws/storyden/app/resources/mark"
	"github.com/Southclaws/storyden/app/resources/mq"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
	"github.com/rs/xid"
)

type summariseConsumer struct {
	l *zap.Logger

	nodeWriter *node_writer.Writer

	qnode pubsub.Topic[mq.SummariseNode]

	summariser semdex.Summariser
}

func newSummariseConsumer(
	l *zap.Logger,

	nodeWriter *node_writer.Writer,

	qnode pubsub.Topic[mq.SummariseNode],

	summariser semdex.Summariser,
) *summariseConsumer {
	return &summariseConsumer{
		l: l,

		nodeWriter: nodeWriter,

		qnode:      qnode,
		summariser: summariser,
	}
}

func (i *summariseConsumer) summariseNode(ctx context.Context, id library.NodeID) error {
	summary, err := i.summariser.Summarise(ctx, &library.Node{Mark: library.NewMark(xid.ID(id), "")})
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	qk := library.QueryKey{mark.NewQueryKeyID(xid.ID(id))}
	_, err = i.nodeWriter.Update(ctx, qk, node_writer.WithDescription(summary))
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}
