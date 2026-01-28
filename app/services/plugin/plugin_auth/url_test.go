package plugin_auth

import (
	"net/url"
	"testing"

	resource_plugin "github.com/Southclaws/storyden/app/resources/plugin"
	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildAndParseConnectionURL(t *testing.T) {
	baseURL := "http://localhost:8000/rpc"
	pluginID := resource_plugin.InstallationID(xid.New())

	secret, err := GenerateSecret()
	require.NoError(t, err)

	connectionURL, err := BuildConnectionURL(baseURL, pluginID, secret)
	require.NoError(t, err)
	require.NotNil(t, connectionURL)

	urlStr := connectionURL.String()
	assert.Contains(t, urlStr, "plugin_id=")
	assert.Contains(t, urlStr, "token=")
	assert.Contains(t, urlStr, xid.ID(pluginID).String())

	params, err := ParseConnectionURL(connectionURL)
	require.NoError(t, err)
	require.NotNil(t, params)

	assert.Equal(t, pluginID, params.PluginID)
	assert.NotEmpty(t, params.Token)
}

func TestParseConnectionURLMissingParams(t *testing.T) {
	tests := []struct {
		name    string
		urlStr  string
		wantErr string
	}{
		{
			name:    "missing plugin_id",
			urlStr:  "http://localhost:8000/rpc?token=abc123",
			wantErr: "missing plugin_id query parameter",
		},
		{
			name:    "missing token",
			urlStr:  "http://localhost:8000/rpc?plugin_id=d5sgqg5o2dto5q61jnd0",
			wantErr: "missing token query parameter",
		},
		{
			name:    "invalid plugin_id format",
			urlStr:  "http://localhost:8000/rpc?plugin_id=invalid&token=abc123",
			wantErr: "invalid plugin ID format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u, err := url.Parse(tt.urlStr)
			require.NoError(t, err)

			_, err = ParseConnectionURL(u)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantErr)
		})
	}
}

func TestBuildConnectionURLPreservesSchemeAndHost(t *testing.T) {
	baseURL := "https://api.example.com:9000/rpc"
	pluginID := resource_plugin.InstallationID(xid.New())

	secret, err := GenerateSecret()
	require.NoError(t, err)

	connectionURL, err := BuildConnectionURL(baseURL, pluginID, secret)
	require.NoError(t, err)

	assert.Equal(t, "https", connectionURL.Scheme)
	assert.Equal(t, "api.example.com:9000", connectionURL.Host)
	assert.Equal(t, "/rpc", connectionURL.Path)
}
