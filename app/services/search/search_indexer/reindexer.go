package search_indexer

import (
	"context"
	"log/slog"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/post/reply"
	"github.com/Southclaws/storyden/app/resources/post/thread"
	"github.com/Southclaws/storyden/internal/ent/node"
	ent_node "github.com/Southclaws/storyden/internal/ent/node"
	ent_post "github.com/Southclaws/storyden/internal/ent/post"
)

func (idx *Indexer) ReindexAll(ctx context.Context) error {
	started := time.Now().UTC()
	idx.logger.Info("starting full reindex",
		slog.Int("chunk_size", idx.chunkSize),
	)

	tn, err := idx.reindexThreads(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	rn, err := idx.reindexReplies(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	nn, err := idx.reindexNodes(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	processed := tn + rn + nn

	if processed > 0 {
		idx.logger.Info("reindex completed",
			slog.Duration("duration", time.Since(started)),
			slog.Int("processed", processed),
		)
	} else {
		idx.logger.Info("reindex skipped: nothing to reindex")
	}

	return nil
}

func (idx *Indexer) reindexThreads(ctx context.Context) (int, error) {
	return reindex(ctx, idx, func() ([]datagraph.Item, error) {
		threads, err := idx.db.Post.Query().
			Where(
				ent_post.RootPostIDIsNil(),
				ent_post.VisibilityEQ(ent_post.VisibilityPublished),
				func(s *sql.Selector) {
					s.Where(sql.Or(
						sql.IsNull(ent_post.FieldIndexedAt),
						sql.GT(ent_post.FieldUpdatedAt, sql.Raw(ent_post.FieldIndexedAt)),
					))
				},
			).
			WithTags().
			Order(ent_post.ByUpdatedAt(), ent_post.ByID()).
			Limit(idx.chunkSize).
			All(ctx)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}

		return dt.MapErr(threads, thread.ItemRef)
	}, func(ids []xid.ID, t time.Time) (int, error) {
		i, err := idx.db.Post.Update().
			Where(ent_post.IDIn(ids...)).
			SetIndexedAt(t).
			Save(ctx)
		if err != nil {
			return 0, fault.Wrap(err, fctx.With(ctx))
		}
		return i, nil
	})
}

func (idx *Indexer) reindexReplies(ctx context.Context) (int, error) {
	return reindex(ctx, idx, func() ([]datagraph.Item, error) {
		replies, err := idx.db.Post.Query().
			Where(
				ent_post.RootPostIDNotNil(),
				// NOTE: Prior to version v1.25.12, replies would always be set
				// to visibility draft on creation, as visibility was never used
				// for replies. Since v1.25.12, replies always get visibility
				// published (semantically correct, and we may use it in future
				// for draft replies, etc.) However, to ensure older instances
				// are indexed correctly in 1.25.12+, we do not filter by the
				// visibility. Yet... this will change once replies have a use
				// case for visibility. The hope is, someone won't jump from
				// version < 1.25.11 to whatever that version is that will add
				// draft reply support, and thus removing this filter. Reason
				// being is, we wouldn't want non-published replies to end up
				// in the search index. For now, we leave this commented out and
				// that bridge to be crossed when the time comes.
				// ent_post.VisibilityEQ(ent_post.VisibilityPublished),
				func(s *sql.Selector) {
					s.Where(sql.Or(
						sql.IsNull(ent_post.FieldIndexedAt),
						sql.GT(ent_post.FieldUpdatedAt, sql.Raw(ent_post.FieldIndexedAt)),
					))
				},
			).
			Order(ent_post.ByUpdatedAt(), ent_post.ByID()).
			Limit(idx.chunkSize).
			All(ctx)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}

		return dt.MapErr(replies, reply.ItemRef)
	}, func(ids []xid.ID, t time.Time) (int, error) {
		i, err := idx.db.Post.Update().
			Where(ent_post.IDIn(ids...)).
			SetIndexedAt(t).
			Save(ctx)
		if err != nil {
			return 0, fault.Wrap(err, fctx.With(ctx))
		}
		return i, nil
	})
}

func (idx *Indexer) reindexNodes(ctx context.Context) (int, error) {
	return reindex(ctx, idx, func() ([]datagraph.Item, error) {
		threads, err := idx.db.Node.Query().
			Where(
				node.VisibilityEQ(node.VisibilityPublished),
				func(s *sql.Selector) {
					s.Where(sql.Or(
						sql.IsNull(node.FieldIndexedAt),
						sql.GT(node.FieldUpdatedAt, sql.Raw(node.FieldIndexedAt)),
					))
				},
			).
			Order(node.ByUpdatedAt(), node.ByID()).
			Limit(idx.chunkSize).
			All(ctx)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}

		return dt.MapErr(threads, library.ItemRef)
	}, func(ids []xid.ID, t time.Time) (int, error) {
		i, err := idx.db.Node.Update().
			Where(ent_node.IDIn(ids...)).
			SetIndexedAt(t).
			Save(ctx)
		if err != nil {
			return 0, fault.Wrap(err, fctx.With(ctx))
		}
		return i, nil
	})
}

func reindex[T datagraph.Item](
	ctx context.Context,
	idx *Indexer,
	fetch func() ([]T, error),
	update func([]xid.ID, time.Time) (int, error),
) (int, error) {
	ids := make([]xid.ID, 0, idx.chunkSize)
	processed := 0
	var k datagraph.Kind

	for {
		ids = ids[:0]

		v, err := fetch()
		if err != nil {
			return processed, fault.Wrap(err, fctx.With(ctx))
		}

		if len(v) == 0 {
			break
		}

		for i, item := range v {
			err := idx.searchIndexer.Index(ctx, item)
			if err != nil {
				idx.logger.Error("failed to index item",
					slog.String("kind", item.GetKind().String()),
					slog.String("id", item.GetID().String()),
					slog.String("error", err.Error()),
				)
			}
			ids = append(ids, item.GetID())

			if i == 0 {
				k = item.GetKind()
			}
		}

		// Go 2 seconds forward, to ensure databases with 1 second precision do
		// not round the indexed at time above the updated at time.
		n, err := update(ids, time.Now().Add((2 * time.Second)))
		if err != nil {
			return processed, fault.Wrap(err, fctx.With(ctx))
		}

		processed += n

		idx.logger.Debug("reindexed chunk",
			slog.String("kind", k.String()),
			slog.Int("processed", n),
		)

		if len(v) < idx.chunkSize {
			break
		}
	}

	return processed, nil
}
