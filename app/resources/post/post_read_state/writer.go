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
	"github.com/Southclaws/storyden/internal/ent/postread"
)

type Writer struct {
	db *ent.Client
}

func New(db *ent.Client) *Writer {
	return &Writer{db: db}
}

func (w *Writer) UpsertReadState(ctx context.Context, accountID account.AccountID, threadID post.ID) error {
	_, err := w.db.PostRead.Create().
		SetAccountID(xid.ID(accountID)).
		SetRootPostID(xid.ID(threadID)).
		SetLastSeenAt(time.Now().UTC()).
		OnConflictColumns(postread.FieldRootPostID, postread.FieldAccountID).
		UpdateNewValues().
		ID(ctx)
	if err != nil {
		return fault.Wrap(err, ftag.With(ftag.Internal))
	}

	return nil
}
