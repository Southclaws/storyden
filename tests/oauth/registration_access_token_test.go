package oauth_test

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/rs/xid"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/ent"
	entaccount "github.com/Southclaws/storyden/internal/ent/account"
	"github.com/Southclaws/storyden/internal/ent/oauthclient"
	"github.com/Southclaws/storyden/internal/ent/oauthdynamicregistrationaccesstokens"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/tests"
)

func agentRegistration(t *testing.T) (openapi.OAuthClientRegisterJSONRequestBody, *rsa.PrivateKey) {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)
	return openapi.OAuthClientRegisterJSONRequestBody{
		ClientName:              new("agent-" + xid.New().String()),
		GrantTypes:              &[]string{"client_credentials"},
		TokenEndpointAuthMethod: new("private_key_jwt"),
		Jwks: &map[string]interface{}{"keys": []any{map[string]any{
			"kty": "RSA", "kid": "agent-key", "n": base64.RawURLEncoding.EncodeToString(key.N.Bytes()), "e": "AQAB",
		}}},
	}, key
}

func TestOAuthRegistrationAccessTokens(t *testing.T) {
	t.Parallel()
	cfg := oauthConfig(t)

	integration.Test(t, cfg, e2e.Setup(), withOAuthRegistration(t, "protected"), fx.Invoke(func(lc fx.Lifecycle, root context.Context, cl *openapi.ClientWithResponses, sh *e2e.SessionHelper, aw *account_writer.Writer, db *ent.Client) {
		lc.Append(fx.StartHook(func() {
			adminCtx, admin := e2e.WithAccount(root, aw, seed.Account_001_Odin)
			adminSession := sh.WithSession(adminCtx)
			memberCtx, _ := e2e.WithAccount(root, aw, seed.Account_003_Baldur)
			memberSession := sh.WithSession(memberCtx)
			issue := func(t *testing.T, limit *int) *openapi.OAuthDynamicRegistrationAccessTokenIssued {
				t.Helper()
				response := tests.AssertRequest(cl.AdminOAuthDynamicRegistrationAccessTokenCreateWithResponse(root, openapi.AdminOAuthDynamicRegistrationAccessTokenCreateJSONRequestBody{
					Label: "Trusted team", ExpiresAt: time.Now().Add(time.Hour), MaxRegistrations: limit,
				}, adminSession))(t, http.StatusCreated)
				require.NotNil(t, response.JSON201)
				require.Equal(t, "no-store", response.HTTPResponse.Header.Get("Cache-Control"))
				require.Equal(t, admin.ID.String(), string(response.JSON201.Iat.CreatorAccountId))
				return response.JSON201
			}
			get := func(t *testing.T, id openapi.Identifier) *openapi.OAuthDynamicRegistrationAccessToken {
				t.Helper()
				response := tests.AssertRequest(cl.AdminOAuthDynamicRegistrationAccessTokenGetWithResponse(root, id, adminSession))(t, http.StatusOK)
				require.NotNil(t, response.JSON200)
				return response.JSON200
			}

			t.Run("admin_only_management", func(t *testing.T) {
				issued := issue(t, nil)
				for _, access := range []struct {
					editor openapi.RequestEditorFn
					status int
				}{
					{memberSession, http.StatusForbidden},
					{bearer(issued.Token), http.StatusUnauthorized},
				} {
					tests.AssertRequest(cl.AdminOAuthDynamicRegistrationAccessTokenCreateWithResponse(root, openapi.AdminOAuthDynamicRegistrationAccessTokenCreateJSONRequestBody{Label: "denied", ExpiresAt: time.Now().Add(time.Hour)}, access.editor))(t, access.status)
					tests.AssertRequest(cl.AdminOAuthDynamicRegistrationAccessTokenListWithResponse(root, nil, access.editor))(t, access.status)
					tests.AssertRequest(cl.AdminOAuthDynamicRegistrationAccessTokenGetWithResponse(root, issued.Iat.Id, access.editor))(t, access.status)
					tests.AssertRequest(cl.AdminOAuthDynamicRegistrationAccessTokenRevokeWithResponse(root, issued.Iat.Id, access.editor))(t, access.status)
				}
				tests.AssertRequest(cl.OAuthClientListWithResponse(root, bearer(issued.Token)))(t, http.StatusUnauthorized)
			})

			t.Run("creation_validation", func(t *testing.T) {
				for _, body := range []openapi.AdminOAuthDynamicRegistrationAccessTokenCreateJSONRequestBody{
					{Label: " ", ExpiresAt: time.Now().Add(time.Hour)},
					{Label: "expired", ExpiresAt: time.Now().Add(-time.Hour)},
					{Label: "zero", ExpiresAt: time.Now().Add(time.Hour), MaxRegistrations: new(0)},
					{Label: "negative", ExpiresAt: time.Now().Add(time.Hour), MaxRegistrations: new(-1)},
				} {
					tests.AssertRequest(cl.AdminOAuthDynamicRegistrationAccessTokenCreateWithResponse(root, body, adminSession))(t, http.StatusBadRequest)
				}
			})

			t.Run("invalid_profiles_do_not_consume_admission", func(t *testing.T) {
				issued := issue(t, nil)
				body, _ := agentRegistration(t)
				body.GrantTypes = &[]string{"client_credentials", "authorization_code"}
				body.RedirectUris = &[]string{"https://client.example/callback"}
				tests.AssertRequest(cl.OAuthClientRegisterWithResponse(root, body, bearer(issued.Token)))(t, http.StatusBadRequest)
				body.GrantTypes = &[]string{"client_credentials"}
				body.TokenEndpointAuthMethod = new("client_secret_basic")
				tests.AssertRequest(cl.OAuthClientRegisterWithResponse(root, body, bearer(issued.Token)))(t, http.StatusBadRequest)
				body.TokenEndpointAuthMethod = new("private_key_jwt")
				body.Jwks = &map[string]any{"keys": []any{map[string]any{"kty": "RSA", "kid": "invalid"}}}
				tests.AssertRequest(cl.OAuthClientRegisterWithResponse(root, body, bearer(issued.Token)))(t, http.StatusBadRequest)
				require.Zero(t, get(t, issued.Iat.Id).RegistrationCount)
				exists, err := db.Account.Query().Where(entaccount.Handle(*body.ClientName)).Exist(root)
				require.NoError(t, err)
				require.False(t, exists)
			})

			t.Run("missing_and_invalid_credentials", func(t *testing.T) {
				body, _ := agentRegistration(t)
				response := tests.AssertRequest(cl.OAuthClientRegisterWithResponse(root, body))(t, http.StatusUnauthorized)
				require.Equal(t, `Bearer realm="oauth/register"`, response.HTTPResponse.Header.Get("WWW-Authenticate"))
				for _, header := range []string{"Bearer ", "Basic abc", "Bearer sdiat_unknown", "Bearer one two"} {
					response := tests.AssertRequest(cl.OAuthClientRegisterWithResponse(root, body, func(_ context.Context, req *http.Request) error { req.Header.Set("Authorization", header); return nil }))(t, http.StatusUnauthorized)
					require.Equal(t, "invalid_token", response.JSON401.Error)
					require.Equal(t, "no-store", response.HTTPResponse.Header.Get("Cache-Control"))
				}
				issued := issue(t, nil)
				tests.AssertRequest(cl.OAuthClientRegisterWithResponse(root, body, func(_ context.Context, req *http.Request) error {
					req.Header.Add("Authorization", "Bearer "+issued.Token)
					req.Header.Add("Authorization", "Bearer "+issued.Token)
					return nil
				}))(t, http.StatusUnauthorized)
				require.Equal(t, 0, get(t, issued.Iat.Id).RegistrationCount)
			})

			t.Run("single_use_hash_and_provenance", func(t *testing.T) {
				issued := issue(t, nil)
				require.Equal(t, 1, issued.Iat.MaxRegistrations)
				require.True(t, strings.HasPrefix(issued.Token, "sdiat_"))
				secret, err := base64.RawURLEncoding.DecodeString(strings.TrimPrefix(issued.Token, "sdiat_"))
				require.NoError(t, err)
				require.Len(t, secret, 32)
				rawID, err := xid.FromString(string(issued.Iat.Id))
				require.NoError(t, err)
				stored, err := db.OAuthDynamicRegistrationAccessTokens.Get(root, rawID)
				require.NoError(t, err)
				sum := sha256.Sum256([]byte(issued.Token))
				require.Equal(t, hex.EncodeToString(sum[:]), stored.TokenHash)
				body, key := agentRegistration(t)
				body.Scope = new("ADMINISTRATOR")
				response := tests.AssertRequest(cl.OAuthClientRegisterWithResponse(root, body, bearer(issued.Token)))(t, http.StatusCreated)
				require.Nil(t, response.JSON201.ClientSecret)
				client, err := db.OAuthClient.Query().Where(oauthclient.ClientID(response.JSON201.ClientId)).Only(root)
				require.NoError(t, err)
				require.NotNil(t, client.AccountID)
				require.NotEqual(t, client.ID, *client.AccountID)
				bot, err := db.Account.Get(root, *client.AccountID)
				require.NoError(t, err)
				require.Equal(t, "bot", bot.Kind.String())
				require.Equal(t, *body.ClientName, bot.Handle)
				assertion := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
					"iss": response.JSON201.ClientId,
					"sub": response.JSON201.ClientId,
					"aud": "http://localhost:8000/api/oauth/token",
					"iat": time.Now().Unix(),
					"exp": time.Now().Add(time.Minute).Unix(),
					"jti": uuid.NewString(),
				})
				assertion.Header["kid"] = "agent-key"
				signed, err := assertion.SignedString(key)
				require.NoError(t, err)
				form := url.Values{
					"grant_type":            {"client_credentials"},
					"client_id":             {response.JSON201.ClientId},
					"client_assertion_type": {"urn:ietf:params:oauth:client-assertion-type:jwt-bearer"},
					"client_assertion":      {signed},
				}
				token := tests.AssertRequest(cl.OAuthTokenWithBodyWithResponse(root, "application/x-www-form-urlencoded", strings.NewReader(form.Encode())))(t, http.StatusOK)
				require.NotNil(t, token.JSON200.AccessToken)
				claims := parseClaims(t, *token.JSON200.AccessToken)
				require.Equal(t, bot.ID.String(), claims["sub"])
				require.NotContains(t, claims["scope"], "ADMINISTRATOR")
				require.Contains(t, claims["scope"], "READ_PUBLISHED_LIBRARY")
				require.False(t, bot.Admin)

				managed := tests.AssertRequest(cl.AdminOAuthClientGetWithResponse(root, openapi.Identifier(client.ID.String()), adminSession))(t, http.StatusOK)
				require.Equal(t, issued.Iat.Id, *managed.JSON200.DcrIatId)
				require.Equal(t, 1, get(t, issued.Iat.Id).RegistrationCount)
				body.ClientName = new("another-" + xid.New().String())
				tests.AssertRequest(cl.OAuthClientRegisterWithResponse(root, body, bearer(issued.Token)))(t, http.StatusUnauthorized)
				found, err := db.Account.Query().Where(entaccount.Handle(*body.ClientName)).Exist(root)
				require.NoError(t, err)
				require.False(t, found)
				listing := tests.AssertRequest(cl.AdminOAuthDynamicRegistrationAccessTokenListWithResponse(root, nil, adminSession))(t, http.StatusOK)
				require.NotContains(t, string(listing.Body), issued.Token)
				require.NotContains(t, string(listing.Body), stored.TokenHash)
				encoded, err := json.Marshal(get(t, issued.Iat.Id))
				require.NoError(t, err)
				require.NotContains(t, string(encoded), "token_hash")
				require.NotContains(t, string(encoded), issued.Token)
			})

			t.Run("invalid_metadata_and_collision_do_not_consume", func(t *testing.T) {
				issued := issue(t, nil)
				body, _ := agentRegistration(t)
				body.Jwks = nil
				tests.AssertRequest(cl.OAuthClientRegisterWithResponse(root, body, bearer(issued.Token)))(t, http.StatusBadRequest)
				body, _ = agentRegistration(t)
				body.ClientName = new(admin.Handle)
				response := tests.AssertRequest(cl.OAuthClientRegisterWithResponse(root, body, bearer(issued.Token)))(t, http.StatusBadRequest)
				require.Equal(t, "invalid_client_metadata", response.JSON400.Error)
				require.Equal(t, 0, get(t, issued.Iat.Id).RegistrationCount)
				body.ClientName = new("retry-" + xid.New().String())
				tests.AssertRequest(cl.OAuthClientRegisterWithResponse(root, body, bearer(issued.Token)))(t, http.StatusCreated)
			})

			t.Run("ordinary_dcr_is_separate", func(t *testing.T) {
				issued := issue(t, nil)
				body := openapi.OAuthClientRegisterJSONRequestBody{ClientName: new("ordinary"), TokenEndpointAuthMethod: new("none"), RedirectUris: &[]string{"https://client.example/callback"}}
				tests.AssertRequest(cl.OAuthClientRegisterWithResponse(root, body))(t, http.StatusCreated)
				response := tests.AssertRequest(cl.OAuthClientRegisterWithResponse(root, body, bearer(issued.Token)))(t, http.StatusForbidden)
				require.Equal(t, "insufficient_scope", response.JSON403.Error)
				require.Equal(t, 0, get(t, issued.Iat.Id).RegistrationCount)
			})

			t.Run("revocation_keeps_existing_clients", func(t *testing.T) {
				issued := issue(t, new(2))
				body, _ := agentRegistration(t)
				registered := tests.AssertRequest(cl.OAuthClientRegisterWithResponse(root, body, bearer(issued.Token)))(t, http.StatusCreated)
				tests.AssertRequest(cl.AdminOAuthDynamicRegistrationAccessTokenRevokeWithResponse(root, issued.Iat.Id, adminSession))(t, http.StatusNoContent)
				revoked := get(t, issued.Iat.Id).RevokedAt
				require.NotNil(t, revoked)
				tests.AssertRequest(cl.AdminOAuthDynamicRegistrationAccessTokenRevokeWithResponse(root, issued.Iat.Id, adminSession))(t, http.StatusNoContent)
				require.Equal(t, revoked, get(t, issued.Iat.Id).RevokedAt)
				body.ClientName = new("revoked-" + xid.New().String())
				tests.AssertRequest(cl.OAuthClientRegisterWithResponse(root, body, bearer(issued.Token)))(t, http.StatusUnauthorized)
				client, err := db.OAuthClient.Query().Where(oauthclient.ClientID(registered.JSON201.ClientId)).Only(root)
				require.NoError(t, err)
				exists, err := db.Account.Query().Where(entaccount.ID(*client.AccountID)).Exist(root)
				require.NoError(t, err)
				require.True(t, exists)
				require.Equal(t, 1, get(t, issued.Iat.Id).RegistrationCount)
				tests.AssertRequest(cl.AdminOAuthDynamicRegistrationAccessTokenRevokeWithResponse(root, openapi.Identifier(xid.New().String()), adminSession))(t, http.StatusNotFound)
			})

			t.Run("expired_token", func(t *testing.T) {
				token := "sdiat_" + uuid.NewString()
				sum := sha256.Sum256([]byte(token))
				_, err := db.OAuthDynamicRegistrationAccessTokens.Create().SetLabel("Expired").SetTokenHash(hex.EncodeToString(sum[:])).SetCreatorAccountID(xid.ID(admin.ID)).SetExpiresAt(time.Now().Add(-time.Second)).SetMaxRegistrations(1).Save(root)
				require.NoError(t, err)
				body, _ := agentRegistration(t)
				tests.AssertRequest(cl.OAuthClientRegisterWithResponse(root, body, bearer(token)))(t, http.StatusUnauthorized)
			})

			t.Run("concurrent_requests_respect_limit", func(t *testing.T) {
				issued := issue(t, new(2))
				body, _ := agentRegistration(t)
				start := make(chan struct{})
				var wg sync.WaitGroup
				statuses := make(chan int, 6)
				errors := make(chan error, 6)
				for n := range 6 {
					wg.Add(1)
					go func(n int) {
						defer wg.Done()
						request := body
						request.ClientName = new(fmt.Sprintf("race-%d-%s", n, xid.New().String()))
						<-start
						response, err := cl.OAuthClientRegisterWithResponse(root, request, bearer(issued.Token))
						if err != nil {
							errors <- err
							return
						}
						statuses <- response.StatusCode()
					}(n)
				}
				close(start)
				wg.Wait()
				close(statuses)
				close(errors)
				for err := range errors {
					require.NoError(t, err)
				}
				counts := map[int]int{}
				for status := range statuses {
					counts[status]++
				}
				require.Equal(t, map[int]int{http.StatusCreated: 2, http.StatusUnauthorized: 4}, counts)
				require.Equal(t, 2, get(t, issued.Iat.Id).RegistrationCount)
				id, err := xid.FromString(string(issued.Iat.Id))
				require.NoError(t, err)
				count, err := db.OAuthClient.Query().Where(oauthclient.DcrIatID(id)).Count(root)
				require.NoError(t, err)
				require.Equal(t, 2, count)
				tokenCount, err := db.OAuthDynamicRegistrationAccessTokens.Query().Where(oauthdynamicregistrationaccesstokens.ID(id)).Count(root)
				require.NoError(t, err)
				require.Equal(t, 1, tokenCount)
			})
		}))
	}))
}

func TestOAuthAgentRegistrationModes(t *testing.T) {
	for _, mode := range []string{"", "disabled", "open"} {
		t.Run(mode, func(t *testing.T) {
			t.Parallel()
			cfg := oauthConfig(t)

			integration.Test(t, cfg, e2e.Setup(), withOAuthRegistration(t, mode), fx.Invoke(func(lc fx.Lifecycle, root context.Context, cl *openapi.ClientWithResponses, sh *e2e.SessionHelper, aw *account_writer.Writer) {
				lc.Append(fx.StartHook(func() {
					adminCtx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
					issued := tests.AssertRequest(cl.AdminOAuthDynamicRegistrationAccessTokenCreateWithResponse(root, openapi.AdminOAuthDynamicRegistrationAccessTokenCreateJSONRequestBody{Label: "Mode test", ExpiresAt: time.Now().Add(time.Hour)}, sh.WithSession(adminCtx)))(t, http.StatusCreated)
					body, _ := agentRegistration(t)
					expected := http.StatusCreated
					if mode != "open" {
						expected = http.StatusBadRequest
					}
					tests.AssertRequest(cl.OAuthClientRegisterWithResponse(root, body))(t, expected)
					body.ClientName = new("mode-" + xid.New().String())
					tests.AssertRequest(cl.OAuthClientRegisterWithResponse(root, body, bearer(issued.JSON201.Token)))(t, expected)
					if mode == "open" {
						body.ClientName = new("invalid-" + xid.New().String())
						tests.AssertRequest(cl.OAuthClientRegisterWithResponse(root, body, bearer("sdiat_invalid")))(t, http.StatusUnauthorized)
					}
				}))
			}))
		})
	}
}
