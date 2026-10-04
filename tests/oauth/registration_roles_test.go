package oauth_test

import (
	"context"
	"crypto/rsa"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/oapi-codegen/nullable"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/ent/accountroles"
	"github.com/Southclaws/storyden/internal/ent/oauthclient"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/tests"
)

func agentAccessToken(t *testing.T, ctx context.Context, cl *openapi.ClientWithResponses, clientID string, key *rsa.PrivateKey, scope string) *openapi.OAuthToken {
	t.Helper()
	assertion := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iss": clientID, "sub": clientID, "aud": "http://localhost:8000/api/oauth/token",
		"iat": time.Now().Unix(), "exp": time.Now().Add(time.Minute).Unix(), "jti": uuid.NewString(),
	})
	assertion.Header["kid"] = "agent-key"
	signed, err := assertion.SignedString(key)
	require.NoError(t, err)
	form := url.Values{"grant_type": {"client_credentials"}, "client_id": {clientID}, "client_assertion_type": {"urn:ietf:params:oauth:client-assertion-type:jwt-bearer"}, "client_assertion": {signed}}
	if scope != "" {
		form.Set("scope", scope)
	}
	response := tests.AssertRequest(cl.OAuthTokenWithBodyWithResponse(ctx, "application/x-www-form-urlencoded", strings.NewReader(form.Encode())))(t, http.StatusOK)
	require.NotNil(t, response.JSON200)
	return (*openapi.OAuthToken)(response.JSON200)
}

func TestAutonomousRegistrationRoles(t *testing.T) {
	t.Parallel()
	integration.Test(t, oauthConfig(t), e2e.Setup(), withOAuthRegistration(t, "open"), fx.Invoke(func(lc fx.Lifecycle, root context.Context, cl *openapi.ClientWithResponses, sh *e2e.SessionHelper, aw *account_writer.Writer, db *ent.Client) {
		lc.Append(fx.StartHook(func() {
			adminCtx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
			adminSession := sh.WithSession(adminCtx)
			createRole := func(name string, permission string) openapi.Identifier {
				t.Helper()
				response := tests.AssertRequest(cl.RoleCreateWithResponse(root, openapi.RoleCreateJSONRequestBody{Name: name, Colour: "green", Permissions: openapi.PermissionList{openapi.Permission(permission)}}, adminSession))(t, http.StatusOK)
				return response.JSON200.Id
			}
			update := func(value openapi.OAuthServiceSettings) *openapi.OAuthServiceSettings {
				t.Helper()
				response := tests.AssertRequest(cl.AdminSettingsUpdateWithResponse(root, openapi.AdminSettingsUpdateJSONRequestBody{Services: &openapi.AdminSettingsServiceProps{Oauth: &value}}, adminSession))(t, http.StatusOK)
				return response.JSON200.Services.Oauth
			}
			botRole := createRole("Bot", "MANAGE_LIBRARY")
			update(openapi.OAuthServiceSettings{AutonomousRegistrationRoleId: nullable.NewNullableWithValue(openapi.NullableIdentifier(botRole))})
			preserved := update(openapi.OAuthServiceSettings{RegistrationApprovalTtl: new(120)})
			require.Equal(t, string(botRole), string(preserved.AutonomousRegistrationRoleId.MustGet()))

			body, key := agentRegistration(t)
			registered := tests.AssertRequest(cl.OAuthClientRegisterWithResponse(root, body))(t, http.StatusCreated).JSON201
			client, err := db.OAuthClient.Query().Where(oauthclient.ClientID(registered.ClientId)).Only(root)
			require.NoError(t, err)
			require.Equal(t, "inherit", string(client.ScopePolicy))
			require.Empty(t, client.AllowedScopes)
			require.NotNil(t, client.AccountID)
			held, err := db.AccountRoles.Query().Where(accountroles.AccountID(*client.AccountID)).All(root)
			require.NoError(t, err)
			require.Len(t, held, 1)
			require.Equal(t, string(botRole), held[0].RoleID.String())
			first := agentAccessToken(t, root, cl, registered.ClientId, key, "")
			require.Contains(t, *first.Scope, "READ_PUBLISHED_LIBRARY")
			require.Contains(t, *first.Scope, "MANAGE_LIBRARY")
			require.Contains(t, *first.Scope, "CREATE_POST")
			require.NotContains(t, *first.Scope, "MANAGE_SETTINGS")
			tests.AssertRequest(cl.NodeListWithResponse(root, &openapi.NodeListParams{}, bearer(*first.AccessToken)))(t, http.StatusOK)
			tests.AssertRequest(cl.AdminSettingsGetWithResponse(root, bearer(*first.AccessToken)))(t, http.StatusForbidden)

			managerRole := createRole("Bot manager", "MANAGE_SETTINGS")
			tests.AssertRequest(cl.AccountAddRoleWithResponse(root, *body.ClientName, managerRole, adminSession))(t, http.StatusOK)
			expanded := agentAccessToken(t, root, cl, registered.ClientId, key, "")
			require.Contains(t, *expanded.Scope, "MANAGE_SETTINGS")
			tests.AssertRequest(cl.AdminSettingsGetWithResponse(root, bearer(*expanded.AccessToken)))(t, http.StatusOK)
			tests.AssertRequest(cl.AdminSettingsUpdateWithResponse(root, openapi.AdminSettingsUpdateJSONRequestBody{Services: &openapi.AdminSettingsServiceProps{Oauth: &openapi.OAuthServiceSettings{AutonomousRegistrationRoleId: nullable.NewNullableWithValue(openapi.NullableIdentifier("00000000000000000a00"))}}}, bearer(*expanded.AccessToken)))(t, http.StatusForbidden)

			tests.AssertRequest(cl.AdminSettingsGetWithResponse(root, bearer(*first.AccessToken)))(t, http.StatusForbidden)
			tests.AssertRequest(cl.AccountRemoveRoleWithResponse(root, *body.ClientName, managerRole, adminSession))(t, http.StatusOK)
			tests.AssertRequest(cl.AdminSettingsGetWithResponse(root, bearer(*expanded.AccessToken)))(t, http.StatusForbidden)
			reduced := agentAccessToken(t, root, cl, registered.ClientId, key, "")
			require.NotContains(t, *reduced.Scope, "MANAGE_SETTINGS")

			restricted := agentAccessToken(t, root, cl, registered.ClientId, key, "READ_PUBLISHED_THREADS")
			require.Equal(t, "READ_PUBLISHED_THREADS", *restricted.Scope)
			tests.AssertRequest(cl.NodeListWithResponse(root, &openapi.NodeListParams{}, bearer(*restricted.AccessToken)))(t, http.StatusForbidden)
			attemptedAdmin := agentAccessToken(t, root, cl, registered.ClientId, key, "ADMINISTRATOR")
			require.Empty(t, *attemptedAdmin.Scope)
			tests.AssertRequest(cl.AdminSettingsGetWithResponse(root, bearer(*attemptedAdmin.AccessToken)))(t, http.StatusForbidden)

			tests.AssertRequest(cl.RoleDeleteWithResponse(root, botRole, adminSession))(t, http.StatusOK)
			read := tests.AssertRequest(cl.AdminSettingsGetWithResponse(root, adminSession))(t, http.StatusOK)
			require.Equal(t, string(botRole), string(read.JSON200.Services.Oauth.AutonomousRegistrationRoleId.MustGet()))
			update(openapi.OAuthServiceSettings{RegistrationApprovalTtl: new(180)})
			orphanBody, _ := agentRegistration(t)
			orphan := tests.AssertRequest(cl.OAuthClientRegisterWithResponse(root, orphanBody))(t, http.StatusCreated).JSON201
			orphanClient, err := db.OAuthClient.Query().Where(oauthclient.ClientID(orphan.ClientId)).Only(root)
			require.NoError(t, err)
			held, err = db.AccountRoles.Query().Where(accountroles.AccountID(*orphanClient.AccountID)).All(root)
			require.NoError(t, err)
			require.Empty(t, held)

			update(openapi.OAuthServiceSettings{AutonomousRegistrationMode: new(openapi.OAuthAutonomousRegistrationMode("approval"))})
			pendingBody, _ := agentRegistration(t)
			pendingBody.RegistrationMode = new([]string{"approval"})
			pending := tests.AssertRequest(cl.OAuthClientRegisterWithResponse(root, pendingBody))(t, http.StatusAccepted).JSON202
			update(openapi.OAuthServiceSettings{AutonomousRegistrationRoleId: nullable.NewNullableWithValue(openapi.NullableIdentifier(managerRole))})
			tests.AssertRequest(cl.AdminOAuthRegistrationApprovalSubmitWithResponse(root, openapi.AdminOAuthRegistrationApprovalSubmitJSONRequestBody{VerificationCode: *pending.VerificationCode, Approved: true}, adminSession))(t, http.StatusNoContent)
			require.NoError(t, db.OAuthRegistrationApproval.Update().SetNextPollAt(time.Now().Add(-time.Second)).Exec(root))
			approved := tests.AssertRequest(cl.OAuthClientRegisterWithResponse(root, openapi.OAuthClientRegisterJSONRequestBody{RegistrationCode: pending.RegistrationCode}))(t, http.StatusCreated).JSON201
			approvedClient, err := db.OAuthClient.Query().Where(oauthclient.ClientID(approved.ClientId)).Only(root)
			require.NoError(t, err)
			held, err = db.AccountRoles.Query().Where(accountroles.AccountID(*approvedClient.AccountID)).All(root)
			require.NoError(t, err)
			require.Len(t, held, 1)
			require.Equal(t, string(managerRole), held[0].RoleID.String())

			cleared := update(openapi.OAuthServiceSettings{AutonomousRegistrationRoleId: nullable.NewNullNullable[openapi.NullableIdentifier]()})
			require.True(t, cleared.AutonomousRegistrationRoleId.IsNull())
			read = tests.AssertRequest(cl.AdminSettingsGetWithResponse(root, adminSession))(t, http.StatusOK)
			require.True(t, read.JSON200.Services.Oauth.AutonomousRegistrationRoleId.IsNull())

			update(openapi.OAuthServiceSettings{AutonomousRegistrationMode: new(openapi.OAuthAutonomousRegistrationMode("open")), AutonomousRegistrationRoleId: nullable.NewNullableWithValue(openapi.NullableIdentifier("000000000000000000m0"))})
			memberBody, _ := agentRegistration(t)
			member := tests.AssertRequest(cl.OAuthClientRegisterWithResponse(root, memberBody))(t, http.StatusCreated).JSON201
			memberClient, err := db.OAuthClient.Query().Where(oauthclient.ClientID(member.ClientId)).Only(root)
			require.NoError(t, err)
			held, err = db.AccountRoles.Query().Where(accountroles.AccountID(*memberClient.AccountID)).All(root)
			require.NoError(t, err)
			require.Empty(t, held)
		}))
	}))
}
