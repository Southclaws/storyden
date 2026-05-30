package api

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Southclaws/storyden/cmd/sd/internal/config"
)

func TestNewAuthenticatedClient(t *testing.T) {
	t.Run("returns human readable rate limit errors", func(t *testing.T) {
		r := require.New(t)

		resetAt := time.Now().UTC().Add(2 * time.Minute).Format(time.RFC1123)
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
			switch request.URL.Path {
			case "/api/threads":
				w.Header().Set("X-RateLimit-Limit", "100")
				w.Header().Set("X-RateLimit-Remaining", "0")
				w.Header().Set("X-RateLimit-Reset", resetAt)
				w.Header().Set("Retry-After", resetAt)
				http.Error(w, "Too Many Requests", http.StatusTooManyRequests)

			default:
				http.NotFound(w, request)
			}
		}))
		defer server.Close()

		store := config.NewFileStoreAt(filepath.Join(t.TempDir(), "storyden", "config.yaml"))
		cfg := config.New()
		cfg.UpsertContext("local", config.Context{
			APIURL:   server.URL,
			AuthType: config.AuthStorageFile,
			Auth: &config.Auth{
				Method:      config.AuthMethodAccessKey,
				AccessToken: "sdak_test_access_key",
				TokenType:   "Bearer",
			},
		})
		cfg.SetCurrentContext("local")
		r.NoError(store.Save(cfg))

		client, err := NewAuthenticatedClient(context.Background(), store, WithRateLimitWarnings(io.Discard))
		r.NoError(err)

		response, err := client.OpenAPI.ThreadListWithResponse(context.Background(), nil)

		r.Nil(response)
		r.ErrorContains(err, "Rate limit exceeded.")
		r.ErrorContains(err, "Limit: 100 requests")
		r.ErrorContains(err, "Reset time:")
		r.ErrorContains(err, "(in about")
		r.ErrorContains(err, "Please wait for the reset window")
	})

	t.Run("warns when rate limit is almost exhausted", func(t *testing.T) {
		r := require.New(t)

		resetAt := time.Now().UTC().Add(5 * time.Minute).Format(time.RFC1123)
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
			switch request.URL.Path {
			case "/api/threads":
				w.Header().Set("X-RateLimit-Limit", "100")
				w.Header().Set("X-RateLimit-Remaining", "5")
				w.Header().Set("X-RateLimit-Reset", resetAt)
				w.Header().Set("Content-Type", "application/json")
				r.NoError(json.NewEncoder(w).Encode(map[string]any{
					"current_page": 1,
					"page_size":    25,
					"results":      0,
					"threads":      []any{},
					"total_pages":  1,
				}))

			default:
				http.NotFound(w, request)
			}
		}))
		defer server.Close()

		store := config.NewFileStoreAt(filepath.Join(t.TempDir(), "storyden", "config.yaml"))
		cfg := config.New()
		cfg.UpsertContext("local", config.Context{
			APIURL:   server.URL,
			AuthType: config.AuthStorageFile,
			Auth: &config.Auth{
				Method:      config.AuthMethodAccessKey,
				AccessToken: "sdak_test_access_key",
				TokenType:   "Bearer",
			},
		})
		cfg.SetCurrentContext("local")
		r.NoError(store.Save(cfg))

		var warnings bytes.Buffer
		client, err := NewAuthenticatedClient(context.Background(), store, WithRateLimitWarnings(&warnings))
		r.NoError(err)

		response, err := client.OpenAPI.ThreadListWithResponse(context.Background(), nil)

		r.NoError(err)
		r.NotNil(response.JSON200)
		r.Contains(warnings.String(), "rate limit is getting low")
		r.Contains(warnings.String(), "5/100 requests remaining")
		r.Contains(warnings.String(), "Reset time:")
		r.Contains(warnings.String(), "(in about")
		r.Contains(warnings.String(), "Consider waiting")
	})

	t.Run("uses access key without OAuth discovery or refresh", func(t *testing.T) {
		a := assert.New(t)
		r := require.New(t)

		var sawThreadRequest bool
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
			switch request.URL.Path {
			case "/api/threads":
				sawThreadRequest = true
				a.Equal("Bearer sdak_test_access_key", request.Header.Get("Authorization"))

				w.Header().Set("Content-Type", "application/json")
				r.NoError(json.NewEncoder(w).Encode(map[string]any{
					"current_page": 1,
					"page_size":    25,
					"results":      0,
					"threads":      []any{},
					"total_pages":  1,
				}))

			case "/.well-known/openid-configuration", "/api/oauth/token":
				t.Fatalf("access key auth should not call %s", request.URL.Path)

			default:
				http.NotFound(w, request)
			}
		}))
		defer server.Close()

		store := config.NewFileStoreAt(filepath.Join(t.TempDir(), "storyden", "config.yaml"))
		cfg := config.New()
		cfg.UpsertContext("local", config.Context{
			APIURL:   server.URL,
			AuthType: config.AuthStorageFile,
			Auth: &config.Auth{
				Method:      config.AuthMethodAccessKey,
				AccessToken: "sdak_test_access_key",
				TokenType:   "Bearer",
			},
		})
		cfg.SetCurrentContext("local")
		r.NoError(store.Save(cfg))

		client, err := NewAuthenticatedClient(context.Background(), store)
		r.NoError(err)

		response, err := client.OpenAPI.ThreadListWithResponse(context.Background(), nil)
		r.NoError(err)
		r.NotNil(response.JSON200)
		a.True(sawThreadRequest)
	})

	t.Run("does not retry unauthorized access key requests", func(t *testing.T) {
		a := assert.New(t)
		r := require.New(t)

		var threadRequests int
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
			switch request.URL.Path {
			case "/api/threads":
				threadRequests++
				a.Equal("Bearer invalid-access-key", request.Header.Get("Authorization"))
				http.Error(w, "invalid access key", http.StatusUnauthorized)

			case "/.well-known/openid-configuration", "/api/oauth/token":
				t.Fatalf("access key auth should not call %s", request.URL.Path)

			default:
				http.NotFound(w, request)
			}
		}))
		defer server.Close()

		store := config.NewFileStoreAt(filepath.Join(t.TempDir(), "storyden", "config.yaml"))
		cfg := config.New()
		cfg.UpsertContext("local", config.Context{
			APIURL:   server.URL,
			AuthType: config.AuthStorageFile,
			Auth: &config.Auth{
				Method:      config.AuthMethodAccessKey,
				AccessToken: "invalid-access-key",
				TokenType:   "Bearer",
			},
		})
		cfg.SetCurrentContext("local")
		r.NoError(store.Save(cfg))

		client, err := NewAuthenticatedClient(context.Background(), store)
		r.NoError(err)

		response, err := client.OpenAPI.ThreadListWithResponse(context.Background(), nil)
		r.NoError(err)
		r.Equal(http.StatusUnauthorized, response.StatusCode())
		a.Equal(1, threadRequests)
	})

	t.Run("refreshes expired access token before request and stores rotation", func(t *testing.T) {
		a := assert.New(t)
		r := require.New(t)

		var sawTokenRequest bool
		var sawThreadRequest bool

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
			switch request.URL.Path {
			case "/.well-known/openid-configuration":
				writeDiscovery(t, w, request)

			case "/api/oauth/token":
				sawTokenRequest = true
				r.NoError(request.ParseForm())
				a.Equal("refresh_token", request.Form.Get("grant_type"))
				a.Equal("storyden-cli", request.Form.Get("client_id"))
				a.Equal("old-refresh-token", request.Form.Get("refresh_token"))

				w.Header().Set("Content-Type", "application/json")
				r.NoError(json.NewEncoder(w).Encode(map[string]any{
					"access_token":  "new-access-token",
					"refresh_token": "new-refresh-token",
					"expires_in":    900,
					"scope":         "openid profile offline_access READ_THREADS",
					"token_type":    "Bearer",
				}))

			case "/api/threads":
				sawThreadRequest = true
				a.Equal("Bearer new-access-token", request.Header.Get("Authorization"))

				w.Header().Set("Content-Type", "application/json")
				r.NoError(json.NewEncoder(w).Encode(map[string]any{
					"current_page": 1,
					"page_size":    25,
					"results":      0,
					"threads":      []any{},
					"total_pages":  1,
				}))

			default:
				http.NotFound(w, request)
			}
		}))
		defer server.Close()

		configPath := filepath.Join(t.TempDir(), "storyden", "config.yaml")
		store := config.NewStoreAtWithCredentialStore(configPath, newMemoryCredentialStore())
		cfg := config.New()
		cfg.UpsertContext("local", config.Context{
			APIURL:   server.URL,
			AuthType: config.AuthStorageCredentialStore,
			Auth: &config.Auth{
				AccessToken:  "old-access-token",
				RefreshToken: "old-refresh-token",
				TokenType:    "Bearer",
				ExpiresAt:    time.Now().Add(-time.Minute),
			},
		})
		cfg.SetCurrentContext("local")
		r.NoError(store.Save(cfg))

		client, err := NewAuthenticatedClient(context.Background(), store)
		r.NoError(err)

		response, err := client.OpenAPI.ThreadListWithResponse(context.Background(), nil)
		r.NoError(err)
		r.NotNil(response.JSON200)

		loaded, err := store.Load()
		r.NoError(err)
		auth := loaded.Contexts["local"].Auth
		r.NotNil(auth)

		a.True(sawTokenRequest)
		a.True(sawThreadRequest)
		a.Equal("new-access-token", auth.AccessToken)
		a.Equal("new-refresh-token", auth.RefreshToken)
		a.Equal("openid profile offline_access READ_THREADS", auth.Scope)
		a.True(auth.ExpiresAt.After(time.Now()))

		data, err := os.ReadFile(configPath)
		r.NoError(err)
		a.Contains(string(data), "auth_type: credential_store")
		a.NotContains(string(data), "new-access-token")
		a.NotContains(string(data), "new-refresh-token")
	})

	t.Run("refreshes and retries once after unauthorized response", func(t *testing.T) {
		a := assert.New(t)
		r := require.New(t)

		var threadRequests int

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
			switch request.URL.Path {
			case "/.well-known/openid-configuration":
				writeDiscovery(t, w, request)

			case "/api/oauth/token":
				r.NoError(request.ParseForm())
				a.Equal("refresh_token", request.Form.Get("grant_type"))
				a.Equal("old-refresh-token", request.Form.Get("refresh_token"))

				w.Header().Set("Content-Type", "application/json")
				r.NoError(json.NewEncoder(w).Encode(map[string]any{
					"access_token":  "new-access-token",
					"refresh_token": "new-refresh-token",
					"expires_in":    900,
					"scope":         "openid profile offline_access READ_THREADS",
					"token_type":    "Bearer",
				}))

			case "/api/threads":
				threadRequests++
				switch threadRequests {
				case 1:
					a.Equal("Bearer old-access-token", request.Header.Get("Authorization"))
					http.Error(w, "token revoked", http.StatusUnauthorized)

				case 2:
					a.Equal("Bearer new-access-token", request.Header.Get("Authorization"))
					w.Header().Set("Content-Type", "application/json")
					r.NoError(json.NewEncoder(w).Encode(map[string]any{
						"current_page": 1,
						"page_size":    25,
						"results":      0,
						"threads":      []any{},
						"total_pages":  1,
					}))
				}

			default:
				http.NotFound(w, request)
			}
		}))
		defer server.Close()

		store := config.NewFileStoreAt(filepath.Join(t.TempDir(), "storyden", "config.yaml"))
		cfg := config.New()
		cfg.UpsertContext("local", config.Context{
			APIURL:   server.URL,
			AuthType: config.AuthStorageFile,
			Auth: &config.Auth{
				AccessToken:  "old-access-token",
				RefreshToken: "old-refresh-token",
				TokenType:    "Bearer",
				ExpiresAt:    time.Now().Add(time.Hour),
			},
		})
		cfg.SetCurrentContext("local")
		r.NoError(store.Save(cfg))

		client, err := NewAuthenticatedClient(context.Background(), store)
		r.NoError(err)

		response, err := client.OpenAPI.ThreadListWithResponse(context.Background(), nil)
		r.NoError(err)
		r.NotNil(response.JSON200)

		a.Equal(2, threadRequests)

		loaded, err := store.Load()
		r.NoError(err)
		auth := loaded.Contexts["local"].Auth
		r.NotNil(auth)
		a.Equal("new-access-token", auth.AccessToken)
		a.Equal("new-refresh-token", auth.RefreshToken)
	})
}

type memoryCredentialStore struct {
	auth map[string]config.Auth
}

func newMemoryCredentialStore() *memoryCredentialStore {
	return &memoryCredentialStore{auth: map[string]config.Auth{}}
}

func (m *memoryCredentialStore) SetAuth(contextName string, auth config.Auth) error {
	m.auth[contextName] = auth
	return nil
}

func (m *memoryCredentialStore) GetAuth(contextName string) (config.Auth, bool, error) {
	auth, ok := m.auth[contextName]
	return auth, ok, nil
}

func (m *memoryCredentialStore) DeleteAuth(contextName string) error {
	delete(m.auth, contextName)
	return nil
}

func (m *memoryCredentialStore) Available() bool { return true }
