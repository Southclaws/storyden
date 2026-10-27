package phone

import (
	"context"
	"testing"
	"time"

	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Southclaws/storyden/app/resources/account/authentication"
	"github.com/Southclaws/storyden/internal/infrastructure/cache/cachetest"
)

type stubAuthRepo struct {
	authentication.Repository
	tokens map[authentication.ID]string
}

func (s *stubAuthRepo) Update(ctx context.Context, id authentication.ID, opts ...authentication.Option) (*authentication.Authentication, error) {
	s.tokens[id] = "rotated-" + id.String()
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
		assert.Contains(t, repo.tokens, rec.ID)
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
