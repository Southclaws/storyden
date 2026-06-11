package oauth_test

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/account/role/role_assign"
	"github.com/Southclaws/storyden/app/resources/account/role/role_repo"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/tests"
)

func TestOAuthTokenEndpointResponseSemantics(t *testing.T) {
	t.Parallel()

	integration.Test(t, oauthConfig(t), e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		cl *openapi.ClientWithResponses,
		sh *e2e.SessionHelper,
		aw *account_writer.Writer,
		assignments *role_assign.Assignment,
		roles *role_repo.Repository,
		ts *httptest.Server,
	) {
		lc.Append(fx.StartHook(func() {
			memberCtx, member := e2e.WithAccount(root, aw, seed.Account_001_Odin)
			grantOAuthClientUse(t, root, roles, assignments, member.ID)
			memberSession := sh.WithSession(memberCtx)

			registered := tests.AssertRequest(cl.OAuthClientRegisterWithResponse(root, openapi.OAuthClientRegisterJSONRequestBody{
				ClientName:              ptr("Token Semantics Client"),
				RedirectUris:            &[]string{"https://semantics.example/cb"},
				TokenEndpointAuthMethod: ptr("client_secret_basic"),
			}, memberSession))(t, http.StatusCreated)
			require.NotNil(t, registered.JSON201)
			clientID := registered.JSON201.ClientId

			postToken := func(form url.Values, basicID, basicSecret string) *http.Response {
				httpReq, err := http.NewRequestWithContext(root, http.MethodPost, ts.URL+"/api/oauth/token", strings.NewReader(form.Encode()))
				require.NoError(t, err)
				httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
				if basicID != "" {
					basic := base64.StdEncoding.EncodeToString([]byte(basicID + ":" + basicSecret))
					httpReq.Header.Set("Authorization", "Basic "+basic)
				}
				resp, err := http.DefaultClient.Do(httpReq)
				require.NoError(t, err)
				return resp
			}

			t.Run("invalid_client_via_basic_auth_returns_401_with_challenge", func(t *testing.T) {
				a := assert.New(t)
				r := require.New(t)

				form := url.Values{}
				form.Set("grant_type", oauthGrantAuthorizationCode)
				form.Set("code", "does-not-exist")
				form.Set("redirect_uri", "https://semantics.example/cb")
				form.Set("code_verifier", strings.Repeat("x", 43))

				resp := postToken(form, clientID, "the-wrong-secret")
				defer resp.Body.Close()

				r.Equal(http.StatusUnauthorized, resp.StatusCode)
				a.Equal(`Basic realm="localhost"`, resp.Header.Get("WWW-Authenticate"), "challenge scheme must match the client's Basic auth, realm is the API host")
				a.Equal("no-store", resp.Header.Get("Cache-Control"))
				a.Equal("no-cache", resp.Header.Get("Pragma"))

				var errResp openapi.OAuthError
				r.NoError(json.NewDecoder(resp.Body).Decode(&errResp))
				a.Equal("invalid_client", errResp.Error)
			})

			t.Run("missing_grant_type_is_invalid_request", func(t *testing.T) {
				a := assert.New(t)
				r := require.New(t)

				form := url.Values{}
				form.Set("client_id", clientID)

				resp := postToken(form, "", "")
				defer resp.Body.Close()

				r.Equal(http.StatusBadRequest, resp.StatusCode)

				var errResp openapi.OAuthError
				r.NoError(json.NewDecoder(resp.Body).Decode(&errResp))
				a.Equal("invalid_request", errResp.Error)
			})
		}))
	}))
}

func TestOAuthRegistrationResponseIsNotCacheable(t *testing.T) {
	t.Parallel()

	integration.Test(t, oauthConfig(t), e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		ts *httptest.Server,
	) {
		lc.Append(fx.StartHook(func() {
			a := assert.New(t)
			r := require.New(t)

			body, err := json.Marshal(openapi.OAuthClientRegisterJSONRequestBody{
				ClientName:              ptr("Cacheability Client"),
				RedirectUris:            &[]string{"https://cache.example/cb"},
				TokenEndpointAuthMethod: ptr("client_secret_basic"),
			})
			r.NoError(err)

			req, err := http.NewRequestWithContext(root, http.MethodPost, ts.URL+"/api/oauth/register", strings.NewReader(string(body)))
			r.NoError(err)
			req.Header.Set("Content-Type", "application/json")

			resp, err := http.DefaultClient.Do(req)
			r.NoError(err)
			defer resp.Body.Close()

			r.Equal(http.StatusCreated, resp.StatusCode)
			a.Equal("no-store", resp.Header.Get("Cache-Control"))
			a.Equal("no-cache", resp.Header.Get("Pragma"))
		}))
	}))
}
