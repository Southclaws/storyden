package resolve_test

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Southclaws/storyden/app/resources/resolve"
)

func TestURL(t *testing.T) {
	for _, tt := range []struct {
		name, base, resource, identifier, want string
	}{
		{"admin", "https://example.com", "admin", "oauth-dcr-approval", "https://example.com/_/resolve/admin/oauth-dcr-approval"},
		{"base_path", "https://example.com/community", "admin", "oauth-dcr-approval", "https://example.com/community/_/resolve/admin/oauth-dcr-approval"},
		{"trailing_slash", "https://example.com/community/", "thread", "welcome", "https://example.com/community/_/resolve/thread/welcome"},
		{"escaped_base", "https://example.com/my%20community/", "node", "with space", "https://example.com/my%20community/_/resolve/node/with%20space"},
		{"literal_percent", "https://example.com", "node", "100%", "https://example.com/_/resolve/node/100%25"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			base, err := url.Parse(tt.base)
			require.NoError(t, err)

			require.Equal(t, tt.want, resolve.URL(*base, tt.resource, tt.identifier).String())
			require.Equal(t, tt.base, base.String())
		})
	}
}
