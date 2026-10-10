package password_reset

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewLinkTemplateOriginValidation(t *testing.T) {
	t.Parallel()

	allowed := url.URL{Scheme: "https", Host: "community.example.com"}

	t.Run("accepts matching origin and builds link", func(t *testing.T) {
		lt, err := NewLinkTemplate("https://community.example.com/reset", "token", allowed)
		require.NoError(t, err)
		assert.Equal(t, "https://community.example.com/reset?token=abc", lt.GetURL("abc"))
	})

	t.Run("accepts matching host case-insensitively", func(t *testing.T) {
		_, err := NewLinkTemplate("https://Community.Example.COM/reset", "token", allowed)
		require.NoError(t, err)
	})

	t.Run("accepts host including port", func(t *testing.T) {
		local := url.URL{Scheme: "http", Host: "localhost:3000"}
		_, err := NewLinkTemplate("http://localhost:3000/reset", "token", local)
		require.NoError(t, err)
	})

	t.Run("rejects scheme downgrade to http", func(t *testing.T) {
		_, err := NewLinkTemplate("http://community.example.com/reset", "token", allowed)
		require.ErrorIs(t, err, ErrResetURLOriginMismatch)
	})

	t.Run("rejects attacker host", func(t *testing.T) {
		_, err := NewLinkTemplate("https://evil.example.com/reset", "token", allowed)
		require.ErrorIs(t, err, ErrResetURLOriginMismatch)
	})

	t.Run("rejects userinfo-smuggled host", func(t *testing.T) {
		_, err := NewLinkTemplate("https://community.example.com@evil.com/reset", "token", allowed)
		require.ErrorIs(t, err, ErrResetURLOriginMismatch)
	})

	t.Run("rejects suffixed lookalike host", func(t *testing.T) {
		_, err := NewLinkTemplate("https://community.example.com.evil.com/reset", "token", allowed)
		require.ErrorIs(t, err, ErrResetURLOriginMismatch)
	})

	t.Run("rejects relative url with no host", func(t *testing.T) {
		_, err := NewLinkTemplate("/reset", "token", allowed)
		require.ErrorIs(t, err, ErrResetURLOriginMismatch)
	})

	t.Run("rejects javascript scheme", func(t *testing.T) {
		_, err := NewLinkTemplate("javascript:alert(1)", "token", allowed)
		require.ErrorIs(t, err, ErrResetURLOriginMismatch)
	})
}