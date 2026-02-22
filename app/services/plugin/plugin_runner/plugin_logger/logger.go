package plugin_logger

import (
	"path/filepath"

	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/plugin"
)

func Build() fx.Option {
	return fx.Provide(newWriter, newReader)
}

func getPluginLogDirectory(pluginDataPath string, pluginID plugin.InstallationID) string {
	return filepath.Join(pluginDataPath, pluginID.String(), "logs")
}

func getOutputPath(pluginDataPath string, pluginID plugin.InstallationID) string {
	return filepath.Join(getPluginLogDirectory(pluginDataPath, pluginID), "output.log")
}
