package webauthn

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Southclaws/storyden/internal/config"
)

func TestNewSeparatesRelyingPartyDomainFromOrigin(t *testing.T) {
	for _, tc := range []struct {
		origin string
		rpID   string
	}{
		{"http://localhost:3000", "localhost"},
		{"https://community.example.com:8443", "community.example.com"},
		{"https://community.example.com", "community.example.com"},
	} {
		t.Run(tc.origin, func(t *testing.T) {
			address, err := url.Parse(tc.origin)
			require.NoError(t, err)
			provider, err := New(config.Config{PublicWebAddress: *address})
			require.NoError(t, err)
			require.Equal(t, tc.rpID, provider.Config.RPID)
			require.Equal(t, []string{tc.origin}, provider.Config.RPOrigins)
		})
	}
}
