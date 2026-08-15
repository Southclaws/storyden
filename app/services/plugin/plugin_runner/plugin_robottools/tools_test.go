package plugin_robottools

import (
	"testing"

	"github.com/Southclaws/opt"
	"github.com/rs/xid"
	"github.com/stretchr/testify/require"

	"github.com/Southclaws/storyden/app/resources/plugin"
	"github.com/Southclaws/storyden/internal/utils"
	"github.com/Southclaws/storyden/lib/plugin/rpc"
)

func TestNewToolsForProvider(t *testing.T) {
	t.Run("builds plugin tool definition", func(t *testing.T) {
		installationID := plugin.InstallationID(xid.New())
		tools, err := NewToolsForProvider(installationID, rpc.RobotToolProviderCapabilityConfig{
			ID:      "provider",
			Type:    "robot.tool_provider",
			Version: "v1",
			Tools: []rpc.RobotToolProviderToolConfig{
				{
					ID:                "lookup",
					Name:              "Lookup",
					Description:       "Looks up a thing.",
					RequiresWorkspace: opt.New(true),
					InputSchema: map[string]any{
						"type":       "object",
						"properties": map[string]any{},
					},
				},
			},
		}, nil)

		require.NoError(t, err)
		require.Len(t, tools, 1)
		require.Equal(t, FullyQualifiedName(installationID, "provider", "lookup"), tools[0].Definition.Name)
		require.Equal(t, "plugin", tools[0].Source)
		require.Equal(t, "Lookup", tools[0].Definition.Title)
		require.Equal(t, "Looks up a thing.", tools[0].Definition.Description)
		require.Equal(t, "object", tools[0].Definition.InputSchema.Type)
		require.True(t, tools[0].Definition.RequiresWorkspace)
	})

	t.Run("rejects non-object input schema", func(t *testing.T) {
		_, err := NewToolsForProvider(plugin.InstallationID(xid.New()), rpc.RobotToolProviderCapabilityConfig{
			ID:      "provider",
			Type:    "robot.tool_provider",
			Version: "v1",
			Tools: []rpc.RobotToolProviderToolConfig{
				{
					ID:          "lookup",
					Name:        "Lookup",
					Description: "Looks up a thing.",
					InputSchema: map[string]any{"type": "string"},
				},
			},
		}, nil)

		require.Error(t, err)
		require.ErrorContains(t, err, "schema type must be object")
	})
}

func TestToolsetDefinitions(t *testing.T) {
	installationID := plugin.InstallationID(utils.Must(xid.FromString("d9sc9v5o2dtjr7a9gee0")))
	definitions, err := ToolsetDefinitions(installationID, rpc.RobotToolProviderCapabilityConfig{
		ID: "search",
		Tools: []rpc.RobotToolProviderToolConfig{
			{ID: "web", Name: "Web search", Description: "Search the web.", InputSchema: map[string]any{"type": "object"}},
			{ID: "fetch", Name: "Fetch page", Description: "Fetch a page.", InputSchema: map[string]any{"type": "object"}},
		},
		Toolsets: []rpc.RobotToolProviderToolsetConfig{{
			ID:                "research",
			Name:              "Web research",
			Description:       "Researches current information.",
			Instruction:       opt.New("Search first, then inspect primary sources."),
			RequiresWorkspace: opt.New(true),
			Tools:             []string{"web", "fetch"},
		}},
	})

	require.NoError(t, err)
	require.Len(t, definitions, 1)
	require.Equal(t, "plugin.d9sc9v5o2dtjr7a9gee0.search.research", definitions[0].ID)
	require.Equal(t, "Search first, then inspect primary sources.", definitions[0].Instruction)
	require.True(t, definitions[0].RequiresWorkspace)
	require.Equal(t, []string{
		"plugin__d9sc9v5o2dtjr7a9gee0__search__web",
		"plugin__d9sc9v5o2dtjr7a9gee0__search__fetch",
	}, definitions[0].ToolNames)
}

func TestToolsetDefinitionsRejectsUnknownTool(t *testing.T) {
	_, err := ToolsetDefinitions(plugin.InstallationID(xid.New()), rpc.RobotToolProviderCapabilityConfig{
		ID:       "search",
		Toolsets: []rpc.RobotToolProviderToolsetConfig{{ID: "research", Name: "Research", Description: "Research", Tools: []string{"missing"}}},
	})
	require.ErrorContains(t, err, `references unknown tool "missing"`)
}
