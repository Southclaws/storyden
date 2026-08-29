package tools_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/robot/robot_memory"
	robottools "github.com/Southclaws/storyden/app/services/semdex/robot/tools"
	"github.com/Southclaws/storyden/internal/integration"
	"github.com/Southclaws/storyden/lib/mcp"
)

func TestMemoryToolsRegisterSingleNodeSurfaceAndScopeByRobot(t *testing.T) {
	t.Parallel()

	integration.Test(t, nil, fx.Invoke(func(
		lc fx.Lifecycle,
		ctx context.Context,
		_ *robot_memory.Repository,
		registry *robottools.Registry,
	) {
		lc.Append(fx.StartHook(func() {
			for _, name := range []string{"memory_create", "memory_list", "memory_move", "memory_open", "memory_search", "memory_update"} {
				assert.Contains(t, registry.AllToolIDs(ctx), name)
			}
			assert.False(t, registry.HasTool("memory_fact_search"))

			createTool, err := registry.GetTool(ctx, "memory_create")
			require.NoError(t, err)
			searchTool, err := registry.GetTool(ctx, "memory_search")
			require.NoError(t, err)
			updateTool, err := registry.GetTool(ctx, "memory_update")
			require.NoError(t, err)
			assert.Contains(t, searchTool.Definition.Description, "active knowledge graph")
			assert.Contains(t, searchTool.Definition.Description, "Every supplied filter is ANDed")
			assert.Contains(t, createTool.Definition.Description, "Duplicates are acceptable")

			accountA := robottools.ContextWithRunContext(ctx, robottools.RunContext{RobotRef: "shared-robot", AccountID: "account-a"})
			created := callMemoryTool[mcp.ToolMemoryCreateOutput](t, createTool, accountA, mcp.ToolMemoryCreateInput{
				Content: "Freyja is owned by Southclaws.", Subject: stringPointerForTool("Freyja"), Predicate: stringPointerForTool("owned by"), Object: stringPointerForTool("Southclaws"),
			})
			assert.Equal(t, "freyja", *created.Subject)
			assert.Equal(t, "owned_by", *created.Predicate)
			assert.Equal(t, "Memory saved.", created.Message)
			assert.Equal(t, "Continue the current task; organize memories only when necessary.", created.NextAction)

			partialArguments, err := json.Marshal(mcp.ToolMemoryCreateInput{Content: "Invalid partial triple.", Subject: stringPointerForTool("freyja")})
			require.NoError(t, err)
			_, err = createTool.Handler(accountA, partialArguments)
			require.ErrorContains(t, err, "must be provided together")

			// Memory is shared across accounts using the same Robot.
			accountB := robottools.ContextWithRunContext(ctx, robottools.RunContext{RobotRef: "shared-robot", AccountID: "account-b"})
			shared := callMemoryTool[mcp.ToolMemorySearchOutput](t, searchTool, accountB, mcp.ToolMemorySearchInput{Subject: stringPointerForTool("*rey*")})
			require.Len(t, shared.Results, 1)
			assert.Equal(t, created.MemoryId, shared.Results[0].MemoryId)

			// A different Robot cannot discover the node.
			otherRobot := robottools.ContextWithRunContext(ctx, robottools.RunContext{RobotRef: "other-robot", AccountID: "account-a"})
			isolated := callMemoryTool[mcp.ToolMemorySearchOutput](t, searchTool, otherRobot, mcp.ToolMemorySearchInput{Subject: stringPointerForTool("*rey*")})
			assert.Empty(t, isolated.Results)

			empty := ""
			updated := callMemoryTool[mcp.ToolMemoryUpdateOutput](t, updateTool, accountB, mcp.ToolMemoryUpdateInput{
				Id: created.MemoryId, Subject: &empty, Predicate: &empty, Object: &empty,
			})
			assert.Nil(t, updated.Memory.Subject)
			assert.Nil(t, updated.Memory.Predicate)
			assert.Nil(t, updated.Memory.Object)
		}))
	}))
}

func callMemoryTool[Output any](t *testing.T, tool *robottools.Tool, ctx context.Context, input any) Output {
	t.Helper()
	arguments, err := json.Marshal(input)
	require.NoError(t, err)
	result, err := tool.Handler(ctx, arguments)
	require.NoError(t, err)
	var output Output
	require.NoError(t, json.Unmarshal(result, &output))
	return output
}

func stringPointerForTool(value string) *string { return &value }
