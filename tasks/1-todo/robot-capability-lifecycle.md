# Robot capability lifecycle

Status: deferred

## Context

Denbot progressively activates individual Tools and reusable Toolsets by writing
their IDs into session state. Activation is currently additive: a capability
remains active for the rest of the conversation once loaded.

That is a useful first implementation because it keeps discovery predictable,
but long-running conversations can accumulate executable schemas and Toolset
instructions that are no longer relevant. Shared custom Toolsets and Robots are
also live configuration today: edits affect every consumer without versioning or
pinning.

Robot capability references intentionally remain flexible string IDs for now.
Plugin removal or configuration changes may therefore leave unavailable
references. The runtime should warn and continue where possible instead of
requiring immediate reconfiguration or enforcing database foreign keys.

## Current decision

- Persist loaded Tool and Toolset IDs for the conversation.
- Do not add unload or reset semantics during the initial Robots v2 milestone.
- Treat custom Toolset instructions as trusted administrator-authored prompt
  extensions; Robot access is already restricted to operational users.
- Keep Robot capability references as string IDs and surface unavailable
  capabilities as recoverable runtime issues.
- Do not introduce Robot or Toolset versioning yet.

## Deferred design questions

### Capability release

Decide whether capabilities should be removed through explicit `tool_unload` and
`toolset_unload` operations, a reset-to-base operation, turn-scoped leases, or a
combination of these.

The design must define what happens when:

- a directly loaded Tool is also owned by an active Toolset;
- two active Toolsets share the same Tool;
- unloading a Toolset removes its instruction but another Toolset still owns
  some of its Tools;
- a delegated Robot changes shared session capability state;
- an unloaded capability appears in older event history.

### Robot and Toolset versioning

Explore immutable revisions and optional pinning for both Robots and Toolsets.
The design should preserve fast live editing while allowing important Robots to
run against a reviewed revision and making historical session behavior
explainable.

### Persistence shape

Re-evaluate normalized Robot-to-Tool and Robot-to-Toolset assignment tables when
usage queries, auditing, versioning or reverse-reference operations justify the
extra schema. Database constraints are not a goal by themselves; graceful
handling of unavailable plugin capabilities remains required.

### Durable runtime diagnostics

Capability construction failures are currently isolated and exposed through a
structured runtime issue tool. The issue list is reconstructed while the
provider remains broken; it is not a permanent incident ledger. Consider a
durable diagnostics store and administrator UI if operators need to inspect old
failures after a plugin is disabled, repaired or removed.

## Acceptance criteria for future work

- Capability state remains deterministic across user turns and delegation.
- Releasing a capability cannot accidentally remove a Tool still owned by
  another active source.
- Prompt instructions and executable Tools follow the same lifecycle.
- The agent and management UI can explain active, unavailable and released
  capabilities without exposing raw internal state.
- Existing sessions remain operable when a plugin capability disappears.
