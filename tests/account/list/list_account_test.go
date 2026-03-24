package account_list_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/tests"
)

func TestAdminAccountList(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		cl *openapi.ClientWithResponses,
		sh *e2e.SessionHelper,
		aw *account_writer.Writer,
	) {
		lc.Append(fx.StartHook(func() {
			r := require.New(t)
			a := assert.New(t)

			adminCtx, admin := e2e.WithAccount(root, aw, seed.Account_001_Odin)
			adminSession := sh.WithSession(adminCtx)

			supportCtx, support := e2e.WithAccount(root, aw, seed.Account_003_Baldur)
			supportSession := sh.WithSession(supportCtx)
			grant(t, cl, adminSession, support.Handle, openapi.PermissionList{openapi.VIEWACCOUNTS})

			group := "acctlist-" + xid.New().String()[:8]
			plain := signupMember(t, root, cl, adminSession, group+"-plain", nil)

			invitation := tests.AssertRequest(cl.InvitationCreateWithResponse(
				root,
				openapi.InvitationCreateJSONRequestBody{},
				adminSession,
			))(t, http.StatusOK)
			invitationID := openapi.InvitationIDQueryParam(invitation.JSON200.Id)

			invited := signupMember(t, root, cl, adminSession, group+"-invited", &invitationID)
			suspended := signupMember(t, root, cl, adminSession, group+"-suspended", nil)

			roleID := createRole(t, root, cl, adminSession, "role-account-list-"+xid.New().String(), openapi.PermissionList{})
			tests.AssertRequest(cl.AccountAddRoleWithResponse(
				root,
				invited.Handle,
				roleID,
				adminSession,
			))(t, http.StatusOK)

			tests.AssertRequest(cl.AdminAccountBanCreateWithResponse(
				root,
				suspended.Handle,
				adminSession,
			))(t, http.StatusOK)

			t.Run("admin_can_search_accounts_by_email", func(t *testing.T) {
				resp, err := cl.AccountListWithResponse(root, &openapi.AccountListParams{
					Q: &plain.Email,
				}, adminSession)
				tests.Ok(t, err, resp)
				r.Len(resp.JSON200.Accounts, 1)
				a.Equal(plain.Handle, resp.JSON200.Accounts[0].Handle)
				a.Equal(plain.Email, resp.JSON200.Accounts[0].EmailAddresses[0].EmailAddress)
				a.Contains(resp.JSON200.Accounts[0].AuthServices, openapi.AuthProviderIdentifier("password"))
			})

			t.Run("admin_can_filter_administrators", func(t *testing.T) {
				isAdmin := true
				resp, err := cl.AccountListWithResponse(root, &openapi.AccountListParams{
					Admin: &isAdmin,
				}, adminSession)
				tests.Ok(t, err, resp)
				r.NotEmpty(resp.JSON200.Accounts)
				a.Contains(accountHandles(resp.JSON200.Accounts), admin.Handle)
				for _, acc := range resp.JSON200.Accounts {
					a.True(acc.Admin)
				}
			})

			t.Run("admin_can_filter_non_administrators", func(t *testing.T) {
				isAdmin := false
				resp, err := cl.AccountListWithResponse(root, &openapi.AccountListParams{
					Admin: &isAdmin,
				}, adminSession)
				tests.Ok(t, err, resp)
				r.NotEmpty(resp.JSON200.Accounts)
				a.Contains(accountHandles(resp.JSON200.Accounts), plain.Handle)
				a.NotContains(accountHandles(resp.JSON200.Accounts), admin.Handle)
				for _, acc := range resp.JSON200.Accounts {
					a.False(acc.Admin)
				}
			})

			t.Run("admin_can_filter_suspended_accounts", func(t *testing.T) {
				suspendedOnly := true
				resp, err := cl.AccountListWithResponse(root, &openapi.AccountListParams{
					Suspended: &suspendedOnly,
				}, adminSession)
				tests.Ok(t, err, resp)
				r.NotEmpty(resp.JSON200.Accounts)
				a.Contains(accountHandles(resp.JSON200.Accounts), suspended.Handle)
				for _, acc := range resp.JSON200.Accounts {
					a.NotNil(acc.Suspended)
				}
			})

			t.Run("admin_can_filter_active_accounts", func(t *testing.T) {
				activeOnly := false
				resp, err := cl.AccountListWithResponse(root, &openapi.AccountListParams{
					Suspended: &activeOnly,
				}, adminSession)
				tests.Ok(t, err, resp)
				r.NotEmpty(resp.JSON200.Accounts)
				a.Contains(accountHandles(resp.JSON200.Accounts), plain.Handle)
				a.NotContains(accountHandles(resp.JSON200.Accounts), suspended.Handle)
				for _, acc := range resp.JSON200.Accounts {
					a.Nil(acc.Suspended)
				}
			})

			t.Run("admin_can_filter_by_auth_service", func(t *testing.T) {
				services := openapi.AdminAccountsAuthServicesQuery{
					openapi.AuthProviderIdentifier("password"),
				}
				resp, err := cl.AccountListWithResponse(root, &openapi.AccountListParams{
					AuthService: &services,
				}, adminSession)
				tests.Ok(t, err, resp)
				r.NotEmpty(resp.JSON200.Accounts)
				a.Contains(accountHandles(resp.JSON200.Accounts), plain.Handle)
				for _, acc := range resp.JSON200.Accounts {
					a.Contains(acc.AuthServices, openapi.AuthProviderIdentifier("password"))
				}
			})

			t.Run("admin_can_filter_by_roles", func(t *testing.T) {
				roles := openapi.ProfilesRoleFilterQuery{roleID}
				resp, err := cl.AccountListWithResponse(root, &openapi.AccountListParams{
					Roles: &roles,
				}, adminSession)
				tests.Ok(t, err, resp)
				r.Len(resp.JSON200.Accounts, 1)
				a.Equal(invited.Handle, resp.JSON200.Accounts[0].Handle)
			})

			t.Run("admin_can_filter_by_invited_by", func(t *testing.T) {
				invitedBy := openapi.ProfilesInvitedByQuery{admin.Handle}
				resp, err := cl.AccountListWithResponse(root, &openapi.AccountListParams{
					InvitedBy: &invitedBy,
				}, adminSession)
				tests.Ok(t, err, resp)
				r.NotEmpty(resp.JSON200.Accounts)
				a.Contains(accountHandles(resp.JSON200.Accounts), invited.Handle)
				for _, acc := range resp.JSON200.Accounts {
					r.NotNil(acc.InvitedBy)
					a.Equal(admin.Handle, acc.InvitedBy.Handle)
				}
			})

			t.Run("admin_can_filter_by_joined_range", func(t *testing.T) {
				future := openapi.ProfilesJoinRangeQuery(time.Now().UTC().Add(24 * time.Hour).Format(time.DateOnly))
				resp, err := cl.AccountListWithResponse(root, &openapi.AccountListParams{
					Joined: &future,
				}, adminSession)
				tests.Ok(t, err, resp)
				a.Len(resp.JSON200.Accounts, 0)
			})

			t.Run("admin_can_sort_results", func(t *testing.T) {
				q := openapi.SearchQuery(group)
				sortAsc := openapi.ProfilesSortByQuery("handle")
				respAsc, err := cl.AccountListWithResponse(root, &openapi.AccountListParams{
					Q:    &q,
					Sort: &sortAsc,
				}, adminSession)
				tests.Ok(t, err, respAsc)
				r.GreaterOrEqual(len(respAsc.JSON200.Accounts), 3)

				sortDesc := openapi.ProfilesSortByQuery("-handle")
				respDesc, err := cl.AccountListWithResponse(root, &openapi.AccountListParams{
					Q:    &q,
					Sort: &sortDesc,
				}, adminSession)
				tests.Ok(t, err, respDesc)
				r.GreaterOrEqual(len(respDesc.JSON200.Accounts), 3)

				firstAsc := string(respAsc.JSON200.Accounts[0].Handle)
				firstDesc := string(respDesc.JSON200.Accounts[0].Handle)
				a.NotEqual(firstAsc, firstDesc)
			})

			t.Run("admin_can_apply_combined_filters", func(t *testing.T) {
				adminOnly := false
				activeOnly := false
				roles := openapi.ProfilesRoleFilterQuery{roleID}
				invitedBy := openapi.ProfilesInvitedByQuery{admin.Handle}
				services := openapi.AdminAccountsAuthServicesQuery{
					openapi.AuthProviderIdentifier("password"),
				}
				q := openapi.SearchQuery(group)
				resp, err := cl.AccountListWithResponse(root, &openapi.AccountListParams{
					Q:           &q,
					Admin:       &adminOnly,
					Suspended:   &activeOnly,
					Roles:       &roles,
					InvitedBy:   &invitedBy,
					AuthService: &services,
				}, adminSession)
				tests.Ok(t, err, resp)
				r.Len(resp.JSON200.Accounts, 1)
				a.Equal(invited.Handle, resp.JSON200.Accounts[0].Handle)
			})

			t.Run("view_accounts_can_see_admins_in_list", func(t *testing.T) {
				resp, err := cl.AccountListWithResponse(root, nil, supportSession)
				tests.Ok(t, err, resp)
				a.Contains(accountHandles(resp.JSON200.Accounts), admin.Handle)
			})
		}))
	}))
}

type createdAccount struct {
	Id     openapi.Identifier
	Handle openapi.AccountHandle
	Email  string
}

func signupMember(
	t *testing.T,
	root context.Context,
	cl *openapi.ClientWithResponses,
	adminSession openapi.RequestEditorFn,
	emailGroup string,
	invitationID *openapi.InvitationIDQueryParam,
) createdAccount {
	t.Helper()

	email := fmt.Sprintf("%s-%s@example.com", emailGroup, xid.New().String())
	params := &openapi.AuthEmailPasswordSignupParams{
		InvitationId: invitationID,
	}

	signup := tests.AssertRequest(cl.AuthEmailPasswordSignupWithResponse(
		root,
		params,
		openapi.AuthEmailPasswordSignupJSONRequestBody{
			Email:    email,
			Password: "password",
		},
	))(t, http.StatusOK)

	view := tests.AssertRequest(cl.AccountViewWithResponse(
		root,
		signup.JSON200.Id,
		adminSession,
	))(t, http.StatusOK)

	return createdAccount{
		Id:     signup.JSON200.Id,
		Handle: view.JSON200.Handle,
		Email:  email,
	}
}

func accountHandles(in []openapi.Account) []openapi.AccountHandle {
	out := make([]openapi.AccountHandle, 0, len(in))
	for _, acc := range in {
		out = append(out, acc.Handle)
	}
	return out
}

func grant(
	t *testing.T,
	cl *openapi.ClientWithResponses,
	adminSession openapi.RequestEditorFn,
	targetHandle openapi.AccountHandle,
	permissions openapi.PermissionList,
) {
	t.Helper()

	name := "role-view-accounts-" + xid.New().String()
	colour := "blue"

	role := tests.AssertRequest(
		cl.RoleCreateWithResponse(
			t.Context(),
			openapi.RoleCreateJSONRequestBody{
				Name:        name,
				Colour:      colour,
				Permissions: permissions,
			}, adminSession),
	)(t, http.StatusOK)

	tests.AssertRequest(
		cl.AccountAddRoleWithResponse(
			t.Context(),
			targetHandle,
			role.JSON200.Id,
			adminSession),
	)(t, http.StatusOK)
}

func createRole(
	t *testing.T,
	root context.Context,
	cl *openapi.ClientWithResponses,
	adminSession openapi.RequestEditorFn,
	name string,
	permissions openapi.PermissionList,
) openapi.Identifier {
	t.Helper()

	colour := "blue"
	role := tests.AssertRequest(
		cl.RoleCreateWithResponse(
			root,
			openapi.RoleCreateJSONRequestBody{
				Name:        name,
				Colour:      colour,
				Permissions: permissions,
			}, adminSession),
	)(t, http.StatusOK)

	return role.JSON200.Id
}
