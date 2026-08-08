package token

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Southclaws/opt"
	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/internal/infrastructure/cache/cachetest"
)

type fakeRepo struct {
	session *Session
	err     error
	calls   int
}

func (f *fakeRepo) Issue(context.Context, account.AccountID) (*Session, error) { return f.session, nil }
func (f *fakeRepo) Revoke(context.Context, Token) error                        { return nil }

func (f *fakeRepo) Validate(context.Context, Token) (*Validated, error) {
	f.calls++
	if f.err != nil {
		return nil, f.err
	}
	return (*Validated)(f.session), nil
}

func liveSession(t Token) Session {
	return Session{
		Token:     t,
		AccountID: account.AccountID(xid.New()),
		ExpiresAt: time.Now().Add(time.Hour),
	}
}

func seed(t *testing.T, store *cachetest.Store, s Session) {
	t.Helper()

	payload, err := s.Serialise()
	require.NoError(t, err)
	require.NoError(t, store.Set(context.Background(), cacheKey(s.Token), string(payload), time.Hour))
}

func TestCachedValidateRejectsRevokedCachedSession(t *testing.T) {
	t.Parallel()

	store := cachetest.New()
	inner := &fakeRepo{}
	repo := NewCachedRepository(inner, store)

	tok := Generate()
	revoked := liveSession(tok)
	revoked.RevokedAt = opt.New(time.Now())
	seed(t, store, revoked)

	v, err := repo.Validate(context.Background(), tok)

	require.Error(t, err, "a revoked cached session must not validate")
	assert.True(t, errors.Is(err, ErrTokenRevoked))
	assert.Nil(t, v)
	assert.Zero(t, inner.calls, "a revoked cached session must not fall through to the database")
}

func TestCachedValidateRejectsExpiredCachedSession(t *testing.T) {
	t.Parallel()

	store := cachetest.New()
	inner := &fakeRepo{}
	repo := NewCachedRepository(inner, store)

	tok := Generate()
	expired := liveSession(tok)
	expired.ExpiresAt = time.Now().Add(-time.Minute)
	seed(t, store, expired)

	v, err := repo.Validate(context.Background(), tok)

	require.Error(t, err)
	assert.True(t, errors.Is(err, ErrTokenExpired))
	assert.Nil(t, v)
}

func TestCachedValidateFallsBackWhenEntryIsUnreadable(t *testing.T) {
	t.Parallel()

	store := cachetest.New()

	tok := Generate()
	stored := liveSession(tok)
	inner := &fakeRepo{session: &stored}
	repo := NewCachedRepository(inner, store)

	require.NoError(t, store.Set(context.Background(), cacheKey(tok), "not json", time.Hour))

	v, err := repo.Validate(context.Background(), tok)

	require.NoError(t, err)
	require.NotNil(t, v, "an unreadable entry must never yield a nil session with a nil error")
	assert.Equal(t, stored.AccountID, v.AccountID)
	assert.Equal(t, 1, inner.calls)
}

func TestCachedValidateServesValidCachedSession(t *testing.T) {
	t.Parallel()

	store := cachetest.New()
	inner := &fakeRepo{}
	repo := NewCachedRepository(inner, store)

	tok := Generate()
	stored := liveSession(tok)
	seed(t, store, stored)

	v, err := repo.Validate(context.Background(), tok)

	require.NoError(t, err)
	require.NotNil(t, v)
	assert.Equal(t, stored.AccountID, v.AccountID)
	assert.Zero(t, inner.calls, "a usable cached session must not hit the database")
}
