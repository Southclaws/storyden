package index_job

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"go.uber.org/zap"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/account_querier"
	"github.com/Southclaws/storyden/app/resources/library/node_querier"
	"github.com/Southclaws/storyden/app/resources/mq"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/post/reply"
	"github.com/Southclaws/storyden/app/resources/post/thread"
	"github.com/Southclaws/storyden/app/resources/profile"
	"github.com/Southclaws/storyden/app/services/semdex"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
)

type indexerConsumer struct {
	l *zap.Logger

	threadRepo   thread.Repository
	replyRepo    reply.Repository
	nodeQuerier  *node_querier.Querier
	accountQuery *account_querier.Querier

	qnode   pubsub.Topic[mq.IndexNode]
	qthread pubsub.Topic[mq.IndexThread]
	qreply  pubsub.Topic[mq.IndexReply]

	indexer semdex.Mutator
}

func newIndexConsumer(
	l *zap.Logger,

	threadRepo thread.Repository,
	replyRepo reply.Repository,
	nodeQuerier *node_querier.Querier,
	accountQuery *account_querier.Querier,

	qnode pubsub.Topic[mq.IndexNode],
	qthread pubsub.Topic[mq.IndexThread],
	qreply pubsub.Topic[mq.IndexReply],
	qprofile pubsub.Topic[mq.IndexProfile],

	indexer semdex.Mutator,
) *indexerConsumer {
	return &indexerConsumer{
		l:            l,
		threadRepo:   threadRepo,
		replyRepo:    replyRepo,
		nodeQuerier:  nodeQuerier,
		accountQuery: accountQuery,
		qnode:        qnode,

		qthread: qthread,
		qreply:  qreply,
		indexer: indexer,
	}
}

func (i *indexerConsumer) indexReply(ctx context.Context, id post.ID) error {
	p, err := i.replyRepo.Get(ctx, id)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	_, err = i.indexer.Index(ctx, p)
	return err
}

func (i *indexerConsumer) indexProfile(ctx context.Context, id account.AccountID) error {
	p, err := i.accountQuery.GetByID(ctx, id)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	_, err = i.indexer.Index(ctx, profile.ProfileFromAccount(p))
	return err
}
