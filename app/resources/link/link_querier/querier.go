package link_querier

import (
	"context"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account/role/role_hydrate"
	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/link"
	"github.com/Southclaws/storyden/app/resources/link/link_ref"
	"github.com/Southclaws/storyden/app/resources/pagination"
	"github.com/Southclaws/storyden/internal/ent"
	link_ent "github.com/Southclaws/storyden/internal/ent/link"
)

type LinkQuerier struct {
	db          *ent.Client
	roleQuerier *role_hydrate.Hydrator
}

func New(db *ent.Client, roleQuerier *role_hydrate.Hydrator) *LinkQuerier {
	return &LinkQuerier{
		db:          db,
		roleQuerier: roleQuerier,
	}
}

type Filter func(*ent.LinkQuery)

func WithURL(s string) Filter {
	return func(lq *ent.LinkQuery) {
		lq.Where(link_ent.URLEQ(s))
	}
}

func WithKeyword(s string) Filter {
	return func(lq *ent.LinkQuery) {
		lq.Where(link_ent.Or(
			link_ent.TitleContainsFold(s),
			link_ent.DescriptionContainsFold(s),
			link_ent.URLContainsFold(s),
		))
	}
}

func (d *LinkQuerier) Get(ctx context.Context, slug string) (*link.Link, error) {
	query := d.db.Link.Query().
		Where(link_ent.SlugEqualFold(slug)).
		WithAssets().
		WithPrimaryImage().
		WithFaviconImage().
		WithPosts(func(pq *ent.PostQuery) {
			pq.WithAuthor()
			pq.WithCategory()
			pq.WithRoot()
		}).
		WithNodes(func(nq *ent.NodeQuery) {
			nq.WithOwner()
		})

	r, err := query.First(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if err := d.roleQuerier.HydrateRoleEdges(ctx, roleHydrationTargets(r)...); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	link, err := link.Map(r)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return link, nil
}

func roleHydrationTargets(r *ent.Link) []*ent.Account {
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

func (d *LinkQuerier) GetByID(ctx context.Context, id link.LinkID) (*link_ref.LinkRef, error) {
	r, err := d.db.Link.Query().
		WithAssets().
		Where(link_ent.ID(xid.ID(id))).
		First(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	link := link_ref.Map(r)

	return link, nil
}

func (d *LinkQuerier) Search(ctx context.Context, params pagination.Parameters, filters ...Filter) (*pagination.Result[*link_ref.LinkRef], error) {
	query := d.db.Link.Query()

	for _, fn := range filters {
		fn(query)
	}

	total, err := query.Clone().Count(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	r, err := query.
		WithPrimaryImage().
		WithFaviconImage().
		WithAssets().
		Limit(params.Limit()).
		Offset(params.Offset()).
		Order(ent.Desc(link_ent.FieldCreatedAt), ent.Desc(link_ent.FieldID)).
		All(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	links := dt.Map(r, link_ref.Map)

	result := pagination.NewPageResult(params, total, links)

	return &result, nil
}
