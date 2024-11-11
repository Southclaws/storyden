package node_semdex

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
	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/mq"
	"github.com/Southclaws/storyden/internal/ent"
	ent_node "github.com/Southclaws/storyden/internal/ent/node"
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
	nodes, err := r.db.Node.Query().
		Select(
			ent_node.FieldID,
			ent_node.FieldVisibility,
			ent_node.FieldDeletedAt,
		).
		Where(
			ent_node.Or(
				ent_node.IndexedAtIsNil(),
				ent_node.IndexedAtLT(time.Now().Add(-reindexThreshold)),
			),
		).
		Limit(reindexChunk).
		All(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	keep, discard := lo.FilterReject(nodes, func(p *ent.Node, _ int) bool {
		return p.Visibility == ent_node.VisibilityPublished && p.DeletedAt == nil
	})

	keepIDs := dt.Map(keep, func(p *ent.Node) xid.ID { return p.ID })
	discardIDs := dt.Map(discard, func(p *ent.Node) xid.ID { return p.ID })

	indexed, err := r.retriever.GetMany(ctx, uint(reindexChunk), keepIDs...)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	indexedIDs := dt.Map(indexed, func(p *datagraph.Ref) xid.ID { return p.ID })

	updated := diff(keepIDs, indexedIDs)
	deleted := lo.Intersect(indexedIDs, discardIDs)

	r.logger.Debug("reindexing nodes",
		zap.Int("all", len(nodes)),
		zap.Int("updated", len(updated)),
		zap.Int("deleted", len(deleted)),
	)

	toIndex := dt.Map(updated, func(id xid.ID) mq.IndexNode {
		return mq.IndexNode{ID: library.NodeID(id)}
	})
	toDelete := dt.Map(deleted, func(id xid.ID) mq.DeleteNode {
		return mq.DeleteNode{ID: library.NodeID(id)}
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
