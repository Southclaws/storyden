package plugin_manager

import (
	"testing"

	"github.com/Southclaws/opt"
	"github.com/stretchr/testify/require"

	"github.com/Southclaws/storyden/lib/plugin/rpc"
)

func TestValidateManifestConfiguration(t *testing.T) {
	t.Parallel()

	manifest := rpc.Manifest{
		ConfigurationSchema: opt.New(rpc.ManifestConfigurationSchema{
			Fields: []rpc.PluginConfigurationFieldSchema{
				{
					PluginConfigurationFieldUnion: &rpc.PluginConfigurationFieldString{
						Type: "string",
						ID:   "name",
					},
				},
				{
					PluginConfigurationFieldUnion: &rpc.PluginConfigurationFieldBoolean{
						Type: "boolean",
						ID:   "enabled",
					},
				},
				{
					PluginConfigurationFieldUnion: &rpc.PluginConfigurationFieldNumber{
						Type: "number",
						ID:   "threshold",
					},
				},
			},
		}),
	}

	t.Run("valid configuration", func(t *testing.T) {
		err := validateManifestConfiguration(manifest, map[string]any{
			"name":      "plugin",
			"enabled":   true,
			"threshold": 2.5,
			"unknown":   "allowed",
		})

		require.NoError(t, err)
	})

	t.Run("invalid value type", func(t *testing.T) {
		err := validateManifestConfiguration(manifest, map[string]any{
			"name":      "plugin",
			"enabled":   "true",
			"threshold": 2.5,
		})

		require.Error(t, err)
		require.ErrorContains(t, err, `configuration field "enabled" expects boolean`)
	})

	t.Run("number accepts integer", func(t *testing.T) {
		err := validateManifestConfiguration(manifest, map[string]any{
			"name":      "plugin",
			"enabled":   true,
			"threshold": 2,
		})

		require.NoError(t, err)
	})
}
