package post

import (
	"context"

	"4d63.com/optional"
	"github.com/pkg/errors"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/internal/errmeta"
	"github.com/Southclaws/storyden/internal/infrastructure/db/model"
	"github.com/Southclaws/storyden/internal/infrastructure/db/model/post"
	"github.com/Southclaws/storyden/pkg/resources/account"
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
) (*Post, error) {
	short := MakeShortBody(body)

	thread, err := d.db.Post.Get(ctx, xid.ID(parentID))
	if err != nil {
		return nil, errmeta.Wrap(errors.Wrap(err, "failed to get parent thread"),
			"authorID", authorID.String(),
			"parentID", parentID.String())
	}

	if thread.First == false {
		return nil, errors.New("attempt to create post under non-thread post")
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
			return nil, errors.Wrap(err, "constraint error: check parent ID or reply ID")
		}

		return nil, err
	}

	p, err = d.db.Post.Query().
		Where(post.IDEQ(p.ID)).
		WithAuthor().
		WithRoot(func(pq *model.PostQuery) {
			pq.WithAuthor()
		}).
		Only(ctx)
	if err != nil {
		return nil, err
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
