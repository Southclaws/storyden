# Local Development

To run Storyden locally via the repository, it's pretty easy! You can use this approach for testing, experimenting and contributing.

For full development and contribution documentation, please visit the [GitHub repository](https://github.com/Southclaws/storyden).

## Prerequisites

To run Storyden locally, you need to have the following installed:

- **Go** - the API is written in Go!
- **Node.js** - the frontend is built with Next.js
- **Yarn** - the package manager for the frontend

If anything is missing from this list, please open an issue!

## Setup Instructions

### 1. Clone the Repository

First, clone the repository:

```bash
git clone https://github.com/Southclaws/storyden.git
cd storyden
```

### 2. Run the Backend (Go API)

From inside the Storyden directory, you can run the API service:

```bash
go run ./cmd/backend
```

This will start the API server with default configuration. You'll get:

- `./data/data.db` SQLite database
- `./data/assets` to store assets (avatars, images, files, etc.)
- A local server running at `http://localhost:8000`
- OpenAPI documentation at `http://localhost:8000/api/docs`
- CORS and cookie rules configured to support localhost

### 3. Run the Frontend (Next.js)

You can also run the frontend service:

```bash
cd web
yarn
yarn dev
```

**Note:** Currently, Storyden's web frontend does not support Next.js `--turbo` due to some minor issues with Panda CSS.

The frontend will be available at `http://localhost:3000` and will by default automatically connect to the API at `http://localhost:8000`.

## Development Workflow

### Running Tests

```bash
# Run all Go tests
go test ./...

# Run specific test package
go test ./tests/thread/...

# Run tests with verbose output
go test -v ./...
```

### Code Generation

Storyden uses heavy code generation. After modifying schemas or API specs, run:

```bash
# Regenerate everything (requires Task installed)
task generate

# Or manually:
# Generate database bindings
task generate:db

# Generate OpenAPI code (backend + frontend + docs)
task generate:openapi
```

### Database Management

```bash
# Seed database with test data
go run ./cmd/seed

# Clean database
go run ./cmd/clean

# Run migrations
go run ./cmd/migrate
```

### Frontend Development

```bash
cd web

# Install dependencies
yarn install

# Start development server
yarn dev

# Build for production
yarn build

# Lint code
yarn lint

# Type check
yarn tsc --noEmit

# Regenerate API client from OpenAPI spec
yarn openapi
```

## Environment Variables

Create a `.env` file in the root directory for custom configuration:

```bash
# Database
DATABASE_URL=sqlite://data/data.db

# Server
LISTEN_ADDR=0.0.0.0:8000
PUBLIC_WEB_ADDRESS=http://localhost:3000
PUBLIC_API_ADDRESS=http://localhost:8000

# Log Level
LOG_LEVEL=debug
LOG_FORMAT=dev

# Development helpers
DEV_CHAOS_SLOW_MODE=0s
DEV_CHAOS_FAIL_RATE=0
```

For a full list of configuration options, see `./internal/config/config.yaml`.

## Troubleshooting

### Port Already in Use

If port 8000 or 3000 is already in use:

```bash
# Change backend port
LISTEN_ADDR=0.0.0.0:8080 go run ./cmd/backend

# Change frontend port
cd web
PORT=3001 yarn dev
```

### Database Issues

If you encounter database issues, you can reset it:

```bash
rm -rf ./data/data.db
go run ./cmd/backend  # Will create fresh database
go run ./cmd/seed     # Optional: add test data
```

### Code Generation Errors

If you see errors about missing generated files:

```bash
# Make sure Task is installed
brew install go-task/tap/go-task  # macOS
# or
go install github.com/go-task/task/v3/cmd/task@latest

# Then regenerate
task generate
```
### API Testing

The OpenAPI documentation is available at:
- Local: `http://localhost:8000/api/docs`
- Interactive testing via Swagger UI

You can also use tools like:
- **Postman** - import from `api/openapi.yaml`
- **curl** - for quick testing
- **httpie** - for human-friendly HTTP requests

## Hot Reload

Both backend and frontend support hot reload:

- **Backend**: Use `air` for hot reload (optional)
  ```bash
  go install github.com/cosmtrek/air@latest
  air
  ```

- **Frontend**: Next.js has built-in hot reload with `yarn dev`

## Additional Resources

- [Official Documentation](https://www.storyden.org/docs)
- [GitHub Repository](https://github.com/Southclaws/storyden)
- [API Reference](http://localhost:8000/api/docs) (when running locally)
- [Community](https://makeroom.club)

## Getting Help

If you encounter issues:

1. Check the [documentation](https://www.storyden.org/docs)
2. Search [existing issues](https://github.com/Southclaws/storyden/issues)
3. Open a new issue with details about your environment and the problem

---

**Source:** Based on [Storyden Local Development Guide](https://www.storyden.org/docs/introduction/local-dev)
