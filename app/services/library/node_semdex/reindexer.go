package node_semdex

import (
	"context"
	"log/slog"
	"time"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/rs/xid"
	"github.com/samber/lo"

	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/message"
	"github.com/Southclaws/storyden/internal/ent"
	ent_node "github.com/Southclaws/storyden/internal/ent/node"
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

	updated := keepIDs
	deleted := discardIDs

	r.logger.Debug("reindexing nodes",
		slog.Int("all", len(nodes)),
		slog.Int("updated", len(updated)),
		slog.Int("deleted", len(deleted)),
	)

	toIndex := dt.Map(updated, func(id xid.ID) any {
		return message.CommandNodeIndex{ID: library.NodeID(id)}
	})
	toDelete := dt.Map(deleted, func(id xid.ID) any {
		return message.CommandNodeDeindex{ID: library.NodeID(id)}
	})

	if err := r.bus.MustPublishMany(ctx, toIndex...); err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	if err := r.bus.MustPublishMany(ctx, toDelete...); err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}
