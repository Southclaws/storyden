package oauth_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/ent/account"
	"github.com/Southclaws/storyden/internal/ent/oauthclient"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/tests"
)

func TestOAuthRegistrationProfiles(t *testing.T) {
	for _, mode := range []string{"disabled", "protected", "open", "approval"} {
		t.Run(mode, func(t *testing.T) {
			t.Parallel()
			cfg := oauthConfig(t)

			integration.Test(t, cfg, e2e.Setup(), withOAuthRegistration(t, mode), fx.Invoke(func(lc fx.Lifecycle, root context.Context, cl *openapi.ClientWithResponses, db *ent.Client) {
				lc.Append(fx.StartHook(func() {
					body, _ := agentRegistration(t)
					body.ClientName = new("Delegated MCP Connector")
					body.RedirectUris = &[]string{"https://client.example/callback"}
					body.RegistrationMode = &[]string{"approval"}
					for _, method := range []string{"none", "client_secret_basic", "client_secret_post", "private_key_jwt"} {
						t.Run(method, func(t *testing.T) {
							for _, grants := range []*[]string{nil, new([]string{"authorization_code", "refresh_token"})} {
								body.TokenEndpointAuthMethod = &method
								body.GrantTypes = grants
								registered := tests.AssertRequest(cl.OAuthClientRegisterWithResponse(root, body))(t, http.StatusCreated)
								require.NotNil(t, registered.JSON201)
								require.Nil(t, registered.JSON202)
								client, err := db.OAuthClient.Query().Where(oauthclient.ClientID(registered.JSON201.ClientId)).Only(root)
								require.NoError(t, err)
								require.Nil(t, client.AccountID)
								require.Nil(t, client.DcrIatID)
								require.Nil(t, client.RegistrationApprovedByAccountID)
								require.Equal(t, method, client.TokenEndpointAuthMethod)
								require.Equal(t, []string{"authorization_code", "refresh_token"}, client.AllowedGrants)
								if method == "none" || method == "private_key_jwt" {
									require.Nil(t, registered.JSON201.ClientSecret)
									require.Nil(t, client.ClientSecretHash)
								} else {
									require.NotNil(t, registered.JSON201.ClientSecret)
								}
							}
						})
					}
					t.Run("mixed_grants_cannot_provision", func(t *testing.T) {
						body, _ := agentRegistration(t)
						body.GrantTypes = &[]string{"authorization_code", "client_credentials"}
						body.RedirectUris = &[]string{"https://client.example/callback"}
						body.RegistrationMode = &[]string{"approval"}
						result := tests.AssertRequest(cl.OAuthClientRegisterWithResponse(root, body))(t, http.StatusBadRequest)
						require.Equal(t, "invalid_client_metadata", result.JSON400.Error)
						exists, err := db.Account.Query().Where(account.Handle(*body.ClientName)).Exist(root)
						require.NoError(t, err)
						require.False(t, exists)
					})
				}))
			}))
		})
	}
}
