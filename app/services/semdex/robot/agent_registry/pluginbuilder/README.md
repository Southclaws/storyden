# Plugin Builder Robot

The Plugin Builder Robot is Storyden's built-in agent for creating and maintaining
supervised Storyden plugins from natural language. It is exposed as the built-in
Robot ID `plugin_builder`, so API clients and the UI can treat it like any other
Robot while the runtime selects a preconfigured ADK agent instead of loading a
database-backed Robot definition.

Its job is intentionally narrow: build Go plugins for Storyden, validate them,
package them, and install or update them through the supervised plugin manager.
It is not a general-purpose coding agent and it must not grow generic shell,
desktop, or host-machine administration capabilities.

## Runtime Model

`Agent` declares Plugin Builder's identity, instructions, capabilities, and
toolset, then registers that definition with the Robot agent registry during
service startup. The central Robot runner owns ADK agent construction, session
history, memory, artifact services, callbacks, model selection, workspace
mounting, and SSE delivery.

The registry entry declares that this Robot requires an active workspace, so
runs fail before tool execution when no workspace is mounted. The only special
behavior is the built-in Robot definition and its constrained plugin-building
toolset.

## Workspace Model

Plugin Builder tools operate on the active Robot workspace from ADK session
state. Tool inputs do not include workspace IDs because workspace selection is
execution context, not model intent.

The tools resolve an opaque `workspaceprovider.Workspace` through the provider
registry. Providers hide whether the workspace is backed by the local filesystem,
a remote sandbox, Sprites, or another environment. Tool code should call the
workspace interface for file and command operations rather than inspecting
provider state or local paths.

The workspace abstraction provides:

- file listing, reading, writing, and searching
- generic command execution for typed tools to call
- provider-owned path confinement and platform details

This boundary keeps Plugin Builder portable across local, Windows, Unix, macOS,
and remote sandbox providers.

Each Plugin Builder chat edits one plugin. For a new plugin, the agent creates
starter files in an empty active workspace. For an existing supervised plugin,
the agent lists editable installed plugins, imports the selected plugin archive
into an empty active workspace, and keeps working with that exact installation.
If the user wants to work on a different plugin, they must start a new chat.

## Tool Model

Tools are organized by file as `tool_<name>.go`. Each tool file owns its ADK
registration, input and output schemas, and implementation helpers.

Current tool groups:

- workspace starter and info tools
- installed supervised plugin list/import tools for existing-plugin editing
- workspace-relative file list/read/write/search tools
- patch application
- Go-specific workflows: `gofmt`, `go mod tidy`, `go vet`, and `go test`
- Go AST inspection from raw file bytes
- Storyden SDK event and symbol discovery
- Go package discovery for workspace modules and dependencies
- constrained Storyden documentation lookup
- Storyden concept and gotcha reference lookup
- package build and supervised plugin install/update
- supervised plugin log reading

Tools should expose task-level inputs only. Environment details such as
workspace ID, provider type, filesystem path roots, operating system, or command
runtime are owned by the runner and workspace provider layers.

The toolset is filtered by chat state:

1. Before a plugin is selected, only workspace inspection, new-plugin scaffold,
   installed-plugin listing, and installed-plugin import tools are visible.
2. Once a plugin is selected, editing, discovery, validation, package, and
   install tools become visible; create/import tools are hidden.
3. Once the selected plugin has an installation, runtime log reading is visible.

Runtime guards still enforce the same one-chat-one-plugin rule even if a stale
or malformed state exposes an unexpected path.

### Go Discovery

Plugin Builder learns SDK surfaces from the workspace module graph instead of
relying on hard-coded reference text. Discovery is intentionally progressive:

1. `plugin_go_packages` lists packages by Go package pattern. Use `./...` for
   workspace code or an import-prefix pattern such as
   `github.com/Southclaws/storyden/lib/plugin/...` to discover subpackages.
2. `plugin_go_package_symbols` lists symbols for one exact import path.
3. `plugin_go_symbol_detail` expands one symbol with fields, methods, and
   signatures where available.
4. `plugin_go_symbol_search` searches names and docs by literal substring when
   the agent knows a concept but not the API name.

Every discovery response includes full package import paths so the agent can
drill into subpackages and related imports without loading an entire dependency
graph up front.

Symbol search is not regex or glob search. Queries such as `Event.*Reply` are
invalid; the agent should search for plain terms such as `Event`, `Reply`, or
`reaction`, then inspect matching packages and symbols.

### Storyden SDK Discovery

Storyden SDK discovery is a narrower Go-powered layer over the same package
graph. It exists because Plugin Builder frequently needs a specific Storyden
surface, such as plugin event names, event payload fields, manifest types, or
generated host HTTP API methods. The general Go discovery tools remain useful
for arbitrary dependencies, but their broad responses can include unrelated
imports and symbols.

Use:

- `plugin_storyden_sdk_events` to find valid `manifest.yaml` event names and
  matching `rpc.Event...` payload structs.
- `plugin_storyden_sdk_search` to search known Storyden SDK/API areas:
  `events`, `manifest`, `rpc`, `http_api`, `operations`, or `all`.

These tools are still derived from Go symbols, not hand-maintained prose. They
provide a targeted directory into Storyden's SDK while `plugin_go_packages`,
`plugin_go_package_symbols`, and `plugin_go_symbol_detail` remain available for
deeper inspection or third-party packages.

`plugin_storyden_docs` provides a constrained documentation fetcher for
Storyden-owned documentation. It can read `https://www.storyden.org/llms.txt`
and subpaths under `https://www.storyden.org/docs/`, requesting Markdown/text
responses for agent-friendly output. It is intentionally not a general web
browser. When choosing manifest `access.permissions`, fetch
`/docs/introduction/members/permissions` and treat that page as the permission
vocabulary.

During discovery, Plugin Builder should decide whether the requested behavior
can be implemented from delivered event payloads alone, or whether it needs to
call Storyden host HTTP APIs to read or write installation data. If it needs
host API calls, the implementation must use `BuildAPIClient` and the manifest
must include `access` with a stable bot account handle, display name, and the
narrow Storyden permission names required by those API operations. Use the
permissions docs to choose those names, and avoid `ADMINISTRATOR` unless the
user explicitly asked for broad administrator capability and no narrower
permission satisfies the requested behavior. Do not add `access` for plugins
that only consume events or use plugin RPC helpers that do not need host HTTP
API credentials.

## Runtime Logs

`plugin_logs_read` reads recent output from the installed supervised plugin
bound to the current workspace. It is the runtime evidence tool: use it when
debugging whether a plugin started, handled an event, emitted an error, or
produced the expected output.

Logs are separate from source and API discovery. If a user asks to check logs,
the agent should call `plugin_logs_read`; it should not substitute Go symbol
search, SDK reference text, or documentation lookup.

## Guardrails

- Plugins are written in Go.
- There is no generic Bash, PowerShell, or arbitrary command tool.
- Command execution is available only through typed tools that construct known
  command specs.
- File access is workspace-relative and confined by the workspace provider.
- AST parsing reads source bytes through the workspace abstraction; it does not
  require local filesystem paths.
- Install operations go through the supervised plugin manager and apply to the
  plugin represented by the active workspace.
- Existing-plugin editing starts by importing an installed supervised plugin;
  subsequent install and log tools continue working with that same plugin.
- The agent must not claim a plugin is installed or active unless the install
  tool returns success.

## Editing Strategy

Plugin Builder should act as the maintainer of one Storyden plugin for the
chat. Its operating loop is:

1. Understand the requested Storyden outcome.
2. Inspect the existing plugin when one exists.
3. Plan the behaviour in terms of Storyden concepts and plugin boundaries.
4. Discover the required Storyden SDK, event, API, manifest, configuration, and
   permission details before relying on them.
5. Implement the behaviour in Go and keep `manifest.yaml` aligned.
6. Validate, fix failures where possible, then install or update only when the
   plugin is ready.
7. Use runtime logs to check startup and user-visible behaviour where possible.
8. Compare the implementation against the original request before reporting
   completion.

The agent should not stop at a proposal when the available tools can complete
the plugin change. If a core requested behaviour remains incomplete, it should
continue working or clearly say what is not ready in user-visible terms.

The default editing loop is:

1. Create starter files for a new plugin, or import an installed supervised
   plugin into an empty workspace for an existing-plugin edit.
2. Search or outline files to locate the relevant behavior, then read a small
   line range around the exact code or text being changed.
3. Update `manifest.yaml` with the structured manifest writer. Edit other
   existing files with exact text replacement, or use complete file writes only
   for new files and full rewrites.
4. Run `plugin_validate` as the holistic readiness check.
5. Use focused Go tools only to repair or recheck specific failed validation
   areas.
6. Install/update the supervised plugin. Installation packages internally; the
   agent does not need a separate package artifact.

File inspection is line-oriented. Reads return the selected range, total line
count, and a content revision. Search returns contextual snippets with the same
revision information, and Go outlines return compact import, type, function,
and method ranges so the agent can decide what to read next without loading
entire files.

Focused edits use exact text replacement. The edit tool reads file bytes through
the workspace abstraction, verifies an optional expected revision, chooses the
replacement location from the old text and optional line hint, and writes the
transformed bytes back through the same abstraction. This works for Go source,
module files, README files, and other text assets without requiring external
tools such as `git` or GNU patch.

Manifest edits are structured. The manifest writer uses the generated plugin
manifest schema as its tool input shape, then decodes and validates through
`rpc.ManifestFromMap` before writing `manifest.yaml`. This keeps field names
such as `configuration_schema` aligned with the runtime manifest contract.

Validation is centralized through `plugin_validate`. It checks manifest schema,
manifest/code consistency, incomplete implementation markers, Go formatting,
dependency tidiness, vet/lint, tests, and package archive validity. The
individual Go tools remain available as repair instruments after a specific
validation check fails.

The active workspace is the source of truth for the plugin. Before broad
changes, Plugin Builder should inspect the existing implementation and preserve
the plugin's established architecture, style, and behavior unless the user
requested otherwise. New work should fit into the plugin as a product: extend
existing behavior where possible, avoid parallel implementations and duplicate
configuration, and keep behavior manifest-driven.

Plugin Builder is allowed to iterate with the user. It does not need to solve
every product detail in one turn. It must not, however, present partial, stubbed,
simulated, dry-run-only, or placeholder behavior as complete. Installed plugins
should not contain TODOs, fake success paths, "would do this" behavior, or
intentionally incomplete implementations. If a requested behavior is not ready,
the agent should say what is ready and what still needs work in user-visible
terms.

Configurable plugins should read current stored configuration during startup and
also handle later configuration updates. A configuration callback alone is not
enough for behaviours that must start after install, update, or restart when the
setting was already present. Missing required configuration should be logged
clearly and disable the dependent behaviour where possible, without making the
plugin look successful when it is idle.

Plugin Builder should treat the plugin as software it owns. When touching nearby
code, it may improve clarity, consolidate duplicated logic, remove obsolete
unreachable code, and fix clearly safe correctness, validation, documentation,
naming, or complexity issues. These improvements must not change unrelated
user-visible behavior.

Because Plugin Builder source is primarily maintained by robots, comments and
lightweight documentation are useful state compression for later runs. Comments
should preserve intent, assumptions, architectural decisions, Storyden concept
mappings, capabilities, non-obvious SDK/API usage, and safety assumptions. They
should explain why rather than narrating obvious syntax, and stale comments
should be removed. Substantial plugins should maintain a concise `README.md`
covering purpose, major behavior, affected Storyden concepts, user-visible
changes, external integrations, and safety or moderation assumptions.

## Runtime Logging

Plugin Builder should add useful runtime logging with the Go standard library
logger. Logs are for future robot maintenance and for non-technical user
support, so they should describe what the plugin is doing in product terms.

Every failure path should log what failed before returning the error, unless the
caller already logs the same failure with better context. Important lifecycle and
behavior points should have concise info logs: startup, configuration decisions,
event receipt, skipped actions, external calls, successful user-visible actions,
and shutdown when relevant.

Logs must not leak secrets, tokens, credentials, or private user content. The
goal is a healthy amount of "this is what's happening" signal, not noisy branch
tracing.

## Package Boundaries

- `agent.go` owns the built-in ADK agent construction, registration, runtime
  dependencies, and workspace resolution.
- `tool_*.go` files own the LLM-facing tool contracts and behavior.
- Workspace provider implementations own filesystem, command, path, and platform
  details.

Keep new functionality on the provider-backed workspace boundary unless it is
truly Plugin Builder-specific behavior.
