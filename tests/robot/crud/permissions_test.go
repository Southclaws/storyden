package crud_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/rs/xid"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/account_ref"
	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/account/role/role_assign"
	"github.com/Southclaws/storyden/app/resources/account/role/role_repo"
	"github.com/Southclaws/storyden/app/resources/rbac"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/tests"
)

func grantRobotPerms(t *testing.T, ctx context.Context, roles *role_repo.Repository, assignments *role_assign.Assignment, accountID account.AccountID, perms ...rbac.Permission) {
	t.Helper()
	created, err := roles.Create(ctx, "robot-test-"+uuid.NewString(), "blue", rbac.PermissionList(perms))
	require.NoError(t, err)
	err = assignments.UpdateRoles(ctx, account_ref.ID(accountID), role_assign.Add(created.ID))
	require.NoError(t, err)
}

// TestRobotPermissions covers all permission combinations for robot CRUD endpoints.
//
// Permission matrix:
//
//	                  | Unauthed | Member | USE_ROBOTS | MANAGE_ROBOTS |
//	RobotsList        | 401      | 403    | 200        | 200           |
//	RobotCreate       | 401      | 403    | 403        | 200           |
//	RobotGet          | 401      | 403    | 200        | 200           |
//	RobotUpdate       | 401      | 403    | 403        | 200           |
//	RobotSessionsList | 401      | 403    | 200        | 200           |
func TestRobotPermissions(t *testing.T) {
	t.Parallel()

	integration.Test(t,
		&config.Config{},
		e2e.Setup(),
		fx.Invoke(func(
			lc fx.Lifecycle,
			root context.Context,
			cl *openapi.ClientWithResponses,
			sh *e2e.SessionHelper,
			aw *account_writer.Writer,
			roles *role_repo.Repository,
			assignments *role_assign.Assignment,
		) {
			lc.Append(fx.StartHook(func() {
				// Admin creates a robot used as a target for get/update permission tests.
				adminCtx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
				adminSession := sh.WithSession(adminCtx)

				robot := tests.AssertRequest(cl.RobotCreateWithResponse(root,
					openapi.RobotCreateJSONRequestBody{
						Name:        "perm-test-robot-" + uuid.NewString(),
						Description: "Permissions test robot",
						Playbook:    "Permissions test playbook.",
						Model:       robotModel(testModel),
					},
					adminSession,
				))(t, http.StatusOK)
				robotID := robot.JSON200.Id

				// Plain member: signed in but no robot permissions.
				memberCtx, memberAcc := e2e.WithAccount(root, aw, seed.Account_003_Baldur)
				memberSession := sh.WithSession(memberCtx)
				_ = memberAcc

				// User with only USE_ROBOTS.
				userCtx, userAcc := e2e.WithAccount(root, aw, seed.Account_004_Loki)
				userSession := sh.WithSession(userCtx)
				grantRobotPerms(t, root, roles, assignments, userAcc.ID, rbac.PermissionUseRobots)

				// User with MANAGE_ROBOTS (implies USE_ROBOTS by convention).
				managerCtx, managerAcc := e2e.WithAccount(root, aw, seed.Account_005_Þórr)
				managerSession := sh.WithSession(managerCtx)
				grantRobotPerms(t, root, roles, assignments, managerAcc.ID, rbac.PermissionUseRobots, rbac.PermissionManageRobots)

				t.Run("RobotsList", func(t *testing.T) {
					r := require.New(t)

					// Unauthenticated
					unauthed, err := cl.RobotsListWithResponse(root, &openapi.RobotsListParams{})
					r.NoError(err)
					r.Equal(http.StatusUnauthorized, unauthed.StatusCode())

					// Member without USE_ROBOTS
					forbidden, err := cl.RobotsListWithResponse(root, &openapi.RobotsListParams{}, memberSession)
					r.NoError(err)
					r.Equal(http.StatusForbidden, forbidden.StatusCode())

					// USE_ROBOTS
					ok, err := cl.RobotsListWithResponse(root, &openapi.RobotsListParams{}, userSession)
					r.NoError(err)
					r.Equal(http.StatusOK, ok.StatusCode())

					// MANAGE_ROBOTS
					okMgr, err := cl.RobotsListWithResponse(root, &openapi.RobotsListParams{}, managerSession)
					r.NoError(err)
					r.Equal(http.StatusOK, okMgr.StatusCode())
				})

				t.Run("RobotCreate", func(t *testing.T) {
					body := openapi.RobotCreateJSONRequestBody{
						Name:        "perm-create-" + uuid.NewString(),
						Description: "Permission test create",
						Playbook:    "Playbook.",
						Model:       robotModel(testModel),
					}

					r := require.New(t)

					// Unauthenticated
					unauthed, err := cl.RobotCreateWithResponse(root, body)
					r.NoError(err)
					r.Equal(http.StatusUnauthorized, unauthed.StatusCode())

					// Member without USE_ROBOTS
					forbidden, err := cl.RobotCreateWithResponse(root, body, memberSession)
					r.NoError(err)
					r.Equal(http.StatusForbidden, forbidden.StatusCode())

					// USE_ROBOTS only (not enough — needs MANAGE_ROBOTS)
					forbiddenUser, err := cl.RobotCreateWithResponse(root, body, userSession)
					r.NoError(err)
					r.Equal(http.StatusForbidden, forbiddenUser.StatusCode())

					// MANAGE_ROBOTS
					okMgr, err := cl.RobotCreateWithResponse(root, body, managerSession)
					r.NoError(err)
					r.Equal(http.StatusOK, okMgr.StatusCode())
				})

				t.Run("RobotGet", func(t *testing.T) {
					r := require.New(t)

					// Unauthenticated
					unauthed, err := cl.RobotGetWithResponse(root, robotID)
					r.NoError(err)
					r.Equal(http.StatusUnauthorized, unauthed.StatusCode())

					// Member without USE_ROBOTS
					forbidden, err := cl.RobotGetWithResponse(root, robotID, memberSession)
					r.NoError(err)
					r.Equal(http.StatusForbidden, forbidden.StatusCode())

					// USE_ROBOTS
					ok, err := cl.RobotGetWithResponse(root, robotID, userSession)
					r.NoError(err)
					r.Equal(http.StatusOK, ok.StatusCode())

					// MANAGE_ROBOTS
					okMgr, err := cl.RobotGetWithResponse(root, robotID, managerSession)
					r.NoError(err)
					r.Equal(http.StatusOK, okMgr.StatusCode())
				})

				t.Run("RobotUpdate", func(t *testing.T) {
					body := openapi.RobotUpdateJSONRequestBody{
						Description: strPtr("Updated via perm test"),
					}

					r := require.New(t)

					// Unauthenticated
					unauthed, err := cl.RobotUpdateWithResponse(root, robotID, body)
					r.NoError(err)
					r.Equal(http.StatusUnauthorized, unauthed.StatusCode())

					// Member without USE_ROBOTS
					forbidden, err := cl.RobotUpdateWithResponse(root, robotID, body, memberSession)
					r.NoError(err)
					r.Equal(http.StatusForbidden, forbidden.StatusCode())

					// USE_ROBOTS only (not enough — needs MANAGE_ROBOTS)
					forbiddenUser, err := cl.RobotUpdateWithResponse(root, robotID, body, userSession)
					r.NoError(err)
					r.Equal(http.StatusForbidden, forbiddenUser.StatusCode())

					// MANAGE_ROBOTS
					okMgr, err := cl.RobotUpdateWithResponse(root, robotID, body, managerSession)
					r.NoError(err)
					r.Equal(http.StatusOK, okMgr.StatusCode())
				})

				t.Run("RobotSessionsList", func(t *testing.T) {
					r := require.New(t)

					// Unauthenticated
					unauthed, err := cl.RobotSessionsListWithResponse(root, &openapi.RobotSessionsListParams{})
					r.NoError(err)
					r.Equal(http.StatusUnauthorized, unauthed.StatusCode())

					// Member without USE_ROBOTS
					forbidden, err := cl.RobotSessionsListWithResponse(root, &openapi.RobotSessionsListParams{}, memberSession)
					r.NoError(err)
					r.Equal(http.StatusForbidden, forbidden.StatusCode())

					// USE_ROBOTS
					ok, err := cl.RobotSessionsListWithResponse(root, &openapi.RobotSessionsListParams{}, userSession)
					r.NoError(err)
					r.Equal(http.StatusOK, ok.StatusCode())

					// MANAGE_ROBOTS
					okMgr, err := cl.RobotSessionsListWithResponse(root, &openapi.RobotSessionsListParams{}, managerSession)
					r.NoError(err)
					r.Equal(http.StatusOK, okMgr.StatusCode())
				})

				t.Run("RobotSessionGet_nonexistent", func(t *testing.T) {
					r := require.New(t)

					// Use a valid xid-format ID that doesn't exist in the database.
					missingID := xid.New().String()

					// Unauthenticated — auth gate fires before DB lookup
					unauthed, err := cl.RobotSessionGetWithResponse(root, missingID, &openapi.RobotSessionGetParams{})
					r.NoError(err)
					r.Equal(http.StatusUnauthorized, unauthed.StatusCode())

					// Member without USE_ROBOTS
					forbidden, err := cl.RobotSessionGetWithResponse(root, missingID, &openapi.RobotSessionGetParams{}, memberSession)
					r.NoError(err)
					r.Equal(http.StatusForbidden, forbidden.StatusCode())

					// USE_ROBOTS — session doesn't exist so expect 404
					notFound, err := cl.RobotSessionGetWithResponse(root, missingID, &openapi.RobotSessionGetParams{}, userSession)
					r.NoError(err)
					r.Equal(http.StatusNotFound, notFound.StatusCode())
				})
			}))
		}),
	)
}

// TestRobotsDisabled verifies behaviour when robots are disabled in settings.
// Permissions still gate access: only accounts with USE_ROBOTS/MANAGE_ROBOTS
// can interact with robots regardless of the feature flag. Members without
// robot permissions always receive 403 — the flag does not change this.
func TestRobotsDisabled(t *testing.T) {
	t.Parallel()

	integration.Test(t,
		&config.Config{},
		e2e.Setup(),
		fx.Invoke(func(
			lc fx.Lifecycle,
			root context.Context,
			cl *openapi.ClientWithResponses,
			sh *e2e.SessionHelper,
			aw *account_writer.Writer,
			roles *role_repo.Repository,
			assignments *role_assign.Assignment,
		) {
			lc.Append(fx.StartHook(func() {
				adminCtx, _ := e2e.WithAccount(root, aw, seed.Account_001_Odin)
				adminSession := sh.WithSession(adminCtx)

				// Plain member — no robot permissions.
				memberCtx, _ := e2e.WithAccount(root, aw, seed.Account_003_Baldur)
				memberSession := sh.WithSession(memberCtx)

				t.Run("admin_can_list_when_disabled", func(t *testing.T) {
					r := require.New(t)
					list, err := cl.RobotsListWithResponse(root, &openapi.RobotsListParams{}, adminSession)
					r.NoError(err)
					r.Equal(http.StatusOK, list.StatusCode())
				})

				t.Run("admin_can_create_when_disabled_with_explicit_model", func(t *testing.T) {
					r := require.New(t)
					create, err := cl.RobotCreateWithResponse(root,
						openapi.RobotCreateJSONRequestBody{
							Name:        "disabled-bot-" + uuid.NewString(),
							Description: "Created while robots disabled",
							Playbook:    "Playbook.",
							Model:       robotModel(testModel),
						},
						adminSession,
					)
					r.NoError(err)
					r.Equal(http.StatusOK, create.StatusCode())
				})

				t.Run("member_cannot_list_when_disabled", func(t *testing.T) {
					r := require.New(t)
					list, err := cl.RobotsListWithResponse(root, &openapi.RobotsListParams{}, memberSession)
					r.NoError(err)
					r.Equal(http.StatusForbidden, list.StatusCode())
				})

				t.Run("member_cannot_create_when_disabled", func(t *testing.T) {
					r := require.New(t)
					create, err := cl.RobotCreateWithResponse(root,
						openapi.RobotCreateJSONRequestBody{
							Name:        "disabled-member-bot-" + uuid.NewString(),
							Description: "Should be forbidden",
							Playbook:    "Playbook.",
							Model:       robotModel(testModel),
						},
						memberSession,
					)
					r.NoError(err)
					r.Equal(http.StatusForbidden, create.StatusCode())
				})

				t.Run("sessions_list_empty_when_disabled", func(t *testing.T) {
					if tests.IsSharedPostgresDatabase() {
						t.Skip("skipping empty robot sessions assertion on shared postgres database")
					}

					r := require.New(t)
					list, err := cl.RobotSessionsListWithResponse(root, &openapi.RobotSessionsListParams{}, adminSession)
					r.NoError(err)
					r.Equal(http.StatusOK, list.StatusCode())
					r.Zero(list.JSON200.Results)
				})
			}))
		}),
	)
}
