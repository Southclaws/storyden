package post

import (
	"context"

	"4d63.com/optional"
	"github.com/google/uuid"

	"github.com/Southclaws/storyden/backend/internal/infrastructure/db/model"
	"github.com/Southclaws/storyden/backend/internal/infrastructure/db/model/post"
	"github.com/Southclaws/storyden/backend/pkg/resources/user"
)

type database struct {
	db *model.Client
}

func New(db *model.Client) Repository {
	return &database{db}
}

func (d *database) CreatePost(
	ctx context.Context,
	body string,
	authorID user.UserID,
	parentID PostID,
	replyToID optional.Optional[PostID],
) (*Post, error) {
	short := MakeShortBody(body)

	q := d.db.Post.
		Create().
		SetBody(body).
		SetShort(short).
		SetFirst(false).
		// SetRootID(uuid.UUID(parentID)).
		SetAuthorID(uuid.UUID(authorID))

	replyToID.If(func(value PostID) {
		q.SetReplyToID(uuid.UUID(value))
	})

	p, err := q.Save(ctx)
	if err != nil {
		if model.IsNotFound(err) {
			return nil, nil
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
