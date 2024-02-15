package link

import (
	"context"
	"fmt"
	"math"
	"net/url"
	"strings"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/opt"
	"github.com/gosimple/slug"

	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/ent/link"
)

type database struct {
	db *ent.Client
}

func New(db *ent.Client) Repository {
	return &database{db}
}

func (d *database) Store(ctx context.Context, address, title, description string, opts ...Option) (*Link, error) {
	u, err := url.Parse(address)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	slug, domain := getLinkAttrs(*u)

	create := d.db.Link.Create()
	mutate := create.Mutation()

	mutate.SetURL(address)
	mutate.SetSlug(slug)
	mutate.SetDomain(domain)
	mutate.SetTitle(title)
	mutate.SetDescription(description)

	for _, fn := range opts {
		fn(mutate)
	}

	create.OnConflictColumns("url").UpdateNewValues()
	create.OnConflictColumns("slug").UpdateNewValues()

	r, err := create.Save(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	r, err = d.db.Link.Query().
		WithAssets().
		Where(link.ID(r.ID)).
		First(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	link := Map(r)

	return link, nil
}

func (d *database) Search(ctx context.Context, page int, size int, filters ...Filter) (*Result, error) {
	total, err := d.db.Link.Query().Count(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	query := d.db.Link.Query().
		Limit(size + 1).
		Offset(page * size).
		Order(ent.Desc(link.FieldCreatedAt))

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

	links := MapA(r)

	return &Result{
		PageSize:    size,
		Results:     len(links),
		TotalPages:  int(math.Ceil(float64(total) / float64(size))),
		CurrentPage: page,
		NextPage:    nextPage,
		Links:       links,
	}, nil
}

func getLinkAttrs(u url.URL) (string, string) {
	host := strings.TrimPrefix(u.Hostname(), "www.")

	full := fmt.Sprintf("%s-%s", host, u.Path)

	slugified := slug.Make(full)
	domain := u.Hostname()

	return slugified, domain
}
