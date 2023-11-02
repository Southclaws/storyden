package link

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"

	"github.com/Southclaws/storyden/internal/ent"
)

type database struct {
	db *ent.Client
}

func New(db *ent.Client) Repository {
	return &database{db}
}

func (d *database) Store(ctx context.Context, url, title, description string, opts ...Option) (*Link, error) {
	create := d.db.Link.Create()
	mutate := create.Mutation()

	mutate.SetURL(url)
	mutate.SetTitle(title)
	mutate.SetDescription(description)

	for _, fn := range opts {
		fn(mutate)
	}

	create.OnConflictColumns("url").UpdateNewValues()

	r, err := create.Save(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	link := Map(r)

	return link, nil
}
