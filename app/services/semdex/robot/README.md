# Robot Tool Registry

This directory contains the tool registry system for Storyden's AI agents (Robots). The architecture is schema-driven, with all tool definitions centralized in `/api/robots.yaml` and code generated from that single source of truth.

## Architecture Overview

### Schema-Driven Design

All tool definitions live in `/api/robots.yaml` as JSON Schema definitions. This provides:

- Single source of truth for tool schemas
- Type-safe Go bindings generated via `go-jsonschema`
- LLM-friendly OpenAPI-compatible schemas
- Easy validation and documentation

### ToolResult Pattern

Tool execution functions return `ToolResult[T]` instead of `T` directly because the Google ADK `functiontool` type doesn't allow error returns. This wrapper provides:

```go
type ToolResult[T any] struct {
    Result T      // The actual result data
    Error  string // LLM-friendly error message
}
```

Helper constructors:

- `NewSuccess[T](v T) ToolResult[T]` - For successful results
- `NewError[T](err error) ToolResult[T]` - For errors (uses err.Error())
- `NewErrorMsg[T](msg string) ToolResult[T]` - For custom error messages

### Why This Pattern Exists

The Google ADK has two error handling paths:

1. **Schema validation errors** (wrong argument names, types, missing required fields): The ADK validates arguments before calling our function. When validation fails, it wraps the error as a Go `error` type which doesn't serialize properly to JSON, resulting in `{"error":{}}` sent to the LLM.

2. **Tool execution errors** (our ToolResult): By returning `ToolResult[T]` with a string `Error` field, we ensure errors serialize correctly as `{"Result":...,"Error":"message here"}`.

**Impact**: Schema validation failures appear as empty error objects to the LLM, while tool execution errors have clear messages. We work around this by providing detailed schema descriptions and validation constraints that help the LLM avoid schema errors.

### Error Accumulation for Validation

When validating array inputs (like kinds, authors, tools), we **accumulate errors** instead of failing early. This gives the LLM feedback about which specific arguments were valid vs invalid.

Example from search tool:

```go
var validationErrors []string

// Validate kinds filter
var invalidKinds []string
for _, k := range args.Kinds {
    kind, err := datagraph.NewKind(string(k))
    if err != nil {
        invalidKinds = append(invalidKinds, string(k))
        continue
    }
    kinds = append(kinds, kind)
}
if len(invalidKinds) > 0 {
    validationErrors = append(validationErrors, fmt.Sprintf("invalid kinds: %v", invalidKinds))
}

// Continue with valid kinds, accumulate all errors, return at end
```

This allows partial success - the search can proceed with valid filters while informing the LLM about invalid ones.

## Tool Registry Components

### 1. Schema Definition (`/api/robots.yaml`)

Each tool has three schema definitions:

- `Tool{Name}`: Top-level definition with title and description
- `Tool{Name}Input`: Input parameters schema
- `Tool{Name}Output`: Output result schema

Example:

```yaml
ToolSearch:
  title: search
  description: "Search the Storyden knowledge base..."
  type: object
  required: [input, output]
  properties:
    input:
      $ref: "#/definitions/ToolSearchInput"
    output:
      $ref: "#/definitions/ToolSearchOutput"

ToolSearchInput:
  type: object
  required: [query]
  properties:
    query:
      type: string
      description: The search query text
    # ... more properties
```

### 2. Code Generation (`/api/embed.go`, `/mcp/bindings.go`)

`embed.go` just exists to load the schema YAML into Go's address space at build time.

Then, the `mcp/bindings` code will make this more useful for us:

The `init()` function:

1. Parses `schema.yaml` into Go structs
2. Calls `initAllTools()` to create tool definitions
3. Builds the tool names enum for validation

Key functions:

- `GetSearchTool()` - Returns search tool definition
- `GetRobotCreateTool()` - Returns robot create tool definition
- `AllToolNames()` - Returns slice of all tool names
- `InjectToolNamesEnum(schema, propertyName)` - Injects dynamic enum into schema

### 3. Tool Provider (`tools/provider.go`)

The `DefaultToolProvider` implements the `ToolProvider` interface:

```go
type ToolProvider interface {
    GetTool(ctx context.Context, name string) (tool.Tool, error)
    GetDefaultTools(ctx context.Context) ([]tool.Tool, error)
    GetRobotTools(ctx context.Context, robotID xid.ID) ([]tool.Tool, error)
}
```

Constructor signature:

```go
func NewToolProvider(
    logger *slog.Logger,
    db *ent.Client,
    searcher searcher.Searcher,
    accountQuerier *account_querier.Querier,
    categoryRepo *category.Repository,
) ToolProvider
```

## Available Tools

### Search (`tool_search.go`)

Searches the Storyden knowledge base with filters:

- `query` (required): Search text
- `kinds`: Filter by content type (post, thread, reply, node, collection, profile, event)
- `authors`: Filter by author handles (username strings, looked up to IDs)
- `categories`: Filter by category names (case-insensitive match to IDs)
- `tags`: Filter by tag names
- `max_results`: Limit results (default 10)

**Validation**: Accumulates errors for invalid kinds, unfound authors, and unfound categories while proceeding with valid values.

### Robot Switch (`tool_switch.go`)

Switches the conversation to a different Robot agent.

- Dynamically injects available robot IDs into enum
- Validates robot exists before switching

### Robot CRUD (`tool_robots.go`)

#### Create

Creates a new Robot with:

- `name` (required): Robot name
- `playbook` (required): System prompt/directive
- `description` (optional): Human-readable description
- `tools` (optional): Array of tool names

**Validation**: Validates tool names against `AllToolNames()` enum, accumulates invalid tool errors.

#### List

Lists all robots with optional limit.

#### Get

Retrieves a specific robot by ID.

**Validation**: Validates ID format.

#### Update

Updates robot fields (name, description, playbook, tools).

**Validation**: Same tool name validation as Create.

#### Delete

Permanently deletes a robot.

**Validation**: Validates ID format and robot existence.

## Dynamic Schema Injection

Some tool schemas need runtime data:

### Robot ID Enum (Switch Tool)

```go
robots, err := p.db.Robot.Query().Order(robot.ByCreatedAt()).All(ctx)
robotIDs := make([]any, len(robots))
for i, r := range robots {
    robotIDs[i] = r.ID.String()
}
inputSchema.Properties["robot_id"].Enum = robotIDs
```

### Tool Names Enum (Create/Update Robot)

```go
inputSchema := toolDef.InputSchema
mcp.InjectToolNamesEnum(inputSchema, "tools")
```

The `InjectToolNamesEnum` helper finds the array property and injects the enum of all available tool names.

## Adding a New Tool

1. **Define schema in `/api/robots.yaml`**:

```yaml
ToolMyNewTool:
  title: my_new_tool
  description: "What this tool does..."
  type: object
  required: [input, output]
  properties:
    input:
      $ref: "#/definitions/ToolMyNewToolInput"
    output:
      $ref: "#/definitions/ToolMyNewToolOutput"

ToolMyNewToolInput:
  type: object
  required: [some_field]
  properties:
    some_field:
      type: string
      description: Description here

ToolMyNewToolOutput:
  type: object
  required: [result]
  properties:
    result:
      type: string
```

2. **Update `/api/embed.go`**:

```go
var toolMyNewTool *ToolDefinition

func initAllTools() {
    // ... existing tools
    toolMyNewTool = initTool("MyNewTool")
    // Update AllToolNames() return value
}

func GetMyNewToolTool() *ToolDefinition {
    return toolMyNewTool
}

func AllToolNames() []string {
    return []string{
        // ... existing tools
        toolMyNewTool.Name,
    }
}
```

3. **Create tool file `tool_mynew.go`**:

```go
package tools

func (p *DefaultToolProvider) newMyNewToolTool(ctx context.Context) (tool.Tool, error) {
    toolDef := mcp.GetMyNewToolTool()

    tool, err := functiontool.New(
        functiontool.Config{
            Name:        toolDef.Name,
            Description: toolDef.Description,
            InputSchema: toolDef.InputSchema,
        },
        p.executeMyNewTool,
    )
    return tool, err
}

func (p *DefaultToolProvider) executeMyNewTool(ctx tool.Context, args mcp.ToolMyNewToolInput) ToolResult[mcp.ToolMyNewToolOutput] {
    // Validate inputs, accumulate errors
    var validationErrors []string

    // Do the work
    result, err := doSomething(args.SomeField)
    if err != nil {
        return NewError[mcp.ToolMyNewToolOutput](err)
    }

    // Return success or partial success with validation errors
    output := mcp.ToolMyNewToolOutput{Result: result}
    if len(validationErrors) > 0 {
        return ToolResult[mcp.ToolMyNewToolOutput]{
            Result: output,
            Error:  strings.Join(validationErrors, "; "),
        }
    }
    return NewSuccess(output)
}
```

4. **Register in `provider.go` `GetTool` switch**:

```go
case mcp.GetMyNewToolTool().Name:
    return p.newMyNewToolTool(ctx)
```

5. **Generate code**:

```bash
task generate
```

## Testing Tools

Tools are tested via integration tests that verify:

- Schema validation
- Error handling and accumulation
- Successful execution paths
- LLM-friendly error messages

See `tests/robot/` for examples.

## Best Practices

1. **Always use ToolResult pattern** - Even for infallible operations
2. **Accumulate validation errors** - Don't fail fast, collect all issues
3. **Provide helpful error messages** - Remember the LLM reads these
4. **Validate IDs early** - Parse xid.IDs before database operations
5. **Use schema-driven types** - Never hand-write request/response structs
6. **Update AllToolNames()** - Keep the enum synchronized
7. **Document filter behavior** - Explain how name→ID lookups work

## Dependency Injection

The tool provider is wired via Uber FX in the main application. Add new dependencies to:

1. `DefaultToolProvider` struct in `provider.go`
2. `NewToolProvider` constructor signature
3. FX provider in `app/services/semdex/robot/fx.go`
