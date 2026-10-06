# Storyden CLI

`sd` uses [OpenCLI](https://github.com/opencli-dev/opencli) to generate its Cobra command tree from [`opencli.yaml`](opencli.yaml).

The specification owns commands, aliases, usage, help, arguments, flags, defaults, choices, mutually exclusive groups, and typed handler parameters. Generated files live in `internal/cligen`; do not edit them directly.

Handlers in `internal/commands` implement operations using the generated parameter types and supplied input/output streams. HTTP response types remain generated from Storyden's OpenAPI contract. The existing renderers preserve full JSON responses, readable YAML projections, tables, and streaming JSONL without converting through a second set of response models.

Most management commands use `internal/apiexec`, which preserves the full API response and defaults to JSON. An OpenCLI `operationId` of `api<OpenAPIOperationID>` binds the command to that operation in `api/openapi.yaml`. Path arguments use uppercase API parameter names; query flags replace underscores with hyphens and use `trackChanged: true`. `internal/apigen` validates those mappings against OpenAPI and generates typed handler adapters, an Fx `Build()` provider, and completion input metadata. It does not construct Cobra commands. Both generators run through `go generate ./cmd/sd`; `task check:cli` checks both outputs.

The API executor validates parameters and JSON request bodies before authentication, reuses the CLI's credential refresh and rate-limit handling, and makes one request. `--data` supplies inline JSON; `--file -` reads stdin. It preserves explicit false, null, empty strings, and large JSON numbers without decoding and re-encoding the outgoing body or JSON response. `--output raw` supports text and SSE; `--output-file` saves successful response bytes. Empty JSON responses print `null`, and HEAD requests return status and headers. Mutations run immediately. Files can be partially written if a download is interrupted.

## Available workflows

| Area | Commands |
| --- | --- |
| Instance administration | `admin settings`, `admin theme`, `admin icon`, `admin banner`, `admin audit`, `admin email` |
| Identity and security | `account`, `admin accounts`, `admin access-keys`, `admin oauth`, `role`, `invitation` |
| Moderation and participation | `admin moderation`, `report`, `notification`, `profile`, `post` |
| Content | `page`, `page versions`, `thread`, `category`, `tag`, `collection`, `link`, `event`, `asset` |
| Robots | `robot chat`, `robot sessions`, `robot toolsets`, `robot providers`, `robot models`, `robot mcp`, `robot workspaces` |
| Scheduled and event automation | `trail`, `trail preview`, `trail runs` |
| Plugins | `plugin`, `plugin configuration`, `plugin manifest`, `plugin dev` |
| API discovery and advanced use | `api operations`, `api schema`, `api request` |

The new API-backed commands expose one response at a time. Use the available page/cursor flags and response metadata to continue. They do not offer automatic `--all` traversal. `thread page` provides reply pagination while the existing `thread get` retains its current interface. `config` refers to local CLI configuration; `admin settings` manages stored server settings. Deployment environment variables are not writable through the settings API.

Use `sd api schema OPERATION_ID` to obtain a self-contained OpenAPI excerpt with complete nested request/response types, required fields, enums, and accepted content types. It works offline and describes the API contract compiled into this CLI. `api request` provides access to every compiled operation, including those without a named command. Browser-dependent OAuth and WebAuthn ceremonies still require their normal client/browser interaction; a generic HTTP request does not replace that interaction.

```sh
sd admin settings get --output json
sd admin settings update --data '{"title":"Our community"}'
sd category create --data '{"name":"General","description":"Welcome","colour":"#336699"}'
sd robot list --output json
sd robot chat 'Summarise the latest discussions' --output json
sd robot sessions events SESSION_ID --live sse --output raw --timeout 0s
sd api schema TrailCreate
sd trail create --file trail.json
sd trail runs start TRAIL_ID
sd trail runs list TRAIL_ID --output json
```

Robot chat acknowledges queued input; it does not wait for a completed answer. Read session events/history to observe results. Disconnecting an event stream does not cancel the Robot. `trail update` replaces the mutable definition, so supply name, status, trigger, actions, and any description to retain. Archiving preserves history and does not cancel actions already running.

## Shell completions

Carapace completes commands and flags, output formats, API enums and boolean values,
local contexts, files, directories, and comma-separated filters. Instance lookups
suggest pages, threads, profiles/accounts, categories, tags, collections, events,
roles, invitations, notifications, reports, links, plugins, Robots, sessions,
Toolsets, providers, workspaces, MCP servers, Trails, and administrative resource IDs.
Page versions/assets, moderation warnings/notes, and Trail runs/actions use the
parent arguments already entered. `api request` completes operation IDs, parameter
names, enum values, and supported resource identifiers after `--path key=` or
`--query key=`. Property schema tokens complete types and sort directions.

Generate setup for Bash, Zsh, Fish, Nushell, PowerShell, or the other shells listed
by `sd completion --help`:

```sh
source <(sd completion zsh)
# Fish: sd completion fish | source
# Nushell: sd completion nushell
```

For Nushell, the emitted snippet defines `$sd_completer`; dispatch `sd` spans to
that closure from your external completer. If you already use the `carapace`
multi-command completer, run `sd completion carapace` and save its output as
`sd.yaml` in the specs directory reported by `carapace --help`. This uses the
[Carapace bridge](https://carapace-sh.github.io/carapace-bin/spec/bridge.html#carapace)
and delegates to this binary, so the bridge needs no regeneration when commands
change. The `sd` executable must be on `PATH`.

Remote suggestions honor `--context` anywhere in the command being completed and
reuse normal credential refresh without changing the saved default context.
Each lookup has a 1.5-second network timeout and reads at most 2 MiB, returning up
to 100 suggestions from one response; paginated lists are not exhausted. Network,
authentication, and response errors produce no suggestions. Lookups only call
explicitly selected GET operations; completing a mutation never executes it.
Free-form text, JSON bodies, secrets, and IDs without a supported discovery source
remain manual inputs. No response cache is shared between identities.

`internal/apigen` derives input names, choices, and list behavior from OpenCLI and
API query enums/content types from OpenAPI. `internal/completion` supplies the
resource discovery policy and Carapace actions. Regenerate through the normal
`go generate ./cmd/sd` workflow after changing either contract.

## Development

1. Edit `opencli.yaml` when changing the CLI interface. Use `trackChanged: true` for flags whose presence matters independently of their value.
2. Run `task generate:cli` (or `go generate ./cmd/sd` from the repository root).
3. Implement the generated handler interface and register its provider in `Build()` in `provider.go`.
4. Run `task check:cli` to check generated output and run the CLI tests.

OpenCLI is declared in the `tool` section of `go.mod`, which pins its version. The `go:generate go tool opencli` directive and `check:cli` use that shared dependency. Upgrade it with `go get -tool github.com/opencli-dev/opencli/tools/opencli@VERSION`, regenerate, and review the resulting command behavior. `task check:cli` verifies that generated code matches the specification.

The root applies the generated global `--context` flag before invoking handlers. Each execution builds its own command tree through Fx. Formatting and terminal integration are configured at the root; command help comes from the specification.

## Validation

```sh
task check:cli
go test ./cmd/sd/...
go vet ./cmd/sd/...
go build -o /tmp/sd ./cmd/sd
```

Tests exercise generated command handlers, validation before API calls, context selection, registration, content creation, update flag presence, pagination, and response rendering using local HTTP fixtures.

## Automation and help

Use `--context NAME` for each operation and select structured output explicitly:

```sh
sd --context my-bot info --output json
sd --context my-bot page list --all --output jsonl
sd --context my-bot page create --name Notes --slug notes --visibility draft --markdown --content-file notes.md --output json
```

`--output` / `-o` selects the response format where offered. `--columns wide` selects list columns; `--output-file PATH` selects download/package destinations. Commands with a fixed JSON response say so in their help. Batch page mutations return JSONL with per-item `ok`/`error` outcomes and can partially succeed. Keep stdout and stderr separate and check exit status, including when using pipelines.

Treat help as the command's tool description. Every command and group has Markdown `longDescription` guidance in `opencli.yaml`, with leaf examples in `examples` and group examples in the Markdown description:

- Explain when to choose the command and how to discover its inputs (slugs, IDs, marks, parent/child scope).
- State effects, interactive or long-running behavior, response shape, and useful follow-up commands.
- Describe pagination, replacement semantics, partial failure, and recovery where they apply.
- Use valid flags and realistic workflows. Content is HTML unless `--markdown` is supplied. Do not promise remote validation from a dry run.
- Put accepted choices in flag descriptions so users can discover them without provoking an error.

The help renderer assembles descriptions, examples, usage, commands, and flags into one Markdown document. PTYs render it with terminal styling; pipes receive raw Markdown with no ANSI escapes. Update the specification, regenerate, then exercise changed examples against the actual CLI and local HTTP fixtures. Do not maintain snapshots or tests that merely repeat generated help text.

`sd info` shows the selected context and server-reported authentication status. Its JSON output includes an `auth` object with status, credential method when configured, and account ID/handle when authenticated. `auth status` remains the local credential inspection command. Library commands use `page` (and `collection pages`); API operation IDs, paths, and JSON keys retain their API names such as `NodeGet` and `nodes`.
