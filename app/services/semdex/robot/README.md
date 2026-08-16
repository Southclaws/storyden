# Robots Runtime

Storyden Robots is an agent runtime built on Google ADK v2. New
conversations use the built-in **Denbot** unless a custom Robot is selected
explicitly for testing or a focused interaction. Denbot can discover
reusable Toolsets, load capabilities for the current session, search the
specialist Robot catalogue, and delegate bounded tasks without changing the
user's active conversation.

Database-backed Robots are focused specialists and reusable delegation targets,
but remain valid chat entrypoints so their behaviour can be exercised directly.

## Runtime Contract

- `/api/robots/sessions` accepts an optional Robot selector when creating a session and
  defaults to `denbot`. The chosen root Robot is immutable for that session.
- `robot_switch` and active-Robot session state do not exist.
- Custom Robots can also be invoked directly by internal/plugin `robot_run`
  workflows when an explicit specialist run is the intended API operation.
- Denbot and delegated Robots share one ADK session and one mounted workspace.
- Delegated events retain their database Robot attribution. Storyden assigns
  the parent function-call ID as their authoritative isolation scope.
- Branched specialist output is projected as collapsed reasoning in the UI.
  Denbot's synthesis remains the visible assistant answer.

## Agent Topology

`runner.go` resolves the built-in Denbot definition and exposes every available
database-backed Robot as an asynchronous function tool:

- Denbot uses `llmagent.ModeChat` and owns the user conversation.
- Specialists run in separately leased unattended session turns.
- Each specialist tool is named `robot_<database-xid>` and immediately returns
  a pending asynchronous result.
- `robot_search` returns that callable delegation name with the Robot's purpose
  and configured Toolsets.
- A specialist receives the delegated request as its active task alongside the
  shared session history, its playbook, its model, and its assigned direct
  tools and Toolsets. Its instruction limits it to that bounded task. It
  finishes with `robot_run_finish`; that result updates the original
  delegation and starts a later Denbot turn for synthesis.
- Every Robot receives shared execution guidance that prevents unchanged tool
  calls from being repeated when neither their inputs nor underlying state have
  changed. This belongs to the runtime identity instruction rather than any
  individual Toolset, so custom Toolsets and direct tools inherit it too.

Delegation never hands the visible conversation to the specialist. The initial
call returns `pending`, the specialist transcript arrives asynchronously on the
session stream, and Denbot remains the visible coordinator when the completion
starts a later turn.

## Toolsets

A Toolset is a reusable capability package:

- stable ID and human-readable name
- description used for discovery
- optional specialist instruction injected while active
- a set of raw registered tools
- source metadata: `system`, `custom`, or `plugin`

System Toolsets are registered at startup. Custom Toolsets are Ent resources
owned by an account and managed through `/robots/toolsets`. Plugins may publish
Toolset definitions alongside their tools; plugin host registration and
unregistration is atomic for both surfaces.

Robots store both direct tool names and Toolset references. Direct tools support
narrow configurations where assigning an entire Toolset would grant unrelated
capabilities or inject unnecessary guidance. A custom Toolset cannot be deleted
while any Robot references it.

Tool assignment is intentionally disjoint. Adding a Toolset removes direct
tools already provided by that Toolset. Adding a direct tool already provided
by an assigned Toolset is rejected with the owning Toolset's name.

Denbot starts with the Library management, Discussion management, Content
search, and Robot studio Toolsets. Plugin Studio is a prompt-bearing system
Toolset, not a built-in Robot.

### Toolset discovery and loading

Denbot always receives `toolset_search`, `toolset_load`, and `tool_load`.
`toolset_search` returns only Toolset ID, name, and description. The Robot uses
`toolset_get` to inspect the selected bundle's tool IDs and instruction before
loading it. Loading records active Toolset references in ADK session state and
makes both their prompt and tools available on the next model step.

Individual capability discovery is deliberately separate. `tool_search`
returns only tool ID, name, and description; `tool_get` returns the selected
tool's full input/output schema, confirmation requirement, and workspace
precondition. `tool_get` is
inspection-only. `tool_load` is the explicit state-changing operation that
activates one or more inspected tools for the conversation.

Workspace availability is derived from the session's authoritative
`robot_workspace` mount state. Toolsets can declare that the whole bundle
requires a workspace, and individual schema-backed tools can declare
`x-storyden-tool.requires_workspace`. ADK's dynamic `Toolset.Tools` boundary
hides these capabilities when no workspace is mounted. Inspection remains
available so Denbot can explain the precondition, while `toolset_load` and
`tool_load` reject activation without mutating session state.

ADK v2.1 caches `Toolset.Tools` once per flow: `toolProcessor` returns as soon as
`Flow.Tools` is non-nil, while the same Flow performs each model step in an
invocation. To support same-invocation loading without exposing every schema to
the model, Storyden's discovery Toolset:

1. resolves the complete executable tool pool once for ADK dispatch;
2. records which Toolsets own each callable tool name;
3. prunes inactive Toolset schemas from every model request;
4. re-evaluates session state after `toolset_load` or `tool_load`, allowing the
   newly active schemas through on the following model step.

This keeps normal ADK tool execution, RBAC callbacks, confirmation gates, and
client-side tool handling intact while preserving progressive disclosure.

## Tool Registry

Raw tool contracts are schema-driven. Tool schemas live under `api/robots.yaml`
and `api/robots/tools/`; generated bindings live in `lib/mcp`. The runtime
registry in `tools/` owns native, MCP, and plugin tool implementations and maps
stable IDs to callable ADK names.

Capability discovery follows progressive disclosure. Search responses never
dump schemas, Toolset membership, or provenance:

- `tool_search` -> ID, name, description
- `tool_get` -> full individual schema and runtime preconditions
- `tool_load` -> conversation activation
- `toolset_search` -> ID, name, description
- `toolset_get` -> full bundle configuration, tool IDs, instruction, runtime preconditions, editability
- `toolset_load` -> conversation activation

`source` remains management/UI provenance and is excluded from agent-facing
discovery. `editable` remains on `toolset_get` because it determines which
management actions are possible.

When adding a tool:

1. add or update the schema source;
2. run `task generate:mcp`;
3. register the implementation in `tools/` or a plugin provider;
4. place it in an appropriate Toolset, or leave it catalogue-only for custom
   composition;
5. test the model-visible schema and real execution path.

## Sessions and presentation

`session_storage.go` adapts Storyden persistence to ADK's session service. ADK
runs with streaming disabled, so each event is stored whole in the ordered
session event log, together with indexed attribution fields:

- invocation ID
- branch
- isolation scope
- human author, built-in Robot, or database Robot

Model parts marked `Thought` and text emitted on delegated branches are
projected as AI SDK reasoning parts. They are hidden by default behind an
expandable disclosure. Tool calls and Denbot's final prose keep their normal
presentation.

## Workspaces

Workspaces belong to the root session, not to Plugin Studio. A request may
mount a workspace template or an existing workspace instance. The mount is
stored in session state and is therefore available to Denbot, loaded Toolsets,
and delegated specialists.

Workspace providers own confinement and execution policy. Toolsets must use the
workspace abstraction rather than infer host paths. Arbitrary command tools
must continue to respect the workspace's `allow_untrusted_commands` policy.

## Safety boundaries

- Tool permission checks run immediately before execution using the current
  account session.
- Confirmation-required tools block their turn and resume from a later
  client-supplied approval input.
- Unattended runs disable confirmation-only side effects and block tools that
  need live client input.
- Specialist construction failures are logged and skipped so one unavailable
  model or Toolset does not take down Denbot.
- Plugin Toolsets are validated against the tools registered by that plugin
  before they become discoverable.

The current Robot catalogue is instance-visible, matching the existing Robots
API. Custom Toolset mutations are author-controlled. Per-user Robot visibility
is a separate product boundary and should be designed explicitly rather than
inferred from Toolset ownership.

## Validation

Focused runtime coverage lives in:

- `app/services/semdex/robot/toolsets`
- `app/resources/robot/robot_toolset`
- `app/resources/robot/robot_session`
- `tests/robot/chat`
- `tests/plugin/plugin_robot_run_test.go`

The chat integration tests exercise direct-tool execution, same-session
delegation, branch attribution, collapsed specialist output, and
search-load-execute behavior for custom Toolsets.
