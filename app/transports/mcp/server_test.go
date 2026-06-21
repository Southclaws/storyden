package mcp

import (
	"context"
	"encoding/json"
	"testing"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	robot_tools "github.com/Southclaws/storyden/app/services/semdex/robot/tools"
	storydenmcp "github.com/Southclaws/storyden/lib/mcp"
)

func TestBindToolDefaultsMissingSchemas(t *testing.T) {
	server := sdkmcp.NewServer(&sdkmcp.Implementation{
		Name:    "test",
		Version: "v1",
	}, nil)

	tool := &robot_tools.Tool{
		Definition: &storydenmcp.ToolDefinition{
			Name:        "test_tool",
			Title:       "Test tool",
			Description: "Tool used for schema default testing.",
		},
		Handler: func(context.Context, json.RawMessage) (json.RawMessage, error) {
			return json.RawMessage(`{"ok":true}`), nil
		},
	}

	require.NotPanics(t, func() {
		bindTool(server, tool)
	})
}

func TestNormaliseToolSchemaRequiresObjectSchema(t *testing.T) {
	assert.JSONEq(t,
		`{"type":"object","additionalProperties":true}`,
		string(normaliseToolSchema(map[string]any{"type": "string"})),
	)

	assert.JSONEq(t,
		`{"type":"object","properties":{"query":{"type":"string"}}}`,
		string(normaliseToolSchema(map[string]any{
			"type": "object",
			"properties": map[string]any{
				"query": map[string]any{"type": "string"},
			},
		})),
	)
}

func TestShouldExportToolSkipsExternalMCPTools(t *testing.T) {
	nativeTool := &robot_tools.Tool{
		Definition: &storydenmcp.ToolDefinition{Name: "library_search"},
		Handler: func(context.Context, json.RawMessage) (json.RawMessage, error) {
			return json.RawMessage(`{}`), nil
		},
	}

	externalTool := &robot_tools.Tool{
		Definition: &storydenmcp.ToolDefinition{Name: "mcp:notion-mcp-beta:notion-create-comment"},
		Handler: func(context.Context, json.RawMessage) (json.RawMessage, error) {
			return json.RawMessage(`{}`), nil
		},
	}

	assert.True(t, shouldExportTool(nativeTool))
	assert.False(t, shouldExportTool(externalTool))
}
