package crud_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/google/uuid"
	nullable "github.com/oapi-codegen/nullable"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/seed"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	"github.com/Southclaws/storyden/tests"
)

func TestRobotWorkspaceCRUDAndRobotDefaultWorkspace(t *testing.T) {
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
		) {
			lc.Append(fx.StartHook(func() {
				adminCtx, adminAcc := e2e.WithAccount(root, aw, seed.Account_001_Odin)
				adminSession := sh.WithSession(adminCtx)

				workspaceName := "workspace-" + uuid.NewString()
				description := "Workspace for Robot CRUD tests"
				provider := openapi.RobotWorkspaceProvider("local")
				createWorkspace := tests.AssertRequest(cl.RobotWorkspaceCreateWithResponse(root,
					openapi.RobotWorkspaceCreateJSONRequestBody{
						Name:        workspaceName,
						Description: description,
						Provider:    &provider,
					},
					adminSession,
				))(t, http.StatusOK)

				assert.Equal(t, workspaceName, createWorkspace.JSON200.Name)
				assert.Equal(t, description, createWorkspace.JSON200.Description)
				assert.Equal(t, provider, createWorkspace.JSON200.Provider)
				assert.Equal(t, adminAcc.ID.String(), createWorkspace.JSON200.CreatedBy.Id)

				getWorkspace := tests.AssertRequest(cl.RobotWorkspaceGetWithResponse(root,
					openapi.RobotWorkspaceIDParam(createWorkspace.JSON200.Id),
					adminSession,
				))(t, http.StatusOK)
				assert.Equal(t, createWorkspace.JSON200.Id, getWorkspace.JSON200.Id)

				createdRobot := tests.AssertRequest(cl.RobotCreateWithResponse(root,
					openapi.RobotCreateJSONRequestBody{
						Name:        "workspace-bot-" + uuid.NewString(),
						Description: "Uses a default workspace",
						Playbook:    "Use your workspace.",
						Model:       robotModel(testModel),
						WorkspaceId: (*openapi.Identifier)(&createWorkspace.JSON200.Id),
					},
					adminSession,
				))(t, http.StatusOK)
				require.True(t, createdRobot.JSON200.WorkspaceId.IsSpecified())
				assert.Equal(t, string(createWorkspace.JSON200.Id), createdRobot.JSON200.WorkspaceId.MustGet())

				getRobot := tests.AssertRequest(cl.RobotGetWithResponse(root,
					createdRobot.JSON200.Id,
					adminSession,
				))(t, http.StatusOK)
				require.True(t, getRobot.JSON200.WorkspaceId.IsSpecified())
				assert.Equal(t, string(createWorkspace.JSON200.Id), getRobot.JSON200.WorkspaceId.MustGet())

				clearedRobot := tests.AssertRequest(cl.RobotUpdateWithResponse(root,
					createdRobot.JSON200.Id,
					openapi.RobotUpdateJSONRequestBody{
						WorkspaceId: nullable.NewNullNullable[openapi.NullableIdentifier](),
					},
					adminSession,
				))(t, http.StatusOK)
				assert.False(t, clearedRobot.JSON200.WorkspaceId.IsSpecified())

				listWorkspaces := tests.AssertRequest(cl.RobotWorkspacesListWithResponse(root,
					&openapi.RobotWorkspacesListParams{},
					adminSession,
				))(t, http.StatusOK)
				assert.GreaterOrEqual(t, listWorkspaces.JSON200.Results, 1)

				tests.AssertRequest(cl.RobotWorkspaceDeleteWithResponse(root,
					openapi.RobotWorkspaceIDParam(createWorkspace.JSON200.Id),
					adminSession,
				))(t, http.StatusOK)
			}))
		}),
	)
}
