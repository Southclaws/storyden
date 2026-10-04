# Storyden CLI

`sd` uses [OpenCLI](https://github.com/opencli-dev/opencli) to generate its Cobra command tree from [`opencli.yaml`](opencli.yaml).

The specification owns commands, aliases, usage, help, arguments, flags, defaults, choices, mutually exclusive groups, and typed handler parameters. Generated files live in `internal/cligen`; do not edit them directly.

Handlers in `internal/commands` implement operations using the generated parameter types and supplied input/output streams. HTTP response types remain generated from Storyden's OpenAPI contract. The existing renderers preserve full JSON responses, readable YAML projections, tables, and streaming JSONL without converting through a second set of response models.

## Development

1. Edit `opencli.yaml` when changing the CLI interface. Use `trackChanged: true` for flags whose presence matters independently of their value.
2. Run `task generate:cli` (or `go generate ./cmd/sd` from the repository root).
3. Implement the generated handler interface and register its provider in `main.go`.
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

Tests exercise the assembled generated tree, validation before API calls, context selection, registration, content creation, update flag presence, pagination, and response rendering using local HTTP fixtures.

## Automation and help

Use `--context NAME` for each operation and select structured output explicitly:

```sh
sd --context my-bot info --output json
sd --context my-bot node list --all --output jsonl
sd --context my-bot node create --name Notes --slug notes --visibility draft --markdown --content-file notes.md --output json
```

`--output` / `-o` selects the response format where offered. `--columns wide` selects list columns; `--output-file PATH` selects download/package destinations. Commands with a fixed JSON response say so in their help. Batch node mutations return JSONL with per-item `ok`/`error` outcomes and can partially succeed. Keep stdout and stderr separate and check exit status, including when using pipelines.

Treat help as the command's tool description. Every command and group has Markdown `longDescription` guidance in `opencli.yaml`, with shell examples in `examples`:

- Explain when to choose the command and how to discover its inputs (slugs, IDs, marks, parent/child scope).
- State effects, interactive or long-running behavior, response shape, and useful follow-up commands.
- Describe pagination, replacement semantics, partial failure, and recovery where they apply.
- Use valid flags and realistic workflows. Content is HTML unless `--markdown` is supplied. Do not promise remote validation from a dry run.
- Put accepted choices in flag descriptions so users can discover them without provoking an error.

The help renderer assembles descriptions, examples, usage, commands, and flags into one Markdown document. PTYs render it with terminal styling; pipes receive raw Markdown with no ANSI escapes. Update the specification, regenerate, then exercise changed examples against the actual CLI and local HTTP fixtures. Do not maintain snapshots or tests that merely repeat generated help text.
