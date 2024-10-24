package tag_querier

import (
	"context"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"

	"github.com/Southclaws/storyden/app/resources/tag"
	"github.com/Southclaws/storyden/app/resources/tag/tag_ref"
	"github.com/Southclaws/storyden/internal/ent"
	ent_tag "github.com/Southclaws/storyden/internal/ent/tag"
)

type Querier struct {
	db *ent.Client
}

func New(db *ent.Client) *Querier {
	return &Querier{db}
}

func (q *Querier) List(ctx context.Context) (tag_ref.Tags, error) {
	r, err := q.db.Tag.Query().All(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	tags := dt.Map(r, tag_ref.Map)

	return tags, nil
}

func (q *Querier) Get(ctx context.Context, name tag_ref.Name) (*tag.Tag, error) {
	r, err := q.db.Tag.Query().
		Where(ent_tag.Name(string(name))).
		WithAccounts().
		WithPosts().
		WithNodes().
		First(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	tag, err := tag.Map(r)
	if err != nil {
		return nil, err
	}

	return tag, nil
}
