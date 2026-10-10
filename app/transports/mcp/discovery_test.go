package mcp

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"regexp"
	"strings"
	"testing"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/santhosh-tekuri/jsonschema/v6"
	"github.com/stretchr/testify/require"

	"github.com/Southclaws/storyden/internal/config"
)

func discoveryConfig(enabled bool) config.Config {
	web, _ := url.Parse("https://community.example.org")
	api, _ := url.Parse("https://api.example.org")
	return config.Config{PublicWebAddress: *web, PublicAPIAddress: *api, MCPEnabled: enabled}
}

func TestDiscoveryDocuments(t *testing.T) {
	card, ai, api := discoveryDocuments(discoveryConfig(true), "Example community")
	require.Equal(t, serverCardSchema, card.Schema)
	require.Regexp(t, regexp.MustCompile(`^[a-zA-Z0-9.-]+/[a-zA-Z0-9._-]+$`), card.Name)
	require.Equal(t, "org.example.community/mcp", card.Name)
	require.Equal(t, "Example community", card.Title)
	require.Equal(t, "https://community.example.org", card.WebsiteURL)
	require.Equal(t, "streamable-http", card.Remotes[0].Type)
	require.Equal(t, "https://api.example.org/mcp", card.Remotes[0].URL)
	require.Equal(t, sdkmcp.SupportedProtocolVersions(), card.Remotes[0].SupportedProtocolVersions)
	require.Equal(t, "1.0", ai.SpecVersion)
	require.Equal(t, []aiCatalogEntry{{
		Identifier: "urn:air:community.example.org:mcp:storyden",
		Type:       serverCardMediaType,
		URL:        "https://api.example.org/mcp/server-card",
	}}, ai.Entries)
	require.Equal(t, "https://community.example.org/.well-known/api-catalog", api.Linkset[0].Anchor)
	require.Equal(t, "https://api.example.org/api", api.Linkset[0].Item[0].Href)
	require.Equal(t, "https://api.example.org/mcp", api.Linkset[0].Item[1].Href)
	require.Equal(t, "https://api.example.org/api/openapi.json", api.Linkset[1].ServiceDesc[0].Href)

	cardJSON, err := json.Marshal(card)
	require.NoError(t, err)
	for _, forbidden := range []string{`"tools"`, `"resources"`, `"prompts"`, `"headers"`} {
		require.NotContains(t, string(cardJSON), forbidden)
	}
}

func TestDiscoveryDisabled(t *testing.T) {
	card, ai, api := discoveryDocuments(discoveryConfig(false), "Storyden")
	require.Empty(t, card.Name)
	require.Empty(t, ai.Entries)
	require.Len(t, api.Linkset[0].Item, 1)
	require.Equal(t, "https://api.example.org/api", api.Linkset[0].Item[0].Href)
}

func TestServerCardMatchesPinnedV1Schema(t *testing.T) {
	data, err := os.ReadFile("testdata/server-card-v1.schema.json")
	require.NoError(t, err)
	var schemaDocument any
	require.NoError(t, json.Unmarshal(data, &schemaDocument))
	compiler := jsonschema.NewCompiler()
	require.NoError(t, compiler.AddResource("server-card-v1.schema.json", schemaDocument))
	schema, err := compiler.Compile("server-card-v1.schema.json#/$defs/ServerCard")
	require.NoError(t, err)

	for _, tc := range []struct {
		name  string
		cfg   config.Config
		title string
	}{
		{"separate origins", discoveryConfig(true), "Example community"},
		{"long title", discoveryConfig(true), strings.Repeat("A", 120)},
		{"localhost", func() config.Config { c := discoveryConfig(true); c.PublicWebAddress.Host = "localhost:3000"; return c }(), "Local"},
		{"ip address", func() config.Config { c := discoveryConfig(true); c.PublicWebAddress.Host = "127.0.0.1:3000"; return c }(), "Local"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			card, _, _ := discoveryDocuments(tc.cfg, tc.title)
			encoded, err := json.Marshal(card)
			require.NoError(t, err)
			var document any
			require.NoError(t, json.Unmarshal(encoded, &document))
			require.NoError(t, schema.Validate(document))
			if len(tc.title) > 100 {
				require.Empty(t, card.Title)
			}
		})
	}
}

func TestDiscoveryHandlerCacheAndCORS(t *testing.T) {
	card, _, _ := discoveryDocuments(discoveryConfig(true), "Example community")
	handler := discoveryHandler(card, serverCardMediaType, true)
	get := httptest.NewRecorder()
	handler.ServeHTTP(get, httptest.NewRequest(http.MethodGet, "/mcp/server-card", nil))
	require.Equal(t, http.StatusOK, get.Code)
	require.Equal(t, serverCardMediaType, get.Header().Get("Content-Type"))
	require.Equal(t, "public, max-age=3600", get.Header().Get("Cache-Control"))
	require.Equal(t, "*", get.Header().Get("Access-Control-Allow-Origin"))
	require.Equal(t, "GET", get.Header().Get("Access-Control-Allow-Methods"))
	require.Equal(t, "Content-Type, If-None-Match", get.Header().Get("Access-Control-Allow-Headers"))
	require.Equal(t, "ETag", get.Header().Get("Access-Control-Expose-Headers"))
	require.True(t, strings.HasPrefix(get.Body.String(), "{"))
	etag := get.Header().Get("ETag")
	require.NotEmpty(t, etag)

	conditional := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/mcp/server-card", nil)
	req.Header.Set("If-None-Match", etag)
	handler.ServeHTTP(conditional, req)
	require.Equal(t, http.StatusNotModified, conditional.Code)
	require.Empty(t, conditional.Body.String())
	require.Equal(t, etag, conditional.Header().Get("ETag"))
	weak := httptest.NewRecorder()
	weakReq := httptest.NewRequest(http.MethodGet, "/mcp/server-card", nil)
	weakReq.Header.Set("If-None-Match", `"previous", W/`+etag)
	handler.ServeHTTP(weak, weakReq)
	require.Equal(t, http.StatusNotModified, weak.Code)

	head := httptest.NewRecorder()
	handler.ServeHTTP(head, httptest.NewRequest(http.MethodHead, "/mcp/server-card", nil))
	require.Equal(t, http.StatusOK, head.Code)
	require.Empty(t, head.Body.String())
}

func TestDiscoveryRoutesWhenMCPDisabled(t *testing.T) {
	mux := http.NewServeMux()
	mountDiscovery(mux, discoveryConfig(false), "Storyden")

	card := httptest.NewRecorder()
	mux.ServeHTTP(card, httptest.NewRequest(http.MethodGet, "/mcp/server-card", nil))
	require.Equal(t, http.StatusNotFound, card.Code)

	ai := httptest.NewRecorder()
	mux.ServeHTTP(ai, httptest.NewRequest(http.MethodGet, "/.well-known/ai-catalog.json", nil))
	require.Equal(t, http.StatusOK, ai.Code)
	require.JSONEq(t, `{"specVersion":"1.0","entries":[]}`, ai.Body.String())

	api := httptest.NewRecorder()
	mux.ServeHTTP(api, httptest.NewRequest(http.MethodHead, "/.well-known/api-catalog", nil))
	require.Equal(t, http.StatusOK, api.Code)
	require.Equal(t, apiCatalogMediaType, api.Header().Get("Content-Type"))
	require.Equal(t, "<https://community.example.org/.well-known/api-catalog>; rel=api-catalog", api.Header().Get("Link"))
	require.Empty(t, api.Body.String())
}

func TestDiscoveryRoutesAllowBrowserPreflight(t *testing.T) {
	mux := http.NewServeMux()
	mountDiscovery(mux, discoveryConfig(true), "Example community")
	for _, path := range []string{"/mcp/server-card", "/.well-known/ai-catalog.json", "/.well-known/api-catalog"} {
		t.Run(path, func(t *testing.T) {
			preflight := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodOptions, path, nil)
			req.Header.Set("Origin", "https://community.example.org")
			req.Header.Set("Access-Control-Request-Method", "GET")
			req.Header.Set("Access-Control-Request-Headers", "if-none-match")
			mux.ServeHTTP(preflight, req)
			require.Equal(t, http.StatusNoContent, preflight.Code)
			require.Equal(t, "*", preflight.Header().Get("Access-Control-Allow-Origin"))
			require.Contains(t, preflight.Header().Get("Access-Control-Allow-Headers"), "If-None-Match")
			require.Empty(t, preflight.Body.String())
		})
	}
}
