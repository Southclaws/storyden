package post_search

import (
	"context"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"

	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/reply"
	"github.com/Southclaws/storyden/internal/ent"
	post_model "github.com/Southclaws/storyden/internal/ent/post"
)

type database struct {
	db *ent.Client
}

func New(db *ent.Client) Repository {
	return &database{db}
}

func (d *database) Search(ctx context.Context, filters ...Filter) ([]*reply.Reply, error) {
	if len(filters) == 0 {
		return []*reply.Reply{}, nil
	}

	q := d.db.Post.
		Query().
		WithAuthor().
		WithReacts().
		WithTags().
		WithRoot().
		Order(ent.Asc(post_model.FieldCreatedAt))

	for _, fn := range filters {
		fn(q)
	}

	posts, err := q.All(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	transform := func(v *ent.Post) (*reply.Reply, error) {
		// hydrate the thread-specific info here. post.FromModel cannot do this
		// as this info is only available in the context of a thread of posts.
		dto, err := reply.FromModel(v)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}

		if v.First {
			dto.RootThreadMark = v.Slug
			dto.RootPostID = post.ID(v.ID)
		} else {
			dto.RootThreadMark = v.Edges.Root.Slug
			dto.RootPostID = post.ID(v.Edges.Root.ID)
		}

		return dto, nil
	}

	replies, err := dt.MapErr(posts, transform)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return replies, nil
}
