package pluginbuilder

import (
	"context"
	"iter"
	"testing"

	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/adk/session"
	adktool "google.golang.org/adk/tool"
	"google.golang.org/genai"

	"github.com/Southclaws/storyden/app/services/semdex/robot/workspacestate"
)

func TestBuildToolsetContainsPluginTools(t *testing.T) {
	builder := &Agent{}
	toolset, err := builder.buildToolset()
	require.NoError(t, err)

	dynamic, ok := toolset.(*pluginBuilderToolset)
	require.True(t, ok)

	names := map[string]bool{}
	for _, tool := range dynamic.tools {
		names[tool.Name()] = true
	}

	for _, name := range []string{
		"plugin_workspace_create",
		"plugin_workspace_info",
		"plugin_installed_list",
		"plugin_workspace_import_installation",
		"plugin_manifest_write",
		"plugin_file_list",
		"plugin_file_read",
		"plugin_file_outline",
		"plugin_file_write",
		"plugin_file_search",
		"plugin_file_edit",
		"plugin_go_ast",
		"plugin_go_fmt",
		"plugin_go_vet",
		"plugin_go_tidy",
		"plugin_go_test",
		"plugin_validate",
		"plugin_storyden_sdk_events",
		"plugin_storyden_sdk_search",
		"plugin_go_packages",
		"plugin_go_package_symbols",
		"plugin_go_symbol_detail",
		"plugin_go_symbol_search",
		"plugin_storyden_docs",
		"plugin_install",
		"plugin_logs_read",
		"plugin_run_bash",
	} {
		assert.True(t, names[name], "missing tool %s", name)
	}
	assert.False(t, names["plugin_apply_patch"])

	require.Len(t, dynamic.tools, len(names), "tool names should be unique")
}

func TestBuildToolsetTemporarilyExposesAllToolsWhenUnbound(t *testing.T) {
	toolset := newTestPluginBuilderToolset(t)

	tools, err := toolset.Tools(newPluginBuilderReadonlyTestContext(nil))
	require.NoError(t, err)

	assert.Equal(t, pluginBuilderAllToolNames, pluginToolNames(tools))
}

func TestBuildToolsetAddsBashOnlyWhenWorkspaceAllowsUntrustedCommands(t *testing.T) {
	toolset := newTestPluginBuilderToolset(t)

	tools, err := toolset.Tools(newPluginBuilderReadonlyTestContext(nil))
	require.NoError(t, err)
	assert.NotContains(t, pluginToolNames(tools), "plugin_run_bash")

	tools, err = toolset.Tools(newPluginBuilderReadonlyTestContext(map[string]any{
		workspacestate.WorkspaceStateKey: map[string]any{
			"workspace_id":             xid.New().String(),
			"workspace_instance_id":    xid.New().String(),
			"provider":                 "sprites",
			"provider_state":           map[string]any{},
			"allow_untrusted_commands": true,
			"metadata":                 map[string]any{},
		},
	}))
	require.NoError(t, err)
	assert.Contains(t, pluginToolNames(tools), "plugin_run_bash")
}

func TestBuildToolsetTemporarilyExposesAllToolsWhenBoundToNewPlugin(t *testing.T) {
	toolset := newTestPluginBuilderToolset(t)

	tools, err := toolset.Tools(newPluginBuilderReadonlyTestContext(map[string]any{
		pluginBuildTargetStateKey: pluginBuildTarget{
			Mode:       pluginBuildTargetModeNew,
			ManifestID: "welcome-plugin",
		},
	}))
	require.NoError(t, err)

	names := pluginToolNames(tools)
	assert.Equal(t, pluginBuilderAllToolNames, names)
	assert.Contains(t, names, "plugin_workspace_create")
	assert.Contains(t, names, "plugin_installed_list")
	assert.Contains(t, names, "plugin_workspace_import_installation")
	assert.Contains(t, names, "plugin_logs_read")
}

func TestBuildToolsetTemporarilyExposesAllToolsWhenBoundToInstalledPlugin(t *testing.T) {
	toolset := newTestPluginBuilderToolset(t)

	tools, err := toolset.Tools(newPluginBuilderReadonlyTestContext(map[string]any{
		pluginBuildTargetStateKey: pluginBuildTarget{
			Mode:           pluginBuildTargetModeExisting,
			ManifestID:     "welcome-plugin",
			InstallationID: xid.New().String(),
		},
	}))
	require.NoError(t, err)

	names := pluginToolNames(tools)
	assert.Equal(t, pluginBuilderAllToolNames, names)
	assert.Contains(t, names, "plugin_logs_read")
	assert.Contains(t, names, "plugin_workspace_create")
	assert.Contains(t, names, "plugin_installed_list")
	assert.Contains(t, names, "plugin_workspace_import_installation")
}

func TestInstructionIncludesUnboundState(t *testing.T) {
	agent := &Agent{}

	instruction, err := agent.instruction(newPluginBuilderReadonlyTestContext(nil))
	require.NoError(t, err)

	assert.Contains(t, instruction, "Current chat state: no plugin selected yet")
	assert.Contains(t, instruction, "creating a new plugin or importing an installed plugin")
	assert.Contains(t, instruction, "Some tools are only valid after this chat has created")
	assert.Contains(t, instruction, "plugin_file_edit for focused changes")
	assert.Contains(t, instruction, "plugin_manifest_write for manifest.yaml changes")
	assert.Contains(t, instruction, "Treat the plugin as Storyden software you own")
	assert.Contains(t, instruction, "Do not stop at a proposal")
	assert.Contains(t, instruction, "compare the implemented plugin against the user’s original requested outcome")
	assert.Contains(t, instruction, "Never ask the user to perform development work")
	assert.Contains(t, instruction, "Before writing code, decide whether the requested behaviour only reacts to delivered events")
	assert.Contains(t, instruction, "If plugin code uses BuildAPIClient, manifest.yaml must include access")
	assert.Contains(t, instruction, "stable bot account handle, display name, and the narrow Storyden permissions")
	assert.Contains(t, instruction, "/docs/introduction/members/permissions")
	assert.Contains(t, instruction, "Use plugin_validate while iterating")
	assert.Contains(t, instruction, "Use plugin_install to compile, package, upload or update")
	assert.Contains(t, instruction, "plugin_install packages internally")
	assert.Contains(t, instruction, "It is never acceptable to claim requested behaviour is complete")
	assert.Contains(t, instruction, "Do not leave TODOs, placeholders, stubs")
	assert.Contains(t, instruction, "Read the current stored configuration during startup")
	assert.Contains(t, instruction, "Do not rely only on a configuration callback")
	assert.NotContains(t, instruction, "plugin_apply_patch")
	assert.NotContains(t, instruction, "plugin_package")
	assert.NotContains(t, instruction, "plugin_sdk_reference")
	assert.Contains(t, instruction, "Treat the plugin as your responsibility.")
	assert.Contains(t, instruction, "Prefer comments explaining why rather than what.")
	assert.Contains(t, instruction, "There is no shell, terminal")
	assert.NotContains(t, instruction, "plugin_run_bash")
}

func TestInstructionIncludesUntrustedCommandAccessWhenWorkspaceAllowsIt(t *testing.T) {
	agent := &Agent{}

	instruction, err := agent.instruction(newPluginBuilderReadonlyTestContext(map[string]any{
		workspacestate.WorkspaceStateKey: map[string]any{
			"workspace_id":             xid.New().String(),
			"workspace_instance_id":    xid.New().String(),
			"provider":                 "sprites",
			"provider_state":           map[string]any{},
			"allow_untrusted_commands": true,
			"metadata":                 map[string]any{},
		},
	}))
	require.NoError(t, err)

	assert.Contains(t, instruction, "plugin_run_bash")
	assert.NotContains(t, instruction, "There is no shell, terminal")
}

func TestInstructionIncludesInstalledState(t *testing.T) {
	agent := &Agent{}

	instruction, err := agent.instruction(newPluginBuilderReadonlyTestContext(map[string]any{
		pluginBuildTargetStateKey: pluginBuildTarget{
			Mode:           pluginBuildTargetModeExisting,
			ManifestID:     "welcome-plugin",
			InstallationID: xid.New().String(),
		},
	}))
	require.NoError(t, err)

	assert.Contains(t, instruction, `Current chat state: working on installed plugin "welcome-plugin"`)
	assert.Contains(t, instruction, "for this plugin only")
}

func newTestPluginBuilderToolset(t *testing.T) *pluginBuilderToolset {
	t.Helper()

	builder := &Agent{}
	toolset, err := builder.buildToolset()
	require.NoError(t, err)
	dynamic, ok := toolset.(*pluginBuilderToolset)
	require.True(t, ok)
	return dynamic
}

func pluginToolNames(tools []adktool.Tool) []string {
	names := make([]string, 0, len(tools))
	for _, tool := range tools {
		names = append(names, tool.Name())
	}
	return names
}

type pluginBuilderReadonlyTestContext struct {
	context.Context
	state pluginBuilderReadonlyTestState
}

func newPluginBuilderReadonlyTestContext(values map[string]any) *pluginBuilderReadonlyTestContext {
	state := pluginBuilderReadonlyTestState{}
	for key, value := range values {
		state[key] = value
	}
	return &pluginBuilderReadonlyTestContext{
		Context: context.Background(),
		state:   state,
	}
}

func (c *pluginBuilderReadonlyTestContext) UserContent() *genai.Content { return nil }
func (c *pluginBuilderReadonlyTestContext) InvocationID() string        { return "" }
func (c *pluginBuilderReadonlyTestContext) AgentName() string           { return AgentName }
func (c *pluginBuilderReadonlyTestContext) ReadonlyState() session.ReadonlyState {
	return c.state
}
func (c *pluginBuilderReadonlyTestContext) UserID() string    { return "" }
func (c *pluginBuilderReadonlyTestContext) AppName() string   { return AgentName }
func (c *pluginBuilderReadonlyTestContext) SessionID() string { return "" }
func (c *pluginBuilderReadonlyTestContext) Branch() string    { return "" }

type pluginBuilderReadonlyTestState map[string]any

func (s pluginBuilderReadonlyTestState) Get(key string) (any, error) {
	value, ok := s[key]
	if !ok {
		return nil, session.ErrStateKeyNotExist
	}
	return value, nil
}

func (s pluginBuilderReadonlyTestState) All() iter.Seq2[string, any] {
	return func(yield func(string, any) bool) {
		for key, value := range s {
			if !yield(key, value) {
				return
			}
		}
	}
}
