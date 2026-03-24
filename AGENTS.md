# AGENTS.md

## Purpose

This file is the generic entrypoint for finding project documentation relevant to the current task.
Use this as an index, then open only the docs needed for the work at hand.

## Quick Orientation

- Product and scope:
  - `README.md` for product overview and goals.
  - `web/` is the product application.
  - `home/` is the website/docs app (only use when task targets docs/marketing site).
- Backend high-level direction (Go):
  - `cmd/backend/main.go` is the backend runtime entrypoint.
  - `app/resources/**`, `app/services/**`, `app/transports/**` are the core 3 layers (data, business logic, external interfaces).
  - `internal/**` contains infrastructure and platform integration code.
  - `api/openapi.yaml` and `api/rpc/**` are contract/spec roots for generated transport code.
  - `internal/ent/schema/**` is the DB schema source of truth.
- Local setup and workflow:
  - `Taskfile.yml` for canonical project task commands.
- Architecture and backend:
  - `CLAUDE.md` for high-level architecture map, layering rules, and command guardrails.
  - `docs/architecture/` for deeper architecture notes when needed.
  - `internal/config/config.yaml` for runtime configuration surface.
- Frontend and design system (`web`):
  - `web/panda.config.ts` for Panda setup, tokens, recipes, patterns, and conditions.
  - `web/src/theme/` and `web/src/app/global.css` for semantic tokens and global styling rules.
  - `web/src/app/(dashboard)/**/README.md` for route/domain intent.
- API and SDK:
  - `api/openapi.yaml` as API contract source of truth.
  - `api/rpc/README.md` and `sdk/**/README.md` for transport/SDK usage.

## Backend Task Routing (Go)

Storyden backend is a Go application using Uber Fx dependency injection and a 3-layer app structure.

- Runtime composition and startup:
  - `cmd/backend/main.go` wires `config`, `infrastructure`, `resources`, `services`, and `transports`.
  - `app/resources/resources.go`, `app/services/services.go`, and `app/transports/transports.go` define dependency module composition.
- Layer responsibilities:
  - `app/resources/**`: persistence/data access and query/write repositories.
  - `app/services/**`: business logic and orchestration.
  - `app/transports/**`: external interfaces (HTTP OpenAPI, RPC, MCP).
  - `internal/**`: infrastructure and platform concerns (db, cache, mail, pubsub, vector, config, etc.).
- If your task is API endpoint behavior:
  - Start in `app/transports/http/bindings/**` for handlers/bindings.
  - Check generated types/server in `app/transports/http/openapi/**`.
  - Update API contract in `api/openapi.yaml` (and related `api/common/**` schemas) when request/response surface changes.
- If your task is domain/business logic:
  - Implement changes in `app/services/<domain>/**`.
  - Use `app/resources/<domain>/**` for data access changes needed by that logic.
- If your task is database/schema/persistence:
  - Source of truth is `internal/ent/schema/**`.
  - Regenerated artifacts are under `internal/ent/**`.
- If your task is infra/config wiring:
  - `internal/infrastructure/**` for adapters and integrations.
  - `internal/config/config.yaml` for environment-backed runtime configuration.
- If your task is plugins/RPC:
  - RPC schemas: `api/rpc/**` and `api/plugin.yaml`.
  - Transport/runtime implementation: `app/transports/rpc/**`, `app/services/plugin/**`, `app/resources/plugin/**`.
- If your task is MCP tools:
  - `app/transports/mcp/**` and `app/transports/mcp/tools/**`.

## Backend Testing and Validation

- Unit/integration:
  - `go test ./...` for full suite.
  - Prefer targeted package runs while iterating, then finish with broader test coverage.
- End-to-end:
  - `task test:e2e`, `task test:e2e:fresh`, `task test:e2e:ci` (see `Taskfile.yml`).
  - E2E harness entrypoint: `cmd/e2etest/main.go`.
- Test fixtures/data:
  - Most backend tests live under `tests/**` with per-domain data in sibling `data/` directories.

## Runtime Guardrails

- Never start backend or frontend runtime processes from the agent.
- Do not run Go app entrypoints (for example `go run ./cmd/backend`).
- Do not run Next.js app servers (for example `yarn dev`, `yarn start`, `next dev`, `next start`).
- Assume the user already has required services running; use tests, linters, codegen, and static checks instead.

## Codegen Routing

- Code generation is a core workflow in this repository.
- When you edit a source-of-truth spec/schema, always run the matching generator before finishing.
- Prefer targeted generation while iterating; run broader generation when cross-surface contracts change.

- Common source-of-truth -> command mapping:
  - `internal/ent/schema/**` -> `task generate:db` or `go generate ./internal/ent`
  - `api/openapi.yaml`, `api/common/**`, `api/plugin.yaml`, `api/rpc/**` -> `task generate:openapi` or `go generate ./api`
  - OpenAPI frontend client only -> `task generate:openapi:frontend` (runs `yarn openapi` in `web`)
  - OpenAPI docs site only -> `task generate:openapi:docs` (runs `yarn openapi` in `home`)
  - Full regen when unsure -> `task generate`

- Go generators are typically invoked with `go generate ./<path>`.
- Generator entrypoints/tooling live in `api/generate.go` and `internal/tools/**`.
- After generation, run relevant tests/typechecks for changed surfaces (backend and/or frontend).

## Design Context

For UI/UX/visual/frontend tasks, use `.impeccable.md` as the source of truth for design direction and principles.
For non-frontend tasks, treat `.impeccable.md` as optional context.
