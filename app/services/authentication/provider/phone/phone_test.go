package phone

import (
	"context"
	"testing"
	"time"

	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Southclaws/storyden/app/resources/account/authentication"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/infrastructure/cache/cachetest"
)

type stubAuthRepo struct {
	authentication.Repository
	tokens map[authentication.ID]string
}

func (s *stubAuthRepo) Update(ctx context.Context, id authentication.ID, opts ...authentication.Option) (*authentication.Authentication, error) {
	m := &ent.AuthenticationMutation{}
	for _, opt := range opts {
		opt(m)
	}
	token, _ := m.Token()
	s.tokens[id] = token
	return nil, nil
}

func newTestProvider() (*Provider, *stubAuthRepo) {
	repo := &stubAuthRepo{tokens: map[authentication.ID]string{}}
	return &Provider{auth: repo, cache: cachetest.New()}, repo
}

func record(code string, age time.Duration) *authentication.Authentication {
	return &authentication.Authentication{
		ID:      xid.New(),
		Created: time.Now().Add(-age),
		Token:   code,
	}
}

func TestVerifyCode(t *testing.T) {
	ctx := context.Background()

	t.Run("accepts the correct code and rotates it", func(t *testing.T) {
		p, repo := newTestProvider()
		rec := record("123456", 0)

		require.NoError(t, p.verifyCode(ctx, rec, "123456"))
		rotated, ok := repo.tokens[rec.ID]
		require.True(t, ok)
		assert.NotEmpty(t, rotated)
		assert.NotEqual(t, "123456", rotated)
	})

	t.Run("rejects an expired code", func(t *testing.T) {
		p, repo := newTestProvider()
		rec := record("123456", codeLifespan+time.Minute)

		assert.Error(t, p.verifyCode(ctx, rec, "123456"))
		assert.Empty(t, repo.tokens)
	})

	t.Run("locks out after too many attempts even with the right code", func(t *testing.T) {
		p, repo := newTestProvider()
		rec := record("123456", 0)

		for i := 0; i < maxCodeAttempts; i++ {
			assert.Error(t, p.verifyCode(ctx, rec, "000000"))
		}

		assert.Error(t, p.verifyCode(ctx, rec, "123456"))
		assert.Empty(t, repo.tokens)
	})
}
