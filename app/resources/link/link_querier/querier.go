package link_querier

import (
	"context"
	"math"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account/role/role_repo"
	"github.com/Southclaws/storyden/app/resources/link"
	"github.com/Southclaws/storyden/app/resources/link/link_ref"
	"github.com/Southclaws/storyden/internal/ent"
	link_ent "github.com/Southclaws/storyden/internal/ent/link"
)

type LinkQuerier struct {
	db          *ent.Client
	roleQuerier *role_repo.Repository
}

func New(db *ent.Client, roleQuerier *role_repo.Repository) *LinkQuerier {
	return &LinkQuerier{
		db:          db,
		roleQuerier: roleQuerier,
	}
}

type Result struct {
	PageSize    int
	Results     int
	TotalPages  int
	CurrentPage int
	NextPage    opt.Optional[int]
	Links       []*link_ref.LinkRef
}

type Filter func(*ent.LinkQuery)

func WithURL(s string) Filter {
	return func(lq *ent.LinkQuery) {
		lq.Where(link_ent.URLContainsFold(s))
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

	roleHydrator, err := d.roleQuerier.BuildMultiHydrator(ctx, roleHydrationTargets(r))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	link, err := link.Map(r, roleHydrator.Hydrate)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return link, nil
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

func (d *LinkQuerier) Search(ctx context.Context, page int, size int, filters ...Filter) (*Result, error) {
	total, err := d.db.Link.Query().Count(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	query := d.db.Link.Query().
		WithPrimaryImage().
		WithFaviconImage().
		Limit(size + 1).
		Offset(page * size).
		Order(ent.Desc(link_ent.FieldCreatedAt))

	for _, fn := range filters {
		fn(query)
	}

	query.WithAssets()

	r, err := query.All(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	isNextPage := len(r) >= size
	nextPage := opt.NewSafe(page+1, isNextPage)

	if isNextPage {
		r = r[:len(r)-1]
	}

	links := dt.Map(r, link_ref.Map)

	return &Result{
		PageSize:    size,
		Results:     len(links),
		TotalPages:  int(math.Ceil(float64(total) / float64(size))),
		CurrentPage: page,
		NextPage:    nextPage,
		Links:       links,
	}, nil
}

func roleHydrationTargets(in *ent.Link) []*ent.Account {
	targets := map[xid.ID]*ent.Account{}

	for _, node := range in.Edges.Nodes {
		if node == nil || node.Edges.Owner == nil {
			continue
		}

		targets[node.Edges.Owner.ID] = node.Edges.Owner
	}

	for _, p := range in.Edges.Posts {
		if p == nil || p.Edges.Author == nil {
			continue
		}

		targets[p.Edges.Author.ID] = p.Edges.Author
	}

	out := make([]*ent.Account, 0, len(targets))
	for _, acc := range targets {
		out = append(out, acc)
	}

	return out
}
