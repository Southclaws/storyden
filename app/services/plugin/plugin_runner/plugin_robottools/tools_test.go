package plugin_robottools

import (
	"testing"

	"github.com/rs/xid"
	"github.com/stretchr/testify/require"

	"github.com/Southclaws/storyden/app/resources/plugin"
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
					ID:          "lookup",
					Name:        "Lookup",
					Description: "Looks up a thing.",
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
