# Asynchronous Robot delegation

Status: deferred

## Context

Robot delegation currently runs synchronously inside the browser-initiated chat
request. Denbot invokes a database-backed Robot through ADK `ModeSingleTurn`,
the specialist completes inside that invocation, and its result returns to
Denbot for synthesis.

ADK v2.2 does not propagate a single-turn delegation's parent function-call ID
as its isolation scope. Storyden temporarily infers that scope from the most
recent unresolved Robot function call so delegated events can be attributed and
grouped in persisted history and the UI.

This inference relies on sequential execution. It is ambiguous as soon as two
calls to the same Robot are active concurrently and is not a foundation for
asynchronous tasks.

## Current decision

- Keep synchronous specialist delegation on `ModeSingleTurn` for the current
  Robots milestone.
- Keep `inferDelegationIsolationScope` only as a temporary compatibility shim.
- Do not expand the inference algorithm to support parallel or asynchronous
  execution.
- Remove the shim when ADK propagates single-turn isolation scopes, or when
  Storyden replaces single-turn delegation with an explicitly scoped dispatch
  path.

## What ADK task mode provides

`ModeTask` provides useful execution semantics for work that may span multiple
agent steps or user turns:

- the parent function-call ID is used as a stable run ID and isolation scope;
- the task owns a private conversational history within the shared session;
- an unresolved task can pause and be resumed on a later runner invocation;
- the task explicitly completes through `finish_task`;
- an output schema can validate the completion payload before it returns to the
  coordinator;
- long-running tool interruptions and human confirmation can suspend the task
  without closing its delegation call.

Task mode is resumable orchestration, not a background job runtime. ADK still
runs a task when Storyden invokes the runner. It does not schedule detached work,
persist an application-level queue, wake an idle browser, retry abandoned jobs,
or provide a bidirectional client transport.

## Required Storyden architecture

Before asynchronous delegation is exposed, introduce an application-level task
resource keyed by a stable task ID rather than reconstructing ownership from
session transcript order. It should record at least:

- parent Robot session and originating message/function-call IDs;
- delegated Robot identity and immutable execution input;
- lifecycle state such as queued, running, waiting, completed, failed and
  cancelled;
- ADK run/isolation scope and resume information;
- timestamps, progress and terminal result or failure;
- delivery cursor or event sequence for reconnecting clients.

Execution should be owned by a durable worker/scheduler rather than the HTTP
request lifecycle. A bidirectional or server-push transport can then subscribe
the browser to task and session events, while persisted event sequences allow
replay after disconnects. The transport must be an observation and command
channel, not the owner of task execution.

Timers should enqueue or resume a task by its durable ID. They must not scan chat
history to guess which delegation is active.

## ADK integration questions

- Use `ModeTask` where a specialist must ask follow-up questions, wait for human
  input, survive across invocations, or return validated structured output.
- Decide whether fire-and-forget work is represented by the same task agent
  behind a Storyden scheduler or by an explicit workflow node controlled by the
  task service.
- Verify concurrent dispatch of the same configured Robot before enabling
  parallel fan-out; ADK Go currently has an open shared-agent state race for
  parallel `single_turn` calls.
- Define cancellation and deadline propagation through Storyden workers, ADK
  contexts and tools.
- Define whether a resumed task uses the Robot definition captured at creation
  or the latest edited definition.
- Keep task progress distinct from model reasoning, tool calls and final chat
  messages in both persistence and UI projection.

## Removal criteria for the temporary shim

Delete `inferDelegationIsolationScope` and the split live/stored event handling
when all delegated runs receive their scope from an authoritative execution ID
before the first child event is produced. Add coverage for two simultaneous
delegations to the same Robot to prove attribution does not depend on event
order.
