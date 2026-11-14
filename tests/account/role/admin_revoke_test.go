package role_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/rs/xid"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/account/role"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/tests"
)

func TestRoleAdministratorRevocation(t *testing.T) {
	integration.Test(t, nil, e2e.Setup(), fx.Invoke(func(
		lc fx.Lifecycle,
		root context.Context,
		cl *openapi.ClientWithResponses,
		sh *e2e.SessionHelper,
		aw *account_writer.Writer,
	) {
		lc.Append(fx.StartHook(func() {
			adminCtx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
			adminSession := sh.WithSession(adminCtx)
			victimCtx, victim := e2e.WithAccount(root, aw, seed.Account_002_Frigg)
			victimSession := sh.WithSession(victimCtx)
			managerCtx, manager := e2e.WithAccount(root, aw, seed.Account_004_Loki)
			managerSession := sh.WithSession(managerCtx)
			holderCtx, holder := e2e.WithAccount(root, aw, seed.Account_003_Baldur)
			holderSession := sh.WithSession(holderCtx)

			createRole := func(perms openapi.PermissionList) openapi.Role {
				response := tests.AssertRequest(cl.RoleCreateWithResponse(adminCtx, openapi.RoleCreateJSONRequestBody{
					Name:        "revoke-test-" + xid.New().String(),
					Colour:      "red",
					Permissions: perms,
				}, adminSession))(t, http.StatusOK)
				require.NotNil(t, response.JSON200)

				return *response.JSON200
			}

			managerRole := createRole(openapi.PermissionList{openapi.MANAGEROLES})
			customAdminRole := createRole(openapi.PermissionList{openapi.ADMINISTRATOR})
			ordinaryRole := createRole(openapi.PermissionList{openapi.MANAGELIBRARY})

			tests.AssertRequest(cl.AccountAddRoleWithResponse(adminCtx, manager.Handle, managerRole.Id, adminSession))(t, http.StatusOK)
			tests.AssertRequest(cl.AccountAddRoleWithResponse(adminCtx, holder.Handle, customAdminRole.Id, adminSession))(t, http.StatusOK)

			requireAdministrator := func(t *testing.T, session openapi.RequestEditorFn) {
				t.Helper()
				tests.AssertRequest(cl.AdminThemeGetWithResponse(root, session))(t, http.StatusOK)
			}

			t.Run("manager_cannot_remove_default_admin_role", func(t *testing.T) {
				requireAdministrator(t, victimSession)

				tests.AssertRequest(cl.AccountRemoveRoleWithResponse(managerCtx, victim.Handle, role.DefaultRoleAdminID.String(), managerSession))(t, http.StatusForbidden)

				stored := tests.AssertRequest(cl.AccountGetWithResponse(root, victimSession))(t, http.StatusOK)
				require.True(t, stored.JSON200.Admin)
				requireAdministrator(t, victimSession)
			})

			t.Run("manager_cannot_remove_custom_admin_role", func(t *testing.T) {
				requireAdministrator(t, holderSession)

				tests.AssertRequest(cl.AccountRemoveRoleWithResponse(managerCtx, holder.Handle, customAdminRole.Id, managerSession))(t, http.StatusForbidden)

				stored := tests.AssertRequest(cl.AccountGetWithResponse(root, holderSession))(t, http.StatusOK)
				require.NotNil(t, findRole(stored.JSON200.Roles, customAdminRole.Id))
				requireAdministrator(t, holderSession)
			})

			t.Run("manager_cannot_strip_admin_permission_from_role", func(t *testing.T) {
				permissions := openapi.PermissionList{openapi.MANAGEROLES}
				tests.AssertRequest(cl.RoleUpdateWithResponse(managerCtx, customAdminRole.Id, openapi.RoleUpdateJSONRequestBody{
					Permissions: &permissions,
				}, managerSession))(t, http.StatusForbidden)

				stored := tests.AssertRequest(cl.RoleGetWithResponse(adminCtx, customAdminRole.Id, adminSession))(t, http.StatusOK)
				require.Contains(t, stored.JSON200.Permissions, openapi.ADMINISTRATOR)
				requireAdministrator(t, holderSession)
			})

			t.Run("manager_cannot_rename_admin_role", func(t *testing.T) {
				for _, roleID := range []string{customAdminRole.Id, role.DefaultRoleAdminID.String()} {
					name := "renamed-" + xid.New().String()
					tests.AssertRequest(cl.RoleUpdateWithResponse(managerCtx, roleID, openapi.RoleUpdateJSONRequestBody{
						Name: &name,
					}, managerSession))(t, http.StatusForbidden)

					stored := tests.AssertRequest(cl.RoleGetWithResponse(adminCtx, roleID, adminSession))(t, http.StatusOK)
					require.NotEqual(t, name, stored.JSON200.Name)
				}
			})

			t.Run("manager_cannot_delete_admin_role", func(t *testing.T) {
				tests.AssertRequest(cl.RoleDeleteWithResponse(managerCtx, customAdminRole.Id, managerSession))(t, http.StatusForbidden)

				tests.AssertRequest(cl.RoleGetWithResponse(adminCtx, customAdminRole.Id, adminSession))(t, http.StatusOK)
				requireAdministrator(t, holderSession)
			})

			t.Run("manager_can_manage_ordinary_roles", func(t *testing.T) {
				tests.AssertRequest(cl.AccountAddRoleWithResponse(managerCtx, victim.Handle, ordinaryRole.Id, managerSession))(t, http.StatusOK)

				removed := tests.AssertRequest(cl.AccountRemoveRoleWithResponse(managerCtx, victim.Handle, ordinaryRole.Id, managerSession))(t, http.StatusOK)
				require.Nil(t, findRole(removed.JSON200.Roles, ordinaryRole.Id))

				name := "renamed-" + xid.New().String()
				tests.AssertRequest(cl.RoleUpdateWithResponse(managerCtx, ordinaryRole.Id, openapi.RoleUpdateJSONRequestBody{
					Name: &name,
				}, managerSession))(t, http.StatusOK)

				tests.AssertRequest(cl.RoleDeleteWithResponse(managerCtx, ordinaryRole.Id, managerSession))(t, http.StatusOK)
			})

			t.Run("admin_can_revoke_admin_roles", func(t *testing.T) {
				removed := tests.AssertRequest(cl.AccountRemoveRoleWithResponse(adminCtx, holder.Handle, customAdminRole.Id, adminSession))(t, http.StatusOK)
				require.Nil(t, findRole(removed.JSON200.Roles, customAdminRole.Id))

				removed = tests.AssertRequest(cl.AccountRemoveRoleWithResponse(adminCtx, victim.Handle, role.DefaultRoleAdminID.String(), adminSession))(t, http.StatusOK)
				require.False(t, removed.JSON200.Admin)

				permissions := openapi.PermissionList{openapi.MANAGELIBRARY}
				tests.AssertRequest(cl.RoleUpdateWithResponse(adminCtx, customAdminRole.Id, openapi.RoleUpdateJSONRequestBody{
					Permissions: &permissions,
				}, adminSession))(t, http.StatusOK)

				tests.AssertRequest(cl.RoleDeleteWithResponse(adminCtx, customAdminRole.Id, adminSession))(t, http.StatusOK)
			})
		}))
	}))
}
