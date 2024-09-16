package post_writer

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/internal/ent"
	ent_post "github.com/Southclaws/storyden/internal/ent/post"
)

// PostWriter provides a way to make changes to posts, regardless of the kind.
// Sometimes you may have a reference to a thread or a reply but you don't know
// which it is and finding that out would require a database call. If you only
// need to update a shared field such as the content, you should use this type.
type PostWriter struct {
	db *ent.Client
}

func New(db *ent.Client) *PostWriter {
	return &PostWriter{db: db}
}

type Option func(*ent.PostMutation)

func WithContent(v datagraph.Content) Option {
	return func(pm *ent.PostMutation) {
		pm.SetBody(v.HTML())
		pm.SetShort(v.Short())
	}
}

func WithContentLinks(ids ...xid.ID) Option {
	return func(pm *ent.PostMutation) {
		pm.AddContentLinkIDs(ids...)
	}
}

func (p *PostWriter) Update(ctx context.Context, id post.ID, opts ...Option) (*post.Post, error) {
	update := p.db.Post.UpdateOneID(xid.ID(id))
	mutate := update.Mutation()

	for _, fn := range opts {
		fn(mutate)
	}

	err := update.Exec(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	r, err := p.db.Post.
		Query().
		Where(ent_post.IDEQ(xid.ID(id))).
		WithAuthor(func(aq *ent.AccountQuery) {
			aq.WithRoles()
		}).
		WithCategory().
		WithTags().
		WithAssets().
		WithLink(func(lq *ent.LinkQuery) {
			lq.WithFaviconImage().WithPrimaryImage()
		}).
		Only(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Internal))
	}

	return post.Map(r)
}
