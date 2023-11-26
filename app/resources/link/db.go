package link

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
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

func (d *database) Search(ctx context.Context, filters ...Filter) ([]*Link, error) {
	query := d.db.Link.Query()

	for _, fn := range filters {
		fn(query)
	}

	query.WithAssets()

	r, err := query.All(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	links := MapA(r)

	return links, nil
}

func getLinkAttrs(u url.URL) (string, string) {
	host := strings.TrimPrefix(u.Hostname(), "www.")

	full := fmt.Sprintf("%s-%s", host, u.Path)

	slugified := slug.Make(full)
	domain := u.Hostname()

	return slugified, domain
}
