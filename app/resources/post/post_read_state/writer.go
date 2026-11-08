package post_read_state

import (
	"context"
	"database/sql"
	"errors"
	"time"

	entsql "entgo.io/ent/dialect/sql"
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
	return w.MarkRead(ctx, accountID, threadID, time.Now().UTC())
}

func (w *Writer) MarkRead(ctx context.Context, accountID account.AccountID, threadID post.ID, at time.Time) error {
	return w.update(ctx, accountID, threadID, at, false)
}

func (w *Writer) MarkUnread(ctx context.Context, accountID account.AccountID, threadID post.ID, at time.Time) error {
	return w.update(ctx, accountID, threadID, at, true)
}

func (w *Writer) update(ctx context.Context, accountID account.AccountID, threadID post.ID, at time.Time, unread bool) error {
	table := entsql.Table(postread.Table)
	markerBefore := entsql.LT(table.C(postread.FieldMarkedUnreadAt), at)
	create := w.db.PostRead.Create().
		SetAccountID(xid.ID(accountID)).
		SetRootPostID(xid.ID(threadID)).
		SetLastSeenAt(at)

	if unread {
		create.SetMarkedUnreadAt(at)
		markerBefore = entsql.LTE(table.C(postread.FieldMarkedUnreadAt), at)
	}

	err := create.OnConflict(
		entsql.ConflictColumns(postread.FieldRootPostID, postread.FieldAccountID),
		entsql.UpdateWhere(entsql.And(
			entsql.LTE(table.C(postread.FieldLastSeenAt), at),
			entsql.Or(
				entsql.IsNull(table.C(postread.FieldMarkedUnreadAt)),
				markerBefore,
			),
		)),
	).
		Update(func(update *ent.PostReadUpsert) {
			if unread {
				update.SetMarkedUnreadAt(at)
			} else {
				update.SetLastSeenAt(at)
				update.ClearMarkedUnreadAt()
			}
		}).
		Exec(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return nil
	}

	if err != nil {
		return fault.Wrap(err, ftag.With(ftag.Internal))
	}

	return nil
}
