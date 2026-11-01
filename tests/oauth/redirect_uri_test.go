package oauth_test

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/Southclaws/opt"
	"github.com/rs/xid"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	oauthresource "github.com/Southclaws/storyden/app/resources/oauth"
	"github.com/Southclaws/storyden/app/resources/oauth/oauth_writer"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/tests"
)

func TestOAuthClientRedirectPolicy(t *testing.T) {
	t.Parallel()

	integration.Test(t, oauthConfig(t), e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		cl *openapi.ClientWithResponses,
		sh *e2e.SessionHelper,
		aw *account_writer.Writer,
		ow *oauth_writer.Writer,
		ts *httptest.Server,
	) {
		lc.Append(fx.StartHook(func() {
			adminCtx, admin := e2e.WithAccount(root, aw, seed.Account_001_Odin)
			adminSession := sh.WithSession(adminCtx)

			for _, tc := range []struct {
				name string
				uri  string
				want int
			}{
				{"https", "https://client.example/callback?next=%2Fsettings&label=a%20b", http.StatusOK},
				{"loopback", "http://127.0.0.1:8080/callback", http.StatusOK},
				{"script", "javascript:alert(document.domain)", http.StatusBadRequest},
				{"custom_scheme", "myapp://callback", http.StatusOK},
				{"private_scheme", "com.example.app:/callback", http.StatusOK},
				{"remote_http", "http://client.example/callback", http.StatusBadRequest},
				{"credentials", "https://user:password@client.example/callback", http.StatusBadRequest},
			} {
				t.Run(tc.name, func(t *testing.T) {
					r := require.New(t)
					uris := []string{tc.uri}
					created := tests.AssertRequest(cl.AdminOAuthClientCreateWithResponse(adminCtx, openapi.AdminOAuthClientCreateJSONRequestBody{
						AccountId:     admin.ID.String(),
						ClientId:      "redirect-admin-" + xid.New().String(),
						Name:          "Redirect policy",
						Type:          openapi.OAuthClientTypePublic,
						RedirectUris:  uris,
						AllowedScopes: []string{"READ_PUBLISHED_THREADS"},
						AllowedGrants: []string{oauthGrantAuthorizationCode},
					}, adminSession))(t, tc.want)
					selfCreated := tests.AssertRequest(cl.OAuthClientCreateWithResponse(adminCtx, openapi.OAuthClientCreateJSONRequestBody{
						Name:          "Redirect policy",
						Type:          openapi.OAuthClientTypePublic,
						RedirectUris:  &uris,
						AllowedScopes: []string{"READ_PUBLISHED_THREADS"},
						AllowedGrants: []string{oauthGrantAuthorizationCode},
						PkceRequired:  true,
					}, adminSession))(t, tc.want)
					if tc.want == http.StatusOK {
						r.Equal(uris, created.JSON200.RedirectUris)
						r.Equal(uris, selfCreated.JSON200.Client.RedirectUris)

						verifier := strings.Repeat("r", 43)
						location := authorizeRedirect(t, root, ts, adminSession, authorizeRequest{
							ClientID:      selfCreated.JSON200.Client.ClientId,
							RedirectURI:   tc.uri,
							Scope:         "READ_PUBLISHED_THREADS",
							State:         "redirect-policy",
							CodeChallenge: codeChallenge(verifier),
						})
						consentURL, err := url.Parse(location)
						r.NoError(err)
						submit := tests.AssertRequest(cl.OAuthAuthoriseConsentSubmitWithResponse(root, openapi.OAuthAuthoriseConsentSubmitJSONRequestBody{
							RequestId: consentURL.Query().Get("request_id"),
							Decision:  openapi.OAuthAuthoriseDecisionApprove,
						}, adminSession))(t, http.StatusOK)
						callback, err := url.Parse(submit.JSON200.Location)
						r.NoError(err)
						registered, err := url.Parse(tc.uri)
						r.NoError(err)
						r.Equal(registered.Scheme, callback.Scheme)
						r.Equal(registered.Host, callback.Host)
						r.Equal(registered.Path, callback.Path)
						query := callback.Query()
						r.Equal("redirect-policy", query.Get("state"))
						code := query.Get("code")
						r.NotEmpty(code)
						query.Del("code")
						query.Del("state")
						r.Equal(registered.Query(), query)
						tests.AssertRequest(oauthToken(t, root, cl, oauthTokenRequest{
							GrantType:    oauthGrantAuthorizationCode,
							ClientId:     selfCreated.JSON200.Client.ClientId,
							Code:         &code,
							RedirectUri:  &tc.uri,
							CodeVerifier: &verifier,
						}))(t, http.StatusOK)
					}

					client := createClient(t, root, ow, admin.ID, "redirect-update-"+xid.New().String(), oauthresource.ClientTypePublic, oauthresource.ScopePolicyExplicit, opt.NewEmpty[string](), standardScopes(), []string{oauthGrantAuthorizationCode})
					tests.AssertRequest(cl.AdminOAuthClientUpdateWithResponse(adminCtx, xid.ID(client.ID).String(), openapi.AdminOAuthClientUpdateJSONRequestBody{
						RedirectUris: &uris,
					}, adminSession))(t, tc.want)
					tests.AssertRequest(cl.OAuthClientUpdateWithResponse(adminCtx, xid.ID(client.ID).String(), openapi.OAuthClientUpdateJSONRequestBody{
						RedirectUris: &uris,
					}, adminSession))(t, tc.want)

					stored := tests.AssertRequest(cl.AdminOAuthClientGetWithResponse(adminCtx, xid.ID(client.ID).String(), adminSession))(t, http.StatusOK)
					want := []string{"https://client.example/callback"}
					if tc.want == http.StatusOK {
						want = uris
					}
					r.Equal(want, stored.JSON200.RedirectUris)
				})
			}
		}))
	}))
}

func TestOAuthRejectsStoredUnsafeRedirects(t *testing.T) {
	t.Parallel()

	integration.Test(t, oauthConfig(t), e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		cl *openapi.ClientWithResponses,
		sh *e2e.SessionHelper,
		aw *account_writer.Writer,
		ow *oauth_writer.Writer,
		ts *httptest.Server,
	) {
		lc.Append(fx.StartHook(func() {
			adminCtx, admin := e2e.WithAccount(root, aw, seed.Account_001_Odin)
			adminSession := sh.WithSession(adminCtx)

			for _, uri := range []string{
				"javascript:alert(document.domain)",
				"data:text/html,unsafe",
				"vbscript:MsgBox(1)",
				"file:///etc/passwd",
				"blob:https://client.example/id",
				"about:blank",
				"filesystem:https://client.example/temporary/file",
				"http://client.example/callback",
				"https://user:password@client.example/callback",
			} {
				t.Run(uri, func(t *testing.T) {
					r := require.New(t)
					// Write through the resource layer to represent pre-existing rows.
					client := createClient(t, root, ow, admin.ID, "redirect-legacy-"+xid.New().String(), oauthresource.ClientTypePublic, oauthresource.ScopePolicyExplicit, opt.NewEmpty[string](), standardScopes(), []string{oauthGrantAuthorizationCode})
					_, err := ow.UpdateClient(root, client.ID, oauth_writer.ClientUpdate{RedirectURIs: opt.New([]string{uri})})
					r.NoError(err)

					for _, scope := range []string{"openid", "not-allowed"} {
						resp := authorizeHTTPResponse(t, root, ts, adminSession, authorizeRequest{
							ClientID:      client.ClientID,
							RedirectURI:   uri,
							Scope:         scope,
							CodeChallenge: codeChallenge(strings.Repeat("r", 43)),
						})
						resp.Body.Close()
						r.Equal(http.StatusBadRequest, resp.StatusCode)
						r.Empty(resp.Header.Get("Location"))
					}

					requestID := xid.New().String()
					hash := sha256.Sum256([]byte(requestID))
					_, err = ow.CreateAuthorisationRequest(root, oauth_writer.AuthorisationRequestCreate{
						ClientID:            client.ID,
						AccountID:           admin.ID,
						RequestIDHash:       hex.EncodeToString(hash[:]),
						RedirectURI:         uri,
						Scope:               "openid",
						CodeChallenge:       codeChallenge(strings.Repeat("r", 43)),
						CodeChallengeMethod: "S256",
						ExpiresAt:           time.Now().Add(time.Minute),
					})
					r.NoError(err)

					consent := tests.AssertRequest(cl.OAuthAuthoriseConsentWithResponse(root, &openapi.OAuthAuthoriseConsentParams{
						RequestId: (*openapi.OAuthAuthorizationRequestIDQuery)(&requestID),
					}, adminSession))(t, http.StatusBadRequest)
					r.Equal("invalid_request", consent.JSON400.Error)

					for _, decision := range []openapi.OAuthAuthoriseDecision{openapi.OAuthAuthoriseDecisionApprove, openapi.OAuthAuthoriseDecisionDeny} {
						submit := tests.AssertRequest(cl.OAuthAuthoriseConsentSubmitWithResponse(root, openapi.OAuthAuthoriseConsentSubmitJSONRequestBody{
							RequestId: requestID,
							Decision:  decision,
						}, adminSession))(t, http.StatusBadRequest)
						r.Equal("invalid_request", submit.JSON400.Error)
						r.Nil(submit.JSON200)
					}
				})
			}
		}))
	}))
}
