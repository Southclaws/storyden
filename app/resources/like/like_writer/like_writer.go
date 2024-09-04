package like_writer

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/ent/likepost"
)

type LikeWriter struct {
	db *ent.Client
}

func New(db *ent.Client) *LikeWriter {
	return &LikeWriter{db: db}
}

func (l *LikeWriter) AddPostLike(ctx context.Context, accountID account.AccountID, postID post.ID) error {
	err := l.db.LikePost.Create().
		SetPostID(xid.ID(postID)).
		SetAccountID(xid.ID(accountID)).
		OnConflict(sql.DoNothing()).
		Exec(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}

func (l *LikeWriter) RemovePostLike(ctx context.Context, accountID account.AccountID, postID post.ID) error {
	_, err := l.db.LikePost.
		Delete().
		Where(
			likepost.AccountIDEQ(xid.ID(accountID)),
			likepost.PostIDEQ(xid.ID(postID)),
		).
		Exec(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}
