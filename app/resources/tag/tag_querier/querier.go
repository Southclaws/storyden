package tag_querier

import (
	"context"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/jmoiron/sqlx"

	"github.com/Southclaws/storyden/app/resources/tag"
	"github.com/Southclaws/storyden/app/resources/tag/tag_ref"
	"github.com/Southclaws/storyden/internal/ent"
	ent_tag "github.com/Southclaws/storyden/internal/ent/tag"
)

type Querier struct {
	db  *ent.Client
	raw *sqlx.DB
}

func New(db *ent.Client, raw *sqlx.DB) *Querier {
	return &Querier{db, raw}
}

func (q *Querier) List(ctx context.Context) (tag_ref.Tags, error) {
	r, err := q.db.Tag.Query().All(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	tags := dt.Map(r, tag_ref.Map(nil))

	return tags, nil
}

const tagItemsCountManyQuery = `select
  t.id tag_id,                              -- tag ID
  count(tp.tag_id) + count(tn.tag_id) items -- number of items,
from
  tags t
  left join tag_posts tp on tp.tag_id = t.id
  left join tag_nodes tn on tn.tag_id = t.id
group by
  t.id
`

func (q *Querier) Search(ctx context.Context, query string) (tag_ref.Tags, error) {
	r, err := q.db.Tag.Query().
		Where(
			ent_tag.NameContainsFold(query),
		).
		All(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	var counts tag_ref.TagItemsResults
	err = q.raw.SelectContext(ctx, &counts, tagItemsCountManyQuery)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	tags := dt.Map(r, tag_ref.Map(counts))

	return tags, nil
}

func (q *Querier) Get(ctx context.Context, name tag_ref.Name) (*tag.Tag, error) {
	r, err := q.db.Tag.Query().
		Where(ent_tag.Name(string(name))).
		WithAccounts().
		WithPosts(func(pq *ent.PostQuery) {
			pq.WithCategory()
			pq.WithAuthor(func(aq *ent.AccountQuery) {
				aq.WithAccountRoles(func(arq *ent.AccountRolesQuery) { arq.WithRole() })
			})
		}).
		WithNodes(func(nq *ent.NodeQuery) {
			nq.WithOwner(func(aq *ent.AccountQuery) {
				aq.WithAccountRoles(func(arq *ent.AccountRolesQuery) { arq.WithRole() })
			})
			nq.WithPrimaryImage()
		}).
		First(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	tag, err := tag.Map(r)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return tag, nil
}
