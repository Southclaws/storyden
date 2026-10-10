package node_querier

import (
	"context"
	"time"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/ent/node"
)

type PublicIndexEntry struct {
	Slug      string
	UpdatedAt time.Time
}

type PublicIndexPage struct {
	Entries    []PublicIndexEntry
	Page       int
	TotalPages int
	PageSize   int
}

// ListPublishedIndex selects only fields needed by an anonymous sitemap.
func (q *Querier) ListPublishedIndex(ctx context.Context, page, size int) (PublicIndexPage, error) {
	if page < 1 || size < 1 || size > 500 {
		return PublicIndexPage{}, fault.New("invalid public index page")
	}

	base := q.db.Node.Query().Where(node.VisibilityEQ(node.VisibilityPublished))
	total, err := base.Clone().Count(ctx)
	if err != nil {
		return PublicIndexPage{}, fault.Wrap(err, fctx.With(ctx))
	}

	rows, err := base.Order(ent.Desc(node.FieldUpdatedAt), ent.Desc(node.FieldID)).
		Limit(size).Offset((page-1)*size).
		Select(node.FieldSlug, node.FieldUpdatedAt).
		All(ctx)
	if err != nil {
		return PublicIndexPage{}, fault.Wrap(err, fctx.With(ctx))
	}

	entries := make([]PublicIndexEntry, len(rows))
	for i, row := range rows {
		entries[i] = PublicIndexEntry{Slug: row.Slug, UpdatedAt: row.UpdatedAt}
	}
	return PublicIndexPage{Entries: entries, Page: page, TotalPages: (total + size - 1) / size, PageSize: size}, nil
}
