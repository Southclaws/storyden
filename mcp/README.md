# MCP Tool Schemas

This directory contains JSON Schema definitions for tool inputs and outputs that are shared between the backend (Go) and frontend (TypeScript).

## Schema Definition

All tool schemas are defined in `schema.yaml` using JSON Schema Draft 07 format.

### Naming Convention

Tool schemas follow a strict naming convention:

- **Tool Definition**: `Tool{Name}` (e.g., `ToolRobotSwitch`)
  - Uses `title` field for the tool name (e.g., `switch_agent`)
  - Uses `description` field for the tool description (shown to LLM)
  - Has `input` and `output` properties that reference the input/output schemas

- **Input Schema**: `Tool{Name}Input` (e.g., `ToolRobotSwitchInput`)
  - Defines the structure of arguments passed to the tool

- **Output Schema**: `Tool{Name}Output` (e.g., `ToolRobotSwitchOutput`)
  - Defines the structure of the result returned by the tool

This convention allows the generic `initTool(name)` function to automatically load and resolve all schemas.

## Code Generation

### Go Types

Go types are automatically generated using `go-jsonschema`:

```bash
go generate ./mcp
```

This generates `mcp_schema.go` with type-safe structs.

### TypeScript Types

TypeScript types can be generated using `json-schema-to-typescript`:

```bash
# From the web directory
npx json-schema-to-typescript ../mcp/schema.yaml > src/lib/agent/tool-schemas.ts
```

Or add to your package.json scripts:

```json
{
  "scripts": {
    "generate:tool-schemas": "json-schema-to-typescript ../mcp/schema.yaml > src/lib/agent/tool-schemas.ts"
  }
}
```

## Usage

### Backend (Go)

```go
import "github.com/Southclaws/storyden/mcp"

func execute(ctx tool.Context, args mcp.SwitchAgentInput) mcp.SwitchAgentOutput {
    return mcp.SwitchAgentOutput{
        Success: true,
        RobotId: args.RobotId,
    }
}
```

### Frontend (TypeScript)

```typescript
import { SwitchAgentInput, SwitchAgentOutput } from "@/lib/agent/tool-schemas";

onToolCall: async ({ toolCall }) => {
  if (toolCall.toolName === "switch_agent") {
    const input = toolCall.input as SwitchAgentInput;
    const output: SwitchAgentOutput = {
      success: true,
      robot_id: input.robot_id,
    };
  }
};
```

## Adding New Tools

1. Define the input and output schemas in `schema.yaml`
2. Run `go generate ./mcp` to generate Go types
3. Run the TypeScript generation command to update frontend types
4. Import and use the generated types in your code
