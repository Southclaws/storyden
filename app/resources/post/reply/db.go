package reply

import (
	"context"
	"time"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/fault/ftag"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/ent/asset"
	ent_post "github.com/Southclaws/storyden/internal/ent/post"
)

type database struct {
	db *ent.Client
}

func New(db *ent.Client) Repository {
	return &database{db}
}

func (d *database) Create(
	ctx context.Context,
	authorID account.AccountID,
	parentID post.ID,
	opts ...Option,
) (*Reply, error) {
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
		SetAuthorID(xid.ID(authorID))

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

	return d.Get(ctx, post.ID(p.ID))
}

func (d *database) Get(ctx context.Context, id post.ID) (*Reply, error) {
	p, err := d.db.Post.
		Query().
		Where(ent_post.IDEQ(xid.ID(id))).
		WithAuthor().
		WithRoot(func(pq *ent.PostQuery) {
			pq.WithAuthor()
		}).
		WithAssets(func(aq *ent.AssetQuery) {
			aq.Order(asset.ByUpdatedAt(), asset.ByCreatedAt())
		}).
		Only(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Internal))
	}

	return Map(p)
}

func (d *database) Update(ctx context.Context, id post.ID, opts ...Option) (*Reply, error) {
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

	return Map(p)
}

func (d *database) Delete(ctx context.Context, id post.ID) error {
	err := d.db.Post.
		UpdateOneID(xid.ID(id)).
		SetDeletedAt(time.Now()).
		Exec(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to archive thread root post"))
	}

	return nil
}
