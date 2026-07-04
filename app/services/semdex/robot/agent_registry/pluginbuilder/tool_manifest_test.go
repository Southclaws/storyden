package pluginbuilder

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	localworkspace "github.com/Southclaws/storyden/app/services/semdex/robot/workspaceprovider/local"
	"github.com/Southclaws/storyden/lib/plugin/rpc"
)

func TestWriteManifestWritesConfigurationSchema(t *testing.T) {
	ctx := context.Background()
	workspace, err := localworkspace.NewWorkspace(t.TempDir())
	require.NoError(t, err)

	agent := &Agent{workspace: workspace}
	state := pluginBuilderTestState{}
	toolCtx := &pluginBuilderTestContext{Context: ctx, state: state}
	require.NoError(t, agent.setPluginBuildTarget(toolCtx, pluginBuildTarget{
		Mode:       pluginBuildTargetModeNew,
		ManifestID: "discord-link-archiver",
	}))

	result, err := agent.WriteManifest(toolCtx, testManifestInput("discord-link-archiver"))
	require.NoError(t, err)
	require.Equal(t, manifestYAMLFilename, result.Path)
	require.Equal(t, "discord-link-archiver", result.ID)
	require.NotEmpty(t, result.Revision)

	read, err := workspace.ReadFile(ctx, manifestYAMLFilename, -1)
	require.NoError(t, err)
	require.Contains(t, string(read.Content), "configuration_schema:")
	require.Contains(t, string(read.Content), "discord_bot_token")
	require.NotContains(t, string(read.Content), "\nconfiguration:")

	mf, err := readProjectManifest(ctx, workspace)
	require.NoError(t, err)
	schema, ok := mf.Manifest.ConfigurationSchema.Get()
	require.True(t, ok)
	require.Len(t, schema.Fields, 1)
	field, ok := schema.Fields[0].PluginConfigurationFieldUnion.(*rpc.PluginConfigurationFieldString)
	require.True(t, ok)
	require.Equal(t, "discord_bot_token", field.ID)
}

func TestWriteManifestForcesManagedLaunchFields(t *testing.T) {
	ctx := context.Background()
	workspace, err := localworkspace.NewWorkspace(t.TempDir())
	require.NoError(t, err)

	agent := &Agent{workspace: workspace}
	toolCtx := &pluginBuilderTestContext{
		Context: ctx,
		state: pluginBuilderTestState{
			pluginBuildTargetStateKey: pluginBuildTarget{
				Mode:       pluginBuildTargetModeNew,
				ManifestID: "discord-link-archiver",
			},
		},
	}

	raw := testManifestInput("discord-link-archiver")
	raw["command"] = "./discord-link-archiver"
	raw["args"] = []any{"--unexpected"}

	_, err = agent.WriteManifest(toolCtx, raw)
	require.NoError(t, err)

	mf, err := readProjectManifest(ctx, workspace)
	require.NoError(t, err)
	require.Equal(t, managedPluginCommand, mf.Manifest.Command)
	require.Equal(t, managedPluginArgs, mf.Manifest.Args)
}

func TestWriteManifestDoesNotRequireManagedLaunchFields(t *testing.T) {
	ctx := context.Background()
	workspace, err := localworkspace.NewWorkspace(t.TempDir())
	require.NoError(t, err)

	agent := &Agent{workspace: workspace}
	toolCtx := &pluginBuilderTestContext{
		Context: ctx,
		state: pluginBuilderTestState{
			pluginBuildTargetStateKey: pluginBuildTarget{
				Mode:       pluginBuildTargetModeNew,
				ManifestID: "discord-link-archiver",
			},
		},
	}

	raw := testManifestInput("discord-link-archiver")
	delete(raw, "command")
	delete(raw, "args")

	_, err = agent.WriteManifest(toolCtx, raw)
	require.NoError(t, err)

	mf, err := readProjectManifest(ctx, workspace)
	require.NoError(t, err)
	require.Equal(t, managedPluginCommand, mf.Manifest.Command)
	require.Equal(t, managedPluginArgs, mf.Manifest.Args)
}

func TestWriteManifestRejectsConfigurationField(t *testing.T) {
	ctx := context.Background()
	workspace, err := localworkspace.NewWorkspace(t.TempDir())
	require.NoError(t, err)

	agent := &Agent{workspace: workspace}
	toolCtx := &pluginBuilderTestContext{
		Context: ctx,
		state: pluginBuilderTestState{
			pluginBuildTargetStateKey: pluginBuildTarget{
				Mode:       pluginBuildTargetModeNew,
				ManifestID: "discord-link-archiver",
			},
		},
	}

	raw := testManifestInput("discord-link-archiver")
	raw["configuration"] = map[string]any{"discord_bot_token": "abc"}
	delete(raw, "configuration_schema")

	_, err = agent.WriteManifest(toolCtx, raw)
	require.Error(t, err)
	require.Contains(t, err.Error(), `unknown manifest field "configuration"`)
	require.Contains(t, err.Error(), `did you mean "configuration_schema"`)
}

func TestWriteManifestRejectsUnknownConfigurationSchemaField(t *testing.T) {
	ctx := context.Background()
	workspace, err := localworkspace.NewWorkspace(t.TempDir())
	require.NoError(t, err)

	agent := &Agent{workspace: workspace}
	toolCtx := &pluginBuilderTestContext{
		Context: ctx,
		state: pluginBuilderTestState{
			pluginBuildTargetStateKey: pluginBuildTarget{
				Mode:       pluginBuildTargetModeNew,
				ManifestID: "discord-link-archiver",
			},
		},
	}

	raw := testManifestInput("discord-link-archiver")
	configurationSchema := raw["configuration_schema"].(map[string]any)
	configurationSchema["extra"] = "not-valid"

	_, err = agent.WriteManifest(toolCtx, raw)
	require.Error(t, err)
	require.Contains(t, err.Error(), "configuration_schema.extra is not a valid field")
}

func TestReadProjectManifestRejectsConfigurationField(t *testing.T) {
	ctx := context.Background()
	workspace, err := localworkspace.NewWorkspace(t.TempDir())
	require.NoError(t, err)

	_, err = workspace.WriteFile(ctx, manifestYAMLFilename, []byte(`id: discord-link-archiver
name: Discord Link Archiver
author: storyden
description: Archives links posted in Discord to the Storyden library.
version: 0.1.0
command: go
args:
  - run
  - .
configuration:
  discord_bot_token: abc
`))
	require.NoError(t, err)

	_, err = readProjectManifest(ctx, workspace)
	require.Error(t, err)
	require.Contains(t, err.Error(), `unknown manifest field "configuration"`)
	require.Contains(t, err.Error(), `did you mean "configuration_schema"`)
}

func TestWriteManifestReportsMissingConfigurationFieldDescription(t *testing.T) {
	ctx := context.Background()
	workspace, err := localworkspace.NewWorkspace(t.TempDir())
	require.NoError(t, err)

	agent := &Agent{workspace: workspace}
	toolCtx := &pluginBuilderTestContext{
		Context: ctx,
		state: pluginBuilderTestState{
			pluginBuildTargetStateKey: pluginBuildTarget{
				Mode:       pluginBuildTargetModeNew,
				ManifestID: "discord-link-archiver",
			},
		},
	}

	raw := testManifestInput("discord-link-archiver")
	configurationSchema := raw["configuration_schema"].(map[string]any)
	fields := configurationSchema["fields"].([]any)
	field := fields[0].(map[string]any)
	delete(field, "description")

	_, err = agent.WriteManifest(toolCtx, raw)
	require.Error(t, err)
	require.Contains(t, err.Error(), "configuration_schema.fields[0].description is required")
}

func TestWriteManifestReportsUnknownConfigurationFieldType(t *testing.T) {
	ctx := context.Background()
	workspace, err := localworkspace.NewWorkspace(t.TempDir())
	require.NoError(t, err)

	agent := &Agent{workspace: workspace}
	toolCtx := &pluginBuilderTestContext{
		Context: ctx,
		state: pluginBuilderTestState{
			pluginBuildTargetStateKey: pluginBuildTarget{
				Mode:       pluginBuildTargetModeNew,
				ManifestID: "discord-link-archiver",
			},
		},
	}

	raw := testManifestInput("discord-link-archiver")
	configurationSchema := raw["configuration_schema"].(map[string]any)
	fields := configurationSchema["fields"].([]any)
	field := fields[0].(map[string]any)
	field["type"] = "password"

	_, err = agent.WriteManifest(toolCtx, raw)
	require.Error(t, err)
	require.Contains(t, err.Error(), "configuration_schema.fields[0].type must be one of:")
	require.Contains(t, err.Error(), "string")
	require.Contains(t, err.Error(), "number")
	require.Contains(t, err.Error(), "boolean")
}

func TestBuildPackageRejectsInvalidManifestBeforeArchive(t *testing.T) {
	ctx := context.Background()
	workspace, err := localworkspace.NewWorkspace(t.TempDir())
	require.NoError(t, err)

	_, err = workspace.WriteFile(ctx, manifestYAMLFilename, []byte(`id: discord-link-archiver
name: Discord Link Archiver
author: storyden
description: Archives links posted in Discord to the Storyden library.
version: 0.1.0
command: go
args:
  - run
  - .
configuration:
  discord_bot_token: abc
`))
	require.NoError(t, err)

	_, err = buildPackage(ctx, workspace)
	require.Error(t, err)
	require.Contains(t, err.Error(), `unknown manifest field "configuration"`)
	require.Contains(t, err.Error(), `did you mean "configuration_schema"`)
}

func TestBuildPackageRejectsBuildAPIClientWithoutAccess(t *testing.T) {
	ctx := context.Background()
	workspace, err := localworkspace.NewWorkspace(t.TempDir())
	require.NoError(t, err)

	_, err = workspace.WriteFile(ctx, manifestYAMLFilename, []byte(`id: reply-bot
name: Reply Bot
author: storyden
description: Replies to new threads.
version: 0.1.0
command: go
args:
  - run
  - .
`))
	require.NoError(t, err)
	_, err = workspace.WriteFile(ctx, "main.go", []byte(`package main

func run() {
	_, _ = pl.BuildAPIClient(ctx)
}
`))
	require.NoError(t, err)

	_, err = buildPackage(ctx, workspace)
	require.Error(t, err)
	require.Contains(t, err.Error(), "manifest access is required because plugin code uses BuildAPIClient")
	require.Contains(t, err.Error(), "stable bot account handle")
	require.Contains(t, err.Error(), "narrow Storyden permissions")
}

func TestBuildPackageAllowsBuildAPIClientWithAccess(t *testing.T) {
	ctx := context.Background()
	workspace, err := localworkspace.NewWorkspace(t.TempDir())
	require.NoError(t, err)

	_, err = workspace.WriteFile(ctx, manifestYAMLFilename, []byte(`id: reply-bot
name: Reply Bot
author: storyden
description: Replies to new threads.
version: 0.1.0
command: go
args:
  - run
  - .
access:
  handle: reply-bot
  name: Reply Bot
  permissions:
    - CREATE_POST
`))
	require.NoError(t, err)
	_, err = workspace.WriteFile(ctx, "main.go", []byte(`package main

func run() {
	_, _ = pl.BuildAPIClient(ctx)
}
`))
	require.NoError(t, err)

	_, err = buildPackage(ctx, workspace)
	require.NoError(t, err)
}

func TestPluginManifestToolInputSchemaRejectsConfiguration(t *testing.T) {
	schema, err := pluginManifestToolInputSchema().Resolve(nil)
	require.NoError(t, err)

	raw := testManifestInput("discord-link-archiver")
	raw["configuration"] = map[string]any{"discord_bot_token": "abc"}
	delete(raw, "configuration_schema")

	err = schema.Validate(raw)
	require.Error(t, err)
}

func TestPluginManifestToolInputSchemaHidesManagedLaunchFields(t *testing.T) {
	schema := pluginManifestToolInputSchema()

	require.NotContains(t, schema.Properties, "command")
	require.NotContains(t, schema.Properties, "args")
	require.NotContains(t, schema.Required, "command")

	resolved, err := schema.Resolve(nil)
	require.NoError(t, err)
	raw := testManifestInput("discord-link-archiver")
	delete(raw, "command")
	delete(raw, "args")
	require.NoError(t, resolved.Validate(raw))
}

func testManifestInput(id string) map[string]any {
	return map[string]any{
		"id":          id,
		"name":        "Discord Link Archiver",
		"author":      "storyden",
		"description": "Archives links posted in Discord to the Storyden library.",
		"version":     "0.1.0",
		"command":     "go",
		"args":        []any{"run", "."},
		"events_consumed": []any{
			"EventThreadPublished",
		},
		"configuration_schema": map[string]any{
			"fields": []any{
				map[string]any{
					"id":          "discord_bot_token",
					"label":       "Discord Bot Token",
					"description": "Token used to connect to Discord.",
					"type":        "string",
				},
			},
		},
	}
}
