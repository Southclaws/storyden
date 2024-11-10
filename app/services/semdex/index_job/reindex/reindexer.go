package reindex

import (
	"context"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/datagraph/semdex"
	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/mq"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/internal/ent"
	entaccount "github.com/Southclaws/storyden/internal/ent/account"
	entpost "github.com/Southclaws/storyden/internal/ent/post"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
	"github.com/rs/xid"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

type reindexer struct {
	l *zap.Logger

	ec *ent.Client

	qnode    pubsub.Topic[mq.IndexNode]
	qthread  pubsub.Topic[mq.IndexThread]
	qreply   pubsub.Topic[mq.IndexReply]
	qprofile pubsub.Topic[mq.IndexProfile]

	indexer   semdex.Indexer
	retriever semdex.Retriever
}

func newReindexer(
	l *zap.Logger,

	ec *ent.Client,

	qnode pubsub.Topic[mq.IndexNode],
	qthread pubsub.Topic[mq.IndexThread],
	qreply pubsub.Topic[mq.IndexReply],
	qprofile pubsub.Topic[mq.IndexProfile],

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

		ec: ec,

		qnode:    qnode,
		qthread:  qthread,
		qreply:   qreply,
		qprofile: qprofile,

		indexer:   indexer,
		retriever: retriever,
	}
}

func (r *reindexer) reindexAll(ctx context.Context) error {
	indexed, err := r.retriever.GetAll(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	threads := dt.Filter(indexed, func(i *datagraph.Ref) bool { return i.Kind == datagraph.KindThread })
	replies := dt.Filter(indexed, func(i *datagraph.Ref) bool { return i.Kind == datagraph.KindReply })
	nodes := dt.Filter(indexed, func(i *datagraph.Ref) bool { return i.Kind == datagraph.KindNode })
	profiles := dt.Filter(indexed, func(i *datagraph.Ref) bool { return i.Kind == datagraph.KindProfile })

	if err := r.reindexThreads(ctx, threads); err != nil {
		return err
	}

	if err := r.reindexReplies(ctx, replies); err != nil {
		return err
	}

	if err := r.reindexNodes(ctx, nodes); err != nil {
		return err
	}

	if err := r.reindexProfiles(ctx, profiles); err != nil {
		return err
	}

	return nil
}

func (r *reindexer) reindexNodes(ctx context.Context, indexed []*datagraph.Ref) error {
	nodes, err := r.ec.Node.Query().Select(entpost.FieldID).All(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	indexedIDs := dt.Map(indexed, func(i *datagraph.Ref) xid.ID { return i.ID })
	postIDs := dt.Map(nodes, func(p *ent.Node) xid.ID { return p.ID })

	intersection := lo.Without(postIDs, indexedIDs...)

	r.l.Debug("reindexing all unindexed nodes",
		zap.Int("all_nodes", len(nodes)),
		zap.Int("indexed_nodes", len(indexed)),
		zap.Int("unindexed_nodes", len(intersection)),
	)

	messages := dt.Map(intersection, func(id xid.ID) mq.IndexNode {
		return mq.IndexNode{ID: library.NodeID(id)}
	})

	if err := r.qnode.Publish(ctx, messages...); err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}

func (r *reindexer) reindexThreads(ctx context.Context, indexed []*datagraph.Ref) error {
	posts, err := r.ec.Post.Query().Select(entpost.FieldID).Where(entpost.RootPostIDIsNil()).All(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	indexedIDs := dt.Map(indexed, func(i *datagraph.Ref) xid.ID { return i.ID })
	postIDs := dt.Map(posts, func(p *ent.Post) xid.ID { return p.ID })

	intersection := lo.Without(postIDs, indexedIDs...)

	r.l.Debug("reindexing all unindexed threads",
		zap.Int("all_posts", len(posts)),
		zap.Int("indexed_posts", len(indexed)),
		zap.Int("unindexed_posts", len(intersection)),
	)

	messages := dt.Map(intersection, func(id xid.ID) mq.IndexThread {
		return mq.IndexThread{ID: post.ID(id)}
	})

	if err := r.qthread.Publish(ctx, messages...); err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}

func (r *reindexer) reindexReplies(ctx context.Context, indexed []*datagraph.Ref) error {
	posts, err := r.ec.Post.Query().Select(entpost.FieldID).Where(entpost.RootPostIDNotNil()).All(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	indexedIDs := dt.Map(indexed, func(i *datagraph.Ref) xid.ID { return i.ID })
	postIDs := dt.Map(posts, func(p *ent.Post) xid.ID { return p.ID })

	intersection := lo.Without(postIDs, indexedIDs...)

	r.l.Debug("reindexing all unindexed replies",
		zap.Int("all_posts", len(posts)),
		zap.Int("indexed_posts", len(indexed)),
		zap.Int("unindexed_posts", len(intersection)),
	)

	messages := dt.Map(intersection, func(id xid.ID) mq.IndexReply {
		return mq.IndexReply{ID: post.ID(id)}
	})

	if err := r.qreply.Publish(ctx, messages...); err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}

func (r *reindexer) reindexProfiles(ctx context.Context, indexed []*datagraph.Ref) error {
	profiles, err := r.ec.Account.Query().Select(entaccount.FieldID).All(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	indexedIDs := dt.Map(indexed, func(i *datagraph.Ref) xid.ID { return i.ID })
	accountIDs := dt.Map(profiles, func(p *ent.Account) xid.ID { return p.ID })

	intersection := lo.Without(accountIDs, indexedIDs...)

	r.l.Debug("reindexing all unindexed profiles",
		zap.Int("all_profiles", len(profiles)),
		zap.Int("indexed_profiles", len(indexed)),
		zap.Int("unindexed_profiles", len(intersection)),
	)

	messages := dt.Map(intersection, func(id xid.ID) mq.IndexProfile {
		return mq.IndexProfile{ID: account.AccountID(id)}
	})

	if err := r.qprofile.Publish(ctx, messages...); err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}
