package reply_writer

import (
	"context"
	"time"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/fault/ftag"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/datagraph"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/post/reply"
	"github.com/Southclaws/storyden/app/resources/post/reply_querier"
	"github.com/Southclaws/storyden/internal/ent"
	ent_post "github.com/Southclaws/storyden/internal/ent/post"
)

type Writer struct {
	db      *ent.Client
	querier *reply_querier.Querier
}

func New(db *ent.Client, querier *reply_querier.Querier) *Writer {
	return &Writer{db: db, querier: querier}
}

type Option func(*ent.PostMutation)

func WithID(id post.ID) Option {
	return func(pm *ent.PostMutation) {
		pm.SetID(xid.ID(id))
	}
}

func WithContent(v datagraph.Content) Option {
	return func(pm *ent.PostMutation) {
		pm.SetBody(v.HTML())
		pm.SetShort(v.Short())
	}
}

func WithReplyTo(v post.ID) Option {
	return func(pm *ent.PostMutation) {
		pm.SetReplyToID(xid.ID(v))
	}
}

func WithMeta(meta map[string]any) Option {
	return func(m *ent.PostMutation) {
		m.SetMetadata(meta)
	}
}

func WithAssets(ids ...asset.AssetID) Option {
	return func(m *ent.PostMutation) {
		m.AddAssetIDs(ids...)
	}
}

func WithLink(id xid.ID) Option {
	return func(pm *ent.PostMutation) {
		pm.SetLinkID(id)
	}
}

func WithContentLinks(ids ...xid.ID) Option {
	return func(pm *ent.PostMutation) {
		pm.AddContentLinkIDs(ids...)
	}
}

func (d *Writer) Create(
	ctx context.Context,
	authorID account.AccountID,
	parentID post.ID,
	opts ...Option,
) (*reply.Reply, error) {
	thread, err := d.db.Post.Get(ctx, xid.ID(parentID))
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fault.Wrap(err, fmsg.With("failed to get parent thread"), fctx.With(ctx), ftag.With(ftag.NotFound))
		}

		return nil, fault.Wrap(err, fmsg.With("failed to get parent thread"), fctx.With(ctx), ftag.With(ftag.Internal))
	}

	if thread.RootPostID != nil {
		return nil, fault.New("attempt to create post under non-thread post")
	}

	q := d.db.Post.
		Create().
		SetUpdatedAt(time.Now()).
		SetLastReplyAt(time.Now()). // Required field, but unused for Reply view
		SetRootID(xid.ID(parentID)).
		SetAuthorID(xid.ID(authorID)).
		SetVisibility(ent_post.VisibilityPublished)

	for _, fn := range opts {
		fn(q.Mutation())
	}

	p, err := q.Save(ctx)
	if err != nil {
		if ent.IsConstraintError(err) {
			return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
		}

		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Internal))
	}

	p, err = d.db.Post.Query().
		Where(ent_post.IDEQ(p.ID)).
		WithAuthor().
		WithRoot(func(pq *ent.PostQuery) {
			pq.WithAuthor()
		}).
		WithAssets().
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.NotFound))
		}

		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Internal))
	}

	err = d.db.Post.
		UpdateOneID(xid.ID(parentID)).
		SetLastReplyAt(time.Now()).
		Exec(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return d.querier.Get(ctx, post.ID(p.ID))
}

func (d *Writer) Update(ctx context.Context, id post.ID, opts ...Option) (*reply.Reply, error) {
	update := d.db.Post.UpdateOneID(xid.ID(id))
	mutate := update.Mutation()

	for _, fn := range opts {
		fn(mutate)
	}

	// Only set the updated_at field if not changing the indexed_at field.
	if _, set := mutate.IndexedAt(); !set {
		mutate.SetUpdatedAt(time.Now())
	}

	err := update.Exec(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	p, err := d.db.Post.
		Query().
		Where(ent_post.IDEQ(xid.ID(id))).
		WithAuthor().
		WithRoot(func(pq *ent.PostQuery) {
			pq.WithAuthor()
		}).
		WithAssets().
		Only(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Internal))
	}

	return reply.Map(p)
}

func (d *Writer) Delete(ctx context.Context, id post.ID) error {
	err := d.db.Post.
		UpdateOneID(xid.ID(id)).
		SetDeletedAt(time.Now()).
		Exec(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to archive thread reply"))
	}

	return nil
}
