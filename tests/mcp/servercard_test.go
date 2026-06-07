package mcp_test

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

	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
)

func TestMCPServerCardDisabled(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		ts *httptest.Server,
	) {
		lc.Append(fx.StartHook(func() {
			r := require.New(t)
			a := assert.New(t)

			req, err := http.NewRequestWithContext(root, http.MethodGet, ts.URL+"/.well-known/mcp/server-card.json", nil)
			r.NoError(err)

			resp, err := http.DefaultClient.Do(req)
			r.NoError(err)
			defer resp.Body.Close()

			a.Equal(http.StatusNotFound, resp.StatusCode)
		}))
	}))
}

func TestMCPServerCardEnabled(t *testing.T) {
	t.Parallel()

	publicWeb, err := url.Parse("http://localhost:3000")
	require.NoError(t, err)
	publicAPI, err := url.Parse("http://localhost:8000")
	require.NoError(t, err)

	cfg := &config.Config{
		MCPEnabled:       true,
		PublicWebAddress: *publicWeb,
		PublicAPIAddress: *publicAPI,
	}

	integration.Test(t, cfg, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		ts *httptest.Server,
	) {
		lc.Append(fx.StartHook(func() {
			r := require.New(t)
			a := assert.New(t)

			req, err := http.NewRequestWithContext(root, http.MethodGet, ts.URL+"/.well-known/mcp/server-card.json", nil)
			r.NoError(err)

			resp, err := http.DefaultClient.Do(req)
			r.NoError(err)
			defer resp.Body.Close()

			r.Equal(http.StatusOK, resp.StatusCode)
			a.Equal("application/json", resp.Header.Get("Content-Type"))
			a.Equal("*", resp.Header.Get("Access-Control-Allow-Origin"))
			a.Equal("public, max-age=3600", resp.Header.Get("Cache-Control"))

			var card map[string]any
			r.NoError(json.NewDecoder(resp.Body).Decode(&card))

			a.Equal("https://static.modelcontextprotocol.io/schemas/v1/server-card.schema.json", card["$schema"])
			a.Equal("org.storyden/storyden", card["name"])
			a.NotEmpty(card["version"])
			a.NotEmpty(card["description"])

			remotes, ok := card["remotes"].([]any)
			r.True(ok, "remotes should be an array")
			r.NotEmpty(remotes)

			remote, ok := remotes[0].(map[string]any)
			r.True(ok, "first remote should be an object")
			a.Equal("sse", remote["type"])
			a.Contains(remote["url"], "/mcp/sse")
		}))
	}))
}
