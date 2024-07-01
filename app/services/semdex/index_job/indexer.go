package index_job

import (
	"context"

	"go.uber.org/zap"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/datagraph/node"
	"github.com/Southclaws/storyden/app/resources/mq"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/reply"
	"github.com/Southclaws/storyden/app/services/semdex"
	"github.com/Southclaws/storyden/internal/pubsub"
)

type indexerConsumer struct {
	l *zap.Logger

	replyRepo   reply.Repository
	nodeRepo    node.Repository
	accountRepo account.Repository

	qnode pubsub.Topic[mq.IndexNode]
	qpost pubsub.Topic[mq.IndexPost]

	indexer   semdex.Indexer
	retriever semdex.Retriever
}

func newIndexConsumer(
	l *zap.Logger,

	replyRepo reply.Repository,
	nodeRepo node.Repository,
	accountRepo account.Repository,

	qnode pubsub.Topic[mq.IndexNode],
	qpost pubsub.Topic[mq.IndexPost],
	qprofile pubsub.Topic[mq.IndexProfile],

	indexer semdex.Indexer,
	retriever semdex.Retriever,
) *indexerConsumer {
	return &indexerConsumer{
		l:           l,
		replyRepo:   replyRepo,
		nodeRepo:    nodeRepo,
		accountRepo: accountRepo,
		qnode:       qnode,
		qpost:       qpost,
		indexer:     indexer,
		retriever:   retriever,
	}
}

func (i *indexerConsumer) indexPost(ctx context.Context, id post.ID) error {
	p, err := i.replyRepo.Get(ctx, id)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return i.indexer.Index(ctx, p)
}

func (i *indexerConsumer) indexNode(ctx context.Context, id datagraph.NodeID) error {
	p, err := i.nodeRepo.GetByID(ctx, id)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return i.indexer.Index(ctx, p)
}

func (i *indexerConsumer) indexProfile(ctx context.Context, id account.AccountID) error {
	p, err := i.accountRepo.GetByID(ctx, id)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return i.indexer.Index(ctx, datagraph.ProfileFromAccount(p))
}
