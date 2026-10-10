package password_reset

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewLinkTemplateHostValidation(t *testing.T) {
	t.Parallel()

	const host = "community.example.com"

	t.Run("accepts matching host and builds link", func(t *testing.T) {
		lt, err := NewLinkTemplate("https://community.example.com/reset", "token", host)
		require.NoError(t, err)
		assert.Equal(t, "https://community.example.com/reset?token=abc", lt.GetURL("abc"))
	})

	t.Run("accepts matching host case-insensitively", func(t *testing.T) {
		_, err := NewLinkTemplate("https://Community.Example.COM/reset", "token", host)
		require.NoError(t, err)
	})

	t.Run("accepts host including port", func(t *testing.T) {
		_, err := NewLinkTemplate("http://localhost:3000/reset", "token", "localhost:3000")
		require.NoError(t, err)
	})

	t.Run("rejects attacker host", func(t *testing.T) {
		_, err := NewLinkTemplate("https://evil.example.com/reset", "token", host)
		require.ErrorIs(t, err, ErrResetURLHostMismatch)
	})

	t.Run("rejects userinfo-smuggled host", func(t *testing.T) {
		_, err := NewLinkTemplate("https://community.example.com@evil.com/reset", "token", host)
		require.ErrorIs(t, err, ErrResetURLHostMismatch)
	})

	t.Run("rejects suffixed lookalike host", func(t *testing.T) {
		_, err := NewLinkTemplate("https://community.example.com.evil.com/reset", "token", host)
		require.ErrorIs(t, err, ErrResetURLHostMismatch)
	})

	t.Run("rejects relative url with no host", func(t *testing.T) {
		_, err := NewLinkTemplate("/reset", "token", host)
		require.ErrorIs(t, err, ErrResetURLHostMismatch)
	})

	t.Run("rejects javascript scheme", func(t *testing.T) {
		_, err := NewLinkTemplate("javascript:alert(1)", "token", host)
		require.ErrorIs(t, err, ErrResetURLHostMismatch)
	})
}