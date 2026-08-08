package github

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/infrastructure/endec"
	endecjwt "github.com/Southclaws/storyden/internal/infrastructure/endec/jwt"
)

func TestLoginRejectsPasswordResetTokenAsState(t *testing.T) {
	t.Parallel()

	cfg := config.Config{GitHubEnabled: true, JWTSecret: []byte("07d422e512b23a056ccc953994d1593f")}

	ed, err := endecjwt.New(cfg)
	require.NoError(t, err)

	p, err := New(cfg, nil, ed)
	require.NoError(t, err)

	reset, err := ed.Encrypt(endec.PurposePasswordReset, endec.Claims{
		"account_id": "cv1l2p1cpetc0gm4nqm0",
	}, time.Hour)
	require.NoError(t, err)

	assert.NotPanics(t, func() {
		_, err := p.Login(context.Background(), reset, "any-code")
		assert.Error(t, err)
	})
}

func TestLoginRejectsStateWithoutRedirect(t *testing.T) {
	t.Parallel()

	cfg := config.Config{GitHubEnabled: true, JWTSecret: []byte("07d422e512b23a056ccc953994d1593f")}

	ed, err := endecjwt.New(cfg)
	require.NoError(t, err)

	p, err := New(cfg, nil, ed)
	require.NoError(t, err)

	state, err := ed.Encrypt(endec.PurposeOAuthState, endec.Claims{}, time.Minute*10)
	require.NoError(t, err)

	assert.NotPanics(t, func() {
		_, err := p.Login(context.Background(), state, "any-code")
		assert.Error(t, err)
	})
}
