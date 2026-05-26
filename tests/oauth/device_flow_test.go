package oauth_test

import (
	"context"
	"net/http"
	"net/url"
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

func TestOAuthDeviceFlowConsentURLConfiguration(t *testing.T) {
	t.Parallel()

	cfg := oauthConfig(t)
	consentURL, err := url.Parse("https://custom.example/oauth/device-consent")
	require.NoError(t, err)
	cfg.OAuthDeviceAuthorisationConsentURL = *consentURL

	integration.Test(t, cfg, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		cl *openapi.ClientWithResponses,
		aw *account_writer.Writer,
		ow *oauth_writer.Writer,
	) {
		lc.Append(fx.StartHook(func() {
			a := assert.New(t)
			r := require.New(t)

			_, admin := e2e.WithAccount(root, aw, seed.Account_001_Odin)
			clientID := "custom-consent-url-" + uuid.NewString()
			createClient(t, root, ow, admin.ID, clientID, oauthresource.ClientTypePublic, oauthresource.ScopePolicyInheritUserPermissions, opt.NewEmpty[string](), standardScopes(), []string{oauthGrantDeviceCode})

			start := tests.AssertRequest(cl.OAuthDeviceAuthorisationWithFormdataBodyWithResponse(root, openapi.OAuthDeviceAuthorisationFormdataRequestBody{
				ClientId: clientID,
				Scope:    ptr("openid profile offline_access"),
			}))(t, http.StatusOK)
			r.NotNil(start.JSON200)
			r.NotNil(start.JSON200.VerificationUri)
			r.NotNil(start.JSON200.VerificationUriComplete)
			a.Equal("https://custom.example/oauth/device-consent", *start.JSON200.VerificationUri)
			a.Contains(*start.JSON200.VerificationUriComplete, "https://custom.example/oauth/device-consent?user_code=")
		}))
	}))
}

func TestOAuthDeviceFlowPermissionPolicies(t *testing.T) {
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
			adminCtx, admin := e2e.WithAccount(root, aw, seed.Account_001_Odin)
			memberCtx, member := e2e.WithAccount(root, aw, seed.Account_004_Loki)
			grantOAuthClientUse(t, root, roles, assignments, member.ID, rbac.PermissionCreatePost)
			adminSession := sh.WithSession(adminCtx)
			memberSession := sh.WithSession(memberCtx)

			t.Run("built_in_storyden_cli_client_starts_without_setup", func(t *testing.T) {
				a := assert.New(t)
				r := require.New(t)

				start := tests.AssertRequest(cl.OAuthDeviceAuthorisationWithFormdataBodyWithResponse(root, openapi.OAuthDeviceAuthorisationFormdataRequestBody{
					ClientId: "storyden-cli",
					Scope:    ptr("openid profile offline_access"),
				}))(t, http.StatusOK)
				r.NotNil(start.JSON200)
				r.NotNil(start.JSON200.DeviceCode)
				r.NotNil(start.JSON200.UserCode)
				r.NotNil(start.JSON200.VerificationUriComplete)
				a.Contains(*start.JSON200.VerificationUriComplete, "user_code=")
			})

			t.Run("non_admin_member_can_authorize_admin_created_client", func(t *testing.T) {
				a := assert.New(t)
				r := require.New(t)

				clientID := "admin-created-device-" + uuid.NewString()
				createClient(t, root, ow, admin.ID, clientID, oauthresource.ClientTypePublic, oauthresource.ScopePolicyExplicit, opt.NewEmpty[string](), append(standardScopes(), rbac.PermissionCreatePost.String()), []string{oauthGrantDeviceCode})

				start := tests.AssertRequest(cl.OAuthDeviceAuthorisationWithFormdataBodyWithResponse(root, openapi.OAuthDeviceAuthorisationFormdataRequestBody{
					ClientId: clientID,
					Scope:    ptr("openid profile CREATE_POST"),
				}))(t, http.StatusOK)
				r.NotNil(start.JSON200)
				r.NotNil(start.JSON200.DeviceCode)
				r.NotNil(start.JSON200.UserCode)

				consent := tests.AssertRequest(cl.OAuthDeviceConsentWithResponse(root, &openapi.OAuthDeviceConsentParams{
					UserCode: start.JSON200.UserCode,
				}, memberSession))(t, http.StatusOK)
				r.NotNil(consent.JSON200)
				a.Equal(clientID, consent.JSON200.ClientId)
				a.Contains(consent.JSON200.GrantedScopes, rbac.PermissionCreatePost.String())

				approve := tests.AssertRequest(cl.OAuthDeviceConsentSubmitWithResponse(root, openapi.OAuthDeviceConsentSubmitJSONRequestBody{
					UserCode: *start.JSON200.UserCode,
					Decision: openapi.OAuthDeviceDecisionApprove,
				}, memberSession))(t, http.StatusOK)
				r.NotNil(approve.JSON200)
				a.Equal(openapi.OAuthDeviceConsentResultStatusApproved, approve.JSON200.Status)

				token := tests.AssertRequest(oauthToken(t, root, cl, oauthTokenRequest{
					GrantType:  oauthGrantDeviceCode,
					ClientId:   clientID,
					DeviceCode: start.JSON200.DeviceCode,
				}))(t, http.StatusOK)
				r.NotNil(token.JSON200)
				r.NotNil(token.JSON200.Scope)
				a.Contains(*token.JSON200.Scope, rbac.PermissionCreatePost.String())
			})

			t.Run("built_in_storyden_cli_client_rejects_missing_scope_with_oauth_error", func(t *testing.T) {
				a := assert.New(t)
				r := require.New(t)

				start := tests.AssertRequest(cl.OAuthDeviceAuthorisationWithFormdataBodyWithResponse(root, openapi.OAuthDeviceAuthorisationFormdataRequestBody{
					ClientId: "storyden-cli",
				}))(t, http.StatusBadRequest)
				r.NotNil(start.JSON400)
				a.Equal("invalid_scope", start.JSON400.Error)
			})

			t.Run("explicit_device_client_allows_omitted_scope", func(t *testing.T) {
				a := assert.New(t)
				r := require.New(t)

				clientID := "empty-scope-device-" + uuid.NewString()
				createClient(t, root, ow, admin.ID, clientID, oauthresource.ClientTypePublic, oauthresource.ScopePolicyExplicit, opt.NewEmpty[string](), standardScopes(), []string{oauthGrantDeviceCode})

				start := tests.AssertRequest(cl.OAuthDeviceAuthorisationWithFormdataBodyWithResponse(root, openapi.OAuthDeviceAuthorisationFormdataRequestBody{
					ClientId: clientID,
				}))(t, http.StatusOK)
				r.NotNil(start.JSON200)
				r.NotNil(start.JSON200.DeviceCode)
				r.NotNil(start.JSON200.UserCode)

				consent := tests.AssertRequest(cl.OAuthDeviceConsentWithResponse(root, &openapi.OAuthDeviceConsentParams{
					UserCode: start.JSON200.UserCode,
				}, adminSession))(t, http.StatusOK)
				r.NotNil(consent.JSON200)
				a.Empty(consent.JSON200.RequestedScopes)
				a.Empty(consent.JSON200.GrantedScopes)
			})

			t.Run("public_client_with_inherit_policy_inherits_effective_user_permissions", func(t *testing.T) {
				a := assert.New(t)
				r := require.New(t)

				clientID := "storyden-cli-" + uuid.NewString()
				createClient(t, root, ow, admin.ID, clientID, oauthresource.ClientTypePublic, oauthresource.ScopePolicyInheritUserPermissions, opt.NewEmpty[string](), standardScopes(), []string{oauthGrantDeviceCode, oauthGrantRefreshToken})

				start := tests.AssertRequest(cl.OAuthDeviceAuthorisationWithFormdataBodyWithResponse(root, openapi.OAuthDeviceAuthorisationFormdataRequestBody{
					ClientId: clientID,
					Scope:    ptr("openid profile offline_access"),
				}))(t, http.StatusOK)
				r.NotNil(start.JSON200)
				r.NotNil(start.JSON200.DeviceCode)
				r.NotNil(start.JSON200.UserCode)
				r.NotNil(start.JSON200.VerificationUri)
				r.NotNil(start.JSON200.VerificationUriComplete)
				a.Equal("http://localhost:3000/oauth/consent", *start.JSON200.VerificationUri)
				a.Contains(*start.JSON200.VerificationUriComplete, "http://localhost:3000/oauth/consent?user_code=")

				consent := tests.AssertRequest(cl.OAuthDeviceConsentWithResponse(root, &openapi.OAuthDeviceConsentParams{
					UserCode: start.JSON200.UserCode,
				}, adminSession))(t, http.StatusOK)
				r.NotNil(consent.JSON200)
				a.Equal(clientID, consent.JSON200.ClientId)
				a.True(consent.JSON200.InheritsUserPermissions)

				approve := tests.AssertRequest(cl.OAuthDeviceConsentSubmitWithResponse(root, openapi.OAuthDeviceConsentSubmitJSONRequestBody{
					UserCode: *start.JSON200.UserCode,
					Decision: openapi.OAuthDeviceDecisionApprove,
				}, adminSession))(t, http.StatusOK)
				r.NotNil(approve.JSON200)
				a.Equal(openapi.OAuthDeviceConsentResultStatusApproved, approve.JSON200.Status)

				doubleApprove := tests.AssertRequest(cl.OAuthDeviceConsentSubmitWithResponse(root, openapi.OAuthDeviceConsentSubmitJSONRequestBody{
					UserCode: *start.JSON200.UserCode,
					Decision: openapi.OAuthDeviceDecisionApprove,
				}, adminSession))(t, http.StatusBadRequest)
				r.NotNil(doubleApprove.JSON400)
				a.Equal("invalid_request", doubleApprove.JSON400.Error)

				token := tests.AssertRequest(oauthToken(t, root, cl, oauthTokenRequest{
					GrantType:  oauthGrantDeviceCode,
					ClientId:   clientID,
					DeviceCode: start.JSON200.DeviceCode,
				}))(t, http.StatusOK)
				r.NotNil(token.JSON200.AccessToken)
				r.NotNil(token.JSON200.RefreshToken)
				r.NotNil(token.JSON200.Scope)
				a.Contains(*token.JSON200.Scope, "openid")
				a.Contains(*token.JSON200.Scope, rbac.PermissionAdministrator.String())

				reuse := tests.AssertRequest(oauthToken(t, root, cl, oauthTokenRequest{
					GrantType:  oauthGrantDeviceCode,
					ClientId:   clientID,
					DeviceCode: start.JSON200.DeviceCode,
				}))(t, http.StatusBadRequest)
				r.NotNil(reuse.JSON400)
				a.Equal("invalid_grant", reuse.JSON400.Error)

				adminSettings := tests.AssertRequest(cl.AdminSettingsGetWithResponse(root, bearer(*token.JSON200.AccessToken)))(t, http.StatusOK)
				r.NotNil(adminSettings.JSON200)
			})

			t.Run("explicit_client_cannot_expand_beyond_member_permissions", func(t *testing.T) {
				a := assert.New(t)
				r := require.New(t)

				clientID := "analytics-bot-" + uuid.NewString()
				createClient(t, root, ow, admin.ID, clientID, oauthresource.ClientTypePublic, oauthresource.ScopePolicyExplicit, opt.NewEmpty[string](), append(standardScopes(), rbac.PermissionCreatePost.String(), rbac.PermissionManageReports.String()), []string{oauthGrantDeviceCode})

				start := tests.AssertRequest(cl.OAuthDeviceAuthorisationWithFormdataBodyWithResponse(root, openapi.OAuthDeviceAuthorisationFormdataRequestBody{
					ClientId: clientID,
					Scope:    ptr("openid profile offline_access CREATE_POST MANAGE_REPORTS"),
				}))(t, http.StatusOK)

				tests.AssertRequest(cl.OAuthDeviceConsentWithResponse(root, &openapi.OAuthDeviceConsentParams{
					UserCode: start.JSON200.UserCode,
				}, memberSession))(t, http.StatusOK)

				approve := tests.AssertRequest(cl.OAuthDeviceConsentSubmitWithResponse(root, openapi.OAuthDeviceConsentSubmitJSONRequestBody{
					UserCode: *start.JSON200.UserCode,
					Decision: openapi.OAuthDeviceDecisionApprove,
				}, memberSession))(t, http.StatusOK)
				r.NotNil(approve.JSON200)
				a.Equal(openapi.OAuthDeviceConsentResultStatusApproved, approve.JSON200.Status)

				token := tests.AssertRequest(oauthToken(t, root, cl, oauthTokenRequest{
					GrantType:  oauthGrantDeviceCode,
					ClientId:   clientID,
					DeviceCode: start.JSON200.DeviceCode,
				}))(t, http.StatusOK)
				r.NotNil(token.JSON200)
				r.NotNil(token.JSON200.Scope)
				a.Contains(*token.JSON200.Scope, rbac.PermissionCreatePost.String())
				a.NotContains(*token.JSON200.Scope, rbac.PermissionManageReports.String())

				adminSettings := tests.AssertRequest(cl.AdminSettingsGetWithResponse(root, bearer(*token.JSON200.AccessToken)))(t, http.StatusForbidden)
				a.NotNil(adminSettings)

				category := tests.AssertRequest(cl.CategoryCreateWithResponse(root, openapi.CategoryInitialProps{
					Name: "oauth-explicit-category-" + uuid.NewString(),
				}, adminSession))(t, http.StatusOK)

				thread := tests.AssertRequest(cl.ThreadCreateWithResponse(root, openapi.ThreadInitialProps{
					Title:      "oauth explicit thread " + uuid.NewString(),
					Body:       ptr("<p>created with an oauth token</p>"),
					Category:   &category.JSON200.Id,
					Visibility: ptr(openapi.Published),
				}, bearer(*token.JSON200.AccessToken)))(t, http.StatusOK)
				r.NotNil(thread.JSON200)
			})
		}))
	}))
}

func TestOAuthDeviceFlowRequiresOAuthClientPermission(t *testing.T) {
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
			adminCtx, admin := e2e.WithAccount(root, aw, seed.Account_001_Odin)
			adminSession := sh.WithSession(adminCtx)

			t.Run("member_without_permission_cannot_claim_or_approve_device_consent", func(t *testing.T) {
				a := assert.New(t)
				r := require.New(t)

				memberCtx, _ := e2e.WithAccount(root, aw, seed.Account_004_Loki)
				memberSession := sh.WithSession(memberCtx)

				clientID := "device-permission-required-" + uuid.NewString()
				createClient(t, root, ow, admin.ID, clientID, oauthresource.ClientTypePublic, oauthresource.ScopePolicyExplicit, opt.NewEmpty[string](), append(standardScopes(), rbac.PermissionCreatePost.String()), []string{oauthGrantDeviceCode})

				start := tests.AssertRequest(cl.OAuthDeviceAuthorisationWithFormdataBodyWithResponse(root, openapi.OAuthDeviceAuthorisationFormdataRequestBody{
					ClientId: clientID,
					Scope:    ptr("openid profile CREATE_POST"),
				}))(t, http.StatusOK)
				r.NotNil(start.JSON200)
				r.NotNil(start.JSON200.DeviceCode)
				r.NotNil(start.JSON200.UserCode)

				consent := tests.AssertRequest(cl.OAuthDeviceConsentWithResponse(root, &openapi.OAuthDeviceConsentParams{
					UserCode: start.JSON200.UserCode,
				}, memberSession))(t, http.StatusBadRequest)
				r.NotNil(consent.JSON400)
				a.Equal("access_denied", consent.JSON400.Error)

				approve := tests.AssertRequest(cl.OAuthDeviceConsentSubmitWithResponse(root, openapi.OAuthDeviceConsentSubmitJSONRequestBody{
					UserCode: *start.JSON200.UserCode,
					Decision: openapi.OAuthDeviceDecisionApprove,
				}, memberSession))(t, http.StatusBadRequest)
				r.NotNil(approve.JSON400)
				a.Equal("access_denied", approve.JSON400.Error)

				token := tests.AssertRequest(oauthToken(t, root, cl, oauthTokenRequest{
					GrantType:  oauthGrantDeviceCode,
					ClientId:   clientID,
					DeviceCode: start.JSON200.DeviceCode,
				}))(t, http.StatusBadRequest)
				r.NotNil(token.JSON400)
				a.Equal("authorization_pending", token.JSON400.Error)
			})

			t.Run("permission_revocation_blocks_pending_device_approval", func(t *testing.T) {
				a := assert.New(t)
				r := require.New(t)

				memberCtx, member := e2e.WithAccount(root, aw, seed.Account_003_Baldur)
				oauthRoleID := grantOAuthClientUse(t, root, roles, assignments, member.ID, rbac.PermissionCreatePost)
				memberSession := sh.WithSession(memberCtx)

				clientID := "device-permission-revoked-" + uuid.NewString()
				createClient(t, root, ow, admin.ID, clientID, oauthresource.ClientTypePublic, oauthresource.ScopePolicyExplicit, opt.NewEmpty[string](), append(standardScopes(), rbac.PermissionCreatePost.String()), []string{oauthGrantDeviceCode})

				start := tests.AssertRequest(cl.OAuthDeviceAuthorisationWithFormdataBodyWithResponse(root, openapi.OAuthDeviceAuthorisationFormdataRequestBody{
					ClientId: clientID,
					Scope:    ptr("openid profile CREATE_POST"),
				}))(t, http.StatusOK)
				r.NotNil(start.JSON200)
				r.NotNil(start.JSON200.DeviceCode)
				r.NotNil(start.JSON200.UserCode)

				claimed := tests.AssertRequest(cl.OAuthDeviceConsentWithResponse(root, &openapi.OAuthDeviceConsentParams{
					UserCode: start.JSON200.UserCode,
				}, memberSession))(t, http.StatusOK)
				r.NotNil(claimed.JSON200)

				revokeOAuthClientUse(t, root, assignments, member.ID, oauthRoleID)

				approve := tests.AssertRequest(cl.OAuthDeviceConsentSubmitWithResponse(root, openapi.OAuthDeviceConsentSubmitJSONRequestBody{
					UserCode: *start.JSON200.UserCode,
					Decision: openapi.OAuthDeviceDecisionApprove,
				}, memberSession))(t, http.StatusBadRequest)
				r.NotNil(approve.JSON400)
				a.Equal("access_denied", approve.JSON400.Error)

				deny := tests.AssertRequest(cl.OAuthDeviceConsentSubmitWithResponse(root, openapi.OAuthDeviceConsentSubmitJSONRequestBody{
					UserCode: *start.JSON200.UserCode,
					Decision: openapi.OAuthDeviceDecisionDeny,
				}, memberSession))(t, http.StatusOK)
				r.NotNil(deny.JSON200)
				a.Equal(openapi.OAuthDeviceConsentResultStatusDenied, deny.JSON200.Status)

				adminConsent := tests.AssertRequest(cl.OAuthDeviceConsentWithResponse(root, &openapi.OAuthDeviceConsentParams{
					UserCode: start.JSON200.UserCode,
				}, adminSession))(t, http.StatusBadRequest)
				r.NotNil(adminConsent.JSON400)
				a.Equal("invalid_request", adminConsent.JSON400.Error)
			})
		}))
	}))
}

func TestOAuthDeviceFlowDefensiveBehaviour(t *testing.T) {
	t.Parallel()

	integration.Test(t, oauthConfig(t), e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		cl *openapi.ClientWithResponses,
		sh *e2e.SessionHelper,
		aw *account_writer.Writer,
		ow *oauth_writer.Writer,
	) {
		lc.Append(fx.StartHook(func() {
			adminCtx, admin := e2e.WithAccount(root, aw, seed.Account_001_Odin)
			adminSession := sh.WithSession(adminCtx)

			t.Run("polling_reports_pending_and_slow_down_before_approval", func(t *testing.T) {
				a := assert.New(t)
				r := require.New(t)

				clientID := "polling-cli-" + uuid.NewString()
				createClient(t, root, ow, admin.ID, clientID, oauthresource.ClientTypePublic, oauthresource.ScopePolicyInheritUserPermissions, opt.NewEmpty[string](), standardScopes(), []string{oauthGrantDeviceCode})

				start := tests.AssertRequest(cl.OAuthDeviceAuthorisationWithFormdataBodyWithResponse(root, openapi.OAuthDeviceAuthorisationFormdataRequestBody{
					ClientId: clientID,
					Scope:    ptr("openid profile offline_access"),
				}))(t, http.StatusOK)
				r.NotNil(start.JSON200)
				r.NotNil(start.JSON200.DeviceCode)

				pending := tests.AssertRequest(oauthToken(t, root, cl, oauthTokenRequest{
					GrantType:  oauthGrantDeviceCode,
					ClientId:   clientID,
					DeviceCode: start.JSON200.DeviceCode,
				}))(t, http.StatusBadRequest)
				r.NotNil(pending.JSON400)
				a.Equal("authorization_pending", pending.JSON400.Error)

				slowDown := tests.AssertRequest(oauthToken(t, root, cl, oauthTokenRequest{
					GrantType:  oauthGrantDeviceCode,
					ClientId:   clientID,
					DeviceCode: start.JSON200.DeviceCode,
				}))(t, http.StatusBadRequest)
				r.NotNil(slowDown.JSON400)
				a.Equal("slow_down", slowDown.JSON400.Error)
			})

			t.Run("denied_device_authorization_cannot_be_approved_or_exchanged", func(t *testing.T) {
				a := assert.New(t)
				r := require.New(t)

				clientID := "denied-device-" + uuid.NewString()
				createClient(t, root, ow, admin.ID, clientID, oauthresource.ClientTypePublic, oauthresource.ScopePolicyInheritUserPermissions, opt.NewEmpty[string](), standardScopes(), []string{oauthGrantDeviceCode})

				start := tests.AssertRequest(cl.OAuthDeviceAuthorisationWithFormdataBodyWithResponse(root, openapi.OAuthDeviceAuthorisationFormdataRequestBody{
					ClientId: clientID,
					Scope:    ptr("openid profile offline_access"),
				}))(t, http.StatusOK)
				r.NotNil(start.JSON200)
				r.NotNil(start.JSON200.DeviceCode)
				r.NotNil(start.JSON200.UserCode)

				tests.AssertRequest(cl.OAuthDeviceConsentWithResponse(root, &openapi.OAuthDeviceConsentParams{
					UserCode: start.JSON200.UserCode,
				}, adminSession))(t, http.StatusOK)

				deny := tests.AssertRequest(cl.OAuthDeviceConsentSubmitWithResponse(root, openapi.OAuthDeviceConsentSubmitJSONRequestBody{
					UserCode: *start.JSON200.UserCode,
					Decision: openapi.OAuthDeviceDecisionDeny,
				}, adminSession))(t, http.StatusOK)
				r.NotNil(deny.JSON200)
				a.Equal(openapi.OAuthDeviceConsentResultStatusDenied, deny.JSON200.Status)

				approveAfterDeny := tests.AssertRequest(cl.OAuthDeviceConsentSubmitWithResponse(root, openapi.OAuthDeviceConsentSubmitJSONRequestBody{
					UserCode: *start.JSON200.UserCode,
					Decision: openapi.OAuthDeviceDecisionApprove,
				}, adminSession))(t, http.StatusBadRequest)
				r.NotNil(approveAfterDeny.JSON400)
				a.Equal("invalid_request", approveAfterDeny.JSON400.Error)

				token := tests.AssertRequest(oauthToken(t, root, cl, oauthTokenRequest{
					GrantType:  oauthGrantDeviceCode,
					ClientId:   clientID,
					DeviceCode: start.JSON200.DeviceCode,
				}))(t, http.StatusBadRequest)
				r.NotNil(token.JSON400)
				a.Equal("access_denied", token.JSON400.Error)
			})

			t.Run("unknown_or_disallowed_scope_is_rejected_at_device_authorization", func(t *testing.T) {
				a := assert.New(t)
				r := require.New(t)

				clientID := "scope-denied-" + uuid.NewString()
				createClient(t, root, ow, admin.ID, clientID, oauthresource.ClientTypePublic, oauthresource.ScopePolicyExplicit, opt.NewEmpty[string](), standardScopes(), []string{oauthGrantDeviceCode})

				denied := tests.AssertRequest(cl.OAuthDeviceAuthorisationWithFormdataBodyWithResponse(root, openapi.OAuthDeviceAuthorisationFormdataRequestBody{
					ClientId: clientID,
					Scope:    ptr("openid MANAGE_REPORTS"),
				}))(t, http.StatusBadRequest)
				r.NotNil(denied.JSON400)
				a.Equal("invalid_scope", denied.JSON400.Error)
			})
		}))
	}))
}
