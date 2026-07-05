package pluginbuilder

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	localworkspace "github.com/Southclaws/storyden/app/services/semdex/robot/workspaceprovider/local"
)

func TestValidateReportsMissingAccessForBuildAPIClient(t *testing.T) {
	ctx := context.Background()
	workspace, err := localworkspace.NewWorkspace(t.TempDir())
	require.NoError(t, err)

	writeWorkspaceFile(t, ctx, workspace, manifestYAMLFilename, `id: reply-bot
name: Reply Bot
author: storyden
description: Replies to new threads.
version: 0.1.0
command: go
args:
  - run
  - .
`)
	writeWorkspaceFile(t, ctx, workspace, "main.go", `package main

func run() {
	_, _ = pl.BuildAPIClient(ctx)
}
`)

	agent := &Agent{workspace: workspace}
	result, err := agent.Validate(ctx, ValidateInput{SkipGo: true})
	require.NoError(t, err)
	require.False(t, result.Success)
	requireValidationCheck(t, result.Checks, "manifest", true)
	requireValidationCheck(t, result.Checks, "manifest_code_consistency", false)
	require.Contains(t, result.Message, "manifest_code_consistency")
	require.Contains(t, result.Message, "BuildAPIClient")
}

func TestValidateRejectsIncompleteImplementationMarkers(t *testing.T) {
	ctx := context.Background()
	workspace, err := localworkspace.NewWorkspace(t.TempDir())
	require.NoError(t, err)

	writeWorkspaceFile(t, ctx, workspace, manifestYAMLFilename, `id: link-archiver
name: Link Archiver
author: storyden
description: Archives links.
version: 0.1.0
command: go
args:
  - run
  - .
`)
	writeWorkspaceFile(t, ctx, workspace, "main.go", `package main

import "log"

func createPage() {
	log.Println("Would create node")
}
`)

	agent := &Agent{workspace: workspace}
	result, err := agent.Validate(ctx, ValidateInput{SkipGo: true})
	require.NoError(t, err)
	require.False(t, result.Success)
	requireValidationCheck(t, result.Checks, "implementation_completeness", false)
	require.Contains(t, result.Message, "implementation_completeness")
	require.Contains(t, result.Message, "would create")
}

func TestValidateReportsUnhandledConfigurationSchemaFields(t *testing.T) {
	ctx := context.Background()
	workspace, err := localworkspace.NewWorkspace(t.TempDir())
	require.NoError(t, err)

	writeWorkspaceFile(t, ctx, workspace, manifestYAMLFilename, `id: discord-bot
name: Discord Bot
author: storyden
description: Uses Discord.
version: 0.1.0
command: go
args:
  - run
  - .
configuration_schema:
  fields:
    - id: discord_token
      label: Discord Bot Token
      description: Bot token for Discord.
      type: string
`)
	writeWorkspaceFile(t, ctx, workspace, "main.go", `package main

type pluginConfig struct {
	values map[string]any
}

func parseConfig(raw map[string]any) (pluginConfig, bool, error) {
	return pluginConfig{values: raw}, true, nil
}
`)

	agent := &Agent{workspace: workspace}
	result, err := agent.Validate(ctx, ValidateInput{SkipGo: true})
	require.NoError(t, err)
	require.False(t, result.Success)
	requireValidationCheck(t, result.Checks, "configuration_implementation", false)
	require.Contains(t, result.Message, "discord_token")
	require.Contains(t, result.NextAction, "Handle every manifest configuration_schema field")
}

func TestValidateRejectsEmptyConfigurationStruct(t *testing.T) {
	ctx := context.Background()
	workspace, err := localworkspace.NewWorkspace(t.TempDir())
	require.NoError(t, err)

	writeWorkspaceFile(t, ctx, workspace, manifestYAMLFilename, `id: discord-bot
name: Discord Bot
author: storyden
description: Uses Discord.
version: 0.1.0
command: go
args:
  - run
  - .
configuration_schema:
  fields:
    - id: discord_token
      label: Discord Bot Token
      description: Bot token for Discord.
      type: string
`)
	writeWorkspaceFile(t, ctx, workspace, "main.go", `package main

type pluginConfig struct {
	// Add manifest configuration_schema fields here as this plugin grows.
}

func parseConfig(raw map[string]any) string {
	return raw["discord_token"].(string)
}
`)

	agent := &Agent{workspace: workspace}
	result, err := agent.Validate(ctx, ValidateInput{SkipGo: true})
	require.NoError(t, err)
	require.False(t, result.Success)
	requireValidationCheck(t, result.Checks, "configuration_implementation", false)
	require.Contains(t, result.Message, "empty configuration structs")
	require.Contains(t, result.Message, "pluginConfig")
}

func TestValidateDoesNotTreatConfigurationFieldStringMentionAsHandled(t *testing.T) {
	ctx := context.Background()
	workspace, err := localworkspace.NewWorkspace(t.TempDir())
	require.NoError(t, err)

	writeWorkspaceFile(t, ctx, workspace, manifestYAMLFilename, `id: discord-bot
name: Discord Bot
author: storyden
description: Uses Discord.
version: 0.1.0
command: go
args:
  - run
  - .
configuration_schema:
  fields:
    - id: discord_token
      label: Discord Bot Token
      description: Bot token for Discord.
      type: string
`)
	writeWorkspaceFile(t, ctx, workspace, "main.go", `package main

import "log"

func parseConfig(raw map[string]any) {
	log.Println("discord_token is required")
}
`)

	agent := &Agent{workspace: workspace}
	result, err := agent.Validate(ctx, ValidateInput{SkipGo: true})
	require.NoError(t, err)
	require.False(t, result.Success)
	requireValidationCheck(t, result.Checks, "configuration_implementation", false)
	require.Contains(t, result.Message, "discord_token")
}

func TestValidateAcceptsHandledConfigurationSchemaFields(t *testing.T) {
	ctx := context.Background()
	workspace, err := localworkspace.NewWorkspace(t.TempDir())
	require.NoError(t, err)

	writeWorkspaceFile(t, ctx, workspace, manifestYAMLFilename, `id: discord-bot
name: Discord Bot
author: storyden
description: Uses Discord.
version: 0.1.0
command: go
args:
  - run
  - .
configuration_schema:
  fields:
    - id: discord_token
      label: Discord Bot Token
      description: Bot token for Discord.
      type: string
`)
	writeWorkspaceFile(t, ctx, workspace, "main.go", `package main

import "errors"

type pluginConfig struct {
	DiscordToken string
}

func parseConfig(raw map[string]any) (pluginConfig, bool, error) {
	token, ok := raw["discord_token"].(string)
	if !ok || token == "" {
		return pluginConfig{}, false, nil
	}
	if token == "invalid" {
		return pluginConfig{}, false, errors.New("invalid token")
	}
	return pluginConfig{DiscordToken: token}, true, nil
}
`)

	agent := &Agent{workspace: workspace}
	result, err := agent.Validate(ctx, ValidateInput{SkipGo: true})
	require.NoError(t, err)
	require.True(t, result.Success, "%#v", result.Checks)
	requireValidationCheck(t, result.Checks, "configuration_implementation", true)
}

func TestValidateAcceptsConfigurationFieldReadThroughConstant(t *testing.T) {
	ctx := context.Background()
	workspace, err := localworkspace.NewWorkspace(t.TempDir())
	require.NoError(t, err)

	writeWorkspaceFile(t, ctx, workspace, manifestYAMLFilename, `id: discord-bot
name: Discord Bot
author: storyden
description: Uses Discord.
version: 0.1.0
command: go
args:
  - run
  - .
configuration_schema:
  fields:
    - id: discord_token
      label: Discord Bot Token
      description: Bot token for Discord.
      type: string
`)
	writeWorkspaceFile(t, ctx, workspace, "main.go", `package main

const configDiscordToken = "discord_token"

type pluginConfig struct {
	DiscordToken string
}

func parseConfig(raw map[string]any) (pluginConfig, bool, error) {
	token, ok := raw[configDiscordToken].(string)
	if !ok || token == "" {
		return pluginConfig{}, false, nil
	}
	return pluginConfig{DiscordToken: token}, true, nil
}
`)

	agent := &Agent{workspace: workspace}
	result, err := agent.Validate(ctx, ValidateInput{SkipGo: true})
	require.NoError(t, err)
	require.True(t, result.Success, "%#v", result.Checks)
	requireValidationCheck(t, result.Checks, "configuration_implementation", true)
}

func TestValidateAcceptsConfigurationFieldReadThroughSwitchOnRangeKey(t *testing.T) {
	ctx := context.Background()
	workspace, err := localworkspace.NewWorkspace(t.TempDir())
	require.NoError(t, err)

	writeWorkspaceFile(t, ctx, workspace, manifestYAMLFilename, `id: discord-bot
name: Discord Bot
author: storyden
description: Uses Discord.
version: 0.1.0
command: go
args:
  - run
  - .
configuration_schema:
  fields:
    - id: discord_token
      label: Discord Bot Token
      description: Bot token for Discord.
      type: string
    - id: discord_channel_id
      label: Discord Channel ID
      description: Channel to post into.
      type: string
`)
	writeWorkspaceFile(t, ctx, workspace, "main.go", `package main

const channelField = "discord_channel_id"

type pluginConfig struct {
	DiscordToken     string
	DiscordChannelID string
}

func parseConfig(raw map[string]any) pluginConfig {
	var cfg pluginConfig
	for key, value := range raw {
		switch key {
		case "discord_token":
			cfg.DiscordToken, _ = value.(string)
		case channelField:
			cfg.DiscordChannelID, _ = value.(string)
		}
	}
	return cfg
}
`)

	agent := &Agent{workspace: workspace}
	result, err := agent.Validate(ctx, ValidateInput{SkipGo: true})
	require.NoError(t, err)
	require.True(t, result.Success, "%#v", result.Checks)
	requireValidationCheck(t, result.Checks, "configuration_implementation", true)
}

func TestValidateAcceptsConfigurationFieldReadThroughJSONUnmarshal(t *testing.T) {
	ctx := context.Background()
	workspace, err := localworkspace.NewWorkspace(t.TempDir())
	require.NoError(t, err)

	writeWorkspaceFile(t, ctx, workspace, manifestYAMLFilename, `id: discord-bot
name: Discord Bot
author: storyden
description: Uses Discord.
version: 0.1.0
command: go
args:
  - run
  - .
configuration_schema:
  fields:
    - id: discord_token
      label: Discord Bot Token
      description: Bot token for Discord.
      type: string
    - id: message_prefix
      label: Message Prefix
      description: Prefix for messages.
      type: string
`)
	writeWorkspaceFile(t, ctx, workspace, "main.go", `package main

import "encoding/json"

type runtimeConfig struct {
	DiscordToken string `+"`json:\"discord_token\"`"+`
	MessagePrefix string `+"`json:\"message_prefix,omitempty\"`"+`
}

func parseConfig(raw map[string]any) (runtimeConfig, error) {
	cfg := &runtimeConfig{}
	data, err := json.Marshal(raw)
	if err != nil {
		return runtimeConfig{}, err
	}
	if err := json.Unmarshal(data, cfg); err != nil {
		return runtimeConfig{}, err
	}
	return *cfg, nil
}
`)

	agent := &Agent{workspace: workspace}
	result, err := agent.Validate(ctx, ValidateInput{SkipGo: true})
	require.NoError(t, err)
	require.True(t, result.Success, "%#v", result.Checks)
	requireValidationCheck(t, result.Checks, "configuration_implementation", true)
}

func TestValidateIgnoresGeneratedIncompleteImplementationMarkers(t *testing.T) {
	ctx := context.Background()
	workspace, err := localworkspace.NewWorkspace(t.TempDir())
	require.NoError(t, err)

	writeWorkspaceFile(t, ctx, workspace, manifestYAMLFilename, `id: generated-plugin
name: Generated Plugin
author: storyden
description: Has generated code.
version: 0.1.0
command: go
args:
  - run
  - .
`)
	writeWorkspaceFile(t, ctx, workspace, "go.mod", "module example.com/generatedplugin\n\ngo 1.24\n")
	writeWorkspaceFile(t, ctx, workspace, "main.go", "package main\n\nfunc main() {}\n")
	writeWorkspaceFile(t, ctx, workspace, "generated.go", `// Code generated by test. DO NOT EDIT.
package main

const marker = "TODO generated placeholder"
`)

	agent := &Agent{workspace: workspace}
	result, err := agent.Validate(ctx, ValidateInput{SkipGo: true})
	require.NoError(t, err)
	require.True(t, result.Success, "%#v", result.Checks)
	requireValidationCheck(t, result.Checks, "implementation_completeness", true)
}

func TestIncompleteImplementationMarkerDoesNotMatchPartialWords(t *testing.T) {
	_, ok := incompleteImplementationMarker(`const service = "todoist"`)
	require.False(t, ok)
}

func TestIncompleteImplementationMarkerRejectsCannedRobotSummaries(t *testing.T) {
	_, ok := incompleteImplementationMarker(`summary := "Moderation triage requested from robot system"`)
	require.True(t, ok)

	_, ok = incompleteImplementationMarker(`// For now, return a generic moderation summary.`)
	require.True(t, ok)
}

func TestValidateRunsFullChecks(t *testing.T) {
	ctx := context.Background()
	workspace, err := localworkspace.NewWorkspace(t.TempDir())
	require.NoError(t, err)

	writeWorkspaceFile(t, ctx, workspace, manifestYAMLFilename, `id: valid-plugin
name: Valid Plugin
author: storyden
description: A small valid plugin.
version: 0.1.0
command: go
args:
  - run
  - .
`)
	writeWorkspaceFile(t, ctx, workspace, "go.mod", "module example.com/plugin\n\ngo 1.24\n")
	writeWorkspaceFile(t, ctx, workspace, "main.go", "package main\n\nfunc main() {}\n")

	agent := &Agent{workspace: workspace}
	result, err := agent.Validate(ctx, ValidateInput{})
	require.NoError(t, err)
	require.True(t, result.Success, "%#v", result.Checks)
	requireValidationCheck(t, result.Checks, "manifest", true)
	requireValidationCheck(t, result.Checks, "go_fmt", true)
	requireValidationCheck(t, result.Checks, "go_tidy", true)
	requireValidationCheck(t, result.Checks, "go_vet", true)
	requireValidationCheck(t, result.Checks, "go_test", true)
}

func TestValidationFailureSummarySkipsGoPackageHeaders(t *testing.T) {
	result := ValidateResult{Checks: []ValidationCheck{
		{
			Name:    "go_vet",
			Success: false,
			Message: "# storyden.local/plugins/discord-boot-tagger\n# [storyden.local/plugins/discord-boot-tagger]\nvet: ./main.go:64:24: sess.UserChannels undefined (type *discordgo.Session has no field or method UserChannels)",
		},
		{
			Name:    "go_test",
			Success: false,
			Output:  "# storyden.local/plugins/discord-boot-tagger\n./main.go:4:2: \"context\" imported and not used\nFAIL\tstoryden.local/plugins/discord-boot-tagger [build failed]\nFAIL\n",
		},
	}}

	summary := validationFailureSummary(result)
	require.Contains(t, summary, "go_vet: vet: ./main.go:64:24: sess.UserChannels undefined")
	require.Contains(t, summary, "go_test: ./main.go:4:2: \"context\" imported and not used")
	require.NotContains(t, summary, "go_vet: # storyden.local")
	require.NotContains(t, summary, "go_test: # storyden.local")
}

func TestValidationNextActionForGoErrorsPointsToDiscovery(t *testing.T) {
	result := ValidateResult{Checks: []ValidationCheck{
		{Name: "go_vet", Success: false},
	}}

	next := validationNextAction(result)
	require.Contains(t, next, "plugin_go_package_symbols")
	require.Contains(t, next, "instead of asking the user")
}

func TestValidationNextActionPrioritisesMissingGoSymbolsOverConfiguration(t *testing.T) {
	result := ValidateResult{Checks: []ValidationCheck{
		{Name: "configuration_implementation", Success: false, Message: "configuration_schema fields are not read from runtime configuration in Go source: discord_token"},
		{Name: "go_test", Success: false, Output: "./main.go:64:24: client.RobotChatSSEWithResponse undefined (type *openapi.ClientWithResponses has no field or method RobotChatSSEWithResponse)"},
	}}

	next := validationNextAction(result)
	require.Contains(t, next, "pl.RunRobot")
	require.NotContains(t, next, "Handle every manifest configuration_schema field")
}

func requireValidationCheck(t *testing.T, checks []ValidationCheck, name string, success bool) {
	t.Helper()

	for _, check := range checks {
		if check.Name == name {
			require.Equal(t, success, check.Success, "%s: %s", name, check.Message)
			return
		}
	}
	require.Fail(t, "missing validation check", "missing %q in %#v", name, checks)
}
