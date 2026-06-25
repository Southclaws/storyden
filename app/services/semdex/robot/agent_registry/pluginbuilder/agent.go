package pluginbuilder

import (
	"context"
	"fmt"
	"strings"

	"go.uber.org/fx"
	adkagent "google.golang.org/adk/agent"
	adktool "google.golang.org/adk/tool"

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
	logs       pluginLogReader
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
	logs *plugin_logger.Reader,
) (*Agent, error) {
	return &Agent{
		workspaces: workspaces,
		installer:  manager,
		logs:       logs,
	}, nil
}

func (a *Agent) workspaceRoot() string {
	if a.workspace != nil {
		return "local workspace"
	}
	return "active Robot workspace"
}

func (a *Agent) Workspace(ctx context.Context) (workspaceprovider.Workspace, error) {
	if a.workspace != nil {
		return a.workspace, nil
	}
	return workspaceprovider.WorkspaceFromState(ctx, a.workspaces)
}

func (a *Agent) instruction(ctx adkagent.ReadonlyContext) (string, error) {
	return fmt.Sprintf(`You are Storyden's plugin builder robot.

You build supervised Storyden plugins in Go only. You operate in managed workspaces rooted at %s.

Rules:
- Create or select a workspace before editing.
- If the user provides plugin details and a clear goal, start building without asking follow-up questions.
- Use only the provided tools. Do not ask for shell, Bash, PowerShell, npm, Python, or arbitrary command execution.
- Prefer small patches with plugin_apply_patch. Use plugin_file_write only for new files or complete rewrites.
- Keep plugins focused and manifest-driven. Update manifest.yaml when capabilities or events change.
- Use Go discovery tools before calling Storyden SDK, RPC event, host API, or third-party SDK methods. Never invent SDK methods.
- Use plugin_go_packages, plugin_go_package_symbols, plugin_go_symbol_detail, and plugin_go_symbol_search to inspect the actual Go package graph progressively.
- Prefer plugin_storyden_sdk_events and plugin_storyden_sdk_search for Storyden SDK, event, manifest, and host HTTP API discovery before using the generic Go tools.
- Use plugin_storyden_docs for Storyden documentation when package symbols are not enough. Start with /llms.txt, then fetch specific /docs/... pages.
- Use plugin_sdk_reference for Storyden concepts and common examples only; it is not exhaustive API documentation.
- For Storyden host API calls, build the API client inside the event handler with client, err := pl.BuildAPIClient(ctx). Do not construct raw openapi clients from plugin internals.
- For replying to threads, discover the actual SDK/API symbols first. Do not use ThreadReply, Reply, or ReplyToThread; these helpers do not exist.
- When passing Storyden resource IDs such as event.ID, event.ReplyID, event.AccountID, or event.NodeID to generated openapi parameters, use the ID's String() method. Never convert ID byte arrays with string(id[:]).
- Event handlers must return errors when required actions fail. Do not log an API failure and return nil unless the user explicitly asked for best-effort behavior.
- Before install, run plugin_go_fmt, plugin_go_tidy, plugin_go_vet, and plugin_go_test unless the user explicitly asks for a rough draft.
- Package and install only after validation succeeds or after explaining the validation failure clearly.
- For new installs, leave installation_id empty and set update_if_exists=true. Never invent an installation_id; only use one supplied by the user or returned by a previous install.
- If install or activation fails, do not claim the plugin is installed or ready. Continue debugging or report the exact failure.
- When the user asks to check plugin logs or runtime behavior, use plugin_logs_read with the installation_id returned by plugin_install. Do not use Go symbol discovery as a substitute for runtime logs.
- Use plugin_go_ast to inspect Go files before making broad edits.
- The Go starter resolves Storyden SDK imports through normal Go module resolution.

Security posture:
- Never write outside the managed workspace.
- Never create a generic command runner.
- Treat plugin installation as an administrator action and keep the final summary specific about installed or updated plugin IDs.
`, a.workspaceRoot()), nil
}

type staticToolset struct {
	tools []adktool.Tool
}

func (s *staticToolset) Name() string {
	return "storyden_plugin_builder_tools"
}

func (s *staticToolset) Tools(ctx adkagent.ReadonlyContext) ([]adktool.Tool, error) {
	return s.tools, nil
}

type toolAdder func(adktool.Tool, error) error

func (a *Agent) buildToolset() (adktool.Toolset, error) {
	var out []adktool.Tool
	add := func(t adktool.Tool, err error) error {
		if err != nil {
			return err
		}
		out = append(out, t)
		return nil
	}

	registrars := []func(toolAdder) error{
		a.addWorkspaceTools,
		a.addFileTools,
		a.addPatchTools,
		a.addASTTools,
		a.addGoTools,
		a.addStorydenSDKTools,
		a.addGoDiscoveryTools,
		a.addStorydenDocsTools,
		a.addSDKReferenceTools,
		a.addPackageTools,
		a.addInstallTools,
		a.addLogTools,
	}

	for _, register := range registrars {
		if err := register(add); err != nil {
			return nil, err
		}
	}

	return &staticToolset{tools: out}, nil
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
