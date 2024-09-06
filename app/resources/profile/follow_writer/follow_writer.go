package follow_writer

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/ent/accountfollow"
)

type Writer struct {
	db *ent.Client
}

func New(db *ent.Client) *Writer {
	return &Writer{db}
}

func (w *Writer) Follow(ctx context.Context, follower, following account.AccountID) error {
	err := w.db.AccountFollow.Create().
		SetFollowerAccountID(xid.ID(follower)).
		SetFollowingAccountID(xid.ID(following)).
		Exec(ctx)
	if err != nil {
		if ent.IsConstraintError(err) {
			return nil
		}
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}

func (w *Writer) Unfollow(ctx context.Context, follower, following account.AccountID) error {
	_, err := w.db.AccountFollow.Delete().
		Where(
			accountfollow.FollowerAccountID(xid.ID(follower)),
			accountfollow.FollowingAccountID(xid.ID(following)),
		).
		Exec(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}
