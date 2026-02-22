package role_test

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

func TestRoleAssignment(t *testing.T) {
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

			tests.AssertRequest(cl.CategoryCreateWithResponse(guestCtx, openapi.CategoryCreateJSONRequestBody{
				Name:        xid.New().String(),
				Description: "d",
				Colour:      "c",
			}, guest1Session))(t, http.StatusForbidden)

			getRole := tests.AssertRequest(cl.RoleGetWithResponse(adminCtx, role.JSON200.Id, adminSession))(t, http.StatusOK)
			a.Equal(name, getRole.JSON200.Name)
			a.Equal(colour, getRole.JSON200.Colour)
			a.Equal(permissions, getRole.JSON200.Permissions)

			tests.AssertRequest(cl.AccountAddRoleWithResponse(adminCtx, guest1.Handle, role.JSON200.Id, adminSession))(t, http.StatusOK)

			guest1Get := tests.AssertRequest(cl.AccountGetWithResponse(guestCtx, guest1Session))(t, http.StatusOK)

			r.Len(guest1Get.JSON200.Roles, 2, "1 default roles, 1 new role")
			role1 := findRole(guest1Get.JSON200.Roles, "000000000000000000m0")
			r.NotNil(role1, "default role not found")
			a.Equal("green", role1.Colour)
			a.Equal([]openapi.Permission{
				"CREATE_POST",
				"READ_PUBLISHED_THREADS",
				"CREATE_REACTION",
				"READ_PUBLISHED_LIBRARY",
				"SUBMIT_LIBRARY_NODE",
				"UPLOAD_ASSET",
				"LIST_PROFILES",
				"READ_PROFILE",
				"CREATE_COLLECTION",
				"LIST_COLLECTIONS",
				"READ_COLLECTION",
				"COLLECTION_SUBMIT",
			}, role1.Permissions)
			role2 := findRole(guest1Get.JSON200.Roles, role.JSON200.Id)
			r.NotNil(role2, "new role not found")
			a.Equal(role.JSON200.Id, role2.Id)
			a.Equal(role.JSON200.Colour, role2.Colour)
			a.Equal(role.JSON200.Permissions, role2.Permissions)

			tests.AssertRequest(cl.CategoryCreateWithResponse(guestCtx, openapi.CategoryCreateJSONRequestBody{
				Name:        xid.New().String(),
				Description: "d",
				Colour:      "c",
			}, guest1Session))(t, http.StatusOK)

			remove := tests.AssertRequest(cl.AccountRemoveRoleWithResponse(adminCtx, guest1.Handle, role.JSON200.Id, adminSession))(t, http.StatusOK)

			r.Len(remove.JSON200.Roles, 1)

			tests.AssertRequest(cl.CategoryCreateWithResponse(guestCtx, openapi.CategoryCreateJSONRequestBody{
				Name:        xid.New().String(),
				Description: "d",
				Colour:      "c",
			}, guest1Session))(t, http.StatusForbidden)
		}))
	}))
}
