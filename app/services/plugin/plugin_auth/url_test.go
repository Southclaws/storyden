package plugin_auth

import (
	"net/url"
	"strings"
	"testing"

	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	resource_plugin "github.com/Southclaws/storyden/app/resources/plugin"
)

func TestBuildAndParseConnectionURL(t *testing.T) {
	baseURL, _ := url.Parse("http://localhost:8000/rpc")
	pluginID := resource_plugin.InstallationID(xid.New())

	secret, err := GenerateSecret()
	require.NoError(t, err)

	connectionURL, err := BuildConnectionURL(*baseURL, pluginID, secret)
	require.NoError(t, err)
	require.NotNil(t, connectionURL)

	urlStr := connectionURL.String()
	assert.Contains(t, urlStr, "plugin_id=")
	assert.Contains(t, urlStr, "token=")
	assert.Contains(t, urlStr, xid.ID(pluginID).String())

	params, err := ParseConnectionURL(connectionURL)
	require.NoError(t, err)
	require.NotNil(t, params)

	id, ok := params.PluginID.Get()
	assert.True(t, ok)
	assert.Equal(t, pluginID, id)
	assert.NotEmpty(t, params.Token)
}

func TestParseConnectionURLMissingParams(t *testing.T) {
	validExternalToken := ExternalTokenPrefix + strings.Repeat("A", secretLength)

	tests := []struct {
		name    string
		urlStr  string
		wantErr string
	}{
		{
			name:    "missing plugin_id is invalid for non-external token",
			urlStr:  "http://localhost:8000/rpc?token=abc123",
			wantErr: "missing plugin_id query parameter",
		},
		{
			name:    "missing plugin_id is valid for external tokens",
			urlStr:  "http://localhost:8000/rpc?token=" + validExternalToken,
			wantErr: "",
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
			if tt.wantErr == "" {
				require.NoError(t, err)
				return
			}

			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantErr)
		})
	}
}

func TestBuildConnectionURLUsesWSSForHTTPS(t *testing.T) {
	baseURL, _ := url.Parse("https://api.example.com:9000")
	pluginID := resource_plugin.InstallationID(xid.New())

	secret, err := GenerateSecret()
	require.NoError(t, err)

	connectionURL, err := BuildConnectionURL(*baseURL, pluginID, secret)
	require.NoError(t, err)

	assert.Equal(t, "wss", connectionURL.Scheme)
	assert.Equal(t, "api.example.com:9000", connectionURL.Host)
	assert.Equal(t, "/rpc", connectionURL.Path)
}

func TestBuildConnectionURLUsesWSForHTTP(t *testing.T) {
	baseURL, _ := url.Parse("http://api.example.com:9000")
	pluginID := resource_plugin.InstallationID(xid.New())

	secret, err := GenerateSecret()
	require.NoError(t, err)

	connectionURL, err := BuildConnectionURL(*baseURL, pluginID, secret)
	require.NoError(t, err)

	assert.Equal(t, "ws", connectionURL.Scheme)
	assert.Equal(t, "api.example.com:9000", connectionURL.Host)
	assert.Equal(t, "/rpc", connectionURL.Path)
}
