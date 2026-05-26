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

func TestOAuthAuthorizationCodeFlow(t *testing.T) {
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
			memberCtx, member := e2e.WithAccount(root, aw, seed.Account_004_Loki)
			grantOAuthClientUse(t, root, roles, assignments, member.ID, rbac.PermissionCreatePost)
			memberSession := sh.WithSession(memberCtx)

			t.Run("valid_authorize_request_requires_consent_then_issues_code", func(t *testing.T) {
				a := assert.New(t)
				r := require.New(t)

				clientID := "auth-code-consent-" + uuid.NewString()
				clientSecret := "secret-" + uuid.NewString()
				redirectURI := "https://client.example/callback"
				createClient(t, root, ow, owner.ID, clientID, oauthresource.ClientTypeConfidential, oauthresource.ScopePolicyExplicit, opt.New(clientSecretHash(t, clientSecret)), append(standardScopes(), rbac.PermissionCreatePost.String()), []string{oauthGrantAuthorizationCode, oauthGrantRefreshToken})
				verifier := strings.Repeat("a", 43)
				state := "state-" + uuid.NewString()

				location := authorizeRedirect(t, root, ts, memberSession, authorizeRequest{
					ClientID:            clientID,
					RedirectURI:         redirectURI,
					Scope:               "openid profile CREATE_POST",
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
				}, memberSession))(t, http.StatusOK)
				r.NotNil(consent.JSON200)
				a.Equal(clientID, consent.JSON200.ClientId)
				a.Equal(redirectURI, consent.JSON200.RedirectUri)
				a.Contains(consent.JSON200.GrantedScopes, rbac.PermissionCreatePost.String())

				submit := tests.AssertRequest(cl.OAuthAuthoriseConsentSubmitWithResponse(root, openapi.OAuthAuthoriseConsentSubmitJSONRequestBody{
					RequestId: requestID,
					Decision:  openapi.OAuthAuthoriseDecisionApprove,
				}, memberSession))(t, http.StatusOK)
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
			})

			t.Run("denied_authorize_request_redirects_with_access_denied", func(t *testing.T) {
				a := assert.New(t)
				r := require.New(t)

				clientID := "auth-code-denied-" + uuid.NewString()
				clientSecret := "secret-" + uuid.NewString()
				redirectURI := "https://client.example/callback"
				createClient(t, root, ow, owner.ID, clientID, oauthresource.ClientTypeConfidential, oauthresource.ScopePolicyExplicit, opt.New(clientSecretHash(t, clientSecret)), standardScopes(), []string{oauthGrantAuthorizationCode})
				state := "state-" + uuid.NewString()

				location := authorizeRedirect(t, root, ts, memberSession, authorizeRequest{
					ClientID:            clientID,
					RedirectURI:         redirectURI,
					Scope:               "openid profile",
					State:               state,
					CodeChallenge:       codeChallenge(strings.Repeat("z", 43)),
					CodeChallengeMethod: "S256",
				})
				consentURL, err := url.Parse(location)
				r.NoError(err)
				requestID := consentURL.Query().Get("request_id")
				r.NotEmpty(requestID)

				submit := tests.AssertRequest(cl.OAuthAuthoriseConsentSubmitWithResponse(root, openapi.OAuthAuthoriseConsentSubmitJSONRequestBody{
					RequestId: requestID,
					Decision:  openapi.OAuthAuthoriseDecisionDeny,
				}, memberSession))(t, http.StatusOK)
				r.NotNil(submit.JSON200)
				a.Equal(openapi.OAuthAuthoriseConsentResultStatusDenied, submit.JSON200.Status)

				redirect, err := url.Parse(submit.JSON200.Location)
				r.NoError(err)
				a.Equal(redirectURI, redirect.Scheme+"://"+redirect.Host+redirect.Path)
				a.Equal("access_denied", redirect.Query().Get("error"))
				a.Equal(state, redirect.Query().Get("state"))
			})

			t.Run("authorization_request_allows_omitted_scope", func(t *testing.T) {
				a := assert.New(t)
				r := require.New(t)

				clientID := "empty-scope-auth-" + uuid.NewString()
				clientSecret := "secret-" + uuid.NewString()
				redirectURI := "https://client.example/callback"
				createClient(t, root, ow, owner.ID, clientID, oauthresource.ClientTypeConfidential, oauthresource.ScopePolicyExplicit, opt.New(clientSecretHash(t, clientSecret)), standardScopes(), []string{oauthGrantAuthorizationCode})

				location := authorizeRedirect(t, root, ts, memberSession, authorizeRequest{
					ClientID:            clientID,
					RedirectURI:         redirectURI,
					State:               "state-" + uuid.NewString(),
					CodeChallenge:       codeChallenge(strings.Repeat("y", 43)),
					CodeChallengeMethod: "S256",
				})

				consentURL, err := url.Parse(location)
				r.NoError(err)
				a.Equal("http", consentURL.Scheme)
				a.Equal("localhost:3000", consentURL.Host)
				a.Equal("/oauth/authorize/consent", consentURL.Path)
				a.NotEmpty(consentURL.Query().Get("request_id"))
			})

			t.Run("authorize_rejects_clients_without_authorization_code_grant", func(t *testing.T) {
				a := assert.New(t)
				r := require.New(t)

				clientID := "auth-code-grant-blocked-" + uuid.NewString()
				clientSecret := "secret-" + uuid.NewString()
				redirectURI := "https://client.example/callback"
				createClient(t, root, ow, owner.ID, clientID, oauthresource.ClientTypeConfidential, oauthresource.ScopePolicyExplicit, opt.New(clientSecretHash(t, clientSecret)), append(standardScopes(), rbac.PermissionCreatePost.String()), []string{oauthGrantRefreshToken})

				resp := authorizeHTTPResponse(t, root, ts, memberSession, authorizeRequest{
					ClientID:            clientID,
					RedirectURI:         redirectURI,
					Scope:               "openid profile CREATE_POST",
					State:               "state-" + uuid.NewString(),
					CodeChallenge:       codeChallenge(strings.Repeat("b", 43)),
					CodeChallengeMethod: "S256",
				})
				defer resp.Body.Close()

				a.Equal(http.StatusBadRequest, resp.StatusCode)
				a.Empty(resp.Header.Get("Location"))
				r.NotNil(resp.Body)
			})

			t.Run("token_exchange_rejects_clients_without_authorization_code_grant", func(t *testing.T) {
				a := assert.New(t)
				r := require.New(t)

				clientID := "auth-code-token-grant-blocked-" + uuid.NewString()
				clientSecret := "secret-" + uuid.NewString()
				createClient(t, root, ow, owner.ID, clientID, oauthresource.ClientTypeConfidential, oauthresource.ScopePolicyExplicit, opt.New(clientSecretHash(t, clientSecret)), standardScopes(), []string{oauthGrantRefreshToken})

				code := "fake-code-" + uuid.NewString()
				redirectURI := "https://client.example/callback"
				verifier := strings.Repeat("c", 43)
				resp := tests.AssertRequest(oauthToken(t, root, cl, oauthTokenRequest{
					GrantType:    oauthGrantAuthorizationCode,
					ClientId:     clientID,
					ClientSecret: &clientSecret,
					Code:         &code,
					RedirectUri:  &redirectURI,
					CodeVerifier: &verifier,
				}))(t, http.StatusBadRequest)
				r.NotNil(resp.JSON400)
				a.Equal("unauthorized_client", resp.JSON400.Error)
			})

			t.Run("token_exchange_rejects_wrong_pkce_verifier", func(t *testing.T) {
				a := assert.New(t)
				r := require.New(t)

				clientID := "auth-code-pkce-mismatch-" + uuid.NewString()
				clientSecret := "secret-" + uuid.NewString()
				redirectURI := "https://client.example/callback"
				createClient(t, root, ow, owner.ID, clientID, oauthresource.ClientTypeConfidential, oauthresource.ScopePolicyExplicit, opt.New(clientSecretHash(t, clientSecret)), standardScopes(), []string{oauthGrantAuthorizationCode})

				verifier := strings.Repeat("e", 43)
				location := authorizeRedirect(t, root, ts, memberSession, authorizeRequest{
					ClientID:            clientID,
					RedirectURI:         redirectURI,
					Scope:               "openid",
					State:               "state-" + uuid.NewString(),
					CodeChallenge:       codeChallenge(verifier),
					CodeChallengeMethod: "S256",
				})
				consentURL, err := url.Parse(location)
				r.NoError(err)
				requestID := consentURL.Query().Get("request_id")
				r.NotEmpty(requestID)

				submit := tests.AssertRequest(cl.OAuthAuthoriseConsentSubmitWithResponse(root, openapi.OAuthAuthoriseConsentSubmitJSONRequestBody{
					RequestId: requestID,
					Decision:  openapi.OAuthAuthoriseDecisionApprove,
				}, memberSession))(t, http.StatusOK)
				r.NotNil(submit.JSON200)
				redirect, err := url.Parse(submit.JSON200.Location)
				r.NoError(err)
				code := redirect.Query().Get("code")
				r.NotEmpty(code)

				wrongVerifier := strings.Repeat("f", 43)
				resp := tests.AssertRequest(oauthToken(t, root, cl, oauthTokenRequest{
					GrantType:    oauthGrantAuthorizationCode,
					ClientId:     clientID,
					ClientSecret: &clientSecret,
					Code:         &code,
					RedirectUri:  &redirectURI,
					CodeVerifier: &wrongVerifier,
				}))(t, http.StatusBadRequest)
				r.NotNil(resp.JSON400)
				a.Equal("invalid_grant", resp.JSON400.Error)
			})

			t.Run("unsupported_authorize_request_is_rejected", func(t *testing.T) {
				a := assert.New(t)
				r := require.New(t)

				clientID := "auth-code-unsupported-" + uuid.NewString()
				clientSecret := "secret-" + uuid.NewString()
				redirectURI := "https://client.example/callback"
				createClient(t, root, ow, owner.ID, clientID, oauthresource.ClientTypeConfidential, oauthresource.ScopePolicyExplicit, opt.New(clientSecretHash(t, clientSecret)), append(standardScopes(), rbac.PermissionCreatePost.String()), []string{oauthGrantAuthorizationCode})

				resp := authorizeHTTPResponse(t, root, ts, memberSession, authorizeRequest{
					ResponseType:        "token",
					ClientID:            clientID,
					RedirectURI:         redirectURI,
					Scope:               "openid profile CREATE_POST",
					State:               "state-" + uuid.NewString(),
					CodeChallenge:       codeChallenge(strings.Repeat("d", 43)),
					CodeChallengeMethod: "S256",
				})
				defer resp.Body.Close()

				a.Equal(http.StatusBadRequest, resp.StatusCode)
				r.NotNil(resp.Body)
			})
		}))
	}))
}

func TestOAuthAuthorizationCodeRequiresOAuthClientPermission(t *testing.T) {
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

			t.Run("member_without_permission_cannot_start_authorization_code_consent", func(t *testing.T) {
				a := assert.New(t)
				r := require.New(t)

				memberCtx, _ := e2e.WithAccount(root, aw, seed.Account_004_Loki)
				memberSession := sh.WithSession(memberCtx)

				clientID := "auth-code-permission-required-" + uuid.NewString()
				clientSecret := "secret-" + uuid.NewString()
				redirectURI := "https://client.example/callback"
				createClient(t, root, ow, owner.ID, clientID, oauthresource.ClientTypeConfidential, oauthresource.ScopePolicyExplicit, opt.New(clientSecretHash(t, clientSecret)), append(standardScopes(), rbac.PermissionCreatePost.String()), []string{oauthGrantAuthorizationCode, oauthGrantRefreshToken})

				resp := authorizeHTTPResponse(t, root, ts, memberSession, authorizeRequest{
					ClientID:            clientID,
					RedirectURI:         redirectURI,
					Scope:               "openid profile CREATE_POST",
					State:               "state-" + uuid.NewString(),
					CodeChallenge:       codeChallenge(strings.Repeat("g", 43)),
					CodeChallengeMethod: "S256",
				})
				defer resp.Body.Close()

				a.Equal(http.StatusBadRequest, resp.StatusCode)
				a.Empty(resp.Header.Get("Location"))
				r.NotNil(resp.Body)
			})

			t.Run("permission_revocation_blocks_pending_authorization_code_approval", func(t *testing.T) {
				a := assert.New(t)
				r := require.New(t)

				memberCtx, member := e2e.WithAccount(root, aw, seed.Account_003_Baldur)
				oauthRoleID := grantOAuthClientUse(t, root, roles, assignments, member.ID, rbac.PermissionCreatePost)
				memberSession := sh.WithSession(memberCtx)

				clientID := "auth-code-permission-revoked-" + uuid.NewString()
				clientSecret := "secret-" + uuid.NewString()
				redirectURI := "https://client.example/callback"
				createClient(t, root, ow, owner.ID, clientID, oauthresource.ClientTypeConfidential, oauthresource.ScopePolicyExplicit, opt.New(clientSecretHash(t, clientSecret)), append(standardScopes(), rbac.PermissionCreatePost.String()), []string{oauthGrantAuthorizationCode, oauthGrantRefreshToken})

				location := authorizeRedirect(t, root, ts, memberSession, authorizeRequest{
					ClientID:            clientID,
					RedirectURI:         redirectURI,
					Scope:               "openid profile CREATE_POST",
					State:               "state-" + uuid.NewString(),
					CodeChallenge:       codeChallenge(strings.Repeat("h", 43)),
					CodeChallengeMethod: "S256",
				})
				consentURL, err := url.Parse(location)
				r.NoError(err)
				requestID := consentURL.Query().Get("request_id")
				r.NotEmpty(requestID)

				revokeOAuthClientUse(t, root, assignments, member.ID, oauthRoleID)

				consent := tests.AssertRequest(cl.OAuthAuthoriseConsentWithResponse(root, &openapi.OAuthAuthoriseConsentParams{
					RequestId: (*openapi.OAuthAuthorizationRequestIDQuery)(&requestID),
				}, memberSession))(t, http.StatusBadRequest)
				r.NotNil(consent.JSON400)
				a.Equal("access_denied", consent.JSON400.Error)

				approve := tests.AssertRequest(cl.OAuthAuthoriseConsentSubmitWithResponse(root, openapi.OAuthAuthoriseConsentSubmitJSONRequestBody{
					RequestId: requestID,
					Decision:  openapi.OAuthAuthoriseDecisionApprove,
				}, memberSession))(t, http.StatusBadRequest)
				r.NotNil(approve.JSON400)
				a.Equal("access_denied", approve.JSON400.Error)

				deny := tests.AssertRequest(cl.OAuthAuthoriseConsentSubmitWithResponse(root, openapi.OAuthAuthoriseConsentSubmitJSONRequestBody{
					RequestId: requestID,
					Decision:  openapi.OAuthAuthoriseDecisionDeny,
				}, memberSession))(t, http.StatusOK)
				r.NotNil(deny.JSON200)
				a.Equal(openapi.OAuthAuthoriseConsentResultStatusDenied, deny.JSON200.Status)
			})
		}))
	}))
}
