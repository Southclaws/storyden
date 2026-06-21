package oauthremote_test

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	oauth_remote "github.com/Southclaws/storyden/app/resources/oauth/remote"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/services/authentication/oauthremote"
	"github.com/Southclaws/storyden/app/services/authentication/oauthremotetoken"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
)

func TestOAuthRemoteClientDiscoveryAndAuthorization(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		service *oauthremote.Service,
		remoteRepo *oauth_remote.Repository,
		tokenService *oauthremotetoken.Service,
		aw *account_writer.Writer,
	) {
		lc.Append(fx.StartHook(func() {
			_, admin := e2e.WithAccount(root, aw, seed.Account_001_Odin)

			t.Run("cimd_discovery_prefers_client_metadata_document", func(t *testing.T) {
				as := newFakeRemoteOAuthAS(t, fakeOAuthASConfig{
					ClientIDMetadataDocumentSupported: true,
					RegistrationEndpoint:              true,
				})

				discovery, err := service.Discover(root, as.resourceURL())
				require.NoError(t, err)

				assert.Equal(t, oauthremote.ModeCIMD, discovery.Mode)
				assert.Equal(t, service.ClientMetadataDocumentURL(), discovery.ClientID)
				assert.Equal(t, as.server.URL+"/issuer", discovery.AuthorizationServer)
				assert.Equal(t, as.server.URL+"/authorize", discovery.AuthorizationServerMetadata.AuthorizationEndpoint)
			})

			t.Run("dcr_create_registers_dynamic_client", func(t *testing.T) {
				as := newFakeRemoteOAuthAS(t, fakeOAuthASConfig{
					RegistrationEndpoint: true,
				})

				connection, err := service.CreateConnection(root, oauthremote.CreateConnectionInput{
					ResourceURL: as.resourceURL(),
					Mode:        oauthremote.ModeDCR,
					Scope:       "read write",
					AddedBy:     admin.ID,
				})
				require.NoError(t, err)

				assert.Equal(t, oauth_remote.ModeDCR, connection.Mode)
				assert.Equal(t, "registered-client", connection.ClientID)
				assert.True(t, connection.HasClientSecret)
				assert.Equal(t, "registered-secret", connection.ClientSecret)
				assert.Equal(t, "client_secret_post", connection.TokenEndpointAuthMethod)
				require.Len(t, connection.RedirectURIs, 1)
				assert.Contains(t, connection.RedirectURIs[0], "/oauth/remote/callback")

				registration := as.lastRegistration()
				assert.Equal(t, 1, as.registrationRequestCount())
				assert.Equal(t, []any{"authorization_code", "refresh_token"}, registration["grant_types"])
				assert.Equal(t, []any{"code"}, registration["response_types"])
				assert.Equal(t, "read write", registration["scope"])
			})

			t.Run("create_replaces_unconnected_connection_when_mode_changes", func(t *testing.T) {
				as := newFakeRemoteOAuthAS(t, fakeOAuthASConfig{
					ClientIDMetadataDocumentSupported: true,
					RegistrationEndpoint:              true,
				})

				cimd, err := service.CreateConnection(root, oauthremote.CreateConnectionInput{
					ResourceURL: as.resourceURL(),
					Mode:        oauthremote.ModeCIMD,
					AddedBy:     admin.ID,
				})
				require.NoError(t, err)
				assert.Equal(t, oauth_remote.ModeCIMD, cimd.Mode)

				dcr, err := service.CreateConnection(root, oauthremote.CreateConnectionInput{
					ResourceURL: as.resourceURL(),
					Mode:        oauthremote.ModeDCR,
					AddedBy:     admin.ID,
				})
				require.NoError(t, err)

				assert.NotEqual(t, cimd.ID, dcr.ID)
				assert.Equal(t, oauth_remote.ModeDCR, dcr.Mode)
				assert.Equal(t, "registered-client", dcr.ClientID)
				assert.Equal(t, 1, as.registrationRequestCount())
			})

			t.Run("create_does_not_replace_connected_connection", func(t *testing.T) {
				as := newFakeRemoteOAuthAS(t, fakeOAuthASConfig{
					ClientIDMetadataDocumentSupported: true,
					RegistrationEndpoint:              true,
				})

				cimd, err := service.CreateConnection(root, oauthremote.CreateConnectionInput{
					ResourceURL: as.resourceURL(),
					Mode:        oauthremote.ModeCIMD,
					AddedBy:     admin.ID,
				})
				require.NoError(t, err)

				authorization, err := service.StartAuthorization(root, cimd.ID)
				require.NoError(t, err)
				authURL, err := url.Parse(authorization.AuthURL)
				require.NoError(t, err)
				as.expectCodeChallenge(authURL.Query().Get("code_challenge"))

				_, err = service.HandleCallback(root, authorization.State, "valid-code")
				require.NoError(t, err)

				_, err = service.CreateConnection(root, oauthremote.CreateConnectionInput{
					ResourceURL: as.resourceURL(),
					Mode:        oauthremote.ModeDCR,
					AddedBy:     admin.ID,
				})
				require.Error(t, err)
				assert.Equal(t, 0, as.registrationRequestCount())
			})

			t.Run("manual_create_skips_discovery", func(t *testing.T) {
				as := newFakeRemoteOAuthAS(t, fakeOAuthASConfig{})

				connection, err := service.CreateConnection(root, oauthremote.CreateConnectionInput{
					ResourceURL: as.resourceURL(),
					Mode:        oauthremote.ModeManual,
					Manual: oauthremote.ManualConfig{
						ClientID:              "manual-client",
						ClientSecret:          "manual-secret",
						AuthorizationEndpoint: as.server.URL + "/authorize",
						TokenEndpoint:         as.server.URL + "/token",
						AuthorizationServer:   as.server.URL + "/manual-issuer",
					},
					Scope:   "manual.scope",
					AddedBy: admin.ID,
				})
				require.NoError(t, err)

				assert.Equal(t, oauth_remote.ModeManual, connection.Mode)
				assert.Equal(t, "manual-client", connection.ClientID)
				assert.Equal(t, "manual-secret", connection.ClientSecret)
				assert.Equal(t, as.server.URL+"/manual-issuer", connection.AuthorizationServer)
				assert.Equal(t, 0, as.protectedResourceRequestCount())
			})

			t.Run("callback_exchanges_code_with_pkce_verifier", func(t *testing.T) {
				as := newFakeRemoteOAuthAS(t, fakeOAuthASConfig{})

				connection, err := service.CreateConnection(root, oauthremote.CreateConnectionInput{
					ResourceURL: as.resourceURL(),
					Mode:        oauthremote.ModeManual,
					Manual: oauthremote.ManualConfig{
						ClientID:              "pkce-client",
						AuthorizationEndpoint: as.server.URL + "/authorize",
						TokenEndpoint:         as.server.URL + "/token",
					},
					Scope:   "profile",
					AddedBy: admin.ID,
				})
				require.NoError(t, err)

				authorization, err := service.StartAuthorization(root, connection.ID)
				require.NoError(t, err)

				authURL, err := url.Parse(authorization.AuthURL)
				require.NoError(t, err)
				assert.Equal(t, as.server.URL+"/authorize", authURL.Scheme+"://"+authURL.Host+authURL.Path)
				assert.Equal(t, "pkce-client", authURL.Query().Get("client_id"))
				callbackURL, err := url.Parse(authURL.Query().Get("redirect_uri"))
				require.NoError(t, err)
				assert.Equal(t, "http://localhost/oauth/remote/callback", callbackURL.String())
				assert.Equal(t, "S256", authURL.Query().Get("code_challenge_method"))
				as.expectCodeChallenge(authURL.Query().Get("code_challenge"))

				result, err := service.HandleCallback(root, authorization.State, "valid-code")
				require.NoError(t, err)

				assert.Equal(t, oauth_remote.StatusConnected, result.Connection.Status)
				assert.Equal(t, "remote-access-token", result.Connection.AccessToken)
				assert.Equal(t, "remote-refresh-token", result.Connection.RefreshToken)
				assert.True(t, result.Connection.HasAccessToken)
				assert.True(t, result.Connection.HasRefreshToken)

				tokenRequest := as.lastTokenRequest()
				assert.Equal(t, "authorization_code", tokenRequest.Get("grant_type"))
				assert.Equal(t, "valid-code", tokenRequest.Get("code"))
				assert.NotEmpty(t, tokenRequest.Get("code_verifier"))
				assert.Equal(t, "pkce-client", tokenRequest.Get("client_id"))
				assert.Empty(t, as.lastTokenAuthorizationHeader())
			})

			t.Run("callback_state_is_single_use", func(t *testing.T) {
				as := newFakeRemoteOAuthAS(t, fakeOAuthASConfig{})

				connection, err := service.CreateConnection(root, oauthremote.CreateConnectionInput{
					ResourceURL: as.resourceURL(),
					Mode:        oauthremote.ModeManual,
					Manual: oauthremote.ManualConfig{
						ClientID:              "single-use-client",
						AuthorizationEndpoint: as.server.URL + "/authorize",
						TokenEndpoint:         as.server.URL + "/token",
					},
					AddedBy: admin.ID,
				})
				require.NoError(t, err)

				authorization, err := service.StartAuthorization(root, connection.ID)
				require.NoError(t, err)
				authURL, err := url.Parse(authorization.AuthURL)
				require.NoError(t, err)
				as.expectCodeChallenge(authURL.Query().Get("code_challenge"))

				_, err = service.HandleCallback(root, authorization.State, "valid-code")
				require.NoError(t, err)

				_, err = service.HandleCallback(root, authorization.State, "valid-code")
				require.Error(t, err)
				assert.Equal(t, 1, as.tokenRequestCount())
			})

			t.Run("dcr_callback_uses_registered_token_auth_method", func(t *testing.T) {
				as := newFakeRemoteOAuthAS(t, fakeOAuthASConfig{
					RegistrationEndpoint: true,
				})

				connection, err := service.CreateConnection(root, oauthremote.CreateConnectionInput{
					ResourceURL: as.resourceURL(),
					Mode:        oauthremote.ModeDCR,
					AddedBy:     admin.ID,
				})
				require.NoError(t, err)

				authorization, err := service.StartAuthorization(root, connection.ID)
				require.NoError(t, err)
				authURL, err := url.Parse(authorization.AuthURL)
				require.NoError(t, err)
				as.expectCodeChallenge(authURL.Query().Get("code_challenge"))

				_, err = service.HandleCallback(root, authorization.State, "valid-code")
				require.NoError(t, err)

				tokenRequest := as.lastTokenRequest()
				assert.Equal(t, "registered-client", tokenRequest.Get("client_id"))
				assert.Equal(t, "registered-secret", tokenRequest.Get("client_secret"))
				assert.Empty(t, as.lastTokenAuthorizationHeader())
			})

			t.Run("manual_secret_callback_uses_basic_auth", func(t *testing.T) {
				as := newFakeRemoteOAuthAS(t, fakeOAuthASConfig{})

				connection, err := service.CreateConnection(root, oauthremote.CreateConnectionInput{
					ResourceURL: as.resourceURL(),
					Mode:        oauthremote.ModeManual,
					Manual: oauthremote.ManualConfig{
						ClientID:              "basic-client",
						ClientSecret:          "basic-secret",
						AuthorizationEndpoint: as.server.URL + "/authorize",
						TokenEndpoint:         as.server.URL + "/token",
					},
					AddedBy: admin.ID,
				})
				require.NoError(t, err)

				authorization, err := service.StartAuthorization(root, connection.ID)
				require.NoError(t, err)
				authURL, err := url.Parse(authorization.AuthURL)
				require.NoError(t, err)
				as.expectCodeChallenge(authURL.Query().Get("code_challenge"))

				_, err = service.HandleCallback(root, authorization.State, "valid-code")
				require.NoError(t, err)

				username, password, ok := as.lastTokenBasicAuth()
				require.True(t, ok)
				assert.Equal(t, "basic-client", username)
				assert.Equal(t, "basic-secret", password)
				assert.Empty(t, as.lastTokenRequest().Get("client_secret"))
			})

			t.Run("access_token_refreshes_expired_token_and_replaces_refresh_token", func(t *testing.T) {
				as := newFakeRemoteOAuthAS(t, fakeOAuthASConfig{
					RefreshTokenReplacement: "rotated-refresh-token",
				})
				connection := createStoredRefreshConnection(t, root, service, remoteRepo, admin.ID, as, refreshConnectionConfig{
					AccessToken:  "expired-access-token",
					RefreshToken: "old-refresh-token",
					Expiry:       time.Now().Add(-time.Hour),
				})

				token, err := tokenService.AccessToken(root, connection.ID)
				require.NoError(t, err)

				assert.Equal(t, "refreshed-access-token", token)
				assert.Equal(t, 1, as.refreshRequestCount())
				refreshRequest := as.lastTokenRequest()
				assert.Equal(t, "refresh_token", refreshRequest.Get("grant_type"))
				assert.Equal(t, "old-refresh-token", refreshRequest.Get("refresh_token"))

				stored, err := remoteRepo.GetConnection(root, connection.ID)
				require.NoError(t, err)
				assert.Equal(t, "refreshed-access-token", stored.AccessToken)
				assert.Equal(t, "rotated-refresh-token", stored.RefreshToken)
				assert.Equal(t, oauth_remote.StatusConnected, stored.Status)
				assert.Nil(t, stored.LastError)
				assert.Nil(t, stored.TokenRefreshStartedAt)
			})

			t.Run("access_token_refreshes_near_expiry_token", func(t *testing.T) {
				as := newFakeRemoteOAuthAS(t, fakeOAuthASConfig{})
				connection := createStoredRefreshConnection(t, root, service, remoteRepo, admin.ID, as, refreshConnectionConfig{
					AccessToken:  "near-expiry-access-token",
					RefreshToken: "near-expiry-refresh-token",
					Expiry:       time.Now().Add(time.Minute),
				})

				token, err := tokenService.AccessToken(root, connection.ID)
				require.NoError(t, err)

				assert.Equal(t, "refreshed-access-token", token)
				assert.Equal(t, 1, as.refreshRequestCount())
			})

			t.Run("access_token_does_not_refresh_unexpired_token", func(t *testing.T) {
				as := newFakeRemoteOAuthAS(t, fakeOAuthASConfig{})
				connection := createStoredRefreshConnection(t, root, service, remoteRepo, admin.ID, as, refreshConnectionConfig{
					AccessToken:  "fresh-access-token",
					RefreshToken: "fresh-refresh-token",
					Expiry:       time.Now().Add(time.Hour),
				})

				token, err := tokenService.AccessToken(root, connection.ID)
				require.NoError(t, err)

				assert.Equal(t, "fresh-access-token", token)
				assert.Equal(t, 0, as.refreshRequestCount())
			})

			t.Run("access_token_preserves_refresh_token_when_response_omits_replacement", func(t *testing.T) {
				as := newFakeRemoteOAuthAS(t, fakeOAuthASConfig{
					OmitRefreshToken: true,
				})
				connection := createStoredRefreshConnection(t, root, service, remoteRepo, admin.ID, as, refreshConnectionConfig{
					AccessToken:  "expired-access-token",
					RefreshToken: "preserved-refresh-token",
					Expiry:       time.Now().Add(-time.Hour),
				})

				_, err := tokenService.AccessToken(root, connection.ID)
				require.NoError(t, err)

				stored, err := remoteRepo.GetConnection(root, connection.ID)
				require.NoError(t, err)
				assert.Equal(t, "preserved-refresh-token", stored.RefreshToken)
			})

			t.Run("access_token_refresh_uses_configured_client_auth_method", func(t *testing.T) {
				t.Run("public", func(t *testing.T) {
					as := newFakeRemoteOAuthAS(t, fakeOAuthASConfig{})
					connection := createStoredRefreshConnection(t, root, service, remoteRepo, admin.ID, as, refreshConnectionConfig{
						ClientID:     "public-client",
						AccessToken:  "expired-access-token",
						RefreshToken: "public-refresh-token",
						Expiry:       time.Now().Add(-time.Hour),
					})

					_, err := tokenService.AccessToken(root, connection.ID)
					require.NoError(t, err)

					assert.Equal(t, "public-client", as.lastTokenRequest().Get("client_id"))
					assert.Empty(t, as.lastTokenRequest().Get("client_secret"))
					assert.Empty(t, as.lastTokenAuthorizationHeader())
				})

				t.Run("client_secret_basic", func(t *testing.T) {
					as := newFakeRemoteOAuthAS(t, fakeOAuthASConfig{})
					connection := createStoredRefreshConnection(t, root, service, remoteRepo, admin.ID, as, refreshConnectionConfig{
						ClientID:     "basic-refresh-client",
						ClientSecret: "basic-refresh-secret",
						AccessToken:  "expired-access-token",
						RefreshToken: "basic-refresh-token",
						Expiry:       time.Now().Add(-time.Hour),
					})

					_, err := tokenService.AccessToken(root, connection.ID)
					require.NoError(t, err)

					username, password, ok := as.lastTokenBasicAuth()
					require.True(t, ok)
					assert.Equal(t, "basic-refresh-client", username)
					assert.Equal(t, "basic-refresh-secret", password)
					assert.Empty(t, as.lastTokenRequest().Get("client_secret"))
				})

				t.Run("client_secret_post", func(t *testing.T) {
					as := newFakeRemoteOAuthAS(t, fakeOAuthASConfig{})
					connection, err := remoteRepo.CreateConnection(root, oauth_remote.ConnectionCreate{
						ResourceURL:             as.resourceURL() + "/post",
						AuthorizationServer:     as.server.URL + "/manual-issuer-post",
						Mode:                    oauth_remote.ModeManual,
						Status:                  oauth_remote.StatusPending,
						ClientID:                "post-refresh-client",
						ClientSecret:            "post-refresh-secret",
						AuthorizationEndpoint:   as.server.URL + "/authorize",
						TokenEndpoint:           as.server.URL + "/token",
						TokenEndpointAuthMethod: "client_secret_post",
						RedirectURI:             service.DefaultRedirectURI(),
						RedirectURIs:            []string{service.DefaultRedirectURI()},
						AddedBy:                 admin.ID,
					})
					require.NoError(t, err)
					expired := time.Now().Add(-time.Hour)
					connection, err = remoteRepo.StoreTokens(root, connection.ID, oauth_remote.TokenUpdate{
						AccessToken:  "expired-access-token",
						RefreshToken: "post-refresh-token",
						TokenType:    "Bearer",
						TokenExpiry:  &expired,
					})
					require.NoError(t, err)

					_, err = tokenService.AccessToken(root, connection.ID)
					require.NoError(t, err)

					assert.Equal(t, "post-refresh-client", as.lastTokenRequest().Get("client_id"))
					assert.Equal(t, "post-refresh-secret", as.lastTokenRequest().Get("client_secret"))
					assert.Empty(t, as.lastTokenAuthorizationHeader())
				})
			})

			t.Run("access_token_refresh_failure_marks_error_without_overwriting_tokens", func(t *testing.T) {
				as := newFakeRemoteOAuthAS(t, fakeOAuthASConfig{
					FailRefresh: true,
				})
				connection := createStoredRefreshConnection(t, root, service, remoteRepo, admin.ID, as, refreshConnectionConfig{
					AccessToken:  "expired-access-token",
					RefreshToken: "old-refresh-token",
					Expiry:       time.Now().Add(-time.Hour),
				})

				_, err := tokenService.AccessToken(root, connection.ID)
				require.Error(t, err)

				stored, err := remoteRepo.GetConnection(root, connection.ID)
				require.NoError(t, err)
				assert.Equal(t, "expired-access-token", stored.AccessToken)
				assert.Equal(t, "old-refresh-token", stored.RefreshToken)
				assert.Equal(t, oauth_remote.StatusError, stored.Status)
				require.NotNil(t, stored.LastError)
				assert.NotEmpty(t, strings.TrimSpace(*stored.LastError))
				assert.Nil(t, stored.TokenRefreshStartedAt)
			})

			t.Run("concurrent_access_token_refresh_only_calls_token_endpoint_once", func(t *testing.T) {
				as := newFakeRemoteOAuthAS(t, fakeOAuthASConfig{
					RefreshDelay: 150 * time.Millisecond,
				})
				connection := createStoredRefreshConnection(t, root, service, remoteRepo, admin.ID, as, refreshConnectionConfig{
					AccessToken:  "expired-access-token",
					RefreshToken: "concurrent-refresh-token",
					Expiry:       time.Now().Add(-time.Hour),
				})

				const callers = 8
				var wg sync.WaitGroup
				tokens := make(chan string, callers)
				errs := make(chan error, callers)
				for range callers {
					wg.Add(1)
					go func() {
						defer wg.Done()
						token, err := tokenService.AccessToken(root, connection.ID)
						if err != nil {
							errs <- err
							return
						}
						tokens <- token
					}()
				}
				wg.Wait()
				close(tokens)
				close(errs)

				for err := range errs {
					require.NoError(t, err)
				}
				for token := range tokens {
					assert.Equal(t, "refreshed-access-token", token)
				}
				assert.Equal(t, 1, as.refreshRequestCount())
			})
		}))
	}))
}

type refreshConnectionConfig struct {
	ClientID     string
	ClientSecret string
	AccessToken  string
	RefreshToken string
	Expiry       time.Time
}

func createStoredRefreshConnection(
	t *testing.T,
	ctx context.Context,
	service *oauthremote.Service,
	repo *oauth_remote.Repository,
	addedBy account.AccountID,
	as *fakeRemoteOAuthAS,
	cfg refreshConnectionConfig,
) oauth_remote.Connection {
	t.Helper()

	clientID := cfg.ClientID
	if clientID == "" {
		clientID = "refresh-client"
	}
	connection, err := service.CreateConnection(ctx, oauthremote.CreateConnectionInput{
		ResourceURL: as.resourceURL() + "/" + xid.New().String(),
		Mode:        oauthremote.ModeManual,
		Manual: oauthremote.ManualConfig{
			ClientID:              clientID,
			ClientSecret:          cfg.ClientSecret,
			AuthorizationEndpoint: as.server.URL + "/authorize",
			TokenEndpoint:         as.server.URL + "/token",
			AuthorizationServer:   as.server.URL + "/manual-issuer-" + xid.New().String(),
		},
		AddedBy: addedBy,
	})
	require.NoError(t, err)

	connection, err = repo.StoreTokens(ctx, connection.ID, oauth_remote.TokenUpdate{
		AccessToken:  cfg.AccessToken,
		RefreshToken: cfg.RefreshToken,
		TokenType:    "Bearer",
		TokenExpiry:  &cfg.Expiry,
	})
	require.NoError(t, err)

	return connection
}

type fakeOAuthASConfig struct {
	ClientIDMetadataDocumentSupported bool
	RegistrationEndpoint              bool
	RefreshTokenReplacement           string
	OmitRefreshToken                  bool
	FailRefresh                       bool
	RefreshDelay                      time.Duration
}

type fakeRemoteOAuthAS struct {
	t      *testing.T
	config fakeOAuthASConfig
	server *httptest.Server

	mu                             sync.Mutex
	protectedResourceMetadataCalls int
	registration                   map[string]any
	registrationRequests           int
	expectedCodeChallenge          string
	tokenRequest                   url.Values
	tokenAuthorizationHeader       string
	tokenRequestTotal              int
	refreshRequestTotal            int
}

func newFakeRemoteOAuthAS(t *testing.T, config fakeOAuthASConfig) *fakeRemoteOAuthAS {
	t.Helper()

	as := &fakeRemoteOAuthAS{t: t, config: config}
	as.server = httptest.NewServer(http.HandlerFunc(as.handle))
	t.Cleanup(as.server.Close)
	return as
}

func (as *fakeRemoteOAuthAS) resourceURL() string {
	return as.server.URL + "/mcp"
}

func (as *fakeRemoteOAuthAS) expectCodeChallenge(challenge string) {
	as.mu.Lock()
	defer as.mu.Unlock()

	as.expectedCodeChallenge = challenge
}

func (as *fakeRemoteOAuthAS) protectedResourceRequestCount() int {
	as.mu.Lock()
	defer as.mu.Unlock()

	return as.protectedResourceMetadataCalls
}

func (as *fakeRemoteOAuthAS) lastRegistration() map[string]any {
	as.mu.Lock()
	defer as.mu.Unlock()

	out := make(map[string]any, len(as.registration))
	for k, v := range as.registration {
		out[k] = v
	}
	return out
}

func (as *fakeRemoteOAuthAS) registrationRequestCount() int {
	as.mu.Lock()
	defer as.mu.Unlock()

	return as.registrationRequests
}

func (as *fakeRemoteOAuthAS) lastTokenRequest() url.Values {
	as.mu.Lock()
	defer as.mu.Unlock()

	out := make(url.Values, len(as.tokenRequest))
	for k, v := range as.tokenRequest {
		out[k] = append([]string(nil), v...)
	}
	return out
}

func (as *fakeRemoteOAuthAS) tokenRequestCount() int {
	as.mu.Lock()
	defer as.mu.Unlock()

	return as.tokenRequestTotal
}

func (as *fakeRemoteOAuthAS) refreshRequestCount() int {
	as.mu.Lock()
	defer as.mu.Unlock()

	return as.refreshRequestTotal
}

func (as *fakeRemoteOAuthAS) lastTokenAuthorizationHeader() string {
	as.mu.Lock()
	defer as.mu.Unlock()

	return as.tokenAuthorizationHeader
}

func (as *fakeRemoteOAuthAS) lastTokenBasicAuth() (string, string, bool) {
	header := as.lastTokenAuthorizationHeader()
	if header == "" {
		return "", "", false
	}

	req := &http.Request{Header: http.Header{"Authorization": []string{header}}}
	return req.BasicAuth()
}

func (as *fakeRemoteOAuthAS) handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.URL.Path {
	case "/.well-known/oauth-protected-resource/mcp":
		as.mu.Lock()
		as.protectedResourceMetadataCalls++
		as.mu.Unlock()

		writeJSON(as.t, w, http.StatusOK, oauthremote.ProtectedResourceMetadata{
			Resource:               as.resourceURL(),
			ResourceName:           "Fake MCP",
			AuthorizationServers:   []string{as.server.URL + "/issuer"},
			BearerMethodsSupported: []string{"header"},
		})

	case "/.well-known/oauth-authorization-server/issuer":
		metadata := oauthremote.AuthorizationServerMetadata{
			Issuer:                            as.server.URL + "/issuer",
			AuthorizationEndpoint:             as.server.URL + "/authorize",
			TokenEndpoint:                     as.server.URL + "/token",
			ResponseTypesSupported:            []string{"code"},
			GrantTypesSupported:               []string{"authorization_code", "refresh_token"},
			TokenEndpointAuthMethodsSupported: []string{"none", "client_secret_post"},
			CodeChallengeMethodsSupported:     []string{"S256"},
			ClientIDMetadataDocumentSupported: as.config.ClientIDMetadataDocumentSupported,
		}
		if as.config.RegistrationEndpoint {
			metadata.RegistrationEndpoint = as.server.URL + "/register"
		}
		writeJSON(as.t, w, http.StatusOK, metadata)

	case "/register":
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var registration map[string]any
		if err := json.NewDecoder(r.Body).Decode(&registration); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		as.mu.Lock()
		as.registration = registration
		as.registrationRequests++
		as.mu.Unlock()

		writeJSON(as.t, w, http.StatusCreated, map[string]any{
			"client_id":                  "registered-client",
			"client_secret":              "registered-secret",
			"token_endpoint_auth_method": "client_secret_post",
			"redirect_uris":              registration["redirect_uris"],
			"grant_types":                []string{"authorization_code", "refresh_token"},
			"response_types":             []string{"code"},
		})

	case "/token":
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		as.mu.Lock()
		as.tokenRequest = r.PostForm
		as.tokenAuthorizationHeader = r.Header.Get("Authorization")
		as.tokenRequestTotal++
		as.mu.Unlock()

		switch r.PostForm.Get("grant_type") {
		case "authorization_code":
			verifier := r.PostForm.Get("code_verifier")
			as.mu.Lock()
			expectedChallenge := as.expectedCodeChallenge
			as.mu.Unlock()

			if r.PostForm.Get("code") != "valid-code" ||
				verifier == "" ||
				pkceChallengeForTest(verifier) != expectedChallenge {
				http.Error(w, "invalid token request", http.StatusBadRequest)
				return
			}

			writeJSON(as.t, w, http.StatusOK, map[string]any{
				"access_token":  "remote-access-token",
				"refresh_token": "remote-refresh-token",
				"token_type":    "Bearer",
				"expires_in":    3600,
				"scope":         strings.TrimSpace(r.PostForm.Get("scope")),
			})

		case "refresh_token":
			as.mu.Lock()
			as.refreshRequestTotal++
			as.mu.Unlock()

			if as.config.RefreshDelay > 0 {
				time.Sleep(as.config.RefreshDelay)
			}
			if as.config.FailRefresh {
				http.Error(w, "refresh failed", http.StatusBadRequest)
				return
			}
			if r.PostForm.Get("refresh_token") == "" {
				http.Error(w, "missing refresh token", http.StatusBadRequest)
				return
			}

			body := map[string]any{
				"access_token": "refreshed-access-token",
				"token_type":   "Bearer",
				"expires_in":   3600,
				"scope":        strings.TrimSpace(r.PostForm.Get("scope")),
			}
			if !as.config.OmitRefreshToken {
				refreshToken := as.config.RefreshTokenReplacement
				if refreshToken == "" {
					refreshToken = r.PostForm.Get("refresh_token")
				}
				body["refresh_token"] = refreshToken
			}

			writeJSON(as.t, w, http.StatusOK, body)

		default:
			http.Error(w, "unsupported grant type", http.StatusBadRequest)
		}

	default:
		http.NotFound(w, r)
	}
}

func pkceChallengeForTest(verifier string) string {
	sum := sha256.Sum256([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(sum[:])
}

func writeJSON(t *testing.T, w http.ResponseWriter, status int, value any) {
	t.Helper()

	w.WriteHeader(status)
	require.NoError(t, json.NewEncoder(w).Encode(value))
}
