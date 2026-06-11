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
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/tests"
)

func TestOAuthAuthorizationCodeNonce(t *testing.T) {
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
			grantOAuthClientUse(t, root, roles, assignments, member.ID)
			memberSession := sh.WithSession(memberCtx)

			issueIDToken := func(t *testing.T, scope, nonce string) *openapi.OAuthToken {
				t.Helper()
				a := assert.New(t)
				r := require.New(t)

				clientID := "nonce-" + uuid.NewString()
				clientSecret := "secret-" + uuid.NewString()
				redirectURI := "https://client.example/callback"
				createClient(t, root, ow, owner.ID, clientID, oauthresource.ClientTypeConfidential, oauthresource.ScopePolicyExplicit, opt.New(clientSecretHash(t, clientSecret)), standardScopes(), []string{oauthGrantAuthorizationCode})
				verifier := strings.Repeat("a", 43)

				location := authorizeRedirect(t, root, ts, memberSession, authorizeRequest{
					ClientID:      clientID,
					RedirectURI:   redirectURI,
					Scope:         scope,
					State:         "state-" + uuid.NewString(),
					Nonce:         nonce,
					CodeChallenge: codeChallenge(verifier),
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

				token := tests.AssertRequest(oauthToken(t, root, cl, oauthTokenRequest{
					GrantType:    oauthGrantAuthorizationCode,
					ClientId:     clientID,
					ClientSecret: &clientSecret,
					Code:         &code,
					RedirectUri:  &redirectURI,
					CodeVerifier: &verifier,
				}))(t, http.StatusOK)
				r.NotNil(token.JSON200)
				a.NotNil(token.JSON200.IdToken)

				return token.JSON200
			}

			t.Run("nonce_is_echoed_into_id_token", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				nonce := "n-" + uuid.NewString()
				tok := issueIDToken(t, "openid profile", nonce)

				claims := parseClaims(t, *tok.IdToken)
				r.Contains(claims, "nonce")
				a.Equal(nonce, claims["nonce"])
			})

			t.Run("absent_nonce_omits_claim", func(t *testing.T) {
				a := assert.New(t)

				tok := issueIDToken(t, "openid profile", "")

				claims := parseClaims(t, *tok.IdToken)
				a.NotContains(claims, "nonce")
			})
		}))
	}))
}
