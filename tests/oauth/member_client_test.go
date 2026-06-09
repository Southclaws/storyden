package oauth_test

import (
	"context"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/rs/xid"
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

func TestOAuthMemberClientManagement(t *testing.T) {
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
	) {
		lc.Append(fx.StartHook(func() {
			firstCtx, first := e2e.WithAccount(root, aw, seed.Account_001_Odin)
			secondCtx, _ := e2e.WithAccount(root, aw, seed.Account_002_Frigg)
			nonAdminCtx, nonAdmin := e2e.WithAccount(root, aw, seed.Account_003_Baldur)
			firstSession := sh.WithSession(firstCtx)
			secondSession := sh.WithSession(secondCtx)
			nonAdminSession := sh.WithSession(nonAdminCtx)

			t.Run("create_lists_reads_updates_and_deletes_owned_client", func(t *testing.T) {
				a := assert.New(t)
				r := require.New(t)

				created := tests.AssertRequest(cl.OAuthClientCreateWithResponse(root, openapi.OAuthClientCreateJSONRequestBody{
					Name:          "Analytics Sync",
					Type:          openapi.OAuthClientTypeConfidential,
					AllowedScopes: []string{"READ_PUBLISHED_THREADS"},
					AllowedGrants: []string{oauthGrantClientCredentials},
					PkceRequired:  false,
				}, firstSession))(t, http.StatusOK)
				r.NotNil(created.JSON200)
				r.NotNil(created.JSON200.ClientSecret)
				a.NotEmpty(*created.JSON200.ClientSecret)
				a.True(strings.HasPrefix(created.JSON200.Client.ClientId, oauthresource.OAuthAccessKeyPrefix))
				a.True(strings.HasPrefix(*created.JSON200.ClientSecret, oauthresource.OAuthAccessSecretPrefix))
				a.Equal("Analytics Sync", created.JSON200.Client.Name)
				a.Equal(openapi.OAuthClientTypeConfidential, created.JSON200.Client.Type)
				a.Equal(openapi.Explicit, created.JSON200.Client.ScopePolicy)
				a.Equal([]string{oauthGrantClientCredentials}, created.JSON200.Client.AllowedGrants)
				a.Empty(created.JSON200.Client.RedirectUris)
				a.Contains(created.JSON200.Client.AllowedScopes, "READ_PUBLISHED_THREADS")

				token := tests.AssertRequest(oauthToken(t, root, cl, oauthTokenRequest{
					GrantType:    oauthGrantClientCredentials,
					ClientId:     created.JSON200.Client.ClientId,
					ClientSecret: created.JSON200.ClientSecret,
					Scope:        ptr("READ_PUBLISHED_THREADS"),
				}))(t, http.StatusOK)
				r.NotNil(token.JSON200)
				r.NotNil(token.JSON200.AccessToken)
				a.Equal("READ_PUBLISHED_THREADS", *token.JSON200.Scope)
				a.Nil(token.JSON200.RefreshToken)
				a.Nil(token.JSON200.IdToken)

				overscoped := tests.AssertRequest(oauthToken(t, root, cl, oauthTokenRequest{
					GrantType:    oauthGrantClientCredentials,
					ClientId:     created.JSON200.Client.ClientId,
					ClientSecret: created.JSON200.ClientSecret,
					Scope:        ptr("MANAGE_REPORTS"),
				}))(t, http.StatusBadRequest)
				r.NotNil(overscoped.JSON400)
				a.Equal("invalid_scope", overscoped.JSON400.Error)

				list := tests.AssertRequest(cl.OAuthClientListWithResponse(root, firstSession))(t, http.StatusOK)
				r.NotNil(list.JSON200)
				r.Len(list.JSON200.Clients, 1)
				a.Equal(created.JSON200.Client.Id, list.JSON200.Clients[0].Id)

				otherList := tests.AssertRequest(cl.OAuthClientListWithResponse(root, secondSession))(t, http.StatusOK)
				r.NotNil(otherList.JSON200)
				a.Empty(otherList.JSON200.Clients)

				getOther := tests.AssertRequest(cl.OAuthClientGetWithResponse(root, created.JSON200.Client.Id, secondSession))(t, http.StatusBadRequest)
				r.NotNil(getOther.JSON400)
				a.Equal("invalid_request", getOther.JSON400.Error)

				updated := tests.AssertRequest(cl.OAuthClientUpdateWithResponse(root, created.JSON200.Client.Id, openapi.OAuthClientUpdateJSONRequestBody{
					Name:          ptr("Reports Sync"),
					AllowedScopes: &[]string{"READ_PUBLISHED_THREADS", "MANAGE_REPORTS"},
				}, firstSession))(t, http.StatusOK)
				r.NotNil(updated.JSON200)
				a.Equal("Reports Sync", updated.JSON200.Name)
				a.Contains(updated.JSON200.AllowedScopes, "MANAGE_REPORTS")

				clientXID, err := xid.FromString(created.JSON200.Client.Id)
				r.NoError(err)
				_, err = ow.CreateRefreshToken(root, oauth_writer.RefreshTokenCreate{
					ClientID:  oauthresource.ClientID(clientXID),
					AccountID: first.ID,
					TokenHash: "delete-client-refresh-token",
					Scope:     "READ_PUBLISHED_THREADS",
					ExpiresAt: time.Now().Add(time.Hour),
				})
				r.NoError(err)

				tests.AssertRequest(cl.OAuthClientDeleteWithResponse(root, created.JSON200.Client.Id, firstSession))(t, http.StatusNoContent)

				afterDelete := tests.AssertRequest(cl.OAuthClientListWithResponse(root, firstSession))(t, http.StatusOK)
				r.NotNil(afterDelete.JSON200)
				a.Empty(afterDelete.JSON200.Clients)

				afterDeleteTokens := tests.AssertRequest(cl.OAuthRefreshTokenListWithResponse(root, firstSession))(t, http.StatusOK)
				r.NotNil(afterDeleteTokens.JSON200)
				a.Empty(afterDeleteTokens.JSON200.Tokens)
			})

			t.Run("create_rejects_unknown_permission_scope", func(t *testing.T) {
				a := assert.New(t)

				resp := tests.AssertRequest(cl.OAuthClientCreateWithResponse(root, openapi.OAuthClientCreateJSONRequestBody{
					Name:          "Bad Scope",
					Type:          openapi.OAuthClientTypeConfidential,
					AllowedScopes: []string{"NOT_A_PERMISSION"},
					AllowedGrants: []string{oauthGrantClientCredentials},
					PkceRequired:  false,
				}, firstSession))(t, http.StatusBadRequest)
				a.NotNil(resp)
			})

			t.Run("create_app_integration_client_with_authorization_code", func(t *testing.T) {
				a := assert.New(t)
				r := require.New(t)

				created := tests.AssertRequest(cl.OAuthClientCreateWithResponse(root, openapi.OAuthClientCreateJSONRequestBody{
					Name:          "Claude MCP Integration",
					Type:          openapi.OAuthClientTypeConfidential,
					AllowedScopes: []string{"READ_PUBLISHED_THREADS", "CREATE_POST"},
					AllowedGrants: []string{"authorization_code", "refresh_token"},
					RedirectUris:  &[]string{"https://claude.ai/api/mcp/auth_callback"},
					PkceRequired:  true,
				}, firstSession))(t, http.StatusOK)
				r.NotNil(created.JSON200)
				r.NotNil(created.JSON200.ClientSecret)
				a.NotEmpty(*created.JSON200.ClientSecret)
				a.Equal("Claude MCP Integration", created.JSON200.Client.Name)
				a.Equal(openapi.OAuthClientTypeConfidential, created.JSON200.Client.Type)
				a.Contains(created.JSON200.Client.AllowedGrants, "authorization_code")
				a.Contains(created.JSON200.Client.AllowedGrants, "refresh_token")
				a.Contains(created.JSON200.Client.RedirectUris, "https://claude.ai/api/mcp/auth_callback")
				a.Contains(created.JSON200.Client.AllowedScopes, "READ_PUBLISHED_THREADS")
				a.Contains(created.JSON200.Client.AllowedScopes, "CREATE_POST")
			})

			t.Run("create_public_app_client", func(t *testing.T) {
				a := assert.New(t)
				r := require.New(t)

				created := tests.AssertRequest(cl.OAuthClientCreateWithResponse(root, openapi.OAuthClientCreateJSONRequestBody{
					Name:          "Public Mobile App",
					Type:          openapi.OAuthClientTypePublic,
					AllowedScopes: []string{"READ_PUBLISHED_THREADS"},
					AllowedGrants: []string{"authorization_code", "refresh_token"},
					RedirectUris:  &[]string{"myapp://callback"},
					PkceRequired:  true,
				}, firstSession))(t, http.StatusOK)
				r.NotNil(created.JSON200)
				a.Nil(created.JSON200.ClientSecret)
				a.Equal("Public Mobile App", created.JSON200.Client.Name)
				a.Equal(openapi.OAuthClientTypePublic, created.JSON200.Client.Type)
				a.Contains(created.JSON200.Client.AllowedGrants, "authorization_code")
				a.Contains(created.JSON200.Client.AllowedGrants, "refresh_token")
				a.Contains(created.JSON200.Client.RedirectUris, "myapp://callback")
			})

			t.Run("create_rejects_authorization_code_without_redirect_uris", func(t *testing.T) {
				resp := tests.AssertRequest(cl.OAuthClientCreateWithResponse(root, openapi.OAuthClientCreateJSONRequestBody{
					Name:          "Bad App Integration",
					Type:          openapi.OAuthClientTypeConfidential,
					AllowedScopes: []string{"READ_PUBLISHED_THREADS"},
					AllowedGrants: []string{"authorization_code"},
					PkceRequired:  true,
				}, firstSession))(t, http.StatusBadRequest)
				require.NotNil(t, resp)
			})

			t.Run("create_rejects_public_client_with_client_credentials", func(t *testing.T) {
				resp := tests.AssertRequest(cl.OAuthClientCreateWithResponse(root, openapi.OAuthClientCreateJSONRequestBody{
					Name:          "Bad Public Client",
					Type:          openapi.OAuthClientTypePublic,
					AllowedScopes: []string{"READ_PUBLISHED_THREADS"},
					AllowedGrants: []string{"client_credentials"},
					PkceRequired:  false,
				}, firstSession))(t, http.StatusBadRequest)
				require.NotNil(t, resp)
			})

			t.Run("create_rejects_public_client_without_pkce", func(t *testing.T) {
				resp := tests.AssertRequest(cl.OAuthClientCreateWithResponse(root, openapi.OAuthClientCreateJSONRequestBody{
					Name:          "Public Client Without PKCE",
					Type:          openapi.OAuthClientTypePublic,
					AllowedScopes: []string{"READ_PUBLISHED_THREADS"},
					AllowedGrants: []string{"authorization_code"},
					RedirectUris:  &[]string{"myapp://callback"},
					PkceRequired:  false,
				}, firstSession))(t, http.StatusBadRequest)
				require.NotNil(t, resp)
			})

			t.Run("create_rejects_machine_client_with_redirect_uris", func(t *testing.T) {
				resp := tests.AssertRequest(cl.OAuthClientCreateWithResponse(root, openapi.OAuthClientCreateJSONRequestBody{
					Name:          "Machine Client With Redirects",
					Type:          openapi.OAuthClientTypeConfidential,
					AllowedScopes: []string{"READ_PUBLISHED_THREADS"},
					AllowedGrants: []string{"client_credentials"},
					RedirectUris:  &[]string{"https://example.com/callback"},
					PkceRequired:  false,
				}, firstSession))(t, http.StatusBadRequest)
				require.NotNil(t, resp)
			})

			t.Run("create_rejects_empty_allowed_grants", func(t *testing.T) {
				resp := tests.AssertRequest(cl.OAuthClientCreateWithResponse(root, openapi.OAuthClientCreateJSONRequestBody{
					Name:          "Client With No Grants",
					Type:          openapi.OAuthClientTypeConfidential,
					AllowedScopes: []string{"READ_PUBLISHED_THREADS"},
					AllowedGrants: []string{},
					PkceRequired:  false,
				}, firstSession))(t, http.StatusBadRequest)
				require.NotNil(t, resp)
			})

			t.Run("management_requires_administrator_permission", func(t *testing.T) {
				dummyID := xid.New().String()

				tests.AssertRequest(cl.OAuthClientListWithResponse(root, nonAdminSession))(t, http.StatusForbidden)
				tests.AssertRequest(cl.OAuthClientCreateWithResponse(root, openapi.OAuthClientCreateJSONRequestBody{
					Name:          "Not Allowed",
					AllowedScopes: []string{"READ_PUBLISHED_THREADS"},
				}, nonAdminSession))(t, http.StatusForbidden)
				tests.AssertRequest(cl.OAuthClientGetWithResponse(root, dummyID, nonAdminSession))(t, http.StatusForbidden)
				tests.AssertRequest(cl.OAuthClientUpdateWithResponse(root, dummyID, openapi.OAuthClientUpdateJSONRequestBody{}, nonAdminSession))(t, http.StatusForbidden)
				tests.AssertRequest(cl.OAuthClientDeleteWithResponse(root, dummyID, nonAdminSession))(t, http.StatusForbidden)
				tests.AssertRequest(cl.OAuthRefreshTokenListWithResponse(root, nonAdminSession))(t, http.StatusForbidden)
				tests.AssertRequest(cl.OAuthRefreshTokenDeleteWithResponse(root, dummyID, nonAdminSession))(t, http.StatusForbidden)
			})

			t.Run("non_admin_can_use_admin_created_device_client", func(t *testing.T) {
				r := require.New(t)

				grantOAuthClientUse(t, root, roles, assignments, nonAdmin.ID)

				clientID := "non-admin-device-" + xid.New().String()

				adminClient := tests.AssertRequest(cl.AdminOAuthClientCreateWithResponse(root, openapi.AdminOAuthClientCreateJSONRequestBody{
					AccountId:     first.ID.String(),
					ClientId:      clientID,
					Name:          "Non-Admin Device Client",
					Type:          openapi.OAuthClientTypePublic,
					ScopePolicy:   ptr(openapi.Inherit),
					AllowedScopes: []string{"openid", "profile"},
					AllowedGrants: []string{oauthGrantDeviceCode},
					RedirectUris:  []string{},
				}, firstSession))(t, http.StatusOK)
				r.NotNil(adminClient.JSON200)

				start := tests.AssertRequest(cl.OAuthDeviceAuthorisationWithFormdataBodyWithResponse(root, openapi.OAuthDeviceAuthorisationFormdataRequestBody{
					ClientId: clientID,
					Scope:    ptr("openid profile"),
				}))(t, http.StatusOK)
				r.NotNil(start.JSON200)
				r.NotNil(start.JSON200.DeviceCode)
				r.NotNil(start.JSON200.UserCode)

				tests.AssertRequest(cl.OAuthDeviceConsentWithResponse(root, &openapi.OAuthDeviceConsentParams{
					UserCode: start.JSON200.UserCode,
				}, nonAdminSession))(t, http.StatusOK)

				tests.AssertRequest(cl.OAuthDeviceConsentSubmitWithResponse(root, openapi.OAuthDeviceConsentSubmitJSONRequestBody{
					UserCode: *start.JSON200.UserCode,
					Decision: openapi.OAuthDeviceDecisionApprove,
				}, nonAdminSession))(t, http.StatusOK)

				token := tests.AssertRequest(oauthToken(t, root, cl, oauthTokenRequest{
					GrantType:  oauthGrantDeviceCode,
					ClientId:   clientID,
					DeviceCode: start.JSON200.DeviceCode,
				}))(t, http.StatusOK)
				r.NotNil(token.JSON200)
				r.NotNil(token.JSON200.AccessToken)
			})
		}))
	}))
}
