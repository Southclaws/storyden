package tag_querier

import (
	"context"
	"fmt"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/opt"
	"github.com/jmoiron/sqlx"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account/role/role_hydrate"
	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/pagination"
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
	roleQuerier *role_hydrate.Hydrator
}

func New(db *ent.Client, raw *sqlx.DB, roleQuerier *role_hydrate.Hydrator) *Querier {
	return &Querier{db, raw, roleQuerier}
}

const tagItemsCountPageQuery = `select
  t.id tag_id,
  (
    select count(*)
    from tag_posts tp
    join posts p on p.id = tp.post_id
    where tp.tag_id = t.id
      and p.visibility = 'published'
      and p.deleted_at is null
  ) + (
    select count(*)
    from tag_nodes tn
    join nodes n on n.id = tn.node_id
    where tn.tag_id = t.id
      and n.visibility = 'published'
      and n.deleted_at is null
  ) items
from
  tags t
%s
order by
  items desc, t.name asc
limit ? offset ?
`

func (q *Querier) List(ctx context.Context, params pagination.Parameters) (*pagination.Result[*tag_ref.Tag], error) {
	return q.page(ctx, params, opt.NewEmpty[string]())
}

func (q *Querier) Search(ctx context.Context, query string, params pagination.Parameters) (*pagination.Result[*tag_ref.Tag], error) {
	return q.page(ctx, params, opt.New(query))
}

// tags are ordered by how much they are used, which only exists in the
// aggregate, so the page is picked there and the rows are hydrated after
func (q *Querier) page(ctx context.Context, params pagination.Parameters, search opt.Optional[string]) (*pagination.Result[*tag_ref.Tag], error) {
	countQuery := q.db.Tag.Query()
	where := ""
	args := []any{}

	if s, ok := search.Get(); ok {
		countQuery = countQuery.Where(ent_tag.NameContainsFold(s))
		where = "where lower(t.name) like lower(?)"
		args = append(args, "%"+s+"%")
	}

	total, err := countQuery.Count(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	args = append(args, params.Limit(), params.Offset())

	var counts tag_ref.TagItemsResults
	if err := q.raw.SelectContext(ctx, &counts, q.raw.Rebind(fmt.Sprintf(tagItemsCountPageQuery, where)), args...); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if len(counts) == 0 {
		empty := pagination.NewPageResult(params, total, []*tag_ref.Tag{})
		return &empty, nil
	}

	ids := make([]xid.ID, 0, len(counts))
	for _, c := range counts {
		ids = append(ids, c.TagID)
	}

	rows, err := q.db.Tag.Query().Where(ent_tag.IDIn(ids...)).All(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	byID := make(map[xid.ID]*tag_ref.Tag, len(rows))
	for _, t := range dt.Map(rows, tag_ref.Map(counts)) {
		byID[xid.ID(t.ID)] = t
	}

	tags := make([]*tag_ref.Tag, 0, len(ids))
	for _, id := range ids {
		if t, ok := byID[id]; ok {
			tags = append(tags, t)
		}
	}

	result := pagination.NewPageResult(params, total, tags)

	return &result, nil
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

	if err := q.roleQuerier.HydrateRoleEdges(ctx, roleHydrationTargets(r)...); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	tag, err := tag.Map(r)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return tag, nil
}

func roleHydrationTargets(r *ent.Tag) []*ent.Account {
	targets := make([]*ent.Account, 0, len(r.Edges.Posts)+len(r.Edges.Nodes))

	for _, p := range r.Edges.Posts {
		if p != nil && p.Edges.Author != nil {
			targets = append(targets, p.Edges.Author)
		}
	}

	for _, n := range r.Edges.Nodes {
		targets = append(targets, library.RoleHydrationTargetsFromNode(n)...)
	}

	return targets
}
