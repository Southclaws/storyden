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

type keyStubRepo struct{ session Session }

func (f *keyStubRepo) Issue(context.Context, account.AccountID) (*Session, error) {
	return &f.session, nil
}
func (f *keyStubRepo) Revoke(context.Context, Token) error { return nil }
func (f *keyStubRepo) Validate(context.Context, Token) (*Validated, error) {
	s := f.session
	return (*Validated)(&s), nil
}

func TestSessionCacheKeysAreNamespaced(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	store := cachetest.New()

	tok := Generate()
	inner := &keyStubRepo{session: Session{
		Token:     tok,
		AccountID: account.AccountID(xid.New()),
		ExpiresAt: time.Now().Add(time.Hour),
	}}

	repo := NewCachedRepository(inner, store)

	_, err := repo.Validate(ctx, tok)
	require.NoError(t, err)

	_, err = store.Get(ctx, tok.String())
	assert.Error(t, err, "the bare token must not be usable as a cache key")

	_, err = store.Get(ctx, cachePrefix+tok.String())
	assert.NoError(t, err)
}

func TestSessionCacheReadWriteAndDeleteAgree(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	store := cachetest.New()

	tok := Generate()
	inner := &keyStubRepo{session: Session{
		Token:     tok,
		AccountID: account.AccountID(xid.New()),
		ExpiresAt: time.Now().Add(time.Hour),
	}}

	repo := NewCachedRepository(inner, store)

	_, err := repo.Validate(ctx, tok)
	require.NoError(t, err)

	require.NoError(t, repo.Revoke(ctx, tok))

	_, err = store.Get(ctx, cachePrefix+tok.String())
	assert.Error(t, err, "revoke must delete the key that cache and get use")
}
