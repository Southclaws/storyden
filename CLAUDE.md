# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Hard Rules

- ALMOST NEVER write comments. We're senior engineers here, not learners.
- NEVER run the backend or frontend manually. The human is already doing this.
- ALWAYS test backend changes by writing either unit tests or end-to-end tests.
- ALWAYS test frontend changes by running Playwright MCP against localhost:3000.

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
  - Uses Park UI components and Panda CSS
  - API client generated from OpenAPI specification

- **Project website**: Located in `home/` directory
  - Built with Next.js and Fumadocs
  - Some documentation is automatically generated alongside code

### Data Layer

- **Database**: SQLite or PostgreSQL/CockroachDB with Ent ORM for schema management
- **Caching**: Optional Redis for session storage, rate limiting and caching
- **Search**: Optional vector databases (Pinecone, Weaviate) for semantic search and recommendations
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
yarn dev
go run
task release
```

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
yarn install

# Check types
yarn tsc --noEmit

# Lint code
yarn lint
```

### Documentation Site

```bash
# In home/ directory
cd home

# Install dependencies
yarn install

# Check types
yarn tsc --noEmit

# Lint code
yarn lint
```

## Testing

### Go Tests

- End-to-end tests in `tests/` directory
- Unit tests alongside source code (e.g., `*_test.go`)
- Run all tests: `go test ./...`
- Run specific test: `go test -run TestName ./path/to/package`
- Run tests with verbose output: `go test -v ./...`

### Test Data

- Tests use temporary SQLite databases for isolation
- Test database files created in `tests/*/data/` directories
- Each test gets a unique timestamped database
- You can inspect per-test databases after a test run with `sqlite3 ./path/to/data.db`

## Code Generation

The project heavily relies on code generation. Before making changes to:

- HTTP endpoints
- Database tables and columns
- Enumerated types in Golang code

You must first edit the source of truth and generate the code.

Run `task generate` to regenerate everything.

**Note**: Panda CSS changes (design tokens, recipes, patterns) in `web/panda.config.ts` do NOT require running `task generate`. Panda CSS generates its output automatically.

## Environment Configuration

The application uses environment variables for configuration. The following file contains all configuration parameters: `./internal/config/config.yaml`

With no environment variables set, running the backend will provide a set of sensible defaults for production.
