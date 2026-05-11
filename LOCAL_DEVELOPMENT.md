# Local Development

# 本地开发

## 中文说明

本文档说明如何在本地运行 Storyden。代码块中的命令、路径、环境变量和端口名称保持英文，请原样复制使用。

本地开发通常需要两个终端：一个运行 Go 后端 API，一个运行 `web` 目录中的 Next.js 前端。默认后端地址是 `http://localhost:8000`，默认前端地址是 `http://localhost:3000`。

To run Storyden locally via the repository, it's pretty easy! You can use this approach for testing, experimenting and contributing.

For full development and contribution documentation, please visit the [GitHub repository](https://github.com/Southclaws/storyden).

## Prerequisites

## 前置依赖

你需要安装以下工具：

- **Go**：后端 API 使用 Go 编写。
- **Node.js**：前端基于 Next.js。
- **Yarn**：前端包管理器。
- **Task**：用于代码生成和开发工作流的任务运行器。

To run Storyden locally, you need to have the following installed:

- **Go** - the API is written in Go!
- **Node.js** - the frontend is built with Next.js
- **Yarn** - the package manager for the frontend
- **[Task](https://taskfile.dev)** - task runner for code generation and development workflows

If anything is missing from this list, please open an issue!

## Setup Instructions

## 设置步骤

### 1. Clone the Repository

### 1. 克隆仓库

First, clone the repository:

```bash
git clone https://github.com/Southclaws/storyden.git
cd storyden
```

### 2. Run the Backend (Go API)

### 2. 运行后端（Go API）

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

这会使用默认配置启动 API 服务，并自动创建：

- `./data/data.db` SQLite 数据库
- `./data/assets` 本地资源目录，用于头像、图片、文件等
- 本地 API 服务：`http://localhost:8000`
- OpenAPI 文档：`http://localhost:8000/api/docs`
- 适配 localhost 的 CORS 和 cookie 规则

### 3. Run the Frontend (Next.js)

### 3. 运行前端（Next.js）

You can also run the frontend service:

```bash
cd web
yarn
yarn dev
```

The frontend will be available at `http://localhost:3000` and will by default automatically connect to the API at `http://localhost:8000`.

前端会运行在 `http://localhost:3000`，默认会自动连接 `http://localhost:8000` 的 API。

## Development Workflow

## 开发流程

### Running Tests

### 运行测试

```bash
# Run all Go tests
go test ./...

# Run specific test package
go test ./tests/thread/...

# Run tests with verbose output
go test -v ./...

# Run end-to-end tests - automatically boots fresh backend and frontend on ports 8001/3001 then shuts them down after the tests finish
task test:e2e
```

### Code Generation

### 代码生成

Storyden 大量使用代码生成。修改 schema 或 API spec 后，请运行：

Storyden uses heavy code generation. After modifying schemas or API specs, run:

```bash
# Regenerate everything
task generate

# Or manually:
# Generate database bindings
task generate:db

# Generate OpenAPI code (backend + frontend + docs)
task generate:openapi
```

### Database Management

### 数据库管理

```bash
# Seed database with test data
go run ./cmd/seed

# Clean database
go run ./cmd/clean
```

注意：Storyden 会在启动时自动处理数据库迁移，因此通常不需要手动运行迁移。

Note: Storyden automatically handles database migrations on startup, so you don't need to run migrations manually.

### Frontend Development

### 前端开发

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

## 环境变量

如需自定义配置，可以在仓库根目录创建 `.env` 文件：

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

完整配置项请查看 `./internal/config/config.yaml`。

For a full list of configuration options, see `./internal/config/config.yaml`.

## Troubleshooting

## 故障排查

### Port Already in Use

### 端口已被占用

如果端口 8000 或 3000 已被占用：

If port 8000 or 3000 is already in use:

```bash
# Change backend port
LISTEN_ADDR=0.0.0.0:8080 go run ./cmd/backend

# Change frontend port
cd web
PORT=3001 yarn dev
```

如果你修改了端口，或需要在不同地址上运行，请同步更新地址配置：

If you change ports or need to run on different addresses, make sure to update the address configuration:

```ini
PUBLIC_WEB_ADDRESS=http://localhost:3001 \
PUBLIC_API_ADDRESS=http://localhost:8080 \
LISTEN_ADDR=0.0.0.0:8080 \
go run ./cmd/backend
```

这些环境变量会控制 CORS 和 cookie 设置。所有可用配置项请查看 `./internal/config/config.yaml`。

These environment variables control CORS and cookie settings. See `./internal/config/config.yaml` for all available configuration options.

### Database Issues

### 数据库问题

如果遇到数据库问题，可以重置数据库：

If you encounter database issues, you can reset it:

```bash
rm -rf ./data/data.db
go run ./cmd/backend  # Will create fresh database
go run ./cmd/seed     # Optional: add test data
```

### Code Generation Errors

### 代码生成错误

如果看到缺少生成文件的错误，请重新生成代码：

If you see errors about missing generated files, regenerate the code:

```bash
task generate
```

### API Testing

### API 测试

OpenAPI 文档地址：

The OpenAPI documentation is available at:

- Local: `http://localhost:8000/api/docs`
- Interactive testing via Scalar UI

你也可以使用以下工具：

You can also use tools like:

- **Postman** - import from `api/openapi.yaml`
- **curl** - for quick testing
- **httpie** - for human-friendly HTTP requests

## Hot Reload

## 热重载

运行 `yarn dev` 时，前端会通过 Next.js 内置能力支持热重载。

The frontend has built-in hot reload with Next.js when running `yarn dev`.

## Additional Resources

## 其他资源

- [官方文档](https://www.storyden.org/docs)
- [GitHub 仓库](https://github.com/Southclaws/storyden)
- [API 参考](http://localhost:8000/api/docs)（本地运行时）
- [社区](https://makeroom.club)

- [Official Documentation](https://www.storyden.org/docs)
- [GitHub Repository](https://github.com/Southclaws/storyden)
- [API Reference](http://localhost:8000/api/docs) (when running locally)
- [Community](https://makeroom.club)

## Getting Help

## 获取帮助

如果遇到问题：

1. 查看[文档](https://www.storyden.org/docs)
2. 搜索[已有 issue](https://github.com/Southclaws/storyden/issues)
3. 打开一个新 issue，并说明你的环境和遇到的问题

If you encounter issues:

1. Check the [documentation](https://www.storyden.org/docs)
2. Search [existing issues](https://github.com/Southclaws/storyden/issues)
3. Open a new issue with details about your environment and the problem
