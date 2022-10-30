package post

import (
	"context"

	"4d63.com/optional"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/fault/ftag"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/internal/infrastructure/db/model"
	"github.com/Southclaws/storyden/internal/infrastructure/db/model/post"
)

type database struct {
	db *model.Client
}

func New(db *model.Client) Repository {
	return &database{db}
}

func (d *database) Create(
	ctx context.Context,
	body string,
	authorID account.AccountID,
	parentID PostID,
	replyToID optional.Optional[PostID],
	meta map[string]any,
) (*Post, error) {
	short := MakeShortBody(body)

	thread, err := d.db.Post.Get(ctx, xid.ID(parentID))
	if err != nil {
		if model.IsNotFound(err) {
			return nil, fault.Wrap(err, fmsg.With("failed to get parent thread"), fctx.With(ctx), ftag.With(ftag.NotFound))
		}

		return nil, fault.Wrap(err, fmsg.With("failed to get parent thread"), fctx.With(ctx), ftag.With(ftag.Internal))
	}

	if thread.First == false {
		return nil, fault.New("attempt to create post under non-thread post")
	}

	q := d.db.Post.
		Create().
		SetBody(body).
		SetShort(short).
		SetFirst(false).
		SetRootID(xid.ID(parentID)).
		SetAuthorID(xid.ID(authorID))

	replyToID.If(func(value PostID) {
		q.SetReplyToID(xid.ID(value))
	})

	p, err := q.Save(ctx)
	if err != nil {
		if model.IsConstraintError(err) {
			return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
		}

		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Internal))
	}

	p, err = d.db.Post.Query().
		Where(post.IDEQ(p.ID)).
		WithAuthor().
		WithRoot(func(pq *model.PostQuery) {
			pq.WithAuthor()
		}).
		Only(ctx)
	if err != nil {
		if model.IsNotFound(err) {
			return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.NotFound))
		}

		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Internal))
	}

	return FromModel(p), nil
}

// func (d *database) EditPost(ctx context.Context, authorID, id string, title *string, body *string) (*Post, error) {
// 	// This could probably be optimised. I am too lazy to do it rn.
// 	post, err := d.db.Post.
// 		FindUnique(
// 			db.Post.ID.Equals(id),
// 		).
// 		With(db.Post.Author.Fetch()).
// 		Exec(ctx)
// 	if err != nil {
// 		if errors.Is(err, db.ErrNotFound) {
// 			return nil, nil
// 		}
// 		return nil, err
// 	}
// 	if post.Author().ID != authorID {
// 		return nil, ErrUnauthorised
// 	}

// 	post, err = d.db.Post.
// 		FindUnique(
// 			db.Post.ID.Equals(id),
// 		).
// 		With(db.Post.Author.Fetch()).
// 		Update(
// 			db.Post.Title.SetIfPresent(title),
// 			db.Post.Body.SetIfPresent(body),
// 		).
// 		Exec(ctx)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return FromModel(post), err
// }

// func (d *database) DeletePost(ctx context.Context, authorID, postID string, force bool) (*Post, error) {
// 	// This could probably be optimised. I am too lazy to do it rn.
// 	post, err := d.db.Post.
// 		FindUnique(
// 			db.Post.ID.Equals(postID),
// 		).
// 		With(
// 			db.Post.Author.Fetch(),
// 			db.Post.Tags.Fetch(),
// 			db.Post.Category.Fetch(),
// 		).
// 		Exec(ctx)
// 	if err != nil {
// 		if errors.Is(err, db.ErrNotFound) {
// 			return nil, nil
// 		}
// 		return nil, err
// 	}
// 	if force == false {
// 		if post.Author().ID != authorID {
// 			return nil, ErrUnauthorised
// 		}
// 	}

// 	_, err = d.db.Post.
// 		FindUnique(db.Post.ID.Equals(postID)).
// 		Update(
// 			db.Post.DeletedAt.Set(time.Now()),
// 		).
// 		Exec(ctx)
// 	if err != nil {
// 		if errors.Is(err, db.ErrNotFound) {
// 			return nil, nil
// 		}
// 		return nil, err
// 	}

// 	return FromModel(post), err
// }
