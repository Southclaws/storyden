package token

import (
	"context"
	"testing"
	"time"

	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/internal/infrastructure/cache/cachetest"
)

type revokeRacingRepo struct {
	session      Session
	revoked      bool
	duringRevoke func()
}

func (f *revokeRacingRepo) Issue(context.Context, account.AccountID) (*Session, error) {
	return &f.session, nil
}

func (f *revokeRacingRepo) Revoke(context.Context, Token) error {
	if f.duringRevoke != nil {
		f.duringRevoke()
	}
	f.revoked = true
	return nil
}

func (f *revokeRacingRepo) Validate(context.Context, Token) (*Validated, error) {
	if f.revoked {
		return nil, ErrTokenRevoked
	}
	s := f.session
	return (*Validated)(&s), nil
}

func TestCachedRevokeEvictsSessionCachedDuringTheWrite(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	store := cachetest.New()

	tok := Generate()
	inner := &revokeRacingRepo{
		session: Session{
			Token:     tok,
			AccountID: account.AccountID(xid.New()),
			ExpiresAt: time.Now().Add(90 * 24 * time.Hour),
		},
	}

	repo := NewCachedRepository(inner, store)

	inner.duringRevoke = func() {
		_, err := repo.Validate(ctx, tok)
		require.NoError(t, err)
	}

	require.NoError(t, repo.Revoke(ctx, tok))
	require.True(t, inner.revoked)

	_, err := store.Get(ctx, tok.String())
	assert.Error(t, err, "the pre-revocation session must not survive in the cache")

	v, err := repo.Validate(ctx, tok)
	require.Error(t, err, "a revoked session must not validate after logout")
	assert.Nil(t, v)
}
