package link_querier

import (
	"context"
	"strings"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account/role/role_hydrate"
	"github.com/Southclaws/storyden/app/resources/library"
	"github.com/Southclaws/storyden/app/resources/link"
	"github.com/Southclaws/storyden/app/resources/link/link_ref"
	"github.com/Southclaws/storyden/app/resources/link/link_url"
	"github.com/Southclaws/storyden/app/resources/pagination"
	"github.com/Southclaws/storyden/internal/ent"
	link_ent "github.com/Southclaws/storyden/internal/ent/link"
	"github.com/Southclaws/storyden/internal/ent/predicate"
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

func WithURL(u link_url.URL) Filter {
	return func(lq *ent.LinkQuery) {
		variants := []string{u.String()}
		raw := u.NetURL()
		if raw.Path == "/" {
			raw.Path = ""
			variants = append(variants, raw.String())
		}

		lq.Where(link_ent.URLIn(variants...))
	}
}

func WithKeyword(s string) Filter {
	return func(lq *ent.LinkQuery) {
		s = strings.TrimSpace(s)
		if s == "" {
			return
		}

		if predicates := urlSearch(s); len(predicates) > 0 {
			lq.Where(link_ent.Or(predicates...))
			return
		}

		lq.Where(link_ent.Or(
			link_ent.TitleContainsFold(s),
			link_ent.DescriptionContainsFold(s),
			link_ent.URLContainsFold(s),
		))
	}
}

func urlSearch(s string) []predicate.Link {
	raw := s
	if !strings.Contains(raw, "://") {
		raw = "https://" + raw
	}

	u, err := link_url.Parse(raw)
	if err != nil || !strings.Contains(u.Hostname(), ".") || strings.ContainsAny(s, " \t\n") {
		return nil
	}

	rawURL := u.NetURL()
	hosts := []string{rawURL.Host}
	if strings.HasPrefix(rawURL.Host, "www.") {
		hosts = append(hosts, strings.TrimPrefix(rawURL.Host, "www."))
	} else {
		hosts = append(hosts, "www."+rawURL.Host)
	}

	var predicates []predicate.Link
	for _, host := range hosts {
		for _, scheme := range []string{"http", "https"} {
			candidate := rawURL
			candidate.Host = host
			candidate.Scheme = scheme
			prefix := candidate.String()
			predicates = append(predicates, link_ent.URLHasPrefix(prefix))
			if candidate.Path == "/" && candidate.RawQuery == "" && candidate.Fragment == "" {
				candidate.Path = ""
				predicates = append(predicates, link_ent.URLEQ(candidate.String()))
			}
		}
	}

	return predicates
}

func (d *LinkQuerier) FindByURL(ctx context.Context, u link_url.URL) (*link_ref.LinkRef, error) {
	query := d.db.Link.Query().WithPrimaryImage().WithFaviconImage()
	WithURL(u)(query)
	r, err := query.Order(ent.Asc(link_ent.FieldID)).First(ctx)
	if ent.IsNotFound(err) {
		return nil, nil
	}

	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return link_ref.Map(r)
}

func (d *LinkQuerier) Get(ctx context.Context, slug string) (*link.Link, error) {
	query := d.db.Link.Query().
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

	r, err := query.Clone().Where(link_ent.SlugEQ(slug)).First(ctx)
	if ent.IsNotFound(err) {
		r, err = query.Where(link_ent.SlugEqualFold(slug)).First(ctx)
	}

	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if _, err := link_url.Parse(r.URL); err != nil {
		return nil, fault.New("link not found", ftag.With(ftag.NotFound))
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
		WithPrimaryImage().
		WithFaviconImage().
		Where(link_ent.ID(xid.ID(id))).
		First(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return link_ref.Map(r)
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

	links := make([]*link_ref.LinkRef, 0, len(r))
	for _, row := range r {
		ref, err := link_ref.Map(row)
		if err != nil {
			continue
		}

		links = append(links, ref)
	}

	result := pagination.NewPageResult(params, total, links)

	return &result, nil
}
