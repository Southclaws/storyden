package like_writer

import (
	"context"
	"database/sql"
	"errors"

	ent_sql "entgo.io/ent/dialect/sql"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/fault/ftag"
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
		OnConflict(ent_sql.DoNothing()).
		Exec(ctx)
	if err != nil {
		// NOTE: err no rows is a conflict-returning edge case.
		// and the actual ent not found is where post/account doesn't exist.
		// it's always the post, as the account is already checked earlier.
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		if ent.IsNotFound(err) {
			err = fault.Wrap(err, ftag.With(ftag.NotFound),
				fmsg.WithDesc("post not found", "The liked post was not found."))
		}
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
