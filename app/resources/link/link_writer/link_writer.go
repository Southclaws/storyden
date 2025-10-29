package link_writer

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/link/link_ref"
	"github.com/Southclaws/storyden/app/resources/mark"
	"github.com/Southclaws/storyden/internal/ent"
	link_ent "github.com/Southclaws/storyden/internal/ent/link"
)

type LinkWriter struct {
	db *ent.Client
}

func New(db *ent.Client) *LinkWriter {
	return &LinkWriter{db}
}

type Option func(*ent.LinkMutation)

func WithPosts(ids ...xid.ID) Option {
	return func(lm *ent.LinkMutation) {
		lm.AddPostIDs(ids...)
	}
}

func WithNodes(ids ...xid.ID) Option {
	return func(lm *ent.LinkMutation) {
		lm.AddNodeIDs(ids...)
	}
}

func WithFaviconImage(id asset.AssetID) Option {
	return func(lm *ent.LinkMutation) {
		lm.SetFaviconImageID(id)
	}
}

func WithPrimaryImage(id asset.AssetID) Option {
	return func(lm *ent.LinkMutation) {
		lm.SetPrimaryImageID(id)
	}
}

func WithAssets(ids ...asset.AssetID) Option {
	return func(lm *ent.LinkMutation) {
		lm.AddAssetIDs(ids...)
	}
}

func (d *LinkWriter) Store(ctx context.Context, address, title, description string, opts ...Option) (*link_ref.LinkRef, error) {
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
		Where(link_ent.ID(r.ID)).
		WithFaviconImage().
		WithPrimaryImage().
		First(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return link_ref.Map(r), nil
}

func getLinkAttrs(u url.URL) (string, string) {
	host := strings.TrimPrefix(u.Hostname(), "www.")

	full := fmt.Sprintf("%s-%s", host, u.Path)

	slugified := mark.Slugify(full)
	domain := u.Hostname()

	return slugified, domain
}
