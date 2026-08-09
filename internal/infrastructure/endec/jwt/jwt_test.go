package jwt

import (
	"encoding/base64"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/infrastructure/endec"
)

func TestEncryptDecrypt(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)

	claims := endec.Claims{
		"sub": "test-subject",
		"exp": float64(time.Now().Add(1 * time.Hour).Unix()),
	}

	ed, err := New(config.Config{
		JWTSecret: []byte("07d422e512b23a056ccc953994d1593f"),
	})
	r.NoError(err)

	t.Run("encrypt and decrypt payload", func(t *testing.T) {
		token, err := ed.Encrypt(endec.PurposePasswordReset, claims, time.Hour)
		a.NoError(err)

		gotClaims, err := ed.Decrypt(endec.PurposePasswordReset, token)
		a.NoError(err)
		a.Equal(claims["sub"], gotClaims["sub"])
		a.Equal(claims["exp"], gotClaims["exp"])
	})

	t.Run("invalid_secret", func(t *testing.T) {
		ed, err := New(config.Config{
			JWTSecret: []byte{},
		})
		r.NoError(err)
		r.Nil(ed)
	})
}

func TestPurposeIsNotInterchangeable(t *testing.T) {
	t.Parallel()

	ed, err := New(config.Config{JWTSecret: []byte("07d422e512b23a056ccc953994d1593f")})
	require.NoError(t, err)

	reset, err := ed.Encrypt(endec.PurposePasswordReset, endec.Claims{"account_id": "cv1l2p1cpetc0gm4nqm0"}, time.Hour)
	require.NoError(t, err)

	_, err = ed.Decrypt(endec.PurposeOAuthState, reset)
	assert.Error(t, err, "a password reset token must not be accepted as an oauth state value")

	state, err := ed.Encrypt(endec.PurposeOAuthState, endec.Claims{"redirect": "https://example.com/callback"}, time.Minute*10)
	require.NoError(t, err)

	_, err = ed.Decrypt(endec.PurposePasswordReset, state)
	assert.Error(t, err, "an oauth state value must not be accepted as a password reset token")
}

func TestPurposeCannotBeForgedFromClaims(t *testing.T) {
	t.Parallel()

	ed, err := New(config.Config{JWTSecret: []byte("07d422e512b23a056ccc953994d1593f")})
	require.NoError(t, err)

	// typ lives in the header, so a caller cannot reach it through Claims
	token, err := ed.Encrypt(endec.PurposePasswordReset, endec.Claims{
		"typ": string(endec.PurposeOAuthState),
	}, time.Hour)
	require.NoError(t, err)

	_, err = ed.Decrypt(endec.PurposeOAuthState, token)
	assert.Error(t, err, "claims must not be able to forge the purpose")

	_, err = ed.Decrypt(endec.PurposePasswordReset, token)
	assert.NoError(t, err)
}

func TestPurposeIsCarriedInTheTypHeader(t *testing.T) {
	t.Parallel()

	ed, err := New(config.Config{JWTSecret: []byte("07d422e512b23a056ccc953994d1593f")})
	require.NoError(t, err)

	token, err := ed.Encrypt(endec.PurposeOAuthState, endec.Claims{"redirect": "/"}, time.Minute)
	require.NoError(t, err)

	segments := strings.Split(token, ".")
	require.Len(t, segments, 3)

	raw, err := base64.RawURLEncoding.DecodeString(segments[0])
	require.NoError(t, err)

	var header map[string]any
	require.NoError(t, json.Unmarshal(raw, &header))

	assert.Equal(t, string(endec.PurposeOAuthState), header["typ"])
}

func TestEncryptRequiresAPurpose(t *testing.T) {
	t.Parallel()

	ed, err := New(config.Config{JWTSecret: []byte("07d422e512b23a056ccc953994d1593f")})
	require.NoError(t, err)

	_, err = ed.Encrypt("", endec.Claims{"sub": "nobody"}, time.Hour)
	assert.Error(t, err)
}
