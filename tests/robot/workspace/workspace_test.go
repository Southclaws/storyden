package workspace_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/Southclaws/opt"
	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	"google.golang.org/genai"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/rbac"
	robotresource "github.com/Southclaws/storyden/app/resources/robot"
	"github.com/Southclaws/storyden/app/resources/robot/model_ref"
	"github.com/Southclaws/storyden/app/resources/robot/robot_querier"
	"github.com/Southclaws/storyden/app/resources/robot/robot_ref"
	"github.com/Southclaws/storyden/app/resources/robot/robot_session"
	"github.com/Southclaws/storyden/app/resources/robot/robot_workspace"
	"github.com/Southclaws/storyden/app/resources/robot/robot_writer"
	"github.com/Southclaws/storyden/app/resources/seed"
	authsession "github.com/Southclaws/storyden/app/services/authentication/session"
	robotservice "github.com/Southclaws/storyden/app/services/semdex/robot"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/internal/integration/e2e"
	robot_tests "github.com/Southclaws/storyden/tests/robot"
)

const mockModelAck = "mock/../scripts/robot-chat-ack.yaml"

func TestRobotWorkspaceMountsFromServiceRun(t *testing.T) {
	t.Parallel()

	workspaceRoot := t.TempDir()

	integration.Test(t,
		&config.Config{
			LanguageModelProvider:  "mock",
			RobotWorkspaceDataPath: workspaceRoot,
		},
		e2e.Setup(),
		robot_tests.WithRobotSettings(mockModelAck),
		fx.Invoke(func(
			lc fx.Lifecycle,
			root context.Context,
			aw *account_writer.Writer,
			workspaceRepo *robot_workspace.Repository,
			robotWriter *robot_writer.Writer,
			robotQuerier *robot_querier.Querier,
			sessionRepo *robot_session.Repository,
			agent *robotservice.Agent,
		) {
			lc.Append(fx.StartHook(func() {
				_, acc := e2e.WithAccount(root, aw, seed.Account_001_Odin)
				ctx := authsession.WithAccountPermissions(root, *acc, rbac.NewList(rbac.PermissionUseRobots, rbac.PermissionManageRobots))
				modelRef, err := model_ref.ParseID(mockModelAck)
				require.NoError(t, err)

				workspace, err := workspaceRepo.Create(
					ctx,
					"service-run-workspace",
					"Workspace mounted by Agent.Run",
					robotresource.WorkspaceProviderLocal,
					acc.ID,
				)
				require.NoError(t, err)

				bot, err := robotWriter.Create(
					ctx,
					"Workspace Runner",
					"Mounts a default workspace",
					"Say okay.",
					modelRef,
					acc.ID,
					robot_writer.WithWorkspaceID(xid.ID(workspace.ID)),
				)
				require.NoError(t, err)

				sessionID := xid.New()
				runRobot(t, ctx, agent, bot.ID, acc.ID, robotresource.SessionID(sessionID), robotservice.RunOptions{
					Mode:   robotservice.ModeUnattended,
					Source: robotservice.SourcePluginRPC,
				})

				mount := sessionWorkspaceMount(t, root, sessionRepo, robotresource.SessionID(sessionID))
				assert.Equal(t, workspace.ID.String(), mount.WorkspaceID.String())
				assert.Equal(t, robotresource.WorkspaceProviderLocal, mount.Provider)
				assert.Equal(t, filepath.Join(workspaceRoot, mount.WorkspaceInstanceID.String()), mount.ProviderState["root_path"])
				assertDirExists(t, mount.ProviderState["root_path"])

				updatedBot, err := robotWriter.Update(ctx, bot.ID, robot_writer.WithoutWorkspaceID())
				require.NoError(t, err)
				assert.False(t, updatedBot.WorkspaceID.Ok())

				queriedBot, err := robotQuerier.Get(ctx, robot_ref.ID(bot.ID))
				require.NoError(t, err)
				assert.False(t, queriedBot.WorkspaceID.Ok())

				nextSessionID := xid.New()
				runRobot(t, ctx, agent, bot.ID, acc.ID, robotresource.SessionID(nextSessionID), robotservice.RunOptions{
					Mode:   robotservice.ModeUnattended,
					Source: robotservice.SourcePluginRPC,
					Workspace: opt.New(robotservice.WorkspaceMountSpec{
						WorkspaceInstanceID: opt.New(mount.WorkspaceInstanceID),
						Metadata:            map[string]any{},
					}),
				})

				nextMount := sessionWorkspaceMount(t, root, sessionRepo, robotresource.SessionID(nextSessionID))
				assert.Equal(t, mount.WorkspaceInstanceID.String(), nextMount.WorkspaceInstanceID.String())
				assert.Equal(t, mount.ProviderState["root_path"], nextMount.ProviderState["root_path"])
			}))
		}),
	)
}

func runRobot(
	t *testing.T,
	ctx context.Context,
	agent *robotservice.Agent,
	robotID robot_ref.ID,
	accountID account.AccountID,
	sessionID robotresource.SessionID,
	options robotservice.RunOptions,
) {
	t.Helper()

	events := agent.Run(
		ctx,
		xid.ID(robotID).String(),
		accountID.String(),
		sessionID.String(),
		genai.NewContentFromText("hello", genai.RoleUser),
		nil,
		options,
	)
	for _, err := range events {
		require.NoError(t, err)
	}
}

func sessionWorkspaceMount(
	t *testing.T,
	ctx context.Context,
	sessionRepo *robot_session.Repository,
	sessionID robotresource.SessionID,
) robotresource.WorkspaceMount {
	t.Helper()

	session, _, err := sessionRepo.Get(ctx, sessionID, robotresource.NewMessageCursorParams(opt.NewEmpty[robotresource.MessageID](), 1))
	require.NoError(t, err)

	mount, ok := robotservice.WorkspaceMountFromState(session.State).Get()
	require.True(t, ok)
	return mount
}

func assertDirExists(t *testing.T, raw any) {
	t.Helper()

	path, ok := raw.(string)
	require.True(t, ok)
	info, err := os.Stat(path)
	require.NoError(t, err)
	assert.True(t, info.IsDir())
}
