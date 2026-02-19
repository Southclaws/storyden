package tag_querier

import (
	"context"
	"sort"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/jmoiron/sqlx"

	"github.com/Southclaws/storyden/app/resources/account/role/role_repo"
	"github.com/Southclaws/storyden/app/resources/tag"
	"github.com/Southclaws/storyden/app/resources/tag/tag_ref"
	"github.com/Southclaws/storyden/internal/ent"
	ent_node "github.com/Southclaws/storyden/internal/ent/node"
	ent_post "github.com/Southclaws/storyden/internal/ent/post"
	ent_tag "github.com/Southclaws/storyden/internal/ent/tag"
)

type Querier struct {
	db          *ent.Client
	raw         *sqlx.DB
	roleQuerier *role_repo.Repository
}

func New(db *ent.Client, raw *sqlx.DB, roleQuerier *role_repo.Repository) *Querier {
	return &Querier{
		db:          db,
		raw:         raw,
		roleQuerier: roleQuerier,
	}
}

const tagItemsCountManyQuery = `select
  t.id tag_id,                    -- tag ID
  count(p.id) + count(n.id) items -- number of items,
from
  tags t
  left join tag_posts tp on tp.tag_id = t.id
  left join posts p on p.id = tp.post_id and p.visibility = 'published' and p.deleted_at is null
  left join tag_nodes tn on tn.tag_id = t.id
  left join nodes n on n.id = tn.node_id and n.visibility = 'published' and n.deleted_at is null
group by
  t.id
`

func (q *Querier) List(ctx context.Context) (tag_ref.Tags, error) {
	r, err := q.db.Tag.Query().All(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	var counts tag_ref.TagItemsResults
	err = q.raw.SelectContext(ctx, &counts, tagItemsCountManyQuery)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	tags := tag_ref.Tags(dt.Map(r, tag_ref.Map(counts)))

	sort.Sort(tags)

	return tags, nil
}

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
		Where(ent_tag.Name(name.String())).
		WithAccounts().
		WithPosts(func(pq *ent.PostQuery) {
			pq.Where(
				ent_post.VisibilityEQ(ent_post.VisibilityPublished),
				ent_post.DeletedAtIsNil(),
			)
			pq.WithCategory()
			pq.WithAuthor()
		}).
		WithNodes(func(nq *ent.NodeQuery) {
			nq.Where(
				ent_node.VisibilityEQ(ent_node.VisibilityPublished),
				ent_node.DeletedAtIsNil(),
			)
			nq.WithOwner()
			nq.WithPrimaryImage()
		}).
		First(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	accs := make([]*ent.Account, 0, len(r.Edges.Posts)+len(r.Edges.Nodes))
	for _, post := range r.Edges.Posts {
		accs = append(accs, post.Edges.Author)
	}
	for _, node := range r.Edges.Nodes {
		accs = append(accs, node.Edges.Owner)
	}

	roleHydrator, err := q.roleQuerier.BuildMultiHydrator(ctx, accs)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	tag, err := tag.Map(r, roleHydrator.Hydrate)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return tag, nil
}
