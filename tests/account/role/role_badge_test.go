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
			adminCtx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
			adminSession := sh.WithSession(adminCtx)

			setupBadgeScenario := func(t *testing.T) (context.Context, string, openapi.RequestEditorFn, openapi.Identifier, openapi.Identifier) {
				t.Helper()

				targetCtx, target := e2e.WithAccount(root, aw, seed.Account_004_Loki)
				targetSession := sh.WithSession(targetCtx)

				permissions := openapi.PermissionList{}

				role1 := tests.AssertRequest(cl.RoleCreateWithResponse(adminCtx, openapi.RoleCreateJSONRequestBody{
					Name:        "badge-role-" + xid.New().String(),
					Colour:      "green",
					Permissions: permissions,
				}, adminSession))(t, http.StatusOK)
				role1ID := role1.JSON200.Id

				role2 := tests.AssertRequest(cl.RoleCreateWithResponse(adminCtx, openapi.RoleCreateJSONRequestBody{
					Name:        "badge-role-" + xid.New().String(),
					Colour:      "green",
					Permissions: permissions,
				}, adminSession))(t, http.StatusOK)
				role2ID := role2.JSON200.Id

				tests.AssertRequest(cl.AccountAddRoleWithResponse(adminCtx, target.Handle, role1ID, adminSession))(t, http.StatusOK)

				tests.AssertRequest(cl.AccountAddRoleWithResponse(adminCtx, target.Handle, role2ID, adminSession))(t, http.StatusOK)

				return targetCtx, target.Handle, targetSession, role1ID, role2ID
			}

			t.Run("set_own_role_as_badge", func(t *testing.T) {
				t.Parallel()

				r := require.New(t)
				a := assert.New(t)

				targetCtx, targetHandle, targetSession, role1ID, role2ID := setupBadgeScenario(t)

				guest1Set := tests.AssertRequest(cl.AccountRoleSetBadgeWithResponse(targetCtx, targetHandle, role1ID, targetSession))(t, http.StatusOK)
				r.Len(guest1Set.JSON200.Roles, 3, "1 default role, 2 new roles")
				role1, _ := lo.Find(guest1Set.JSON200.Roles, func(r openapi.AccountRole) bool { return r.Id == role1ID })
				a.Equal(role1ID, role1.Id)
				a.True(role1.Badge)

				guest2Set := tests.AssertRequest(cl.AccountRoleSetBadgeWithResponse(targetCtx, targetHandle, role2ID, targetSession))(t, http.StatusOK)
				r.Len(guest2Set.JSON200.Roles, 3, "1 default role, 2 new roles")
				role2, _ := lo.Find(guest2Set.JSON200.Roles, func(r openapi.AccountRole) bool { return r.Id == role2ID })
				a.Equal(role2ID, role2.Id)
				a.True(role2.Badge)

				guest1Remove := tests.AssertRequest(cl.AccountRoleRemoveBadgeWithResponse(targetCtx, targetHandle, role1ID, targetSession))(t, http.StatusOK)
				r.Len(guest1Remove.JSON200.Roles, 3, "1 default role, 2 new roles")
				role1, _ = lo.Find(guest1Remove.JSON200.Roles, func(r openapi.AccountRole) bool { return r.Id == role1ID })
				a.Equal(role1ID, role1.Id)
				a.False(role1.Badge)
			})

			t.Run("non_role_manager_cannot_set_other_member_badges", func(t *testing.T) {
				t.Parallel()

				targetCtx, targetHandle, _, role1ID, _ := setupBadgeScenario(t)
				actorCtx, _ := e2e.WithAccount(root, aw, seed.Account_003_Baldur)
				actorSession := sh.WithSession(actorCtx)

				tests.AssertRequest(cl.AccountRoleSetBadgeWithResponse(targetCtx, targetHandle, role1ID, actorSession))(t, http.StatusForbidden)

				tests.AssertRequest(cl.AccountRoleRemoveBadgeWithResponse(targetCtx, targetHandle, role1ID, actorSession))(t, http.StatusForbidden)
			})

			t.Run("role_manager_can_set_other_member_badges", func(t *testing.T) {
				t.Parallel()

				r := require.New(t)
				a := assert.New(t)

				targetCtx, targetHandle, _, role1ID, _ := setupBadgeScenario(t)

				setResp := tests.AssertRequest(cl.AccountRoleSetBadgeWithResponse(targetCtx, targetHandle, role1ID, adminSession))(t, http.StatusOK)
				r.Equal(targetHandle, setResp.JSON200.Handle, "admin badge update should mutate target account, not caller")
				targetRole, found := lo.Find(setResp.JSON200.Roles, func(in openapi.AccountRole) bool { return in.Id == role1ID })
				r.True(found, "target role should be present on target account")
				a.True(targetRole.Badge, "target role should be marked as badge")

				removeResp := tests.AssertRequest(cl.AccountRoleRemoveBadgeWithResponse(targetCtx, targetHandle, role1ID, adminSession))(t, http.StatusOK)
				r.Equal(targetHandle, removeResp.JSON200.Handle, "admin badge removal should mutate target account, not caller")
				targetRole, found = lo.Find(removeResp.JSON200.Roles, func(in openapi.AccountRole) bool { return in.Id == role1ID })
				r.True(found, "target role should still be present on target account")
				a.False(targetRole.Badge, "target role badge should be removed")
			})
		}))
	}))
}
