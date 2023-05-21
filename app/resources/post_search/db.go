package post_search

import (
	"context"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"

	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/internal/ent"
	post_model "github.com/Southclaws/storyden/internal/ent/post"
)

type database struct {
	db *ent.Client
}

func New(db *ent.Client) Repository {
	return &database{db}
}

func (d *database) Search(ctx context.Context, filters ...Filter) ([]*post.Post, error) {
	if len(filters) == 0 {
		return []*post.Post{}, nil
	}

	q := d.db.Post.
		Query().
		WithAuthor().
		WithReacts().
		WithTags().
		Order(ent.Asc(post_model.FieldCreatedAt))

	for _, fn := range filters {
		fn(q)
	}

	posts, err := q.All(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return dt.Map(posts, post.FromModel), nil
}
