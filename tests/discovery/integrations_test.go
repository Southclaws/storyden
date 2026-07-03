package discovery_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	transportmcp "github.com/Southclaws/storyden/app/transports/mcp"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
)

type integrationsSurface struct {
	Slug string `json:"slug"`
	Type string `json:"type"`
	URL  string `json:"url"`
	Spec string `json:"spec"`
}

type integrationsDocument struct {
	Version  int                   `json:"version"`
	Surfaces []integrationsSurface `json:"surfaces"`
}

func TestIntegrationsDiscoveryMCPDisabled(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		ts *httptest.Server,
	) {
		lc.Append(fx.StartHook(func() {
			a := assert.New(t)
			r := require.New(t)

			resp := getJSON(t, root, ts.URL+"/.well-known/integrations.json")
			defer resp.Body.Close()
			r.Equal(http.StatusOK, resp.StatusCode)
			a.Equal("application/json", resp.Header.Get("Content-Type"))
			a.Contains(resp.Header.Get("Cache-Control"), "public")

			var doc integrationsDocument
			r.NoError(json.NewDecoder(resp.Body).Decode(&doc))
			a.Equal(3, doc.Version)
			r.Len(doc.Surfaces, 1)
			a.Equal("storyden-api", doc.Surfaces[0].Slug)
			a.Equal("http", doc.Surfaces[0].Type)

			cardResp := getJSON(t, root, ts.URL+"/.well-known/mcp/server-card.json")
			defer cardResp.Body.Close()
			a.Equal(http.StatusNotFound, cardResp.StatusCode)
		}))
	}))
}

func TestIntegrationsDiscoveryMCPEnabled(t *testing.T) {
	t.Parallel()

	integration.Test(t, mcpConfig(t), e2e.Setup(), transportmcp.Build(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		ts *httptest.Server,
	) {
		lc.Append(fx.StartHook(func() {
			a := assert.New(t)
			r := require.New(t)

			resp := getJSON(t, root, ts.URL+"/.well-known/integrations.json")
			defer resp.Body.Close()
			r.Equal(http.StatusOK, resp.StatusCode)

			var doc integrationsDocument
			r.NoError(json.NewDecoder(resp.Body).Decode(&doc))
			r.Len(doc.Surfaces, 2)
			a.Equal("storyden-mcp", doc.Surfaces[1].Slug)
			a.Equal("mcp", doc.Surfaces[1].Type)

			cardResp := getJSON(t, root, ts.URL+"/.well-known/mcp/server-card.json")
			defer cardResp.Body.Close()
			r.Equal(http.StatusOK, cardResp.StatusCode)
			a.Equal("application/json", cardResp.Header.Get("Content-Type"))

			var card map[string]string
			r.NoError(json.NewDecoder(cardResp.Body).Decode(&card))
			r.Contains(card, "url")
			a.Contains(card["url"], "/mcp")
		}))
	}))
}

func getJSON(t *testing.T, ctx context.Context, target string) *http.Response {
	t.Helper()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, target, nil)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	return resp
}

func mcpConfig(t *testing.T) *config.Config {
	t.Helper()

	publicWebAddress, err := url.Parse("http://localhost:3000")
	require.NoError(t, err)
	publicAPIAddress, err := url.Parse("http://localhost:8000")
	require.NoError(t, err)

	return &config.Config{
		PublicWebAddress: *publicWebAddress,
		PublicAPIAddress: *publicAPIAddress,
		MCPEnabled:       true,
	}
}
