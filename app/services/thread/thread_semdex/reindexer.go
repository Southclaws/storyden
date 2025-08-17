package thread_semdex

import (
	"context"
	"log/slog"
	"time"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/rs/xid"
	"github.com/samber/lo"

	"github.com/Southclaws/storyden/internal/ent"
	ent_post "github.com/Southclaws/storyden/internal/ent/post"
)

func (r *semdexer) schedule(ctx context.Context, schedule time.Duration, reindexThreshold time.Duration, reindexChunk int) {
	for range time.NewTicker(schedule).C {
		err := r.reindex(ctx, reindexThreshold, reindexChunk)
		if err != nil {
			r.logger.Error("failed to run reindex job", slog.String("error", err.Error()))
		}
	}
}

func (r *semdexer) reindex(ctx context.Context, reindexThreshold time.Duration, reindexChunk int) error {
	updated, deleted, err := r.gatherTargets(ctx, reindexThreshold, reindexChunk)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	r.logger.Debug("reindexing threads",
		slog.Int("updated", len(updated)),
		slog.Int("deleted", len(deleted)),
	)

	// toIndex := dt.Map(updated, func(id xid.ID) mq.IndexThread {
	// 	return mq.IndexThread{ID: post.ID(id)}
	// })
	// toDelete := dt.Map(deleted, func(id xid.ID) mq.DeleteThread {
	// 	return mq.DeleteThread{ID: post.ID(id)}
	// })

	// if err := r.indexQueue.Publish(ctx, toIndex...); err != nil {
	// 	return fault.Wrap(err, fctx.With(ctx))
	// }

	// if err := r.deleteQueue.Publish(ctx, toDelete...); err != nil {
	// 	return fault.Wrap(err, fctx.With(ctx))
	// }

	return nil
}

func (r *semdexer) gatherTargets(ctx context.Context, reindexThreshold time.Duration, reindexChunk int) ([]xid.ID, []xid.ID, error) {
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
		).
		Order(ent.Desc(ent_post.FieldCreatedAt)).
		Limit(reindexChunk).
		All(ctx)
	if err != nil {
		return nil, nil, fault.Wrap(err, fctx.With(ctx))
	}

	keepIDs, discardIDs := r.partition(threads)

	return keepIDs, discardIDs, nil
}

func (r *semdexer) partition(threads []*ent.Post) ([]xid.ID, []xid.ID) {
	keep, discard := lo.FilterReject(threads, func(p *ent.Post, _ int) bool {
		return p.Visibility == ent_post.VisibilityPublished && p.DeletedAt == nil
	})

	keepIDs := dt.Map(keep, func(p *ent.Post) xid.ID { return p.ID })
	discardIDs := dt.Map(discard, func(p *ent.Post) xid.ID { return p.ID })

	return keepIDs, discardIDs
}
