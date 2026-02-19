package reaction

import (
	"context"
	"database/sql"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/fault/ftag"
	"github.com/pkg/errors"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/role/role_repo"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/ent/react"
)

var ErrInvalidEmoji = errors.New("invalid emoji codepoint")

func New(db *ent.Client, roleQuerier *role_repo.Repository) (*Querier, *Writer) {
	q := &Querier{db: db, roleQuerier: roleQuerier}
	w := &Writer{db, q}
	return q, w
}

type Querier struct {
	db          *ent.Client
	roleQuerier *role_repo.Repository
}

func (q *Querier) Get(ctx context.Context, reactID ReactID) (*React, error) {
	r, err := q.db.React.Query().
		Where(react.ID(xid.ID(reactID))).
		WithAccount().
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			err = fault.Wrap(err, ftag.With(ftag.NotFound))
		}
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return q.mapSingle(ctx, r)
}

func (q *Querier) Lookup(ctx context.Context, accountID account.AccountID, postID xid.ID, e string) (*React, error) {
	r, err := q.db.React.Query().
		Where(
			react.AccountID(xid.ID(accountID)),
			react.PostID(xid.ID(postID)),
			react.Emoji(e),
		).
		WithAccount().
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			err = fault.Wrap(err, ftag.With(ftag.NotFound))
		}
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return q.mapSingle(ctx, r)
}

func (q *Querier) mapSingle(ctx context.Context, in *ent.React) (*React, error) {
	accountEdge, err := in.Edges.AccountOrErr()
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	roleHydrator, err := q.roleQuerier.BuildSingleHydrator(ctx, accountEdge)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	out, err := Map(in, roleHydrator.Hydrate)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return out, nil
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

		if ent.IsNotFound(err) {
			err = fault.Wrap(err, ftag.With(ftag.NotFound),
				fmsg.WithDesc("post not found", "The post you are trying to react to does not exist."),
			)
		}

		return nil, fault.Wrap(err, fctx.With(ctx))
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
