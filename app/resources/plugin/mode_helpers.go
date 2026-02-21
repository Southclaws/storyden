package plugin

func pluginModeFromSupervised(supervised bool) PluginMode {
	if supervised {
		return PluginModeSupervised
	}
	return PluginModeExternal
}

func (m PluginMode) Supervised() bool {
	return m == PluginModeSupervised
}

func (m PluginMode) External() bool {
	return m == PluginModeExternal
}
