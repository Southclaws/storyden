package reply

import (
	"context"
	"time"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/internal/ent"
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
	body string,
	authorID account.AccountID,
	parentID post.ID,
	replyToID opt.Optional[post.ID],
	meta map[string]any,
	opts ...Option,
) (*Reply, error) {
	short := post.MakeShortBody(string(body))

	thread, err := d.db.Post.Get(ctx, xid.ID(parentID))
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fault.Wrap(err, fmsg.With("failed to get parent thread"), fctx.With(ctx), ftag.With(ftag.NotFound))
		}

		return nil, fault.Wrap(err, fmsg.With("failed to get parent thread"), fctx.With(ctx), ftag.With(ftag.Internal))
	}

	if thread.First == false {
		return nil, fault.New("attempt to create post under non-thread post")
	}

	q := d.db.Post.
		Create().
		SetBody(string(body)).
		SetShort(short).
		SetFirst(false).
		SetRootID(xid.ID(parentID)).
		SetAuthorID(xid.ID(authorID))

	for _, fn := range opts {
		fn(q.Mutation())
	}

	replyToID.Call(func(value post.ID) {
		q.SetReplyToID(xid.ID(value))
	})

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

	return FromModel(p), nil
}

func (d *database) Get(ctx context.Context, id post.ID) (*Reply, error) {
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

	return FromModel(p), nil
}

func (d *database) Update(ctx context.Context, id post.ID, opts ...Option) (*Reply, error) {
	update := d.db.Post.UpdateOneID(xid.ID(id))
	mutate := update.Mutation()

	for _, fn := range opts {
		fn(mutate)
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

	return FromModel(p), nil
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
