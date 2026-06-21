# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Hard Rules

- ALMOST NEVER write comments. We're senior engineers here, not learners.
- NEVER run the backend or frontend manually. The human is already doing this.
- ALWAYS test backend changes by writing either unit tests or end-to-end tests.

## Project Overview

Storyden is a modern community platform combining forum, wiki, and community hub features. It's built with Go for the backend and Next.js for the frontend, with a focus on modern security, deployment, and intelligence features.

## Architecture

Storyden follows a strict zero service dependencies production-ready deployment with additional services as optional enhancements. These include the database (PostgreSQL or CockroachDB instead of SQLite), caching (Redis instead of in-process memory), pub-sub (RabbitMQ instead of in-process channels), filesystem (S3 instead of local disk), etc.

### Core Components

- **Backend (Go)**: Rough "clean architecture" not strict or silly levels of adherence though
  - `cmd/backend/main.go`: Main application entry point
  - `app/`: Core application logic divided into resources, services, and transports
    - `app/resources`: Persistence and data structures/data models. NEVER imports from services or transports.
    - `app/services`: Logic layer separated into small simple Go packages. NEVER imports from transports, may import from other services and resources.
    - `app/transports`: Public interface to the outside world via protocols such as HTTP. Imports from services and resources.
  - `internal/`: Infrastructure components (database, pubsub, config, etc.) NEVER imports anything from `./app/` completely isolated from business logic.

- **Frontend (Next.js)**: Located in `web/` directory
  - Built with Next.js 15, React 19, TypeScript
  - Uses Park UI components and Panda CSS with **strict design tokens** (see Design System section below)
  - API client generated from OpenAPI specification

- **Project website**: Located in `home/` directory
  - Built with Next.js and Fumadocs
  - Some documentation is automatically generated alongside code

### Data Layer

- **Database**: SQLite or PostgreSQL/CockroachDB with Ent ORM for schema management
- **Caching**: Optional Redis for session storage, rate limiting and caching
- **Search**: Optional vector databases (Pinecone) for semantic search and recommendations
- **Storage**: Filesystem or S3-compatible object storage for files/assets
- **Message Queue**: Optional RabbitMQ for pubsub

### Key Patterns

- **Simple 3-layer Architecture**: Separation of concerns with resources, services, and transports
- **Dependency Injection**: Using Uber FX for dependency management
- **Code Generation**: Heavy use of code generation for database schema, API clients, Go enums
- **Event-Driven**: Uses Watermill for pub/sub messaging

## Development Commands

### Commands you MUST NEVER RUN

```bash
pnpm dev
go run ./cmd/backend
task release
```

Exception: the isolated Playwright harness tasks are allowed (`task test:e2e`, `task test:e2e:fresh`, `task test:e2e:ci`). Those tasks may run `go run ./cmd/e2etest` and `pnpm start` internally because they create and tear down their own test-only backend/frontend pair.

### Full-stack

```bash
# Run all code-generation
task generate
```

### Backend (Go)

```bash
# Check code
go vet ./...

# Run tests
go test ./...

# Run specific test package
go test ./tests/account/...

# Seed database with test data
go run ./cmd/seed

# Generate database bindings (Ent)
task generate:db
```

### Frontend (Next.js)

```bash
# In web/ directory
cd web

# Install dependencies
pnpm install

# Check types
pnpm tsc --noEmit

# Lint code
pnpm lint
```

### Documentation Site

```bash
# In home/ directory
cd home

# Install dependencies
pnpm install

# Check types
pnpm tsc --noEmit

# Lint code
pnpm lint
```

## Testing

### Go Tests

- End-to-end tests in `tests/` directory
- Unit tests alongside source code (e.g., `*_test.go`)
- Run all tests: `go test ./...`
- Run specific test: `go test -run TestName ./path/to/package`
- Run tests with verbose output: `go test -v ./...`

#### Test Data

- Tests use temporary SQLite databases for isolation
- Test database files created in `tests/*/data/` directories
- Each test gets a unique timestamped database
- You can inspect per-test databases after a test run with `sqlite3 ./path/to/data.db`

### Frontend Tests

- Edit .spec.ts files in `./web/tests`
- Run full Playwright test suite: `task test:e2e:ci`
- Run a specific Playwright test: `task test:e2e:ci -- mytest.spec.ts`
- Always build shared helpers over duplicate copy-paste

## Code Generation

The project heavily relies on code generation. Before making changes to:

- HTTP endpoints
- Database tables and columns
- Enumerated types in Golang code

You must first edit the source of truth and generate the code.

Run `task generate` to regenerate everything.

**Note**: Panda CSS changes (design tokens, recipes, patterns) in `web/panda.config.ts` do NOT require running `task generate`. Panda CSS generates its output automatically.

## Design System (Panda CSS)

The frontend uses Panda CSS with **STRICT TOKEN ENFORCEMENT**. This is a mature project with a comprehensive design system.

### Critical Rules

- **NEVER use raw CSS values** (no `16rem`, `1px solid`, `#3b82f6`, etc.)
- **ALWAYS use design tokens** for all styling
- The config enforces `strictTokens: true` and `strictPropertyValues: true` - builds will fail with raw values
- When writing CSS modules, use CSS variables like `var(--colors-border-subtle)`, `var(--spacing-6)`, etc.
- GUIDELINES NOT RULES: Sometimes you MAY have to break these a little to get things done, but that's rare.

### Common Design Tokens

**Spacing**: Use token names like `4`, `6`, `8` or CSS vars like `var(--spacing-4)`

```tsx
// Good (TSX)
<Box p="4" gap="6" />

// Good (CSS)
padding: var(--spacing-4);
gap: var(--spacing-6);

// Bad
<Box p="16px" />
padding: 16px;
```

**Colors**: Use semantic tokens like `border.subtle`, `bg.canvas`, `fg.muted`

```tsx
// Good (TSX)
<Box borderColor="border.subtle" bg="bg.canvas" />

// Good (CSS)
border-color: var(--colors-border-subtle);
background: var(--colors-bg-canvas);

// Bad
<Box borderColor="#e5e7eb" />
border-color: #e5e7eb;
```

**Borders**: Use `borderWidth="thin"` or CSS vars like `var(--borders-thin)`

```tsx
// Good (TSX)
<Box borderWidth="thin" borderColor="border.subtle" />

// Good (CSS)
border: var(--borders-thin) solid var(--colors-border-subtle);

// Bad
<Box borderWidth="1px" />
border: 1px solid #ccc;
```

**Sizes**: Use token names like `64`, `72` or CSS vars like `var(--sizes-64)`

```tsx
// Good (TSX)
<Box w="64" maxW="breakpoint-2xl" />

// Good (CSS)
width: var(--sizes-64);
max-width: var(--sizes-breakpoint-2xl);

// Bad
<Box w="256px" />
width: 256px;
```

**Z-Index**: Use semantic tokens like `docked`, `overlay`, `popover`

```tsx
// Good (TSX)
<Box zIndex="docked" />

// Good (CSS)
z-index: var(--z-index-docked);

// Bad
z-index: 50;
```

### Finding Available Tokens

- Check `web/panda.config.ts` for the config
- Look at existing components in `web/src/components/` for examples
- Common patterns: `border.subtle`, `bg.canvas`, `fg.muted`, `spacing-{n}`, `sizes-{n}`
- Use Panda's generated types for autocomplete in TSX

## Environment Configuration

The application uses environment variables for configuration. The following file contains all configuration parameters: `./internal/config/config.yaml`

With no environment variables set, running the backend will provide a set of sensible defaults for production.
