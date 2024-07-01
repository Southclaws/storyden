package reindex

import (
	"context"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/mq"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/services/semdex"
	"github.com/Southclaws/storyden/internal/ent"
	entpost "github.com/Southclaws/storyden/internal/ent/post"
	"github.com/Southclaws/storyden/internal/pubsub"
	"github.com/rs/xid"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

type reindexer struct {
	l *zap.Logger

	ec *ent.Client

	qnode pubsub.Topic[mq.IndexNode]
	qpost pubsub.Topic[mq.IndexPost]

	indexer   semdex.Indexer
	retriever semdex.Retriever
}

func newReindexer(
	l *zap.Logger,

	qnode pubsub.Topic[mq.IndexNode],
	qpost pubsub.Topic[mq.IndexPost],

	indexer semdex.Indexer,
	retriever semdex.Retriever,
) *reindexer {
	// If the indexer is a searcher only, we don't need to reindex anything.
	switch indexer.(type) {
	case *semdex.OnlySearcher:
		return nil
	case *semdex.Empty:
		return nil
	}

	return &reindexer{
		l: l,

		qnode:     qnode,
		qpost:     qpost,
		indexer:   indexer,
		retriever: retriever,
	}
}

func (r *reindexer) reindexAll(ctx context.Context) error {
	indexed, err := r.retriever.GetAll(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	posts := dt.Filter(indexed, func(i *datagraph.NodeReference) bool { return i.Kind == datagraph.KindPost })
	if err := r.reindexPosts(ctx, posts); err != nil {
		return err
	}

	nodes := dt.Filter(indexed, func(i *datagraph.NodeReference) bool { return i.Kind == datagraph.KindNode })
	if err := r.reindexNodes(ctx, nodes); err != nil {
		return err
	}

	return nil
}

func (r *reindexer) reindexNodes(ctx context.Context, indexed []*datagraph.NodeReference) error {
	r.l.Debug("reindexing all unindexed nodes")

	return nil
}

func (r *reindexer) reindexPosts(ctx context.Context, indexed []*datagraph.NodeReference) error {
	posts, err := r.ec.Post.Query().Select(entpost.FieldID).All(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	indexedIDs := dt.Map(indexed, func(i *datagraph.NodeReference) xid.ID { return i.ID })
	postIDs := dt.Map(posts, func(p *ent.Post) xid.ID { return p.ID })

	intersection := lo.Without(postIDs, indexedIDs...)

	r.l.Debug("reindexing all unindexed posts",
		zap.Int("all_posts", len(posts)),
		zap.Int("indexed_posts", len(indexed)),
		zap.Int("unindexed_posts", len(intersection)),
	)

	messages := dt.Map(intersection, func(id xid.ID) mq.IndexPost {
		return mq.IndexPost{ID: post.ID(id)}
	})

	if err := r.qpost.Publish(ctx, messages...); err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}
