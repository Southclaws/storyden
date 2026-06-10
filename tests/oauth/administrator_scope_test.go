package oauth_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/Southclaws/opt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/account_writer"
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

func TestOAuthAdministratorScopeGrant(t *testing.T) {
	t.Parallel()

	integration.Test(t, oauthConfig(t), e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		cl *openapi.ClientWithResponses,
		sh *e2e.SessionHelper,
		aw *account_writer.Writer,
		assignments *role_assign.Assignment,
		ow *oauth_writer.Writer,
		roles *role_repo.Repository,
		ts *httptest.Server,
	) {
		lc.Append(fx.StartHook(func() {
			_, owner := e2e.WithAccount(root, aw, seed.Account_001_Odin)

			t.Run("administrator_account_can_authorize_with_specific_scopes", func(t *testing.T) {
				a := assert.New(t)
				r := require.New(t)

				adminCtx, admin := e2e.WithAccount(root, aw, seed.Account_006_Freyja)
				grantOAuthClientUse(t, root, roles, assignments, admin.ID, rbac.PermissionAdministrator)
				adminSession := sh.WithSession(adminCtx)

				clientID := "admin-scope-test-" + uuid.NewString()
				clientSecret := "secret-" + uuid.NewString()
				redirectURI := "https://client.example/callback"
				requestedScopes := append(standardScopes(),
					rbac.PermissionCreatePost.String(),
					rbac.PermissionReadPublishedThreads.String(),
				)
				createClient(t, root, ow, owner.ID, clientID, oauthresource.ClientTypeConfidential, oauthresource.ScopePolicyExplicit, opt.New(clientSecretHash(t, clientSecret)), requestedScopes, []string{oauthGrantAuthorizationCode, oauthGrantRefreshToken})
				verifier := strings.Repeat("a", 43)
				state := "state-" + uuid.NewString()

				location := authorizeRedirect(t, root, ts, adminSession, authorizeRequest{
					ClientID:            clientID,
					RedirectURI:         redirectURI,
					Scope:               "openid profile email offline_access CREATE_POST READ_PUBLISHED_THREADS",
					State:               state,
					CodeChallenge:       codeChallenge(verifier),
					CodeChallengeMethod: "S256",
				})

				consentURL, err := url.Parse(location)
				r.NoError(err)
				a.Equal("http", consentURL.Scheme)
				a.Equal("localhost:3000", consentURL.Host)
				a.Equal("/oauth/authorize/consent", consentURL.Path)
				requestID := consentURL.Query().Get("request_id")
				r.NotEmpty(requestID)

				consent := tests.AssertRequest(cl.OAuthAuthoriseConsentWithResponse(root, &openapi.OAuthAuthoriseConsentParams{
					RequestId: (*openapi.OAuthAuthorizationRequestIDQuery)(&requestID),
				}, adminSession))(t, http.StatusOK)
				r.NotNil(consent.JSON200)
				a.Equal(clientID, consent.JSON200.ClientId)
				a.Equal(redirectURI, consent.JSON200.RedirectUri)
				a.Contains(consent.JSON200.GrantedScopes, rbac.PermissionCreatePost.String())
				a.Contains(consent.JSON200.GrantedScopes, rbac.PermissionReadPublishedThreads.String())

				submit := tests.AssertRequest(cl.OAuthAuthoriseConsentSubmitWithResponse(root, openapi.OAuthAuthoriseConsentSubmitJSONRequestBody{
					RequestId: requestID,
					Decision:  openapi.OAuthAuthoriseDecisionApprove,
				}, adminSession))(t, http.StatusOK)
				r.NotNil(submit.JSON200)
				a.Equal(openapi.OAuthAuthoriseConsentResultStatusApproved, submit.JSON200.Status)

				redirect, err := url.Parse(submit.JSON200.Location)
				r.NoError(err)
				a.Equal(redirectURI, redirect.Scheme+"://"+redirect.Host+redirect.Path)
				a.Equal(state, redirect.Query().Get("state"))
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
				a.Contains(*token.JSON200.Scope, rbac.PermissionCreatePost.String())
				a.Contains(*token.JSON200.Scope, rbac.PermissionReadPublishedThreads.String())
			})

			t.Run("administrator_can_authorize_client_with_administrator_in_allowed_scopes", func(t *testing.T) {
				a := assert.New(t)
				r := require.New(t)

				adminCtx, admin := e2e.WithAccount(root, aw, seed.Account_007_Freyr)
				grantOAuthClientUse(t, root, roles, assignments, admin.ID, rbac.PermissionAdministrator)
				adminSession := sh.WithSession(adminCtx)

				clientID := "admin-scope-administrator-allowed-" + uuid.NewString()
				clientSecret := "secret-" + uuid.NewString()
				redirectURI := "https://client.example/callback"
				clientAllowedScopes := append(standardScopes(), rbac.PermissionAdministrator.String())
				createClient(t, root, ow, owner.ID, clientID, oauthresource.ClientTypeConfidential, oauthresource.ScopePolicyExplicit, opt.New(clientSecretHash(t, clientSecret)), clientAllowedScopes, []string{oauthGrantAuthorizationCode, oauthGrantRefreshToken})
				verifier := strings.Repeat("b", 43)
				state := "state-" + uuid.NewString()

				location := authorizeRedirect(t, root, ts, adminSession, authorizeRequest{
					ClientID:            clientID,
					RedirectURI:         redirectURI,
					Scope:               "openid profile email offline_access CREATE_POST READ_PUBLISHED_THREADS",
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
				}, adminSession))(t, http.StatusOK)
				r.NotNil(consent.JSON200)
				a.Contains(consent.JSON200.GrantedScopes, rbac.PermissionCreatePost.String())
				a.Contains(consent.JSON200.GrantedScopes, rbac.PermissionReadPublishedThreads.String())
			})
		}))
	}))
}
