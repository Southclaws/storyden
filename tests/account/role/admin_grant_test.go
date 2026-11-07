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

func TestRoleAdministratorGrants(t *testing.T) {
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
			managerCtx, manager := e2e.WithAccount(root, aw, seed.Account_004_Loki)
			managerSession := sh.WithSession(managerCtx)
			targetCtx, target := e2e.WithAccount(root, aw, seed.Account_004_Loki)
			targetSession := sh.WithSession(targetCtx)

			createRole := func(perms openapi.PermissionList) openapi.Role {
				response := tests.AssertRequest(cl.RoleCreateWithResponse(adminCtx, openapi.RoleCreateJSONRequestBody{
					Name:        "grant-test-" + xid.New().String(),
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

			t.Run("manager_cannot_create_admin_role", func(t *testing.T) {
				tests.AssertRequest(cl.RoleCreateWithResponse(managerCtx, openapi.RoleCreateJSONRequestBody{
					Name:        "forbidden-admin-" + xid.New().String(),
					Colour:      "red",
					Permissions: openapi.PermissionList{openapi.MANAGEROLES, openapi.ADMINISTRATOR},
				}, managerSession))(t, http.StatusForbidden)
			})

			t.Run("manager_cannot_add_admin_permission", func(t *testing.T) {
				for _, roleID := range []string{
					managerRole.Id,
					ordinaryRole.Id,
					role.DefaultRoleMemberID.String(),
					role.DefaultRoleGuestID.String(),
				} {
					t.Run(roleID, func(t *testing.T) {
						permissions := openapi.PermissionList{openapi.ADMINISTRATOR}
						tests.AssertRequest(cl.RoleUpdateWithResponse(managerCtx, roleID, openapi.RoleUpdateJSONRequestBody{
							Permissions: &permissions,
						}, managerSession))(t, http.StatusForbidden)

						stored := tests.AssertRequest(cl.RoleGetWithResponse(adminCtx, roleID, adminSession))(t, http.StatusOK)
						require.NotContains(t, stored.JSON200.Permissions, openapi.ADMINISTRATOR)
					})
				}
			})

			t.Run("manager_cannot_assign_admin_roles", func(t *testing.T) {
				for _, roleID := range []string{role.DefaultRoleAdminID.String(), customAdminRole.Id} {
					tests.AssertRequest(cl.AccountAddRoleWithResponse(managerCtx, target.Handle, roleID, managerSession))(t, http.StatusForbidden)
				}

				tests.AssertRequest(cl.AccountAddRoleWithResponse(managerCtx, manager.Handle, customAdminRole.Id, managerSession))(t, http.StatusForbidden)
				tests.AssertRequest(cl.AccountAddRoleWithResponse(managerCtx, manager.Handle, role.DefaultRoleAdminID.String(), managerSession))(t, http.StatusBadRequest)

				for _, accountSession := range []openapi.RequestEditorFn{managerSession, targetSession} {
					stored := tests.AssertRequest(cl.AccountGetWithResponse(root, accountSession))(t, http.StatusOK)
					for _, heldRole := range stored.JSON200.Roles {
						require.NotContains(t, heldRole.Permissions, openapi.ADMINISTRATOR)
					}
				}
			})

			t.Run("manager_can_grant_other_permissions", func(t *testing.T) {
				created := tests.AssertRequest(cl.RoleCreateWithResponse(managerCtx, openapi.RoleCreateJSONRequestBody{
					Name:        "allowed-role-" + xid.New().String(),
					Colour:      "red",
					Permissions: openapi.PermissionList{openapi.MANAGEACCOUNTS},
				}, managerSession))(t, http.StatusOK)

				permissions := openapi.PermissionList{openapi.MANAGELIBRARY, openapi.MANAGEACCOUNTS}
				tests.AssertRequest(cl.RoleUpdateWithResponse(managerCtx, created.JSON200.Id, openapi.RoleUpdateJSONRequestBody{
					Permissions: &permissions,
				}, managerSession))(t, http.StatusOK)

				for _, handle := range []string{manager.Handle, target.Handle} {
					assigned := tests.AssertRequest(cl.AccountAddRoleWithResponse(managerCtx, handle, created.JSON200.Id, managerSession))(t, http.StatusOK)
					heldRole := findRole(assigned.JSON200.Roles, created.JSON200.Id)
					require.NotNil(t, heldRole)
					require.ElementsMatch(t, permissions, heldRole.Permissions)
				}
			})

			t.Run("admin_can_grant_admin_permission", func(t *testing.T) {
				permissions := openapi.PermissionList{openapi.ADMINISTRATOR}
				updated := tests.AssertRequest(cl.RoleUpdateWithResponse(adminCtx, ordinaryRole.Id, openapi.RoleUpdateJSONRequestBody{
					Permissions: &permissions,
				}, adminSession))(t, http.StatusOK)
				require.Equal(t, permissions, updated.JSON200.Permissions)
			})

			t.Run("admin_can_assign_admin_roles", func(t *testing.T) {
				for _, roleID := range []string{customAdminRole.Id, ordinaryRole.Id, role.DefaultRoleAdminID.String()} {
					assigned := tests.AssertRequest(cl.AccountAddRoleWithResponse(adminCtx, target.Handle, roleID, adminSession))(t, http.StatusOK)
					heldRole := findRole(assigned.JSON200.Roles, roleID)
					require.NotNil(t, heldRole)
					require.Contains(t, heldRole.Permissions, openapi.ADMINISTRATOR)
				}
			})
		}))
	}))
}
