package thread_semdex

import (
	"context"
	"time"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/rs/xid"
	"github.com/samber/lo"
	"go.uber.org/zap"

	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/mq"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/internal/ent"
	ent_post "github.com/Southclaws/storyden/internal/ent/post"
)

func (r *semdexer) schedule(ctx context.Context, schedule time.Duration, reindexThreshold time.Duration, reindexChunk int) {
	for range time.NewTicker(schedule).C {
		err := r.reindex(ctx, reindexThreshold, reindexChunk)
		if err != nil {
			r.logger.Error("failed to run reindex job", zap.Error(err))
		}
	}
}

func (r *semdexer) reindex(ctx context.Context, reindexThreshold time.Duration, reindexChunk int) error {
	threads, err := r.db.Post.Query().
		Select(
			ent_post.FieldID,
			ent_post.FieldVisibility,
			ent_post.FieldDeletedAt,
		).
		Where(
			ent_post.Or(
				ent_post.IndexedAtIsNil(),
				ent_post.IndexedAtLT(time.Now().Add(-reindexThreshold)),
			),
			ent_post.RootPostIDIsNil(),
		).
		Limit(reindexChunk).
		All(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	keep, discard := lo.FilterReject(threads, func(p *ent.Post, _ int) bool {
		return p.Visibility == ent_post.VisibilityPublished && p.DeletedAt == nil
	})

	keepIDs := dt.Map(keep, func(p *ent.Post) xid.ID { return p.ID })
	discardIDs := dt.Map(discard, func(p *ent.Post) xid.ID { return p.ID })

	indexed, err := r.retriever.GetMany(ctx, uint(reindexChunk), keepIDs...)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	indexedIDs := dt.Map(indexed, func(p *datagraph.Ref) xid.ID { return p.ID })

	updated := diff(keepIDs, indexedIDs)
	deleted := lo.Intersect(indexedIDs, discardIDs)

	r.logger.Debug("reindexing threads",
		zap.Int("all", len(threads)),
		zap.Int("updated", len(updated)),
		zap.Int("deleted", len(deleted)),
	)

	toIndex := dt.Map(updated, func(id xid.ID) mq.IndexThread {
		return mq.IndexThread{ID: post.ID(id)}
	})
	toDelete := dt.Map(deleted, func(id xid.ID) mq.DeleteThread {
		return mq.DeleteThread{ID: post.ID(id)}
	})

	if err := r.indexQueue.Publish(ctx, toIndex...); err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	if err := r.deleteQueue.Publish(ctx, toDelete...); err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}

func diff(targets []xid.ID, indexed []xid.ID) []xid.ID {
	_, ids := lo.Difference(indexed, targets)
	return ids
}
