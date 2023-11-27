package link_graph

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"

	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/ent/link"
)

type database struct {
	db *ent.Client
}

func New(db *ent.Client) Repository {
	return &database{db}
}

func (d *database) Get(ctx context.Context, slug string) (*WithRefs, error) {
	query := d.db.Link.Query().
		Where(link.SlugEqualFold(slug)).
		WithAssets().
		WithPosts(func(pq *ent.PostQuery) {
			pq.WithAuthor()
		}).
		WithClusters().
		WithItems()

	r, err := query.First(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	link, err := Map(r)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return link, nil
}
