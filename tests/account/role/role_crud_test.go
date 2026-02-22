package role_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/rs/xid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
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

func TestRoleCRUD(t *testing.T) {
	t.Parallel()

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

			t.Run("role_edit", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				name := "test-role-" + xid.New().String()
				colour := "red"
				permissions := openapi.PermissionList{openapi.MANAGECATEGORIES}

				role := tests.AssertRequest(cl.RoleCreateWithResponse(adminCtx, openapi.RoleCreateJSONRequestBody{
					Name:        name,
					Colour:      colour,
					Permissions: permissions,
				}, adminSession))(t, http.StatusOK)

				guestCtx, guest1 := e2e.WithAccount(root, aw, seed.Account_004_Loki)
				guest1Session := sh.WithSession(guestCtx)

				vis := openapi.Published

				tests.AssertRequest(cl.AccountAddRoleWithResponse(adminCtx, guest1.Handle, role.JSON200.Id, adminSession))(t, http.StatusOK)
				tests.AssertRequest(cl.CategoryCreateWithResponse(guestCtx, openapi.CategoryCreateJSONRequestBody{
					Name:        xid.New().String(),
					Description: "d",
					Colour:      "c",
				}, guest1Session))(t, http.StatusOK)
				tests.AssertRequest(cl.NodeCreateWithResponse(guestCtx, openapi.NodeCreateJSONRequestBody{
					Name:       xid.New().String(),
					Visibility: &vis,
				}, guest1Session))(t, http.StatusForbidden)

				newName := "test-role-" + xid.New().String()
				newColour := "red"
				newPermissions := openapi.PermissionList{openapi.MANAGELIBRARY}

				tests.AssertRequest(cl.RoleUpdateWithResponse(adminCtx, role.JSON200.Id, openapi.RoleUpdateJSONRequestBody{
					Name:        &newName,
					Colour:      &newColour,
					Permissions: &newPermissions,
				}, adminSession))(t, http.StatusOK)

				getRole := tests.AssertRequest(cl.RoleGetWithResponse(adminCtx, role.JSON200.Id, adminSession))(t, http.StatusOK)
				a.Equal(newName, getRole.JSON200.Name)
				a.Equal(newColour, getRole.JSON200.Colour)
				a.Equal(openapi.PermissionList{openapi.MANAGELIBRARY}, getRole.JSON200.Permissions)

				tests.AssertRequest(cl.CategoryCreateWithResponse(guestCtx, openapi.CategoryCreateJSONRequestBody{
					Name:        xid.New().String(),
					Description: "d",
					Colour:      "c",
				}, guest1Session))(t, http.StatusForbidden)

				page2 := tests.AssertRequest(cl.NodeCreateWithResponse(guestCtx, openapi.NodeCreateJSONRequestBody{
					Name:       xid.New().String(),
					Visibility: &vis,
				}, guest1Session))(t, http.StatusOK)
				r.NotNil(page2.JSON200)
			})

			t.Run("role_metadata_roundtrip", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				name := "meta-role-" + xid.New().String()
				meta := openapi.Metadata{
					"display": map[string]any{
						"style": "bold",
						"glow":  true,
					},
				}

				create := tests.AssertRequest(cl.RoleCreateWithResponse(adminCtx, openapi.RoleCreateJSONRequestBody{
					Name:        name,
					Colour:      "purple",
					Permissions: openapi.PermissionList{},
					Meta:        &meta,
				}, adminSession))(t, http.StatusOK)

				get := tests.AssertRequest(cl.RoleGetWithResponse(adminCtx, create.JSON200.Id, adminSession))(t, http.StatusOK)

				r.NotNil(get.JSON200.Meta)
				a.Equal(meta, *get.JSON200.Meta)
			})

			t.Run("admin_default_role_rejects_permission_updates", func(t *testing.T) {
				name := "admin-default-role-" + xid.New().String()
				permissions := openapi.PermissionList{openapi.MANAGECATEGORIES}

				tests.AssertRequest(cl.RoleUpdateWithResponse(adminCtx, role.DefaultRoleAdminID.String(), openapi.RoleUpdateJSONRequestBody{
					Name:        &name,
					Permissions: &permissions,
				}, adminSession))(t, http.StatusBadRequest)
			})

			t.Run("role_reorder_custom_roles", func(t *testing.T) {
				r := require.New(t)
				a := assert.New(t)

				createRole := func(name string) openapi.Role {
					resp := tests.AssertRequest(cl.RoleCreateWithResponse(adminCtx, openapi.RoleCreateJSONRequestBody{
						Name:        name,
						Colour:      "orange",
						Permissions: openapi.PermissionList{},
					}, adminSession))(t, http.StatusOK)
					return *resp.JSON200
				}

				role1 := createRole("order-role-" + xid.New().String())
				role2 := createRole("order-role-" + xid.New().String())
				role3 := createRole("order-role-" + xid.New().String())

				nonDefaultCustomIDs := func(roles []openapi.Role) []openapi.Identifier {
					return lo.FilterMap(roles, func(v openapi.Role, _ int) (openapi.Identifier, bool) {
						if v.Id == "0000000000000000000g" || v.Id == "000000000000000000m0" || v.Id == "00000000000000000a00" {
							return "", false
						}
						if v.Id == role1.Id || v.Id == role2.Id || v.Id == role3.Id {
							return "", false
						}
						return openapi.Identifier(v.Id), true
					})
				}

				sameIDSet := func(a, b []openapi.Identifier) bool {
					if len(a) != len(b) {
						return false
					}

					counts := map[openapi.Identifier]int{}
					for _, id := range a {
						counts[id]++
					}
					for _, id := range b {
						counts[id]--
						if counts[id] < 0 {
							return false
						}
					}
					for _, count := range counts {
						if count != 0 {
							return false
						}
					}
					return true
				}

				var remainingIDs []openapi.Identifier
				var reordered bool
				const maxAttempts = 12
				deadline := time.Now().Add(3 * time.Second)

				// This retries around parallel-test interference in shared CI DB:
				// we re-list until remainingIDs/nonDefaultCustomIDs stabilise via
				// sameIDSet, then send RoleUpdateOrderWithResponse with reordered.
				for attempt := 0; attempt < maxAttempts && time.Now().Before(deadline); attempt++ {
					beforeListResp := tests.AssertRequest(cl.RoleListWithResponse(adminCtx, adminSession))(t, http.StatusOK)
					r.NotNil(beforeListResp.JSON200)

					remainingIDs = nonDefaultCustomIDs(beforeListResp.JSON200.Roles)

					stableCheckResp := tests.AssertRequest(cl.RoleListWithResponse(adminCtx, adminSession))(t, http.StatusOK)
					r.NotNil(stableCheckResp.JSON200)

					stableRemainingIDs := nonDefaultCustomIDs(stableCheckResp.JSON200.Roles)
					if !sameIDSet(remainingIDs, stableRemainingIDs) {
						time.Sleep(40 * time.Millisecond)
						continue
					}

					remainingIDs = stableRemainingIDs
					reorderedIDs := append([]openapi.Identifier{role3.Id, role1.Id, role2.Id}, remainingIDs...)
					reorderResp, err := cl.RoleUpdateOrderWithResponse(adminCtx, openapi.RoleUpdateOrderJSONRequestBody{
						RoleIds: reorderedIDs,
					}, adminSession)
					r.NoError(err)
					r.NotNil(reorderResp)

					if reorderResp.StatusCode() == http.StatusOK {
						r.NotNil(reorderResp.JSON200)
						reordered = true
						break
					}

					if reorderResp.StatusCode() != http.StatusBadRequest {
						a.Failf("unexpected role reorder status", "status=%d", reorderResp.StatusCode())
						break
					}

					time.Sleep(40 * time.Millisecond)
				}
				r.True(reordered, "failed to reorder roles with a stable custom role set")

				listResp := tests.AssertRequest(cl.RoleListWithResponse(adminCtx, adminSession))(t, http.StatusOK)
				r.NotNil(listResp.JSON200)

				roles := listResp.JSON200.Roles
				r.True(len(roles) >= 5)
				a.Equal("0000000000000000000g", roles[0].Id, "guest role must remain first")
				a.Equal("000000000000000000m0", roles[1].Id, "member role must remain second")

				ids := lo.Map(roles, func(r openapi.Role, _ int) string { return r.Id })
				idx1 := lo.IndexOf(ids, role1.Id)
				idx2 := lo.IndexOf(ids, role2.Id)
				idx3 := lo.IndexOf(ids, role3.Id)

				r.NotEqual(-1, idx1)
				r.NotEqual(-1, idx2)
				r.NotEqual(-1, idx3)
				a.True(idx3 < idx1 && idx1 < idx2, "custom roles should match reordered precedence")

				tests.AssertRequest(cl.RoleUpdateOrderWithResponse(adminCtx, openapi.RoleUpdateOrderJSONRequestBody{
					RoleIds: append([]openapi.Identifier{
						openapi.Identifier("000000000000000000m0"),
						role1.Id,
						role2.Id,
						role3.Id,
					}, remainingIDs...),
				}, adminSession))(t, http.StatusBadRequest)
			})
		}))
	}))
}
