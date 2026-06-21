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

## Tool Model

Tools are organized by file as `tool_<name>.go`. Each tool file owns its ADK
registration, input and output schemas, and implementation helpers.

Current tool groups:

- workspace starter and info tools
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
graph up front. `plugin_sdk_reference` remains useful for Storyden concepts,
manifest gotchas, and common examples, but it is not exhaustive API
documentation.

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
browser.

## Runtime Logs

`plugin_logs_read` reads recent output from an installed supervised plugin by
installation ID. It is the runtime evidence tool: use it when debugging whether
a plugin started, handled an event, emitted an error, or produced the expected
output.

Logs are separate from source and API discovery. If a user asks to check logs,
the agent should call `plugin_logs_read` with the `installation_id` returned by
`plugin_install`; it should not substitute Go symbol search, SDK reference text,
or documentation lookup.

## Guardrails

- Plugins are written in Go.
- There is no generic Bash, PowerShell, or arbitrary command tool.
- Command execution is available only through typed tools that construct known
  command specs.
- File access is workspace-relative and confined by the workspace provider.
- AST parsing reads source bytes through the workspace abstraction; it does not
  require local filesystem paths.
- Install operations go through the supervised plugin manager.
- The agent must not claim a plugin is installed or active unless the install
  tool returns success.

## Editing Strategy

The default editing loop is:

1. Create starter files or inspect existing workspace state.
2. Read/search files and inspect Go source with AST tools where useful.
3. Edit with patches or complete file writes.
4. Run Go validation tools.
5. Package and install/update the supervised plugin.

Patch application is provider-independent for Go source files. The patch tool
reads file bytes through the workspace abstraction, applies a semantic Go patch
in process, and writes the transformed bytes back through the same abstraction.
Hosted workspaces do not need external tools such as `git` to apply Go source
edits.

## Package Boundaries

- `agent.go` owns the built-in ADK agent construction, registration, runtime
  dependencies, and workspace resolution.
- `tool_*.go` files own the LLM-facing tool contracts and behavior.
- Workspace provider implementations own filesystem, command, path, and platform
  details.

Keep new functionality on the provider-backed workspace boundary unless it is
truly Plugin Builder-specific behavior.
