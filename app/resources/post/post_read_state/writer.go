package post_read_state

import (
	"context"
	"time"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/ftag"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/internal/ent"
)

type Writer struct {
	db *ent.Client
}

func New(db *ent.Client) *Writer {
	return &Writer{db: db}
}

func (w *Writer) UpsertReadState(ctx context.Context, accountID account.AccountID, threadID post.ID) error {
	threadXID := xid.ID(threadID)

	p, err := w.db.Post.Get(ctx, threadXID)
	if err != nil {
		if ent.IsNotFound(err) {
			return fault.Wrap(err, ftag.With(ftag.NotFound))
		}
		return fault.Wrap(err, ftag.With(ftag.Internal))
	}

	_, err = w.db.PostRead.Create().
		SetAccountID(xid.ID(accountID)).
		SetRootPostID(p.ID).
		SetLastSeenAt(time.Now().UTC()).
		OnConflict().
		UpdateNewValues().
		ID(ctx)
	if err != nil {
		return fault.Wrap(err, ftag.With(ftag.Internal))
	}

	return nil
}
