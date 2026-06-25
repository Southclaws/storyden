package robot

import (
	"context"
	"database/sql"
	"log/slog"
	"testing"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Southclaws/storyden/app/services/semdex/robot/agent_registry"
	"github.com/Southclaws/storyden/app/services/semdex/robot/agent_registry/robotbuilder"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/ent/enttest"

	_ "github.com/glebarez/go-sqlite"
)

func TestResolveAgentSpecResolvesDatabaseRobot(t *testing.T) {
	ctx := context.Background()
	db := newResolutionTestDB(t)

	author := db.Account.Create().
		SetHandle("admin").
		SetName("Admin").
		SaveX(ctx)
	dbRobot := db.Robot.Create().
		SetAuthorID(author.ID).
		SetName("Database Robot").
		SetDescription("A database-backed robot.").
		SetPlaybook("Handle database robot turns.").
		SetModel("mock/robot-chat-simple").
		SetTools([]string{"robot_switch"}).
		SaveX(ctx)

	agent := &Agent{
		logger: slog.Default(),
		db:     db,
	}

	spec, err := agent.resolveAgentSpec(ctx, dbRobot.ID.String())
	require.NoError(t, err)

	assert.Equal(t, dbRobot.ID.String(), spec.RobotRef)
	assert.Equal(t, "Database Robot", spec.AppName)
	assert.Equal(t, "Database Robot", spec.AgentName)
	assert.Equal(t, "Database Robot", spec.DisplayName)
	assert.Equal(t, "A database-backed robot.", spec.Description)
	assert.Equal(t, "Handle database robot turns.", spec.Instruction)
	assert.Equal(t, []string{"robot_switch"}, spec.ToolNames)
	assert.Equal(t, []string{"robot_switch"}, spec.Capabilities)
	require.True(t, spec.DatabaseRobotID.Ok())
	assert.True(t, spec.ModelRef.Ok())
	assert.Nil(t, spec.WorkspaceDefinition)
}

func TestResolveAgentSpecResolvesRobotBuilderBuiltin(t *testing.T) {
	registry := agent_registry.New(slog.Default())
	require.NoError(t, robotbuilder.Register(registry))

	agent := &Agent{
		logger: slog.Default(),
		agents: registry,
	}

	spec, err := agent.resolveAgentSpec(context.Background(), agent_registry.RobotBuilderID)
	require.NoError(t, err)

	assert.Equal(t, agent_registry.RobotBuilderID, spec.RobotRef)
	assert.Equal(t, robotbuilder.AppName, spec.AppName)
	assert.Equal(t, robotbuilder.AgentName, spec.AgentName)
	assert.Equal(t, robotbuilder.DisplayName, spec.DisplayName)
	assert.False(t, spec.DatabaseRobotID.Ok())
	assert.Nil(t, spec.WorkspaceDefinition)
	assert.NotEmpty(t, spec.ToolNames)
	assert.Empty(t, spec.Toolsets)
}

func TestResolveAgentSpecResolvesWorkspaceRequiredBuiltin(t *testing.T) {
	registry := agent_registry.New(slog.Default())
	require.NoError(t, registry.Register(agent_registry.Definition{
		ID:                agent_registry.PluginBuilderID,
		Name:              "Plugin Builder",
		Description:       "Builds plugins.",
		RequiresWorkspace: true,
		AppName:           "storyden_plugin_builder",
		AgentName:         "storyden_plugin_builder",
		Instruction:       "Build plugins.",
	}))

	agent := &Agent{
		logger: slog.Default(),
		agents: registry,
	}

	spec, err := agent.resolveAgentSpec(context.Background(), agent_registry.PluginBuilderID)
	require.NoError(t, err)

	assert.Equal(t, agent_registry.PluginBuilderID, spec.RobotRef)
	require.NotNil(t, spec.WorkspaceDefinition)
	assert.Equal(t, "Plugin Builder", spec.WorkspaceDefinition.Name)
	assert.Empty(t, spec.ToolNames)
	assert.Empty(t, spec.Toolsets)
}

func TestResolveAgentSpecRejectsUnknownBuiltin(t *testing.T) {
	agent := &Agent{
		logger: slog.Default(),
		agents: agent_registry.New(slog.Default()),
	}

	_, err := agent.resolveAgentSpec(context.Background(), "missing_builtin")
	require.Error(t, err)
	assert.Contains(t, err.Error(), `unknown robot "missing_builtin"`)
}

func newResolutionTestDB(t *testing.T) *ent.Client {
	t.Helper()

	sqlDB, err := sql.Open("sqlite", "file:"+t.Name()+"?mode=memory&cache=shared&_pragma=foreign_keys(1)")
	require.NoError(t, err)
	t.Cleanup(func() { _ = sqlDB.Close() })

	db := enttest.NewClient(t, enttest.WithOptions(ent.Driver(entsql.OpenDB(dialect.SQLite, sqlDB))))
	t.Cleanup(func() { _ = db.Close() })
	return db
}
