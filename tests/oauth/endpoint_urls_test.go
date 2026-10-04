package oauth_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/tests"
)

func TestOAuthEndpointURLs(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		name, address, tokenEndpoint string
	}{
		{"root", "https://api.example.com", "https://api.example.com/api/oauth/token"},
		{"api_suffix", "https://api.example.com/api/", "https://api.example.com/api/oauth/token"},
		{"prefix", "https://api.example.com/community/", "https://api.example.com/community/api/oauth/token"},
		{"prefix_api_suffix", "https://api.example.com/community/api", "https://api.example.com/community/api/oauth/token"},
		{"escaped_prefix", "https://api.example.com/my%20community/api/", "https://api.example.com/my%20community/api/oauth/token"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cfg := oauthConfig(t)
			address, err := url.Parse(tt.address)
			require.NoError(t, err)
			cfg.PublicAPIAddress = *address

			integration.Test(t, cfg, e2e.Setup(), withOAuthRegistration(t, "open"), fx.Invoke(func(lc fx.Lifecycle, ctx context.Context, cl *openapi.ClientWithResponses, ts *httptest.Server) {
				lc.Append(fx.StartHook(func() {
					response, err := http.Get(ts.URL + "/.well-known/oauth-authorization-server")
					require.NoError(t, err)
					defer response.Body.Close()
					require.Equal(t, http.StatusOK, response.StatusCode)

					var discovery struct {
						TokenEndpoint string `json:"token_endpoint"`
					}
					require.NoError(t, json.NewDecoder(response.Body).Decode(&discovery))
					require.Equal(t, tt.tokenEndpoint, discovery.TokenEndpoint)

					body, key := agentRegistration(t)
					registered := tests.AssertRequest(cl.OAuthClientRegisterWithResponse(ctx, body))(t, http.StatusCreated).JSON201

					for _, audience := range []string{tt.tokenEndpoint, "https://other.example.com/api/oauth/token"} {
						assertion := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.RegisteredClaims{
							Issuer:    registered.ClientId,
							Subject:   registered.ClientId,
							Audience:  jwt.ClaimStrings{audience},
							IssuedAt:  jwt.NewNumericDate(time.Now()),
							ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute)),
							ID:        uuid.NewString(),
						})
						assertion.Header["kid"] = "agent-key"
						signed, err := assertion.SignedString(key)
						require.NoError(t, err)

						form := url.Values{
							"grant_type":            {"client_credentials"},
							"client_id":             {registered.ClientId},
							"client_assertion_type": {"urn:ietf:params:oauth:client-assertion-type:jwt-bearer"},
							"client_assertion":      {signed},
						}
						status := http.StatusOK
						if audience != tt.tokenEndpoint {
							status = http.StatusBadRequest
						}

						token := tests.AssertRequest(cl.OAuthTokenWithBodyWithResponse(ctx, "application/x-www-form-urlencoded", strings.NewReader(form.Encode())))(t, status)
						if status == http.StatusBadRequest {
							require.Equal(t, "invalid_client", token.JSON400.Error)
						}
					}
				}))
			}))
		})
	}
}
