# Plugin Builder Eval Playbook

This playbook describes the manual eval loop for the built-in `plugin_builder`
Robot against the real Storyden backend, Robot SSE transport, and workspace
provider. Use it when changing Plugin Builder instructions, tool descriptions,
tool schemas, tool result shaping, validation, install/update behavior, or
workspace provisioning.

The goal is not to prove one prompt works once. The goal is to find places
where the Robot is forced to guess, then move that burden into tools,
validation, or result shaping.

## What This Eval Covers

A good Plugin Builder eval exercises the full managed flow:

1. Start a backend from the current tree.
2. Authenticate as a local admin (ask the user for a cookie.)
3. Attach a Robot workspace, usually Sprites.
4. Send a realistic `/sse/chat` request to `plugin_builder`.
5. Capture the full SSE stream and backend logs.
6. Inspect the session transcript, tool calls, generated workspace, and
   installed plugin state.
7. Patch tool behavior or validation when the Robot gets confused.
8. Add regression tests for the exact confusion.
9. Restart the backend and rerun the same eval prompt.

Prefer eval prompts that require several Storyden concepts at once: events,
configuration, host API access, `RunRobot`, third-party libraries, validation,
installation, and runtime logs.

## Preconditions

Run from the Storyden repository root:

```sh
cd /Users/southclaws/work/storyden
```

Expected local state:

- `data/data.db` exists and has seeded local accounts.
- A Robot workspace exists, or one can be created through the UI/API.
- The selected workspace provider is able to run Go commands.
- The configured model/provider has enough credit for a multi-minute run.
- Network access is available if the workspace needs module downloads.

Do not read Sprite auth/config secrets while debugging. Keep live Sprite
inspection to metadata, file lists, command availability, and workspace source
unless there is explicit approval to inspect sensitive files.

## Start The Backend

Start the backend from the current working tree in a terminal you can stop:

```sh
go run ./cmd/backend
```

Wait for:

```text
storyden http server starting
```

The server should listen on `http://localhost:8000`. If you change tool code
after starting the backend, stop and restart the backend before rerunning the
eval. Otherwise the Robot will still use the old tool implementation.

## Authenticate

For local seed data, signing in as `odin` is usually enough:

```sh
COOKIE=/tmp/storyden-pluginbuilder-cookies.txt
rm -f "$COOKIE"

curl -sS -i -c "$COOKIE" \
  -H 'Content-Type: application/json' \
  --data '{"identifier":"odin","token":"password"}' \
  http://localhost:8000/api/auth/password/signin
```

Copy the `storyden-session` value from the `Set-Cookie` header. In local HTTP
testing the cookie may be marked `Secure`, so `curl -b "$COOKIE"` may not send
it back to `http://localhost`. Passing the cookie explicitly is usually simpler:

```sh
SESSION_COOKIE='storyden-session=<value-from-signin>'
```

If `/sse/chat` returns `Unauthorized`, the cookie was not sent or the session is
invalid.

## Choose A Workspace

List available Robot workspaces:

```sh
curl -sS \
  -H "Cookie: $SESSION_COOKIE" \
  http://localhost:8000/api/robots/workspaces | jq .
```

Pick a workspace ID. A Sprites workspace is preferred for production-like
coverage:

```sh
WORKSPACE_ID=d92op7do2dtjpvta6jh0
```

If `/sse/chat` returns:

```text
Plugin Builder requires an active Robot workspace
```

then the request body did not include a valid workspace mount, or the backend
could not mount it.

## Generate A Session ID

The SSE endpoint requires a valid XID session ID.

```sh
tmp=$(mktemp -d)
cat > "$tmp/main.go" <<'GO'
package main

import (
	"fmt"

	"github.com/rs/xid"
)

func main() {
	fmt.Println(xid.New().String())
}
GO

(cd "$tmp" && go mod init xidtmp >/dev/null && go get github.com/rs/xid >/dev/null && go run .)
rm -rf "$tmp"
```

Set the output:

```sh
CHAT_SESSION_ID=<xid>
```

## Send The Eval Prompt

Use prompts that are specific enough to require real implementation, but not so
over-specified that they hard-code the answer.

Example prompt:

```sh
PROMPT='Build a real Storyden supervised Go plugin called Discord Weekly Moderation Digest. Every time a Storyden report is created, it should collect the report details, ask Storyden robot d94r50jara0c2aq6dpj0 for a concise moderation-friendly triage summary, and post that summary to a configured Discord channel using a bot token. It must include configuration fields for discord_token and discord_channel_id at minimum, plus optional storyden_robot_id, message_prefix, and minimum_report_status. It must not crash if configuration is missing at first boot. It should avoid duplicate Discord sessions when configuration changes. It should use Storyden event, host API, and robot APIs correctly rather than stubbing. It should install and activate when ready.'
```

Create the request body:

```sh
REQUEST=/tmp/pluginbuilder-${CHAT_SESSION_ID}.json
STREAM=/tmp/pluginbuilder-${CHAT_SESSION_ID}.sse

jq -n \
  --arg sid "$CHAT_SESSION_ID" \
  --arg workspace "$WORKSPACE_ID" \
  --arg prompt "$PROMPT" \
  '{
    id: $sid,
    sessionId: $sid,
    robotId: "plugin_builder",
    workspace: { workspace_id: $workspace },
    messages: [
      {
        id: ($sid + "m"),
        role: "user",
        parts: [{ type: "text", text: $prompt }]
      }
    ]
  }' > "$REQUEST"
```

Run the request and keep the stream:

```sh
curl -sS -N --max-time 900 \
  -H 'Content-Type: application/json' \
  -H "Cookie: $SESSION_COOKIE" \
  --data @"$REQUEST" \
  http://localhost:8000/sse/chat | tee "$STREAM"
```

Long initial silence can be normal while the backend mounts or provisions the
workspace. Check backend logs before assuming the run is stuck.

## Inspect The Run

The stream is the first source of truth:

```sh
less "$STREAM"
```

Look for:

- tool inputs and outputs
- validation failures and `next_action`
- whether the Robot uses authoritative discovery tools
- whether it asks the user instead of fixing code it can inspect
- whether it calls `plugin_install`
- whether it calls `plugin_logs_read` after activation
- whether final text matches the actual install/log results

The database transcript is useful when the stream is truncated:

```sh
sqlite3 ./data/data.db \
  "select id, created_at, substr(event_data,1,1600)
   from robot_session_messages
   where session_id='$CHAT_SESSION_ID'
   order by created_at;"
```

Installed plugin state can be checked through the app UI or database. Runtime
logs are best checked through the Robot tool result when available, because it
is the same evidence the Robot saw.

## Optional Sprite Inspection

If the workspace behavior is suspect, use the Sprite CLI to inspect the exact
workspace instance mentioned in backend logs:

```sh
/Users/southclaws/.local/bin/sprite list
/Users/southclaws/.local/bin/sprite -s <sprite-name> exec sh -lc 'pwd; find /workspace -maxdepth 2 -type f | sort | sed -n "1,120p"'
```

For generated plugin source:

```sh
/Users/southclaws/.local/bin/sprite -s <sprite-name> exec sh -lc 'sed -n "1,260p" /workspace/main.go'
```

Avoid inspecting credentials or auth config. For provisioning bugs, prefer
metadata checks such as PATH, command availability, marker file existence, and
workspace files.

## Success Criteria

A passing eval should show:

- `plugin_workspace_create` or `plugin_workspace_import_installation`, depending
  on the task.
- Storyden event/API discovery before implementation when the task touches
  Storyden events, host APIs, permissions, or Robot APIs.
- Manifest written through `plugin_manifest_write`.
- Runtime fields such as command and binary details left to tools.
- Configuration handled as live state: no crash when first boot has no config.
- Real implementation, not placeholders, TODOs, dry-run logic, canned summaries,
  or "would send" behavior.
- `plugin_validate` passes.
- `plugin_install` succeeds and activates when requested.
- `plugin_logs_read` shows either healthy startup or a user-actionable runtime
  state such as waiting for configuration.
- Post-install file edits are rejected unless the user requested another change.

For Robot integration specifically:

- Correct: `pl.RunRobot(ctx, robotID, message)`.
- Incorrect: generated HTTP `RobotChatSSE`, `RobotChatSSEWithResponse`,
  `RobotRunWithResponse`, or manually parsing UI chat streams from a plugin.
- Manifest access includes `USE_ROBOTS` when the plugin calls `RunRobot`.

## Failure Analysis

Classify failures by what should change.

Tool contract failures:

- The Robot cannot discover a real API with natural search terms.
- A tool exposes fields the Robot should not control.
- A validation error is technically correct but points at the wrong next action.
- The tool result omits the next concrete repair step.
- A retry creates a new installed plugin instead of updating the current one.

Validation failures:

- A semantic check rejects a legitimate implementation pattern.
- A semantic check allows stubs, canned summaries, or fake side effects.
- Go compile errors are hidden behind lower-priority semantic warnings.

Prompt/playbook failures:

- The Robot ignores a workflow stage even though tools expose the right path.
- The Robot asks the user to make an implementation choice it can resolve from
  available APIs.
- The Robot reports technical details to the user instead of user-facing status.

Model behavior failures:

- The Robot sees the right tool result and still chooses a forbidden API.
- The Robot repeatedly makes the same invalid edit after specific feedback.
- The Robot stops before install despite validation success and user request.

Prefer fixing tool contracts, schemas, and result shape before adding another
line to the Robot instructions.

## Improve Tools And Tests

For each failure:

1. Reduce it to the smallest concrete confusion.
2. Patch the relevant tool, validation check, or result shape.
3. Add a regression test that uses the same query, error text, or code pattern
   from the transcript.
4. Run focused tests.
5. Restart the backend.
6. Rerun the same eval prompt.

Useful test commands:

```sh
go test ./app/services/semdex/robot/agent_registry/pluginbuilder
go test ./app/services/semdex/robot/agent_registry/...
```

When a tool change affects generated API schemas or source-of-truth specs, run
the matching generator before rerunning tests.

## Known High-Value Checks

These checks caught real issues and should remain part of future eval review:

- Natural SDK search terms such as `robot`, `run`, `robot run`, and areas like
  `http_api` still steer to `pl.RunRobot`.
- Validation `next_action` prioritizes missing Go methods, fields, and types
  before config semantic cleanup.
- Config validation accepts direct raw map access, switch-on-key parsing, and
  JSON marshal/unmarshal into tagged config structs.
- The Robot does not add dummy `_ = raw["field"]` reads just to satisfy a
  validator.
- The Robot does not call `RobotChatSSE` or any UI streaming endpoint from a
  plugin.
- Optional configuration fields are either implemented or omitted. Do not let
  the Robot add "reserved for future use" fields while claiming the request is
  complete.
- Plugin logs that say "waiting for configuration" are a successful startup
  state when required UI config has not been supplied yet.

## Cleanup

Stop the backend process you started for the eval:

```text
Ctrl-C
```

If the eval installed local plugins, clean them up through the UI or normal
plugin management paths when they are no longer useful. Do not manually delete
database rows or plugin package files unless the cleanup itself is the task.
