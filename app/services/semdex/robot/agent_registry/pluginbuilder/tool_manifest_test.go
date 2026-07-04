package pluginbuilder

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"

	pluginresource "github.com/Southclaws/storyden/app/resources/plugin"
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
	require.Equal(t, packagedPluginCommand, mf.Manifest.Command)
	require.Empty(t, mf.Manifest.Args)
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
	require.Equal(t, packagedPluginCommand, mf.Manifest.Command)
	require.Empty(t, mf.Manifest.Args)
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
	_, err = workspace.WriteFile(ctx, "go.mod", []byte("module example.com/replybot\n\ngo 1.24\n"))
	require.NoError(t, err)
	_, err = workspace.WriteFile(ctx, "main.go", []byte(`package main

// access marker: BuildAPIClient(
func main() {}
`))
	require.NoError(t, err)

	_, err = buildPackage(ctx, workspace)
	require.NoError(t, err)
}

func TestBuildPackageCompilesMainExeAndRewritesRuntimeManifest(t *testing.T) {
	ctx := context.Background()
	workspace, err := localworkspace.NewWorkspace(t.TempDir())
	require.NoError(t, err)

	_, err = workspace.WriteFile(ctx, manifestYAMLFilename, []byte(`id: compiled-plugin
name: Compiled Plugin
author: storyden
description: Compiles to a managed runtime binary.
version: 0.1.0
command: go
args:
  - run
  - .
`))
	require.NoError(t, err)
	_, err = workspace.WriteFile(ctx, "go.mod", []byte("module example.com/compiledplugin\n\ngo 1.24\n"))
	require.NoError(t, err)
	_, err = workspace.WriteFile(ctx, "main.go", []byte("package main\n\nfunc main() {}\n"))
	require.NoError(t, err)

	pkg, err := buildPackageForTarget(ctx, workspace, pluginBuildRuntimeTarget{
		GOOS:   runtime.GOOS,
		GOARCH: runtime.GOARCH,
	})
	require.NoError(t, err)
	require.Contains(t, pkg.Files, pluginresource.ArchiveManifestFileName)
	require.Contains(t, pkg.Files, "go.mod")
	require.Contains(t, pkg.Files, "main.go")
	require.Contains(t, pkg.Files, packagedPluginBinary)

	manifest := readPackagedManifestForTest(t, pkg.Bytes)
	require.Equal(t, "compiled-plugin", manifest["id"])
	require.Equal(t, packagedPluginCommand, manifest["command"])
	require.Empty(t, manifest["args"])
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

func TestPluginManifestToolInputSchemaAcceptsManagedLaunchFields(t *testing.T) {
	schema := pluginManifestToolInputSchema()

	require.Contains(t, schema.Properties, "command")
	require.Contains(t, schema.Properties, "args")
	require.NotContains(t, schema.Required, "command")

	resolved, err := schema.Resolve(nil)
	require.NoError(t, err)
	raw := testManifestInput("discord-link-archiver")
	require.NoError(t, resolved.Validate(raw))

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

func readPackagedManifestForTest(t *testing.T, data []byte) map[string]any {
	t.Helper()

	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	require.NoError(t, err)
	for _, file := range zr.File {
		if file.Name != pluginresource.ArchiveManifestFileName {
			continue
		}
		rc, err := file.Open()
		require.NoError(t, err)
		defer rc.Close()
		content, err := io.ReadAll(rc)
		require.NoError(t, err)
		var manifest map[string]any
		require.NoError(t, json.Unmarshal(content, &manifest))
		return manifest
	}
	require.FailNow(t, "manifest.json not found")
	return nil
}
