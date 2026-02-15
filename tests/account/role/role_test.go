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
				role1 := findRole(guest1Get.JSON200.Roles, "000000000000000000m0")
				r.NotNil(role1, "default role not found")
				a.Equal("green", role1.Colour)
				a.Equal([]openapi.Permission{
					"CREATE_POST", "READ_PUBLISHED_THREADS", "CREATE_REACTION", "READ_PUBLISHED_LIBRARY", "SUBMIT_LIBRARY_NODE", "UPLOAD_ASSET", "LIST_PROFILES", "READ_PROFILE", "CREATE_COLLECTION", "LIST_COLLECTIONS", "READ_COLLECTION", "COLLECTION_SUBMIT",
				}, role1.Permissions)
				role2 := findRole(guest1Get.JSON200.Roles, role.JSON200.Id)
				r.NotNil(role2, "new role not found")
				a.Equal(role.JSON200.Id, role2.Id)
				a.Equal(role.JSON200.Colour, role2.Colour)
				a.Equal(role.JSON200.Permissions, role2.Permissions)

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
				tests.Status(t, err, page, http.StatusForbidden)

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

			t.Run("role_metadata_roundtrip", func(t *testing.T) {
				name := "meta-role-" + xid.New().String()
				meta := openapi.Metadata{
					"display": map[string]any{
						"style": "bold",
						"glow":  true,
					},
				}

				create, err := cl.RoleCreateWithResponse(adminCtx, openapi.RoleCreateJSONRequestBody{
					Name:        name,
					Colour:      "purple",
					Permissions: openapi.PermissionList{},
					Meta:        &meta,
				}, adminSession)
				tests.Ok(t, err, create)

				get, err := cl.RoleGetWithResponse(adminCtx, create.JSON200.Id, adminSession)
				tests.Ok(t, err, get)

				r.NotNil(get.JSON200.Meta)
				a.Equal(meta, *get.JSON200.Meta)
			})

			t.Run("role_reorder_custom_roles", func(t *testing.T) {
				defaultGuestID := openapi.Identifier("0000000000000000000g")
				defaultMemberID := openapi.Identifier("000000000000000000m0")
				defaultAdminID := openapi.Identifier("00000000000000000a00")

				createRole := func(name string) openapi.Role {
					resp := tests.AssertRequest(cl.RoleCreateWithResponse(adminCtx, openapi.RoleCreateJSONRequestBody{
						Name:        name,
						Colour:      "orange",
						Permissions: openapi.PermissionList{},
					}, adminSession))(t, http.StatusOK)
					r.NotNil(resp.JSON200)
					return *resp.JSON200
				}

				role1 := createRole("order-role-" + xid.New().String())
				role2 := createRole("order-role-" + xid.New().String())
				role3 := createRole("order-role-" + xid.New().String())

				getCustomIDs := func(roles []openapi.Role) []openapi.Identifier {
					return lo.FilterMap(roles, func(v openapi.Role, _ int) (openapi.Identifier, bool) {
						if v.Id == defaultGuestID || v.Id == defaultMemberID || v.Id == defaultAdminID {
							return "", false
						}
						return openapi.Identifier(v.Id), true
					})
				}
				targetIDs := []openapi.Identifier{role3.Id, role1.Id, role2.Id}

				var reorderResp *openapi.RoleUpdateOrderResponse
				var reorderedIDs []openapi.Identifier
				for attempt := 0; attempt < 5; attempt++ {
					beforeListResp := tests.AssertRequest(cl.RoleListWithResponse(adminCtx, adminSession))(t, http.StatusOK)
					r.NotNil(beforeListResp.JSON200)

					customIDs := getCustomIDs(beforeListResp.JSON200.Roles)
					r.Contains(customIDs, role1.Id)
					r.Contains(customIDs, role2.Id)
					r.Contains(customIDs, role3.Id)

					remainingIDs := lo.Filter(customIDs, func(id openapi.Identifier, _ int) bool {
						return id != role1.Id && id != role2.Id && id != role3.Id
					})

					reorderedIDs = append([]openapi.Identifier{}, targetIDs...)
					reorderedIDs = append(reorderedIDs, remainingIDs...)

					resp, err := cl.RoleUpdateOrderWithResponse(adminCtx, openapi.RoleUpdateOrderJSONRequestBody{
						RoleIds: reorderedIDs,
					}, adminSession)
					r.NoError(err)
					r.NotNil(resp)

					if resp.StatusCode() == http.StatusOK {
						reorderResp = resp
						break
					}
					if resp.StatusCode() != http.StatusBadRequest {
						t.Fatalf("unexpected role reorder status %d", resp.StatusCode())
					}
				}

				r.NotNil(reorderResp, "failed to reorder roles after retries due concurrent role mutations")
				r.NotNil(reorderResp.JSON200)

				roles := reorderResp.JSON200.Roles
				ids := lo.Map(roles, func(r openapi.Role, _ int) openapi.Identifier { return r.Id })
				idxGuest := lo.IndexOf(ids, defaultGuestID)
				idxMember := lo.IndexOf(ids, defaultMemberID)
				idx1 := lo.IndexOf(ids, role1.Id)
				idx2 := lo.IndexOf(ids, role2.Id)
				idx3 := lo.IndexOf(ids, role3.Id)

				r.NotEqual(-1, idxGuest)
				r.NotEqual(-1, idxMember)
				r.NotEqual(-1, idx1)
				r.NotEqual(-1, idx2)
				r.NotEqual(-1, idx3)
				a.True(idxGuest < idxMember, "default guest role should sort before member role")
				a.True(idx3 < idx1 && idx1 < idx2, "custom roles should match reordered precedence")

				tests.AssertRequest(cl.RoleUpdateOrderWithResponse(adminCtx, openapi.RoleUpdateOrderJSONRequestBody{
					RoleIds: append([]openapi.Identifier{
						defaultMemberID,
					}, reorderedIDs...),
				}, adminSession))(t, http.StatusBadRequest)
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
				tests.Status(t, err, guest1Set, http.StatusForbidden)

				guest1Remove, err := cl.AccountRoleRemoveBadgeWithResponse(guestCtx, admin.Handle, role1ID, guest1Session)
				tests.Status(t, err, guest1Remove, http.StatusForbidden)
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

func findRole(roles []openapi.AccountRole, id string) *openapi.AccountRole {
	for _, role := range roles {
		if role.Id == id {
			return &role
		}
	}
	return nil
}
