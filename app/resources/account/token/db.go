package token

import (
	"context"
	"time"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/ent/session"
)

// Sessions last 90 days.
// TODO: Make this configurable and match the session expiry in the cookie.
var Expiry = 24 * time.Hour * 90

type persistedRepository struct {
	db *ent.Client
}

func New(
	db *ent.Client,
) Repository {
	return &persistedRepository{
		db: db,
	}
}

func (r *persistedRepository) Issue(ctx context.Context, accountID account.AccountID) (*Session, error) {
	token := Token{xid.New()}

	create := r.db.Session.Create().
		SetID(token.ID).
		SetAccountID(xid.ID(accountID)).
		SetExpiresAt(time.Now().Add(Expiry))

	result, err := create.Save(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return Map(result), nil
}

func (r *persistedRepository) Revoke(ctx context.Context, id Token) error {
	update := r.db.Session.Update().Where(session.ID(id.ID))

	update.SetRevokedAt(time.Now())

	err := update.Exec(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}

func (r *persistedRepository) Validate(ctx context.Context, t Token) (*Validated, error) {
	query := r.db.Session.Query().Where(session.ID(t.ID))

	result, err := query.Only(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	v, err := Map(result).Validate()
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return v, nil
}

func Map(s *ent.Session) *Session {
	return &Session{
		Token:     Token{s.ID},
		AccountID: account.AccountID(s.AccountID),
		ExpiresAt: s.ExpiresAt,
		RevokedAt: opt.NewPtr(s.RevokedAt),
	}
}
