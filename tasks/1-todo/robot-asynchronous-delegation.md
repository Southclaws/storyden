# Autonomous Robot invocations

Status: foundation implemented; product decisions remain

## Implemented execution model

`session_coordinator.Enqueue` is the common internal invocation boundary. The
HTTP chat binding, plugin `robot_run`, scheduled work, and delegated specialist
work all ultimately submit durable session inputs through the coordinator.
Session execution remains serial: one replica claims a session lease, performs
one turn, publishes its ordered events, and then signals the next eligible
input.

Plugin invocation is already a non-browser producer, but its RPC method is
synchronous because the caller expects a collected structured result. The
coordinator run wrapper exists for that request/result contract; it is not
needed for browser compatibility and it does not own execution after enqueue.

Scheduled inputs use `not_before`. They remain queued until due and are found by
the normal runnable-session reconciliation pass, so process restarts and
multi-replica deployments do not lose the wake-up. `check_back_later` is the
first producer: it returns a pending asynchronous tool result, stores a due
internal input, and resumes the root Robot in a later turn. The accepting
worker schedules a best-effort local wake-up; the database-backed runnable
session reconciliation pass is the crash-safe fallback.

Specialist delegation uses the same pattern. An asynchronous `robot_<id>` tool:

1. enqueues an unattended specialist input with a stable ID derived from the
   session and function-call ID;
2. immediately returns `pending` to the coordinator Robot;
3. runs the specialist in a separately leased turn;
4. requires the specialist to finish through `robot_run_finish`;
5. enqueues that structured result as an authoritative internal input tied to
   the original function-call ID;
6. starts another root-Robot turn to synthesize the result.

The coordinator assigns the original function-call ID as the specialist's
isolation scope before the first child event is stored. Attribution no longer
depends on transcript order or ADK single-turn delegation behavior.

Manual cancellation is also a coordinator command. The database records a
cancellation request, pubsub wakes the replica that owns the active lease, and
the worker emits a durable `turn_cancelled` event. A database check in the lease
heartbeat covers a missed pubsub notification. Cancelling one turn leaves later
queued inputs intact.

## Decisions still open

### Invocation authority

An autonomous input currently retains the initiating account ID, root Robot,
session, workspace request, source, invocation context, and a stable input ID.
At execution time the worker reloads that account and its current roles; missing
accounts fail the turn rather than inheriting the worker's internal context.
Before adding external schedules or resource listeners, decide whether that
current-authority policy is the final contract or whether some authority should
be snapshotted:

- whether deleted, banned, or permission-reduced accounts invalidate pending
  work;
- whether a scheduled run uses current membership and tool permissions;
- whether integrations act as an account, a service principal, or a distinct
  actor type;
- what attribution a multiplayer session displays for system-originated work.
- whether every member who can view a shared session may cancel its active turn,
  or cancellation requires a narrower session role. The current endpoint uses
  the existing session-view permission.

### Definition and environment drift

Decide whether delayed work uses the Robot definition, model, Toolsets,
workspace mount, and invocation context captured at enqueue time or their latest
values at execution time. The current implementation stores the requested
Robot and workspace identity but resolves the Robot definition when it runs.

### Scheduling contract

`check_back_later` is deliberately small. Product limits still need decisions:

- maximum pending timers per account/session and rate limits;
- required wake-up precision and how much scheduler load is acceptable;
- retention and expiration for very old pending work;
- editing, listing, and cancelling scheduled inputs before they become turns;
- behavior when the target session or Robot is deleted.

### Delegation lifecycle

The first implementation serializes all work in a session. Before parallel
fan-out, define:

- whether sibling delegations may execute concurrently and how they consume a
  shared workspace safely;
- whether cancelling the root request cascades to already-enqueued specialist
  work;
- retry policy after process loss or external side effects;
- progress events distinct from model reasoning and tool calls;
- UI and API inspection of pending/failed specialist tasks;
- whether a specialist may itself delegate or schedule follow-up work.

ADK task mode remains useful for a specialist that must suspend and resume its
private workflow, but it is not the durable scheduler. Storyden's input queue,
leases, event log, and pubsub remain the ownership boundary.

### Plugin result delivery

Plugin `robot_run` remains synchronous even though its execution is detached
from the RPC context after enqueue. A future asynchronous plugin contract needs
a durable invocation ID plus polling, webhook, pubsub, or stream-based result
delivery semantics. That is a separate API decision from the session runtime.
