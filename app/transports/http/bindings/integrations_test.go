package bindings

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	oauthservice "github.com/Southclaws/storyden/app/services/authentication/oauth"
	"github.com/Southclaws/storyden/internal/config"
)

func mustParseURL(t *testing.T, raw string) url.URL {
	t.Helper()

	u, err := url.Parse(raw)
	require.NoError(t, err)

	return *u
}

func TestBuildIntegrationsDocumentMCPDisabled(t *testing.T) {
	t.Parallel()

	cfg := config.Config{
		PublicAPIAddress: mustParseURL(t, "https://example.com"),
		PublicWebAddress: mustParseURL(t, "https://example.com"),
		MCPEnabled:       false,
	}
	oauthBinding := OAuth{oauth: &oauthservice.Service{}}

	doc := buildIntegrationsDocument(cfg, oauthBinding)

	assert.Equal(t, 3, doc.Version)
	require.Len(t, doc.Surfaces, 1)

	api := doc.Surfaces[0]
	assert.Equal(t, "storyden-api", api.Slug)
	assert.Equal(t, "http", api.Type)
	assert.Equal(t, "https://example.com/api", api.URL)
	assert.Equal(t, "https://example.com/api/openapi.json", api.Spec)
	assert.Equal(t, "https://example.com/api/docs", api.Docs)
	assert.Equal(t, "https://example.com/.well-known/integrations.json", api.Basis.Source)

	require.Equal(t, "required", api.Auth.Status)
	require.Len(t, api.Auth.Entries, 1)
	use := api.Auth.Entries[0].Use[0]
	assert.Equal(t, accessKeyCredentialID, use.ID)
	assert.Equal(t, "http", use.Mechanics.Source)
	assert.Equal(t, "header", use.Mechanics.In)
	assert.Equal(t, "Authorization", use.Mechanics.HeaderName)
	assert.Equal(t, "Bearer", use.Mechanics.Scheme)

	require.Contains(t, doc.Credentials, accessKeyCredentialID)
	assert.NotContains(t, doc.Credentials, oauthCredentialID)
}

func TestBuildIntegrationsDocumentMCPEnabled(t *testing.T) {
	t.Parallel()

	cfg := config.Config{
		PublicAPIAddress: mustParseURL(t, "https://api.example.com"),
		PublicWebAddress: mustParseURL(t, "https://example.com"),
		MCPEnabled:       true,
	}
	oauthBinding := OAuth{oauth: &oauthservice.Service{}}

	doc := buildIntegrationsDocument(cfg, oauthBinding)

	require.Len(t, doc.Surfaces, 2)

	mcp := doc.Surfaces[1]
	assert.Equal(t, "storyden-mcp", mcp.Slug)
	assert.Equal(t, "mcp", mcp.Type)
	// /mcp is mounted on the backend's own mux (same as /api), so it must be
	// reachable via PublicAPIAddress, not the frontend's PublicWebAddress.
	assert.Equal(t, "https://api.example.com/mcp", mcp.URL)
	assert.Equal(t, []string{"streamable-http"}, mcp.Transports)
	assert.Equal(t, "https://example.com/.well-known/integrations.json", mcp.Basis.Source)

	require.Equal(t, "required", mcp.Auth.Status)
	require.Len(t, mcp.Auth.Entries, 1)
	assert.Equal(t, accessKeyCredentialID, mcp.Auth.Entries[0].Use[0].ID)
}
