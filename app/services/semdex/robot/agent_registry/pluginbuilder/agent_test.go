package pluginbuilder

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildToolsetExposesPluginTools(t *testing.T) {
	builder := &Agent{}
	toolset, err := builder.buildToolset()
	require.NoError(t, err)

	static, ok := toolset.(*staticToolset)
	require.True(t, ok)

	names := map[string]bool{}
	for _, tool := range static.tools {
		names[tool.Name()] = true
	}

	for _, name := range []string{
		"plugin_workspace_create",
		"plugin_workspace_info",
		"plugin_file_list",
		"plugin_file_read",
		"plugin_file_write",
		"plugin_file_search",
		"plugin_apply_patch",
		"plugin_go_ast",
		"plugin_go_fmt",
		"plugin_go_vet",
		"plugin_go_tidy",
		"plugin_go_test",
		"plugin_storyden_sdk_events",
		"plugin_storyden_sdk_search",
		"plugin_go_packages",
		"plugin_go_package_symbols",
		"plugin_go_symbol_detail",
		"plugin_go_symbol_search",
		"plugin_storyden_docs",
		"plugin_sdk_reference",
		"plugin_package",
		"plugin_install",
		"plugin_logs_read",
	} {
		assert.True(t, names[name], "missing tool %s", name)
	}

	require.Len(t, static.tools, len(names), "tool names should be unique")
}
