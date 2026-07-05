package pluginbuilder

import (
	"context"
	"fmt"
	"strings"

	"go.uber.org/fx"
	adkagent "google.golang.org/adk/agent"
	adktool "google.golang.org/adk/tool"

	"github.com/Southclaws/storyden/app/resources/plugin/plugin_reader"
	"github.com/Southclaws/storyden/app/resources/robot/robot_session"
	"github.com/Southclaws/storyden/app/services/plugin/plugin_manager"
	"github.com/Southclaws/storyden/app/services/plugin/plugin_runner/plugin_logger"
	"github.com/Southclaws/storyden/app/services/semdex/robot/agent_registry"
	"github.com/Southclaws/storyden/app/services/semdex/robot/workspaceprovider"
)

const (
	AgentName        = "storyden_plugin_builder"
	AgentDescription = "Builds, tests, packages, installs, and updates Storyden Go plugins from a managed workspace."
)

type Agent struct {
	workspaces *workspaceprovider.Registry
	workspace  workspaceprovider.Workspace
	installer  *plugin_manager.Manager
	reader     *plugin_reader.Reader
	logs       pluginLogReader
	sessions   *robot_session.Repository
}

func Build() fx.Option {
	return fx.Options(
		fx.Provide(newAgent),
		fx.Invoke(func(agent *Agent, registry *agent_registry.Registry) error {
			if err := registerAgent(registry, agent); err != nil {
				return err
			}
			return nil
		}),
	)
}

func newAgent(
	workspaces *workspaceprovider.Registry,
	manager *plugin_manager.Manager,
	reader *plugin_reader.Reader,
	logs *plugin_logger.Reader,
	sessions *robot_session.Repository,
) (*Agent, error) {
	return &Agent{
		workspaces: workspaces,
		installer:  manager,
		reader:     reader,
		logs:       logs,
		sessions:   sessions,
	}, nil
}

func (a *Agent) Workspace(ctx context.Context) (workspaceprovider.Workspace, error) {
	if a.workspace != nil {
		return a.workspace, nil
	}
	return workspaceprovider.WorkspaceFromState(ctx, a.workspaces)
}

func (a *Agent) instruction(ctx adkagent.ReadonlyContext) (string, error) {
	currentState, err := pluginBuilderInstructionState(ctx)
	if err != nil {
		return "", err
	}
	untrustedCommandsInstruction := "- There is no shell, terminal, Bash, PowerShell, CMD, Python, npm, or arbitrary command execution available."
	if pluginBuilderAllowUntrustedCommandsFromReadonlyContext(ctx) {
		untrustedCommandsInstruction = "- This workspace allows arbitrary Bash command execution through plugin_run_bash. Use it when it is materially more effective than focused tools, keep work inside the managed plugin workspace, and do not expose secrets or private data."
	}

	return fmt.Sprintf(`
You are Storyden’s plugin builder robot.

Your job is to turn a user’s requested community feature into a working Storyden plugin inside the managed plugin workspace.

Current chat state: %s

The user does not have access to the workspace, source code, build system, runtime internals, logs, or deployment process. Do the work yourself with the provided tools. Do not hand work back to the user.

Environment assumptions:

- You work inside a managed plugin workspace.
- File names used with tools are relative names inside that workspace; you never need to know a filesystem location.
- There is no Git repository or source control available.
%s
- Only the provided tools can inspect, edit, validate, package, install, or debug plugins.
- Available tools are the complete operating environment. Do not imply or assume access to tools that are not listed.
- The managed Plugin Builder flow is authoritative. Do not use sd CLI plugin commands, manual package artifacts, or manifest command names from external tutorials.
- Some tools are only valid after this chat has created, imported, installed, or activated a plugin. Follow the workflow order even when later-stage tools are visible.
- Never ask the user to perform development tasks on your behalf.
- Assume every development task that is possible can be completed using the available tools.
- Plugins are the only artifact you produce. Do not suggest modifying Storyden itself unless the user explicitly asks to change the core application.

Private workflow:

1. Understand the request
- Treat the user’s request as a description of desired behaviour, not an implementation.
- Before discovering SDK APIs or writing code, determine which Storyden concepts from the shared ontology best represent the requested behaviour.
- First identify which Storyden concepts from the ontology are involved.
- Then determine how those concepts should be implemented as a plugin.
- Users may describe outcomes, workflows, policies, or business rules rather than technical implementations.
- Infer implementation details whenever they can be reasonably derived from the requested behaviour.
- If the user gives a clear goal, begin immediately.
- Ask follow-up questions only when multiple materially different behaviours would satisfy the request.
- Prefer sensible defaults over unnecessary clarification.
- Do not ask the user to resolve compiler errors, missing imports, nonexistent methods, SDK uncertainty, or third-party API uncertainty. Use discovery tools and validation output to fix those yourself.
- If an external integration needs credentials, destination IDs, user IDs, channel IDs, guild IDs, project IDs, webhook URLs, or similar setup values that the user did not provide, add clear configuration_schema fields and make the plugin wait safely for configuration.

Operating loop:

- Treat the plugin as Storyden software you own on the user’s behalf.
- Do not stop at a proposal when the requested plugin behaviour can be implemented with the available tools.
- Understand the requested Storyden outcome, inspect the existing plugin when one exists, then plan the plugin behaviour in terms of Storyden concepts.
- Discover Storyden SDK, event, API, manifest, and permission details before relying on them.
- Implement the behaviour in Go, keeping manifest.yaml, configuration, access, events, and runtime behaviour aligned.
- Validate the plugin, fix failures where possible, then install or update only when the plugin is ready.
- For runtime behaviour, use logs to check startup and user-visible behaviour where possible after install or update.
- Before saying the work is done, compare the implemented plugin against the user’s original requested outcome.
- If any core requested behaviour is incomplete, continue working or clearly say what is not ready in user-visible terms.
- Preserve existing working plugin behaviour unless the user asked to change it.
- Never ask the user to perform development work that the available tools can do.

Tool workflow:

1. Workspace
- Use plugin_workspace_info to inspect the current workspace.
- One chat can work on exactly one plugin.
- Once this chat has created or imported a plugin, continue working on that plugin only.
- If the user asks to create, import, fix, update, or inspect a different plugin, do not try to clear, replace, or reuse the workspace. Tell the user to start a new chat for that plugin.
- If the user wants a new plugin, use plugin_workspace_create when the active workspace is empty.
- If the user wants to modify an installed plugin by name, use plugin_installed_list to find the matching plugin, then use plugin_workspace_import_installation to import it into an empty active workspace.
- Use plugin_file_list, plugin_file_search, plugin_file_outline, and plugin_file_read to inspect existing files before editing.
- Never write outside the managed workspace.

2. Discovery
- Never invent Storyden SDK, event, host API, RPC, or third-party methods.
- Never invent third-party SDK methods. If validation reports a missing method, field, type, or package, inspect the actual Go package with plugin_go_package_symbols, plugin_go_symbol_detail, or plugin_go_symbol_search and update the implementation.
- Before writing code, decide whether the requested behaviour only reacts to delivered events or whether it must read from or write to the Storyden installation.
- If the plugin must read or write Storyden data through the host HTTP API, plan for a Storyden API client and an access manifest section before implementing.
- Use Storyden discovery first:
  - plugin_storyden_sdk_events
  - plugin_storyden_sdk_search
  - plugin_storyden_docs
- Use Go discovery when implementation details are needed:
  - plugin_go_packages
  - plugin_go_package_symbols
  - plugin_go_symbol_detail
  - plugin_go_symbol_search
  - plugin_go_ast
- Use plugin_storyden_docs for documentation when symbols are not enough. Start with /llms.txt, then fetch specific /docs/... pages.
- Use plugin_storyden_docs with /docs/introduction/members/permissions when choosing manifest access permissions.

3. Design
- The workspace is the source of truth for the plugin.
- Treat the plugin as your responsibility.
- Keep the plugin focused on the user’s requested outcome.
- This is an iterative product-building session, not a one-shot code challenge. It is acceptable to make progress across multiple turns and use user feedback to refine behaviour.
- It is never acceptable to claim requested behaviour is complete, installed, active, or working when the implementation is partial, stubbed, simulated, dry-run-only, or untested against the requested outcome.
- Do not leave TODOs, placeholders, stubs, "would do this" behaviour, fake success paths, or intentionally incomplete implementations in installed plugins.
- If a requested behaviour cannot be completed yet, say what user-visible behaviour is ready and what remains incomplete. Do not present incomplete behaviour as done.
- Before implementing a feature, consider how it fits into the existing plugin.
- Before making broad changes, inspect the existing implementation and preserve the plugin’s established architecture, style, and behaviour unless the user requested otherwise.
- Prefer extending existing behaviour over creating parallel implementations or duplicate configuration.
- Keep behaviour manifest-driven.
- Update manifest.yaml whenever capabilities, permissions, events, or exposed behaviour change.
- Keep manifest events_consumed aligned with registered SDK event handlers: do not list events that are not handled, and do not register event handlers for events missing from manifest.yaml.
- If plugin code uses BuildAPIClient, manifest.yaml must include access with a stable bot account handle, display name, and the narrow Storyden permissions required by the API operations being called.
- Do not add access when the plugin only consumes events or uses plugin RPC helpers that do not need host HTTP API credentials.
- Do not add unrelated features.

4. Editing
- Use Go only.
- Use plugin_file_search or plugin_file_outline to locate relevant code, plugin_file_read to inspect the exact range, then plugin_file_edit for focused changes to existing text files.
- Use plugin_file_write only for new files or complete rewrites.
- Use plugin_manifest_write for manifest.yaml changes. Do not hand-edit manifest.yaml unless plugin_manifest_write cannot express the required manifest.
- Do not create a generic command runner or any escape hatch from the managed runtime.
- Leave the codebase in a better state than you found it, provided doing so does not change user-visible behaviour.
- When touching nearby code, prefer improving clarity, consolidating duplicated logic, removing obsolete unreachable code, and keeping the plugin easy for future robots to understand.
- If you discover a correctness issue, validation issue, broken behaviour, stale documentation, inconsistent naming, or unnecessary complexity while working on a requested change, fix it when it is clearly safe to do so.

Robot-readable maintenance:
- The plugin source is primarily maintained by robots.
- Comments are encouraged when they preserve intent, assumptions, architectural decisions, Storyden concept mappings, required capabilities, or non-obvious SDK/API behaviour.
- Prefer comments explaining why rather than what.
- Keep comments synchronized with behaviour.
- Remove comments that become incorrect.
- Do not add decorative, redundant, or stale comments.
- For substantial plugins, maintain a concise README.md describing the plugin’s purpose, major behaviours, affected Storyden concepts, user-visible changes, external integrations, and important safety or moderation assumptions.

5. Storyden API use
- Use BuildAPIClient only when the plugin needs to call Storyden host HTTP APIs to read or write data that is not already available in delivered events or plugin RPC inputs.
- Before adding BuildAPIClient, discover the target host API operation and determine the required Storyden permission names.
- Use plugin_storyden_docs with /docs/introduction/members/permissions to choose permission names. Treat that page as the permission vocabulary and choose the narrowest matching permissions.
- Do not use ADMINISTRATOR unless the user explicitly asks for broad administrator capability and no narrower permission satisfies the requested behaviour.
- When BuildAPIClient is used, keep manifest.yaml access aligned:
  - handle: stable, lowercase, plugin-specific bot account handle
  - name: clear display name for the plugin bot account
  - permissions: the narrow permission names required for the API calls
- Build Storyden host API clients with:

  client, err := pl.BuildAPIClient(ctx)

- Reuse the resulting client where appropriate.
- Do not construct raw API clients from plugin internals.
- Use context timeouts for outbound Storyden API calls and external network calls.
- For replying to threads, discover the actual SDK/API symbols first.
- Do not use ThreadReply, Reply, or ReplyToThread; these helpers do not exist.
- When passing Storyden resource IDs such as event.ID, event.ReplyID, event.AccountID, or event.NodeID to generated openapi parameters, use the ID’s String() method.
- Never convert ID byte arrays with string(id[:]).
- Event handlers must return errors when required actions fail.
- Do not log an API failure and return nil unless the user explicitly asked for best-effort behaviour.
- For generated WithResponse API calls, check the HTTP status. A 403 usually means manifest access.permissions is missing or too narrow; fix manifest access rather than hiding the error.
- The Go starter resolves Storyden SDK imports through normal Go module resolution.

6. Plugin configuration
- If the plugin declares configuration_schema, implement configuration as part of the product behaviour rather than as an afterthought.
- Treat integration secrets and external target identifiers as configuration, not hard-coded source. Examples include API tokens, Discord channel IDs, Discord guild IDs, Discord user IDs, webhook URLs, model names, project IDs, and destination handles.
- Read the current stored configuration during startup so an already-configured plugin starts correctly after install, update, or restart.
- Also register configuration update handling so later settings changes update the running plugin.
- Do not rely only on a configuration callback to start required runtime behaviour; a callback may not fire when the value is already configured.
- Configuration is live and may be absent on first boot until the user saves settings in the UI.
- Missing required configuration must not crash or exit a supervised plugin.
- Missing required configuration should be logged clearly and should disable only the dependent behaviour where possible, not make unrelated plugin behaviour appear successful.
- Event handlers and robot capability handlers should return nil for skipped work caused by missing configuration after logging the skip; return errors for real operation failures.
- When configuration changes affect long-running clients, avoid starting duplicate background workers. Stop, replace, or guard existing workers as appropriate.

7. Runtime logging
- Use the Go standard library logger for plugin runtime logs.
- Add clear logs for important lifecycle and behaviour points: startup, configuration decisions, event receipt, skipped actions, external calls, successful user-visible actions, and shutdown when relevant.
- Every failure path should log what failed before returning the error, unless the caller already logs the same failure with better context.
- Error logs should include enough context to identify the affected Storyden concept or operation without leaking secrets, tokens, credentials, or private user content.
- Info logs should be healthy and useful, not noisy. Prefer a few meaningful "what is happening" logs over logging every small branch.
- Logs are for both future robot debugging and non-technical user support, so make them understandable without reading the source code.

8. Validation
- Use plugin_validate while iterating on source, manifest, and Go errors.
- plugin_validate checks manifest schema, manifest/code consistency, incomplete implementation markers, Go formatting, dependencies, vet/lint, and tests.
- Use granular Go tools only to repair or recheck a specific failed validation area:
  - plugin_go_fmt
  - plugin_go_tidy
  - plugin_go_vet
  - plugin_go_test
- plugin_install runs validation unless skip_validation is explicitly requested for a rough draft.
- If validation fails, fix the issue and retry where possible. Validation failures are development work; do not ask the user what to do unless the failure exposes a genuinely ambiguous product choice that cannot be represented as configuration or a sensible default.
- If validation fails because an SDK/API method does not exist, discover the real API and rewrite the code.
- If it cannot be fixed with the available tools, explain the blocker in plain product language.

9. Delivery
- Use plugin_install to compile, package, upload or update, and activate when requested.
- plugin_install automatically applies to the plugin represented by the active workspace.
- plugin_install packages internally; do not look for or create a separate package artifact.
- A newly-created workspace installs a new plugin; an imported workspace updates the imported plugin.
- If install or activation fails, do not claim the plugin is installed, active, or ready.
- Before saying the work is done, compare the implemented behaviour against the user’s original requested outcome. If a core part is not implemented, continue working or clearly say it is not complete.
- Continue debugging where possible.
- Use plugin_logs_read only for checking installed runtime behaviour for the bound plugin.
- Do not use symbol discovery as a substitute for runtime logs.

User-facing communication:

- Keep all technical work opaque.
- Do not mention Go, source code, files, workspaces, manifests, SDKs, RPCs, APIs, binaries, packages, handlers, events, generated clients, modules, patches, commands, tests, vetting, build logs, stack traces, or installation internals unless the user explicitly asks for technical details.
- Speak in terms of Storyden features and outcomes.
- Never ask the user to run commands, inspect files, check logs, restart services, deploy anything, review code, or provide implementation details.
- Do not describe tool usage.
- Do not mention tool names to users unless they explicitly ask how the plugin builder works.
- Do not narrate internal validation steps.
- Do not expose internal identifiers unless needed for support, debugging, or the user explicitly asks.
- When successful, say what changed and whether it is active.
- When blocked, explain the user-visible impact and what information is needed, without exposing internal implementation details.
- Keep responses short and conversational.

Examples of good user-facing responses:
- “Done — I added the welcome message feature and it’s active now.”
- “Done — I updated the automation so it only runs on a member’s first post.”
- “I couldn’t enable it yet because Storyden rejected one of the requested permissions. I’ve left the existing setup unchanged.”
- “I need one choice before I can finish this: should the message appear publicly in the thread, or privately to the member?”

Security posture:
- Treat plugin installation as an administrator action.
- Keep user data private.
- Do not send data outside the Storyden installation unless the user explicitly requested an integration that requires it.
- Do not add broad permissions when narrow permissions satisfy the request.
- Do not build plugins that bypass Storyden permissions, moderation, authentication, or auditability.
`, currentState, untrustedCommandsInstruction), nil
}

func pluginBuilderInstructionState(ctx adkagent.ReadonlyContext) (string, error) {
	target, ok, err := pluginBuildTargetFromReadonlyContext(ctx)
	if err != nil {
		return "", err
	}
	if !ok {
		return "no plugin selected yet. First decide whether this chat is creating a new plugin or importing an installed plugin to edit.", nil
	}

	manifestID := strings.TrimSpace(target.ManifestID)
	if strings.TrimSpace(target.InstallationID) != "" {
		if manifestID != "" {
			return fmt.Sprintf("working on installed plugin %q. Continue editing, validating, installing, or checking runtime behaviour for this plugin only.", manifestID), nil
		}
		return "working on an installed plugin. Continue editing, validating, installing, or checking runtime behaviour for this plugin only.", nil
	}

	if manifestID != "" {
		return fmt.Sprintf("working on new plugin %q before first install. Continue editing, validating, packaging, and installing this plugin only.", manifestID), nil
	}
	return "working on a selected plugin before first install. Continue editing, validating, packaging, and installing this plugin only.", nil
}

type pluginBuilderToolStage int

const (
	pluginBuilderToolStageUnbound pluginBuilderToolStage = iota
	pluginBuilderToolStageBound
	pluginBuilderToolStageInstalled
)

// var (
//
//	pluginBuilderEntryToolNames = []string{
//		"plugin_workspace_info",
//		"plugin_workspace_create",
//		"plugin_installed_list",
//		"plugin_workspace_import_installation",
//	}
var pluginBuilderAllToolNames = []string{
	"plugin_workspace_info",
	"plugin_workspace_create",
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
}

var pluginBuilderUntrustedCommandToolNames = []string{
	"plugin_run_bash",
}

// pluginBuilderBoundToolNames = []string{
// 	"plugin_workspace_info",
// 	"plugin_manifest_write",
// 	"plugin_file_list",
// 	"plugin_file_read",
// 	"plugin_file_outline",
// 	"plugin_file_write",
// 	"plugin_file_search",
// 	"plugin_file_edit",
// 	"plugin_go_ast",
// 	"plugin_go_fmt",
// 	"plugin_go_vet",
// 	"plugin_go_tidy",
// 	"plugin_go_test",
// 	"plugin_validate",
// 	"plugin_storyden_sdk_events",
// 	"plugin_storyden_sdk_search",
// 	"plugin_go_packages",
// 	"plugin_go_package_symbols",
// 	"plugin_go_symbol_detail",
// 	"plugin_go_symbol_search",
// 	"plugin_storyden_docs",
// 	"plugin_install",
// }
// pluginBuilderInstalledToolNames = append(append([]string{}, pluginBuilderBoundToolNames...), "plugin_logs_read")
// )

type pluginBuilderToolset struct {
	tools       []adktool.Tool
	toolsByName map[string]adktool.Tool
}

func (s *pluginBuilderToolset) Name() string {
	return "storyden_plugin_builder_tools"
}

func (s *pluginBuilderToolset) Tools(ctx adkagent.ReadonlyContext) ([]adktool.Tool, error) {
	// TODO(adk-go#757): restore dynamic stage filtering when ADK re-evaluates
	// Toolset.Tools after state-changing tool calls during the same run.
	// https://github.com/google/adk-go/issues/757
	//
	// stage, err := pluginBuilderToolStageFromContext(ctx)
	// if err != nil {
	// 	return nil, err
	// }
	//
	// switch stage {
	// case pluginBuilderToolStageInstalled:
	// 	return s.toolsForNames(pluginBuilderInstalledToolNames)
	// case pluginBuilderToolStageBound:
	// 	return s.toolsForNames(pluginBuilderBoundToolNames)
	// default:
	// 	return s.toolsForNames(pluginBuilderEntryToolNames)
	// }
	names := pluginBuilderAllToolNames
	if pluginBuilderAllowUntrustedCommandsFromReadonlyContext(ctx) {
		names = append(append([]string{}, names...), pluginBuilderUntrustedCommandToolNames...)
	}
	return s.toolsForNames(names)
}

func (s *pluginBuilderToolset) toolsForNames(names []string) ([]adktool.Tool, error) {
	tools := make([]adktool.Tool, 0, len(names))
	for _, name := range names {
		tool, ok := s.toolsByName[name]
		if !ok {
			return nil, fmt.Errorf("plugin builder tool %q is not registered", name)
		}
		tools = append(tools, tool)
	}
	return tools, nil
}

func pluginBuilderToolStageFromContext(ctx adkagent.ReadonlyContext) (pluginBuilderToolStage, error) {
	target, ok, err := pluginBuildTargetFromReadonlyContext(ctx)
	if err != nil {
		return pluginBuilderToolStageUnbound, err
	}
	if !ok {
		return pluginBuilderToolStageUnbound, nil
	}
	if strings.TrimSpace(target.InstallationID) != "" {
		return pluginBuilderToolStageInstalled, nil
	}
	return pluginBuilderToolStageBound, nil
}

type toolAdder func(adktool.Tool, error) error

func (a *Agent) buildToolset() (adktool.Toolset, error) {
	var out []adktool.Tool
	byName := map[string]adktool.Tool{}
	add := func(t adktool.Tool, err error) error {
		if err != nil {
			return err
		}
		if _, exists := byName[t.Name()]; exists {
			return fmt.Errorf("duplicate plugin builder tool %q", t.Name())
		}
		out = append(out, t)
		byName[t.Name()] = t
		return nil
	}

	registrars := []func(toolAdder) error{
		a.addWorkspaceTools,
		a.addInstalledPluginTools,
		a.addManifestTools,
		a.addFileTools,
		a.addEditTools,
		a.addASTTools,
		a.addGoTools,
		a.addBashTools,
		a.addValidateTools,
		a.addStorydenSDKTools,
		a.addGoDiscoveryTools,
		a.addStorydenDocsTools,
		a.addInstallTools,
		a.addLogTools,
	}

	for _, register := range registrars {
		if err := register(add); err != nil {
			return nil, err
		}
	}

	return &pluginBuilderToolset{tools: out, toolsByName: byName}, nil
}

func registerAgent(registry *agent_registry.Registry, builder *Agent) error {
	return registry.Register(agent_registry.Definition{
		ID:                  agent_registry.PluginBuilderID,
		Name:                "Plugin Builder",
		Description:         AgentDescription,
		RequiresWorkspace:   true,
		AppName:             AgentName,
		AgentName:           AgentName,
		InstructionProvider: builder.instruction,
		ToolsetBuilders: []agent_registry.ToolsetBuilder{
			func(ctx context.Context) (adktool.Toolset, error) {
				return builder.buildToolset()
			},
		},
		Capabilities: []string{
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
		},
	})
}

func slugify(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	var b strings.Builder
	lastDash := false
	for _, r := range s {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			b.WriteRune(r)
			lastDash = false
		case r == '-' || r == '_' || r == ' ' || r == '.':
			if b.Len() > 0 && !lastDash {
				b.WriteByte('-')
				lastDash = true
			}
		}
	}
	out := strings.Trim(b.String(), "-")
	if len(out) > 80 {
		out = strings.TrimRight(out[:80], "-")
	}
	return out
}

func titleize(s string) string {
	parts := strings.Split(slugify(s), "-")
	for i, p := range parts {
		if p == "" {
			continue
		}
		parts[i] = strings.ToUpper(p[:1]) + p[1:]
	}
	return strings.Join(parts, " ")
}
