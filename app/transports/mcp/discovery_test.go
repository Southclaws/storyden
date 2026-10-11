package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/Southclaws/opt"
	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/santhosh-tekuri/jsonschema/v6"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/settings"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/tests"
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
	require.Equal(t, "https://api.example.org/.well-known/api-catalog", api.Linkset[0].Anchor)
	require.Equal(t, "https://api.example.org/api", api.Linkset[0].Item[0].Href)
	require.Equal(t, "https://api.example.org/mcp", api.Linkset[0].Item[1].Href)
	require.Equal(t, "https://api.example.org/api/openapi.json", api.Linkset[1].ServiceDesc[0].Href)

	cardJSON, err := json.Marshal(card)
	require.NoError(t, err)
	for _, forbidden := range []string{`"tools"`, `"resources"`, `"prompts"`, `"headers"`} {
		require.NotContains(t, string(cardJSON), forbidden)
	}

	for _, tc := range []struct {
		host string
		want string
	}{
		{"", "localhost/mcp"},
		{"localhost", "localhost/mcp"},
		{"127.0.0.1", "localhost/mcp"},
		{"::1", "localhost/mcp"},
		{"Community.Example.ORG", "org.example.community/mcp"},
		{"community.example.org.", "org.example.community/mcp"},
		{"bücher.example", "example.xn--bcher-kva/mcp"},
		{"xn--bcher-kva.example", "example.xn--bcher-kva/mcp"},
		{"bücher。example", "example.xn--bcher-kva/mcp"},
		{"bücher。example。", "example.xn--bcher-kva/mcp"},
		{"invalid_host", "localhost/mcp"},
		{"-invalid.example", "localhost/mcp"},
		{"invalid-.example", "localhost/mcp"},
		{"example..org", "localhost/mcp"},
		{".example.org", "localhost/mcp"},
		{"example.org..", "localhost/mcp"},
		{"xn--.example", "localhost/mcp"},
		{strings.Repeat("a", 63) + ".example", "example." + strings.Repeat("a", 63) + "/mcp"},
		{strings.Repeat("a", 64) + ".example", "localhost/mcp"},
		{strings.Repeat("a", 63) + "." + strings.Repeat("b", 63) + "." + strings.Repeat("c", 63) + ".test", "test." + strings.Repeat("c", 63) + "." + strings.Repeat("b", 63) + "." + strings.Repeat("a", 63) + "/mcp"},
		{strings.Repeat("a.", 98) + "b", "localhost/mcp"},
	} {
		t.Run("server name for "+tc.host, func(t *testing.T) {
			require.Equal(t, tc.want, serverName(tc.host))
		})
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
	require.Equal(t, "no-cache, must-revalidate", get.Header().Get("Cache-Control"))
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
	if tests.IsSharedPostgresDatabase() {
		t.Skip("skipping global discovery settings mutations on shared postgres database")
	}

	cfg := discoveryConfig(false)
	integration.Test(t, &cfg, e2e.Setup(), fx.Invoke(MountDiscovery), fx.Invoke(func(
		lc fx.Lifecycle,
		repo *settings.SettingsRepository,
		ts *httptest.Server,
	) {
		lc.Append(fx.StartHook(func(ctx context.Context) {
			r := require.New(t)

			_, err := repo.Set(ctx, enabledDiscoverySettings())
			r.NoError(err)

			card, err := ts.Client().Get(ts.URL + "/mcp/server-card")
			r.NoError(err)
			defer card.Body.Close()
			r.Equal(http.StatusNotFound, card.StatusCode)

			ai, err := ts.Client().Get(ts.URL + "/.well-known/ai-catalog.json")
			r.NoError(err)
			defer ai.Body.Close()
			r.Equal(http.StatusOK, ai.StatusCode)

			body, err := io.ReadAll(ai.Body)
			r.NoError(err)
			r.JSONEq(`{"specVersion":"1.0","entries":[]}`, string(body))

			api, err := ts.Client().Head(ts.URL + "/.well-known/api-catalog")
			r.NoError(err)
			defer api.Body.Close()
			r.Equal(http.StatusOK, api.StatusCode)
			r.Equal(`application/linkset+json; profile="https://www.rfc-editor.org/info/rfc9727"`, api.Header.Get("Content-Type"))
			r.Equal("<https://api.example.org/.well-known/api-catalog>; rel=api-catalog", api.Header.Get("Link"))

			body, err = io.ReadAll(api.Body)
			r.NoError(err)
			r.Empty(body)
		}))
	}))
}

func TestDiscoveryRoutesAllowBrowserPreflight(t *testing.T) {
	if tests.IsSharedPostgresDatabase() {
		t.Skip("skipping global discovery settings mutations on shared postgres database")
	}

	cfg := discoveryConfig(true)
	integration.Test(t, &cfg, e2e.Setup(), fx.Invoke(MountDiscovery), fx.Invoke(func(
		lc fx.Lifecycle,
		repo *settings.SettingsRepository,
		ts *httptest.Server,
	) {
		lc.Append(fx.StartHook(func(ctx context.Context) {
			_, err := repo.Set(ctx, enabledDiscoverySettings())
			require.NoError(t, err)

			for _, path := range []string{"/mcp/server-card", "/.well-known/ai-catalog.json", "/.well-known/api-catalog"} {
				t.Run(path, func(t *testing.T) {
					r := require.New(t)

					req, err := http.NewRequestWithContext(ctx, http.MethodOptions, ts.URL+path, nil)
					r.NoError(err)
					req.Header.Set("Origin", "https://community.example.org")
					req.Header.Set("Access-Control-Request-Method", "GET")
					req.Header.Set("Access-Control-Request-Headers", "if-none-match")

					resp, err := ts.Client().Do(req)
					r.NoError(err)
					defer resp.Body.Close()
					r.Equal(http.StatusNoContent, resp.StatusCode)
					r.Equal("*", resp.Header.Get("Access-Control-Allow-Origin"))
					r.Equal("GET", resp.Header.Get("Access-Control-Allow-Methods"))
					r.Equal("Content-Type, If-None-Match", resp.Header.Get("Access-Control-Allow-Headers"))

					body, err := io.ReadAll(resp.Body)
					r.NoError(err)
					r.Empty(body)
				})
			}
		}))
	}))
}

func TestAPICatalogLinksetWireFormat(t *testing.T) {
	if tests.IsSharedPostgresDatabase() {
		t.Skip("skipping global discovery settings mutations on shared postgres database")
	}

	for _, enabled := range []bool{false, true} {
		t.Run(fmt.Sprintf("MCP enabled=%t", enabled), func(t *testing.T) {
			cfg := discoveryConfig(enabled)
			integration.Test(t, &cfg, e2e.Setup(), fx.Invoke(MountDiscovery), fx.Invoke(func(
				lc fx.Lifecycle,
				repo *settings.SettingsRepository,
				ts *httptest.Server,
			) {
				lc.Append(fx.StartHook(func(ctx context.Context) {
					r := require.New(t)

					_, err := repo.Set(ctx, enabledDiscoverySettings())
					r.NoError(err)

					resp, err := ts.Client().Get(ts.URL + "/.well-known/api-catalog")
					r.NoError(err)
					defer resp.Body.Close()
					r.Equal(http.StatusOK, resp.StatusCode)
					r.Equal(`application/linkset+json; profile="https://www.rfc-editor.org/info/rfc9727"`, resp.Header.Get("Content-Type"))
					r.Equal("<https://api.example.org/.well-known/api-catalog>; rel=api-catalog", resp.Header.Get("Link"))

					items := `{"href":"https://api.example.org/api"}`
					if enabled {
						items += `, {"href":"https://api.example.org/mcp"}`
					}

					body, err := io.ReadAll(resp.Body)
					r.NoError(err)
					r.JSONEq(fmt.Sprintf(`{
						"linkset": [
							{
								"anchor": "https://api.example.org/.well-known/api-catalog",
								"item": [%s]
							},
							{
								"anchor": "https://api.example.org/api",
								"service-desc": [{
									"href": "https://api.example.org/api/openapi.json",
									"type": "application/vnd.oai.openapi+json"
								}]
							}
						]
					}`, items), string(body))
				}))
			}))
		})
	}
}

func TestDiscoveryRoutesDisabledByDefault(t *testing.T) {
	if tests.IsSharedPostgresDatabase() {
		t.Skip("skipping global discovery defaults assertions on shared postgres database")
	}

	cfg := discoveryConfig(true)
	integration.Test(t, &cfg, e2e.Setup(), fx.Invoke(MountDiscovery), fx.Invoke(func(
		lc fx.Lifecycle,
		ts *httptest.Server,
	) {
		lc.Append(fx.StartHook(func(ctx context.Context) {
			for _, path := range []string{"/mcp/server-card", "/.well-known/ai-catalog.json", "/.well-known/api-catalog"} {
				for _, method := range []string{http.MethodGet, http.MethodHead, http.MethodOptions} {
					t.Run(method+" "+path, func(t *testing.T) {
						r := require.New(t)

						req, err := http.NewRequestWithContext(ctx, method, ts.URL+path, nil)
						r.NoError(err)

						resp, err := ts.Client().Do(req)
						r.NoError(err)
						defer resp.Body.Close()
						r.Equal(http.StatusNotFound, resp.StatusCode)
						r.Equal("no-store", resp.Header.Get("Cache-Control"))
					})
				}
			}
		}))
	}))
}

func TestDiscoverySettingsApplyWithoutRestart(t *testing.T) {
	if tests.IsSharedPostgresDatabase() {
		t.Skip("skipping global discovery settings mutations on shared postgres database")
	}

	cfg := discoveryConfig(true)
	integration.Test(t, &cfg, e2e.Setup(), fx.Invoke(MountDiscovery), fx.Invoke(func(
		lc fx.Lifecycle,
		repo *settings.SettingsRepository,
		ts *httptest.Server,
	) {
		lc.Append(fx.StartHook(func(ctx context.Context) {
			paths := []struct {
				path   string
				policy settings.DiscoveryServiceSettings
			}{
				{"/mcp/server-card", settings.DiscoveryServiceSettings{MCPServerCard: opt.New(false)}},
				{"/.well-known/ai-catalog.json", settings.DiscoveryServiceSettings{AICatalog: opt.New(false)}},
				{"/.well-known/api-catalog", settings.DiscoveryServiceSettings{APICatalog: opt.New(false)}},
			}

			for _, tc := range paths {
				t.Run(tc.path, func(t *testing.T) {
					r := require.New(t)

					_, err := repo.Set(ctx, enabledDiscoverySettings())
					r.NoError(err)

					before, err := ts.Client().Get(ts.URL + tc.path)
					r.NoError(err)
					defer before.Body.Close()
					r.Equal(http.StatusOK, before.StatusCode)
					r.Equal("no-cache, must-revalidate", before.Header.Get("Cache-Control"))
					etag := before.Header.Get("ETag")
					r.NotEmpty(etag)

					conditional, err := http.NewRequestWithContext(ctx, http.MethodGet, ts.URL+tc.path, nil)
					r.NoError(err)
					conditional.Header.Set("If-None-Match", etag)

					unchanged, err := ts.Client().Do(conditional)
					r.NoError(err)
					defer unchanged.Body.Close()
					r.Equal(http.StatusNotModified, unchanged.StatusCode)
					r.Equal(etag, unchanged.Header.Get("ETag"))

					_, err = repo.Set(ctx, settings.Settings{
						Services: opt.New(settings.ServiceSettings{Discovery: opt.New(tc.policy)}),
					})
					r.NoError(err)

					for _, method := range []string{http.MethodGet, http.MethodHead, http.MethodOptions} {
						req, err := http.NewRequestWithContext(ctx, method, ts.URL+tc.path, nil)
						r.NoError(err)
						req.Header.Set("If-None-Match", etag)

						resp, err := ts.Client().Do(req)
						r.NoError(err)
						defer resp.Body.Close()
						r.Equal(http.StatusNotFound, resp.StatusCode, method)
						r.Equal("no-store", resp.Header.Get("Cache-Control"))
					}

					for _, other := range paths {
						if other.path == tc.path {
							continue
						}

						resp, err := ts.Client().Get(ts.URL + other.path)
						r.NoError(err)
						defer resp.Body.Close()
						r.Equal(http.StatusOK, resp.StatusCode, other.path)
					}
				})
			}

			t.Run("card visibility and title update", func(t *testing.T) {
				r := require.New(t)

				_, err := repo.Set(ctx, enabledDiscoverySettings())
				r.NoError(err)

				_, err = repo.Set(ctx, settings.Settings{
					Services: opt.New(settings.ServiceSettings{
						Discovery: opt.New(settings.DiscoveryServiceSettings{MCPServerCard: opt.New(false)}),
					}),
				})
				r.NoError(err)

				catalog, err := ts.Client().Get(ts.URL + "/.well-known/ai-catalog.json")
				r.NoError(err)
				defer catalog.Body.Close()
				r.Equal(http.StatusOK, catalog.StatusCode)

				body, err := io.ReadAll(catalog.Body)
				r.NoError(err)
				r.JSONEq(`{"specVersion":"1.0","entries":[]}`, string(body))

				_, err = repo.Set(ctx, settings.Settings{
					Title: opt.New("After"),
					Services: opt.New(settings.ServiceSettings{
						Discovery: opt.New(settings.DiscoveryServiceSettings{MCPServerCard: opt.New(true)}),
					}),
				})
				r.NoError(err)

				card, err := ts.Client().Get(ts.URL + "/mcp/server-card")
				r.NoError(err)
				defer card.Body.Close()
				r.Equal(http.StatusOK, card.StatusCode)

				var document map[string]any
				r.NoError(json.NewDecoder(card.Body).Decode(&document))
				r.Equal("After", document["title"])
			})
		}))
	}))
}

func enabledDiscoverySettings() settings.Settings {
	return settings.Settings{
		Title: opt.New("Example community"),
		Services: opt.New(settings.ServiceSettings{
			Discovery: opt.New(settings.DiscoveryServiceSettings{
				MCPServerCard: opt.New(true),
				AICatalog:     opt.New(true),
				APICatalog:    opt.New(true),
			}),
		}),
	}
}
