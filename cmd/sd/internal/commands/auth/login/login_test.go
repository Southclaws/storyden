package login

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Southclaws/storyden/app/services/authentication/oauth"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/cmd/sd/internal/config"
	"github.com/Southclaws/storyden/cmd/sd/internal/tui"
)

func TestContextName(t *testing.T) {
	t.Run("uses host slug", func(t *testing.T) {
		a := assert.New(t)

		cfg := config.New()

		a.Equal("forum-example-com", contextName(cfg, "https://forum.example.com/api"))
	})

	t.Run("reuses context for the same API URL", func(t *testing.T) {
		a := assert.New(t)

		cfg := config.New()
		cfg.UpsertContext("forum-example-com", config.Context{APIURL: "https://forum.example.com/api"})

		a.Equal("forum-example-com", contextName(cfg, "https://forum.example.com/api"))
	})

	t.Run("suffixes host when existing context points elsewhere", func(t *testing.T) {
		a := assert.New(t)

		cfg := config.New()
		cfg.UpsertContext("forum-example-com", config.Context{APIURL: "https://forum.example.com/api"})

		a.Equal("forum-example-com-2", contextName(cfg, "https://forum.example.com"))
	})
}

func TestAccessKeyStdinExample(t *testing.T) {
	t.Run("windows uses powershell", func(t *testing.T) {
		shell, command := accessKeyStdinExample("windows")

		assert.Equal(t, "powershell", shell)
		assert.Equal(t, `$env:STORYDEN_ACCESS_KEY | sd auth login http://localhost:8000 --access-key-stdin --auth-storage file`, command)
	})

	t.Run("macos uses zsh", func(t *testing.T) {
		shell, command := accessKeyStdinExample("darwin")

		assert.Equal(t, "zsh", shell)
		assert.Equal(t, `printf '%s' "$STORYDEN_ACCESS_KEY" | sd auth login http://localhost:8000 --access-key-stdin --auth-storage file`, command)
	})

	t.Run("linux uses bash", func(t *testing.T) {
		shell, command := accessKeyStdinExample("linux")

		assert.Equal(t, "bash", shell)
		assert.Equal(t, `printf '%s' "$STORYDEN_ACCESS_KEY" | sd auth login http://localhost:8000 --access-key-stdin --auth-storage file`, command)
	})
}

func TestLoginCommand(t *testing.T) {
	t.Run("stores access key auth without OAuth discovery", func(t *testing.T) {
		a := assert.New(t)
		r := require.New(t)

		configPath := filepath.Join(t.TempDir(), "storyden", "config.yaml")
		store := config.NewFileStoreAt(configPath)
		command := (*cobra.Command)(New(store))
		command.SetArgs([]string{
			"http://localhost:8000",
			"--access-key-stdin",
			"--auth-storage", "file",
		})
		command.SetIn(strings.NewReader("sdak_test_access_key\n"))
		var stdout bytes.Buffer
		command.SetOut(&stdout)

		r.NoError(command.Execute())

		cfg, err := store.Load()
		r.NoError(err)
		a.Equal("localhost-8000", cfg.CurrentContext)

		ctx := cfg.Contexts["localhost-8000"]
		a.Equal("http://localhost:8000", ctx.APIURL)
		a.Equal(config.AuthStorageFile, ctx.AuthType)
		r.NotNil(ctx.Auth)
		a.Equal(config.AuthMethodAccessKey, ctx.Auth.Method)
		a.Equal("sdak_test_access_key", ctx.Auth.AccessToken)
		a.Equal("Bearer", ctx.Auth.TokenType)
		a.Empty(ctx.Auth.RefreshToken)
		a.Contains(stdout.String(), "Authenticated with")
	})

	t.Run("stores API URL as current context", func(t *testing.T) {
		a := assert.New(t)
		r := require.New(t)

		var deviceAuthorizationCalled bool
		var tokenCalled bool
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
			switch request.URL.Path {
			case "/.well-known/openid-configuration":
				writeJSON(t, w, map[string]any{
					"issuer":                                "https://storyden.example",
					"authorization_endpoint":                serverURL(request) + "/api/oauth/authorize",
					"device_authorization_endpoint":         serverURL(request) + "/api/oauth/device_authorization",
					"token_endpoint":                        serverURL(request) + "/api/oauth/token",
					"userinfo_endpoint":                     serverURL(request) + "/api/oauth/userinfo",
					"jwks_uri":                              serverURL(request) + "/api/oauth/jwks",
					"response_types_supported":              []string{"code"},
					"grant_types_supported":                 []string{oauth.GrantTypeDeviceCode},
					"scopes_supported":                      []string{"openid", "profile", "offline_access"},
					"subject_types_supported":               []string{"public"},
					"id_token_signing_alg_values_supported": []string{"RS256"},
					"code_challenge_methods_supported":      []string{"S256"},
				})

			case "/api/oauth/device_authorization":
				deviceAuthorizationCalled = true
				r.NoError(request.ParseForm())
				a.Equal(oauth.StorydenCLIClientID, request.Form.Get("client_id"))
				a.Equal(deviceAuthScope, request.Form.Get("scope"))
				writeJSON(t, w, map[string]any{
					"device_code":               "device-code",
					"user_code":                 "ABCD-EFGH",
					"verification_uri":          serverURL(request) + "/api/oauth/consent",
					"verification_uri_complete": serverURL(request) + "/api/oauth/consent?user_code=ABCD-EFGH",
					"expires_in":                600,
					"interval":                  1,
				})

			case "/api/oauth/token":
				tokenCalled = true
				r.NoError(request.ParseForm())
				a.Equal(oauth.GrantTypeDeviceCode, request.Form.Get("grant_type"))
				a.Equal(oauth.StorydenCLIClientID, request.Form.Get("client_id"))
				a.Equal("device-code", request.Form.Get("device_code"))
				writeJSON(t, w, map[string]any{
					"access_token":  "access-token",
					"refresh_token": "refresh-token",
					"token_type":    "Bearer",
					"expires_in":    900,
					"scope":         "openid profile offline_access CREATE_POST",
				})

			default:
				http.NotFound(w, request)
			}
		}))
		defer server.Close()

		var stdout bytes.Buffer
		var stderr bytes.Buffer
		store := config.NewFileStoreAt(filepath.Join(t.TempDir(), "storyden", "config.yaml"))
		root := &cobra.Command{Use: "sd"}
		root.SetOut(&stdout)
		root.SetErr(&stderr)
		root.AddCommand((*cobra.Command)(New(store)))
		root.SetArgs([]string{"login", server.URL + "/api"})

		r.NoError(root.Execute())

		cfg, err := store.Load()
		r.NoError(err)

		a.True(deviceAuthorizationCalled)
		a.True(tokenCalled)
		a.Equal(slug(strings.TrimPrefix(server.URL, "http://")), cfg.CurrentContext)
		ctx := cfg.Contexts[cfg.CurrentContext]
		a.Equal(server.URL, ctx.APIURL)
		a.Equal(config.AuthStorageFile, ctx.AuthType)
		r.NotNil(ctx.Auth)
		a.Equal("access-token", ctx.Auth.AccessToken)
		a.Equal("refresh-token", ctx.Auth.RefreshToken)
		a.Equal("Bearer", ctx.Auth.TokenType)
		a.Equal("openid profile offline_access CREATE_POST", ctx.Auth.Scope)
		a.Equal("https://storyden.example", ctx.Auth.Issuer)
		a.Equal(oauth.StorydenCLIClientID, ctx.Auth.ClientID)
		a.Contains(stdout.String(), "ABCD-EFGH")
		a.Contains(stdout.String(), server.URL)
		a.Contains(stderr.String(), "credentials will be stored in the config file")
	})
}

func TestDeviceVerificationInstructions(t *testing.T) {
	t.Run("accepts verification URI complete", func(t *testing.T) {
		a := assert.New(t)

		device := &openapi.OAuthDeviceAuthorisation{
			VerificationUriComplete: ptr("http://example.com/complete"),
		}

		a.True(hasDeviceVerification(device))
	})

	t.Run("accepts verification URI with user code", func(t *testing.T) {
		a := assert.New(t)

		device := &openapi.OAuthDeviceAuthorisation{
			VerificationUri: ptr("http://example.com/device"),
			UserCode:        ptr("ABCD-EFGH"),
		}
		var stdout bytes.Buffer

		a.True(hasDeviceVerification(device))
		writeDeviceInstructions(&stdout, device)

		a.Contains(stdout.String(), tui.URL.Render("http://example.com/device"))
		a.Contains(stdout.String(), "ABCD-EFGH")
		a.Contains(stdout.String(), "Enter this code")
	})

	t.Run("rejects verification URI without user code", func(t *testing.T) {
		a := assert.New(t)

		device := &openapi.OAuthDeviceAuthorisation{
			VerificationUri: ptr("http://example.com/device"),
		}

		a.False(hasDeviceVerification(device))
	})
}

func TestResolveAuthStorage(t *testing.T) {
	t.Run("auto uses credential store when available", func(t *testing.T) {
		a := assert.New(t)
		r := require.New(t)

		store := config.NewStoreAtWithCredentialStore(filepath.Join(t.TempDir(), "storyden", "config.yaml"), newMemoryCredentialStore())

		storage, err := resolveAuthStorage("auto", store)
		r.NoError(err)

		a.Equal(config.AuthStorageCredentialStore, storage)
	})

	t.Run("file override forces file storage when credential store is available", func(t *testing.T) {
		a := assert.New(t)
		r := require.New(t)

		store := config.NewStoreAtWithCredentialStore(filepath.Join(t.TempDir(), "storyden", "config.yaml"), newMemoryCredentialStore())

		storage, err := resolveAuthStorage("file", store)
		r.NoError(err)

		a.Equal(config.AuthStorageFile, storage)
	})
}

func TestEndpointFromArgs(t *testing.T) {
	t.Run("reuses current context when no endpoint argument is provided", func(t *testing.T) {
		a := assert.New(t)
		r := require.New(t)

		store := config.NewFileStoreAt(filepath.Join(t.TempDir(), "storyden", "config.yaml"))
		cfg := config.New()
		cfg.UpsertContext("local", config.Context{APIURL: "http://localhost:8000"})
		cfg.SetCurrentContext("local")
		r.NoError(store.Save(cfg))

		endpoint, err := endpointFromArgs(&cobra.Command{}, nil, store)
		r.NoError(err)

		a.Equal("http://localhost:8000", endpoint)
	})
}

func writeJSON(t *testing.T, w http.ResponseWriter, value any) {
	t.Helper()

	w.Header().Set("Content-Type", "application/json")
	require.NoError(t, json.NewEncoder(w).Encode(value))
}

func serverURL(r *http.Request) string {
	return "http://" + r.Host
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
