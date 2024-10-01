package reaction

import (
	"context"
	"database/sql"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/pkg/errors"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/ent/react"
)

var ErrInvalidEmoji = errors.New("invalid emoji codepoint")

func New(db *ent.Client) (*Querier, *Writer) {
	q := &Querier{db}
	w := &Writer{db, q}
	return q, w
}

type Querier struct {
	db *ent.Client
}

func (q *Querier) Get(ctx context.Context, reactID ReactID) (*React, error) {
	r, err := q.db.React.Query().
		Where(react.ID(xid.ID(reactID))).
		WithAccount(func(aq *ent.AccountQuery) { aq.WithAccountRoles(func(arq *ent.AccountRolesQuery) { arq.WithRole() }) }).
		Only(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Internal))
	}

	return Map(r)
}

func (q *Querier) Lookup(ctx context.Context, accountID account.AccountID, postID xid.ID, e string) (*React, error) {
	r, err := q.db.React.Query().
		Where(
			react.AccountID(xid.ID(accountID)),
			react.PostID(xid.ID(postID)),
			react.Emoji(e),
		).
		WithAccount(func(aq *ent.AccountQuery) { aq.WithAccountRoles(func(arq *ent.AccountRolesQuery) { arq.WithRole() }) }).
		Only(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Internal))
	}

	return Map(r)
}

type Writer struct {
	db      *ent.Client
	querier *Querier
}

func (w *Writer) Add(ctx context.Context, accountID account.AccountID, postID xid.ID, emojiID string) (*React, error) {
	e, ok := IsValidEmoji(emojiID)
	if !ok {
		return nil, fault.Wrap(ErrInvalidEmoji, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	reactID, err := w.tryAdd(ctx, accountID, postID, e)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if reactID == nil {
		return w.querier.Lookup(ctx, accountID, postID, e)
	}

	return w.querier.Get(ctx, ReactID(*reactID))
}

func (w *Writer) tryAdd(ctx context.Context, accountID account.AccountID, postID xid.ID, e string) (*xid.ID, error) {
	reactID, err := w.db.React.
		Create().
		SetEmoji(e).
		SetAccountID(xid.ID(accountID)).
		SetPostID(xid.ID(postID)).
		OnConflictColumns(react.FieldPostID, react.FieldAccountID, react.FieldEmoji).DoNothing().
		ID(ctx)
	if err != nil {
		// NOTE: Not found is a red herring here, due to SQL being as weird as
		// it normally is, on-conflict-do-nothing doesn't return anything.
		// NOTE 2: but, ent.IsNotFound(err) does not work, wrong sentinel error.
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Internal))
	}

	return &reactID, nil
}

func (w *Writer) Remove(ctx context.Context, accountID account.AccountID, reactID ReactID) error {
	err := w.db.React.
		DeleteOneID(xid.ID(reactID)).
		Exec(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Internal))
	}

	return nil
}
