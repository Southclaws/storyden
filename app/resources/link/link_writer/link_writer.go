package link_writer

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/link/link_ref"
	"github.com/Southclaws/storyden/app/resources/link/link_url"
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

func (d *LinkWriter) Store(ctx context.Context, u link_url.URL, title, description string, opts ...Option) (*link_ref.LinkRef, error) {
	if u.IsZero() {
		return nil, link_url.ErrEmpty
	}

	slug, domain := getLinkAttrs(u)

	create := d.db.Link.Create()
	mutate := create.Mutation()

	mutate.SetURL(u.String())
	mutate.SetSlug(slug)
	mutate.SetDomain(domain)
	mutate.SetTitle(title)
	mutate.SetDescription(description)

	for _, fn := range opts {
		fn(mutate)
	}

	create.OnConflictColumns(link_ent.FieldURL).UpdateNewValues().Update(func(upsert *ent.LinkUpsert) {
		if title == "" {
			upsert.SetIgnore(link_ent.FieldTitle)
		}

		if description == "" {
			upsert.SetIgnore(link_ent.FieldDescription)
		}

		upsert.SetIgnore(link_ent.FieldLivenessStatus)
		upsert.SetIgnore(link_ent.FieldLivenessFailures)
		upsert.SetIgnore(link_ent.FieldScrapeStatus)
	})

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

	return link_ref.Map(r)
}

func getLinkAttrs(u link_url.URL) (string, string) {
	host := strings.TrimPrefix(u.Hostname(), "www.")

	full := fmt.Sprintf("%s-%s", host, u.NetURL().Path)

	slugified := mark.Slugify(full)
	domain := u.Hostname()

	return slugified + "-" + xid.New().String(), domain
}

type ScrapeResult struct {
	Success     bool
	Title       string
	Description string
	Domain      string
	Options     []Option
}

func (d *LinkWriter) FinishScrape(ctx context.Context, id xid.ID, now time.Time, result ScrapeResult) error {
	update := d.db.Link.UpdateOneID(id)
	mutation := update.Mutation()
	mutation.SetScrapeStatus(link_ent.ScrapeStatusFailed)
	mutation.SetScrapeAttemptedAt(now)

	if result.Success {
		mutation.SetScrapeStatus(link_ent.ScrapeStatusSuccess)
		mutation.SetScrapedAt(now)
		if result.Title != "" {
			mutation.SetTitle(result.Title)
		}

		if result.Description != "" {
			mutation.SetDescription(result.Description)
		}

		if result.Domain != "" {
			mutation.SetDomain(result.Domain)
		}

		for _, fn := range result.Options {
			fn(mutation)
		}
	}

	_, err := update.Save(ctx)
	return fault.Wrap(err, fctx.With(ctx))
}
