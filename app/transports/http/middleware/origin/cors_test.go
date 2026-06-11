package origin

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Southclaws/storyden/internal/config"
)

func mustURL(t *testing.T, raw string) url.URL {
	t.Helper()
	u, err := url.Parse(raw)
	assert.NoError(t, err)
	return *u
}

func TestAllowedOrigins(t *testing.T) {
	t.Run("trusts web and api addresses, rejects arbitrary", func(t *testing.T) {
		cfg := config.Config{
			PublicWebAddress: mustURL(t, "https://app.example.com"),
			PublicAPIAddress: mustURL(t, "https://api.example.com"),
		}

		got := allowedOrigins(cfg)

		assert.Contains(t, got, "https://app.example.com")
		assert.Contains(t, got, "https://api.example.com")
		assert.NotContains(t, got, "https://evil.example.com")
	})

	t.Run("includes configured extras and dedupes", func(t *testing.T) {
		cfg := config.Config{
			PublicWebAddress:   mustURL(t, "https://app.example.com"),
			PublicAPIAddress:   mustURL(t, "https://app.example.com"),
			CORSAllowedOrigins: []string{"https://admin.example.com", " https://app.example.com "},
		}

		got := allowedOrigins(cfg)

		assert.ElementsMatch(t, []string{"https://app.example.com", "https://admin.example.com"}, got)
	})

	t.Run("normalises scheme and host casing, drops path", func(t *testing.T) {
		cfg := config.Config{
			PublicWebAddress:   mustURL(t, "https://app.example.com"),
			CORSAllowedOrigins: []string{"HTTPS://Admin.Example.com/some/path"},
		}

		got := allowedOrigins(cfg)

		assert.Contains(t, got, "https://admin.example.com")
	})

	t.Run("skips empty and unparseable entries", func(t *testing.T) {
		cfg := config.Config{
			PublicWebAddress:   mustURL(t, "https://app.example.com"),
			CORSAllowedOrigins: []string{"", "   ", "://nohost", "notaurl"},
		}

		got := allowedOrigins(cfg)

		assert.Equal(t, []string{"https://app.example.com"}, got)
	})
}

func TestNormaliseOrigin(t *testing.T) {
	assert.Equal(t, "https://app.example.com", normaliseOrigin("https://app.example.com"))
	assert.Equal(t, "http://localhost:3000", normaliseOrigin("http://localhost:3000"))
	assert.Equal(t, "https://app.example.com", normaliseOrigin("HTTPS://App.Example.com/cb?x=1"))
	assert.Equal(t, "", normaliseOrigin(""))
	assert.Equal(t, "", normaliseOrigin("notaurl"))
	assert.Equal(t, "", normaliseOrigin("https://"))
}
