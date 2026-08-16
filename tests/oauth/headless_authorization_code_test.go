package oauth_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/Southclaws/opt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/account/authentication/access_key"
	"github.com/Southclaws/storyden/app/resources/account/role/role_assign"
	"github.com/Southclaws/storyden/app/resources/account/role/role_repo"
	oauthresource "github.com/Southclaws/storyden/app/resources/oauth"
	"github.com/Southclaws/storyden/app/resources/oauth/oauth_writer"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/tests"
)

// TestOAuthHeadlessAuthorizationCodeFlow proves the Authorization Code + PKCE
// flow can be completed end-to-end by a script holding only a personal access
// key: no session cookie, no browser, and no reliance on the Storyden web
// frontend at any point.
func TestOAuthHeadlessAuthorizationCodeFlow(t *testing.T) {
	t.Parallel()

	integration.Test(t, oauthConfig(t), e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		cl *openapi.ClientWithResponses,
		aw *account_writer.Writer,
		assignments *role_assign.Assignment,
		ow *oauth_writer.Writer,
		roles *role_repo.Repository,
		akRepo *access_key.Repository,
		ts *httptest.Server,
	) {
		lc.Append(fx.StartHook(func() {
			_, owner := e2e.WithAccount(root, aw, seed.Account_001_Odin)
			_, member := e2e.WithAccount(root, aw, seed.Account_004_Loki)
			grantOAuthClientUse(t, root, roles, assignments, member.ID, rbac.PermissionCreatePost)

			akr, err := akRepo.Create(root, member.ID, access_key.AccessKeyKindPersonal, "headless-e2e", opt.NewEmpty[time.Time]())
			require.NoError(t, err)
			bearerAuth := bearer(akr.String())

			t.Run("personal_access_key_completes_full_authorization_code_flow", func(t *testing.T) {
				a := assert.New(t)
				r := require.New(t)

				clientID := "headless-auth-code-" + uuid.NewString()
				clientSecret := "secret-" + uuid.NewString()
				redirectURI := "https://client.example/callback"
				createClient(t, root, ow, owner.ID, clientID, oauthresource.ClientTypeConfidential, oauthresource.ScopePolicyExplicit, opt.New(clientSecretHash(t, clientSecret)), append(standardScopes(), rbac.PermissionCreatePost.String()), []string{oauthGrantAuthorizationCode, oauthGrantRefreshToken})
				verifier := strings.Repeat("k", 43)
				state := "state-" + uuid.NewString()

				location := authorizeRedirect(t, root, ts, bearerAuth, authorizeRequest{
					ClientID:            clientID,
					RedirectURI:         redirectURI,
					Scope:               "openid profile CREATE_POST",
					State:               state,
					CodeChallenge:       codeChallenge(verifier),
					CodeChallengeMethod: "S256",
				})

				consentURL, err := url.Parse(location)
				r.NoError(err)
				requestID := consentURL.Query().Get("request_id")
				r.NotEmpty(requestID)

				consent := tests.AssertRequest(cl.OAuthAuthoriseConsentWithResponse(root, &openapi.OAuthAuthoriseConsentParams{
					RequestId: (*openapi.OAuthAuthorizationRequestIDQuery)(&requestID),
				}, bearerAuth))(t, http.StatusOK)
				r.NotNil(consent.JSON200)
				a.Equal(clientID, consent.JSON200.ClientId)

				submit := tests.AssertRequest(cl.OAuthAuthoriseConsentSubmitWithResponse(root, openapi.OAuthAuthoriseConsentSubmitJSONRequestBody{
					RequestId: requestID,
					Decision:  openapi.OAuthAuthoriseDecisionApprove,
				}, bearerAuth))(t, http.StatusOK)
				r.NotNil(submit.JSON200)
				a.Equal(openapi.OAuthAuthoriseConsentResultStatusApproved, submit.JSON200.Status)

				redirect, err := url.Parse(submit.JSON200.Location)
				r.NoError(err)
				code := redirect.Query().Get("code")
				r.NotEmpty(code)

				token := tests.AssertRequest(oauthToken(t, root, cl, oauthTokenRequest{
					GrantType:    oauthGrantAuthorizationCode,
					ClientId:     clientID,
					ClientSecret: &clientSecret,
					Code:         &code,
					RedirectUri:  &redirectURI,
					CodeVerifier: &verifier,
				}))(t, http.StatusOK)
				r.NotNil(token.JSON200)
				r.NotNil(token.JSON200.AccessToken)

				userinfo := tests.AssertRequest(cl.OAuthUserInfoWithResponse(root, bearer(*token.JSON200.AccessToken)))(t, http.StatusOK)
				r.NotNil(userinfo.JSON200)
			})

			t.Run("bot_access_key_is_rejected_for_authorize_and_consent", func(t *testing.T) {
				a := assert.New(t)
				r := require.New(t)

				botAkr, err := akRepo.Create(root, member.ID, access_key.AccessKeyKindBot, "bot-e2e", opt.NewEmpty[time.Time]())
				r.NoError(err)
				botAuth := bearer(botAkr.String())

				clientID := "headless-bot-blocked-" + uuid.NewString()
				clientSecret := "secret-" + uuid.NewString()
				redirectURI := "https://client.example/callback"
				createClient(t, root, ow, owner.ID, clientID, oauthresource.ClientTypeConfidential, oauthresource.ScopePolicyExplicit, opt.New(clientSecretHash(t, clientSecret)), standardScopes(), []string{oauthGrantAuthorizationCode})

				resp := authorizeHTTPResponse(t, root, ts, botAuth, authorizeRequest{
					ClientID:            clientID,
					RedirectURI:         redirectURI,
					State:               "state-" + uuid.NewString(),
					CodeChallenge:       codeChallenge(strings.Repeat("o", 43)),
					CodeChallengeMethod: "S256",
				})
				defer resp.Body.Close()
				a.Equal(http.StatusForbidden, resp.StatusCode)

				requestID := openapi.OAuthAuthorizationRequestIDQuery("request-" + uuid.NewString())
				consent := tests.AssertRequest(cl.OAuthAuthoriseConsentWithResponse(root, &openapi.OAuthAuthoriseConsentParams{
					RequestId: &requestID,
				}, botAuth))(t, http.StatusForbidden)
				r.Nil(consent.JSON200)

				submit := tests.AssertRequest(cl.OAuthAuthoriseConsentSubmitWithResponse(root, openapi.OAuthAuthoriseConsentSubmitJSONRequestBody{
					RequestId: string(requestID),
					Decision:  openapi.OAuthAuthoriseDecisionApprove,
				}, botAuth))(t, http.StatusForbidden)
				r.Nil(submit.JSON200)
			})

			t.Run("no_credential_at_all_still_redirects_to_default_login_url", func(t *testing.T) {
				r := require.New(t)

				resp := authorizeHTTPResponse(t, root, ts, func(ctx context.Context, req *http.Request) error { return nil }, authorizeRequest{
					ClientID:            "anon-" + uuid.NewString(),
					RedirectURI:         "https://client.example/callback",
					State:               "state-" + uuid.NewString(),
					CodeChallenge:       codeChallenge(strings.Repeat("l", 43)),
					CodeChallengeMethod: "S256",
				})
				defer resp.Body.Close()

				r.Equal(http.StatusFound, resp.StatusCode)
				r.Equal("http://localhost:3000/login", resp.Header.Get("Location"))
			})

			t.Run("oauth_access_token_is_rejected_for_authorize_and_consent", func(t *testing.T) {
				a := assert.New(t)
				r := require.New(t)

				ccClientID := "headless-cc-" + uuid.NewString()
				ccClientSecret := "secret-" + uuid.NewString()
				createClient(t, root, ow, member.ID, ccClientID, oauthresource.ClientTypeConfidential, oauthresource.ScopePolicyExplicit, opt.New(clientSecretHash(t, ccClientSecret)), standardScopes(), []string{oauthGrantClientCredentials})

				ccToken := tests.AssertRequest(oauthToken(t, root, cl, oauthTokenRequest{
					GrantType:    oauthGrantClientCredentials,
					ClientId:     ccClientID,
					ClientSecret: &ccClientSecret,
				}))(t, http.StatusOK)
				r.NotNil(ccToken.JSON200)
				r.NotNil(ccToken.JSON200.AccessToken)
				oauthBearer := bearer(*ccToken.JSON200.AccessToken)

				authClientID := "headless-blocked-" + uuid.NewString()
				authClientSecret := "secret-" + uuid.NewString()
				createClient(t, root, ow, owner.ID, authClientID, oauthresource.ClientTypeConfidential, oauthresource.ScopePolicyExplicit, opt.New(clientSecretHash(t, authClientSecret)), standardScopes(), []string{oauthGrantAuthorizationCode})

				resp := authorizeHTTPResponse(t, root, ts, oauthBearer, authorizeRequest{
					ClientID:            authClientID,
					RedirectURI:         "https://client.example/callback",
					State:               "state-" + uuid.NewString(),
					CodeChallenge:       codeChallenge(strings.Repeat("m", 43)),
					CodeChallengeMethod: "S256",
				})
				defer resp.Body.Close()
				a.Equal(http.StatusForbidden, resp.StatusCode)

				requestID := openapi.OAuthAuthorizationRequestIDQuery("request-" + uuid.NewString())
				consent := tests.AssertRequest(cl.OAuthAuthoriseConsentWithResponse(root, &openapi.OAuthAuthoriseConsentParams{
					RequestId: &requestID,
				}, oauthBearer))(t, http.StatusForbidden)
				r.Nil(consent.JSON200)

				submit := tests.AssertRequest(cl.OAuthAuthoriseConsentSubmitWithResponse(root, openapi.OAuthAuthoriseConsentSubmitJSONRequestBody{
					RequestId: string(requestID),
					Decision:  openapi.OAuthAuthoriseDecisionApprove,
				}, oauthBearer))(t, http.StatusForbidden)
				r.Nil(submit.JSON200)
			})
		}))
	}))
}

// TestOAuthAuthorisationLoginURLConfiguration proves the unauthenticated
// redirect target used by the Authorization Code flow's entry point is
// configurable, resolving the TODO that previously hardcoded it to the
// frontend's `/login` route.
func TestOAuthAuthorisationLoginURLConfiguration(t *testing.T) {
	t.Parallel()

	cfg := oauthConfig(t)
	loginURL, err := url.Parse("https://custom.example/sign-in")
	require.NoError(t, err)
	cfg.OAuthAuthorisationLoginURL = *loginURL

	integration.Test(t, cfg, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		ts *httptest.Server,
	) {
		lc.Append(fx.StartHook(func() {
			r := require.New(t)

			resp := authorizeHTTPResponse(t, root, ts, func(ctx context.Context, req *http.Request) error { return nil }, authorizeRequest{
				ClientID:            "anon-" + uuid.NewString(),
				RedirectURI:         "https://client.example/callback",
				State:               "state-" + uuid.NewString(),
				CodeChallenge:       codeChallenge(strings.Repeat("n", 43)),
				CodeChallengeMethod: "S256",
			})
			defer resp.Body.Close()

			r.Equal(http.StatusFound, resp.StatusCode)
			r.Equal("https://custom.example/sign-in", resp.Header.Get("Location"))
		}))
	}))
}
