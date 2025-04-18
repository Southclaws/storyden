package role_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/rs/xid"
	"github.com/samber/lo"
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

func TestRoles(t *testing.T) {
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

			adminCtx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
			adminSession := sh.WithSession(adminCtx)

			t.Run("role_assignment", func(t *testing.T) {
				t.Parallel()

				name := "test-role-" + xid.New().String()
				colour := "red"
				permissions := openapi.PermissionList{openapi.MANAGECATEGORIES}

				role, err := cl.RoleCreateWithResponse(adminCtx, openapi.RoleCreateJSONRequestBody{Name: name, Colour: colour, Permissions: permissions}, adminSession)
				tests.Ok(t, err, role)

				guestCtx, guest1 := e2e.WithAccount(root, aw, seed.Account_004_Loki)
				guest1Session := sh.WithSession(guestCtx)

				// guest1 cannot create categories
				cat1, err := cl.CategoryCreateWithResponse(guestCtx, openapi.CategoryCreateJSONRequestBody{Name: xid.New().String(), Description: "d", Colour: "c"}, guest1Session)
				tests.Status(t, err, cat1, http.StatusForbidden)

				getRole, err := cl.RoleGetWithResponse(adminCtx, role.JSON200.Id, adminSession)
				tests.Ok(t, err, getRole)
				a.Equal(name, getRole.JSON200.Name)
				a.Equal(colour, getRole.JSON200.Colour)
				a.Equal(permissions, getRole.JSON200.Permissions)

				addRole1, err := cl.AccountAddRoleWithResponse(adminCtx, guest1.Handle, role.JSON200.Id, adminSession)
				tests.Ok(t, err, addRole1)

				guest1Get, err := cl.AccountGetWithResponse(guestCtx, guest1Session)
				tests.Ok(t, err, guest1Get)

				r.Len(guest1Get.JSON200.Roles, 2, "1 default roles, 1 new role")
				role1 := guest1Get.JSON200.Roles[0] // custom roles are always first, default roles always last
				a.Equal(role.JSON200.Id, role1.Id)
				a.Equal(role.JSON200.Colour, role1.Colour)
				a.Equal(role.JSON200.Permissions, role1.Permissions)

				// guest1 can now create categories
				cat, err := cl.CategoryCreateWithResponse(guestCtx, openapi.CategoryCreateJSONRequestBody{Name: xid.New().String(), Description: "d", Colour: "c"}, guest1Session)
				tests.Ok(t, err, cat)

				remove, err := cl.AccountRemoveRoleWithResponse(adminCtx, guest1.Handle, role.JSON200.Id, adminSession)
				tests.Ok(t, err, remove)

				r.Len(remove.JSON200.Roles, 1)

				// guest1 now cannot create categories
				cat2, err := cl.CategoryCreateWithResponse(guestCtx, openapi.CategoryCreateJSONRequestBody{Name: xid.New().String(), Description: "d", Colour: "c"}, guest1Session)
				tests.Status(t, err, cat2, http.StatusForbidden)
			})

			t.Run("role_edit", func(t *testing.T) {
				t.Parallel()

				name := "test-role-" + xid.New().String()
				colour := "red"
				permissions := openapi.PermissionList{openapi.MANAGECATEGORIES}

				role, err := cl.RoleCreateWithResponse(adminCtx, openapi.RoleCreateJSONRequestBody{Name: name, Colour: colour, Permissions: permissions}, adminSession)
				tests.Ok(t, err, role)

				guestCtx, guest1 := e2e.WithAccount(root, aw, seed.Account_004_Loki)
				guest1Session := sh.WithSession(guestCtx)

				vis := openapi.Published

				// Give role to guest, has category permissions, but not library permissions
				addRole1, err := cl.AccountAddRoleWithResponse(adminCtx, guest1.Handle, role.JSON200.Id, adminSession)
				tests.Ok(t, err, addRole1)
				cat, err := cl.CategoryCreateWithResponse(guestCtx, openapi.CategoryCreateJSONRequestBody{Name: xid.New().String(), Description: "d", Colour: "c"}, guest1Session)
				tests.Ok(t, err, cat)
				page, err := cl.NodeCreateWithResponse(guestCtx, openapi.NodeCreateJSONRequestBody{Name: xid.New().String(), Visibility: &vis}, guest1Session)
				tests.Status(t, err, page, http.StatusUnauthorized)

				// Rename, re-colour and remove category permissions
				newName := "test-role-" + xid.New().String()
				newColour := "red"
				newPermissions := openapi.PermissionList{openapi.MANAGELIBRARY}

				edit, err := cl.RoleUpdateWithResponse(adminCtx, role.JSON200.Id, openapi.RoleUpdateJSONRequestBody{Name: &newName, Colour: &newColour, Permissions: &newPermissions}, adminSession)
				tests.Ok(t, err, edit)

				getRole, err := cl.RoleGetWithResponse(adminCtx, role.JSON200.Id, adminSession)
				tests.Ok(t, err, getRole)
				a.Equal(newName, getRole.JSON200.Name)
				a.Equal(newColour, getRole.JSON200.Colour)
				a.Equal(openapi.PermissionList{openapi.MANAGELIBRARY}, getRole.JSON200.Permissions)

				// guest1 now cannot create categories
				cat2, err := cl.CategoryCreateWithResponse(guestCtx, openapi.CategoryCreateJSONRequestBody{Name: xid.New().String(), Description: "d", Colour: "c"}, guest1Session)
				tests.Status(t, err, cat2, http.StatusForbidden)

				// guest1 can now create published library pages
				page2, err := cl.NodeCreateWithResponse(guestCtx, openapi.NodeCreateJSONRequestBody{Name: xid.New().String(), Visibility: &vis}, guest1Session)
				tests.Ok(t, err, page2)
			})
		}))
	}))
}

func TestRoleBadges(t *testing.T) {
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

			guestCtx, guest1 := e2e.WithAccount(root, aw, seed.Account_004_Loki)
			guest1Session := sh.WithSession(guestCtx)

			colour := "green"
			permissions := openapi.PermissionList{ /* purely aesthetic role, no permissions*/ }

			role1, err := cl.RoleCreateWithResponse(adminCtx, openapi.RoleCreateJSONRequestBody{Name: "badge-role-" + xid.New().String(), Colour: colour, Permissions: permissions}, adminSession)
			tests.Ok(t, err, role1)
			role1ID := role1.JSON200.Id

			role2, err := cl.RoleCreateWithResponse(adminCtx, openapi.RoleCreateJSONRequestBody{Name: "badge-role-" + xid.New().String(), Colour: colour, Permissions: permissions}, adminSession)
			tests.Ok(t, err, role2)
			role2ID := role2.JSON200.Id

			addRole1, err := cl.AccountAddRoleWithResponse(adminCtx, guest1.Handle, role1ID, adminSession)
			tests.Ok(t, err, addRole1)

			addRole2, err := cl.AccountAddRoleWithResponse(adminCtx, guest1.Handle, role2ID, adminSession)
			tests.Ok(t, err, addRole2)

			t.Run("set_own_role_as_badge", func(t *testing.T) {
				t.Parallel()

				guest1Set, err := cl.AccountRoleSetBadgeWithResponse(guestCtx, guest1.Handle, role1ID, guest1Session)
				tests.Ok(t, err, guest1Set)
				r.Len(guest1Set.JSON200.Roles, 3, "1 default roles, 1 new role")
				role1, _ := lo.Find(guest1Set.JSON200.Roles, func(r openapi.AccountRole) bool { return r.Id == role1ID })
				a.Equal(role1ID, role1.Id)
				a.True(role1.Badge)

				guest2Set, err := cl.AccountRoleSetBadgeWithResponse(guestCtx, guest1.Handle, role2ID, guest1Session)
				tests.Ok(t, err, guest2Set)
				r.Len(guest2Set.JSON200.Roles, 3, "1 default roles, 1 new role")
				role2, _ := lo.Find(guest2Set.JSON200.Roles, func(r openapi.AccountRole) bool { return r.Id == role2ID })
				a.Equal(role2ID, role2.Id)
				a.True(role2.Badge)

				guest1Remove, err := cl.AccountRoleRemoveBadgeWithResponse(guestCtx, guest1.Handle, role1ID, guest1Session)
				tests.Ok(t, err, guest1Remove)
				r.Len(guest1Remove.JSON200.Roles, 3, "1 default roles, 1 new role")
				role1, _ = lo.Find(guest1Remove.JSON200.Roles, func(r openapi.AccountRole) bool { return r.Id == role1ID })
				a.Equal(role1ID, role1.Id)
				a.False(role1.Badge)
			})

			t.Run("non_role_manager_cannot_set_other_member_badges", func(t *testing.T) {
				t.Parallel()

				guest1Set, err := cl.AccountRoleSetBadgeWithResponse(guestCtx, admin.Handle, role1ID, guest1Session)
				tests.Status(t, err, guest1Set, http.StatusUnauthorized)

				guest1Remove, err := cl.AccountRoleRemoveBadgeWithResponse(guestCtx, admin.Handle, role1ID, guest1Session)
				tests.Status(t, err, guest1Remove, http.StatusUnauthorized)
			})

			t.Run("role_manager_can_set_other_member_badges", func(t *testing.T) {
				t.Parallel()

				guest1Set, err := cl.AccountRoleSetBadgeWithResponse(guestCtx, guest1.Handle, role1ID, adminSession)
				tests.Ok(t, err, guest1Set)

				guest1Remove, err := cl.AccountRoleRemoveBadgeWithResponse(guestCtx, guest1.Handle, role1ID, adminSession)
				tests.Ok(t, err, guest1Remove)
			})
		}))
	}))
}
