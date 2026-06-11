package datagraph

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCanonicalResolvePath(t *testing.T) {
	assert.Equal(t, "/_/resolve/thread/my-thread-abc123", CanonicalResolvePath(KindThread, "my-thread-abc123"))
	assert.Equal(t, "/_/resolve/profile/southclaws", CanonicalResolvePath(KindProfile, "southclaws"))
	assert.Equal(t, "/_/resolve/reply/cabcdef0123456789", CanonicalResolvePath(KindReply, "cabcdef0123456789"))
}

func TestCanonicalResolveURL(t *testing.T) {
	t.Run("root path", func(t *testing.T) {
		web, err := url.Parse("https://example.com")
		require.NoError(t, err)
		assert.Equal(t,
			"https://example.com/_/resolve/thread/my-thread-abc123",
			CanonicalResolveURL(*web, KindThread, "my-thread-abc123").String(),
		)
	})

	t.Run("base path with trailing slash", func(t *testing.T) {
		web, err := url.Parse("https://example.com/community/")
		require.NoError(t, err)
		assert.Equal(t,
			"https://example.com/community/_/resolve/node/handbook",
			CanonicalResolveURL(*web, KindNode, "handbook").String(),
		)
	})

	t.Run("escapes the identifier", func(t *testing.T) {
		web, err := url.Parse("https://example.com")
		require.NoError(t, err)
		assert.Equal(t,
			"https://example.com/_/resolve/profile/with%20space",
			CanonicalResolveURL(*web, KindProfile, "with space").String(),
		)
	})
}
