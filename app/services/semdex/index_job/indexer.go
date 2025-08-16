package index_job

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/library/node_querier"
	"github.com/Southclaws/storyden/app/resources/post/reply"
	"github.com/Southclaws/storyden/app/resources/post/thread"
	"github.com/Southclaws/storyden/app/resources/profile/profile_querier"
	"github.com/Southclaws/storyden/app/services/semdex"
)

type indexerConsumer struct {
	threadRepo   thread.Repository
	replyRepo    reply.Repository
	nodeQuerier  *node_querier.Querier
	profileQuery *profile_querier.Querier

	indexer semdex.Mutator
}

func newIndexConsumer(
	threadRepo thread.Repository,
	replyRepo reply.Repository,
	nodeQuerier *node_querier.Querier,
	profileQuery *profile_querier.Querier,

	indexer semdex.Mutator,
) *indexerConsumer {
	return &indexerConsumer{
		threadRepo:   threadRepo,
		replyRepo:    replyRepo,
		nodeQuerier:  nodeQuerier,
		profileQuery: profileQuery,

		indexer: indexer,
	}
}

func (i *indexerConsumer) indexProfile(ctx context.Context, id account.AccountID) error {
	p, err := i.profileQuery.GetByID(ctx, id)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	_, err = i.indexer.Index(ctx, p)
	return err
}
