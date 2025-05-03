package jwt

import (
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
		token, err := ed.Encrypt(claims, time.Hour)
		a.NoError(err)

		gotClaims, err := ed.Decrypt(token)
		a.NoError(err)
		a.Equal(claims["sub"], gotClaims["sub"])
		a.Equal(claims["exp"], gotClaims["exp"])
	})

	t.Run("invalid_secret", func(t *testing.T) {
		ed, err := New(config.Config{
			JWTSecret: []byte{},
		})
		r.NoError(err)

		token, err := ed.Encrypt(claims, time.Hour)
		r.Error(err)
		a.EqualError(err, "no JWT secret provided")

		gotClaims, err := ed.Decrypt(token)
		r.Error(err)
		a.EqualError(err, "token is malformed: token contains an invalid number of segments: token is malformed")
		a.Nil(gotClaims)
	})
}
