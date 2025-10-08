package view_test

import (
	"context"
	"net/http"
	"testing"

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

func TestAccountView(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		cl *openapi.ClientWithResponses,
		sh *e2e.SessionHelper,
		aw *account_writer.Writer,
	) {
		lc.Append(fx.StartHook(func() {
			adminCtx, admin := e2e.WithAccount(root, aw, seed.Account_001_Odin)
			adminSession := sh.WithSession(adminCtx)

			admin2Ctx, admin2 := e2e.WithAccount(root, aw, seed.Account_002_Frigg)
			_ = sh.WithSession(admin2Ctx) // admin2Session not used in tests

			// A support account with VIEW_ACCOUNTS permission
			supportCtx, support := e2e.WithAccount(root, aw, seed.Account_003_Baldur)
			supportSession := sh.WithSession(supportCtx)

			// Grant VIEW_ACCOUNTS permission to support account
			grant(t, cl, adminSession, support.Handle, openapi.PermissionList{openapi.VIEWACCOUNTS})

			// A regular member account
			memberCtx, member := e2e.WithAccount(root, aw, seed.Account_004_Loki)
			memberSession := sh.WithSession(memberCtx)

			t.Run("owner_can_view_own_account", func(t *testing.T) {
				r := require.New(t)

				acc := tests.AssertRequest(
					cl.AccountViewWithResponse(root, member.ID.String(), memberSession),
				)(t, http.StatusOK)

				r.NotNil(acc.JSON200)
				r.Equal(member.Handle, acc.JSON200.Handle)
				r.Equal(member.ID.String(), acc.JSON200.Id)
			})

			t.Run("admin_can_view_another_admin", func(t *testing.T) {
				r := require.New(t)

				// Admin (Odin) viewing another admin (Frigg)
				acc := tests.AssertRequest(
					cl.AccountViewWithResponse(root, admin2.ID.String(), adminSession),
				)(t, http.StatusOK)

				r.NotNil(acc.JSON200)
				r.Equal(admin2.Handle, acc.JSON200.Handle)
				r.Equal(admin2.ID.String(), acc.JSON200.Id)
			})

			t.Run("admin_can_view_regular_member", func(t *testing.T) {
				r := require.New(t)

				// Admin (Odin) viewing regular member (Loki)
				acc := tests.AssertRequest(
					cl.AccountViewWithResponse(root, member.ID.String(), adminSession),
				)(t, http.StatusOK)

				r.NotNil(acc.JSON200)
				r.Equal(member.Handle, acc.JSON200.Handle)
				r.Equal(member.ID.String(), acc.JSON200.Id)
			})

			t.Run("view_accounts_can_view_regular_member", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				// Support (Baldur with VIEW_ACCOUNTS) viewing regular member (Loki)
				acc := tests.AssertRequest(
					cl.AccountViewWithResponse(root, member.ID.String(), supportSession),
				)(t, http.StatusOK)

				r.NotNil(acc.JSON200)
				a.Equal(member.Handle, acc.JSON200.Handle)
				a.Equal(member.ID.String(), acc.JSON200.Id)

				// Verify account details are present (emails, auth methods, etc.)
				r.NotNil(acc.JSON200.EmailAddresses)
			})

			t.Run("view_accounts_cannot_view_admin", func(t *testing.T) {
				// Support (Baldur with VIEW_ACCOUNTS) trying to view admin (Odin)
				tests.AssertRequest(
					cl.AccountViewWithResponse(root, admin.ID.String(), supportSession),
				)(t, http.StatusForbidden)
			})

			t.Run("view_accounts_cannot_view_another_admin", func(t *testing.T) {
				// Support (Baldur with VIEW_ACCOUNTS) trying to view another admin (Frigg)
				tests.AssertRequest(
					cl.AccountViewWithResponse(root, admin2.ID.String(), supportSession),
				)(t, http.StatusForbidden)
			})

			t.Run("regular_member_cannot_view_other_member", func(t *testing.T) {
				// Regular member (Loki) trying to view support (Baldur)
				tests.AssertRequest(
					cl.AccountViewWithResponse(root, support.ID.String(), memberSession),
				)(t, http.StatusForbidden)
			})

			t.Run("regular_member_cannot_view_admin", func(t *testing.T) {
				// Regular member (Loki) trying to view admin (Odin)
				tests.AssertRequest(
					cl.AccountViewWithResponse(root, admin.ID.String(), memberSession),
				)(t, http.StatusForbidden)
			})

			t.Run("unauthenticated_cannot_view_account", func(t *testing.T) {
				// Unauthenticated request
				tests.AssertRequest(
					cl.AccountViewWithResponse(root, member.ID.String()),
				)(t, http.StatusUnauthorized)
			})

			t.Run("invalid_account_id_returns_bad_request", func(t *testing.T) {
				// Invalid account ID format
				tests.AssertRequest(
					cl.AccountViewWithResponse(root, "invalid-id", adminSession),
				)(t, http.StatusBadRequest)
			})

			t.Run("nonexistent_account_returns_not_found", func(t *testing.T) {
				// Valid ID format but doesn't exist
				nonexistentID := xid.New().String()
				tests.AssertRequest(
					cl.AccountViewWithResponse(root, nonexistentID, adminSession),
				)(t, http.StatusNotFound)
			})
		}))
	}))
}

// grant creates a role with the specified permissions and assigns it to the target account
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
