package anthropic

import (
	"encoding/json"
	"testing"

	anthropicapi "github.com/anthropics/anthropic-sdk-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/adk/v2/model"
	"google.golang.org/genai"
)

func TestBuildToolInputSchemaSerializesEmptyParametersAsObject(t *testing.T) {
	inputSchema := buildToolInputSchema(&genai.FunctionDeclaration{
		Name: "system_robot_tool_catalog",
	})

	raw, err := json.Marshal(anthropicapi.ToolParam{
		Name:        "system_robot_tool_catalog",
		InputSchema: inputSchema,
	})
	require.NoError(t, err)

	var payload map[string]any
	require.NoError(t, json.Unmarshal(raw, &payload))

	inputSchemaPayload, ok := payload["input_schema"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "object", inputSchemaPayload["type"])
	assert.Equal(t, map[string]any{}, inputSchemaPayload["properties"])
}

func TestBuildToolInputSchemaFallsBackToJSONSchemaWhenParametersAreEmpty(t *testing.T) {
	inputSchema := buildToolInputSchema(&genai.FunctionDeclaration{
		Name:       "search",
		Parameters: &genai.Schema{},
		ParametersJsonSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"query": map[string]any{
					"type": "string",
				},
			},
			"required": []any{"query"},
		},
	})

	assert.Equal(t, []string{"query"}, inputSchema.Required)

	properties, ok := inputSchema.Properties.(map[string]any)
	require.True(t, ok)
	query, ok := properties["query"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "string", query["type"])
}

func TestConvertToAnthropicMessagesSerializesNilToolArgsAsEmptyObject(t *testing.T) {
	messages := convertToAnthropicMessages(&model.LLMRequest{
		Contents: []*genai.Content{{
			Role: genai.RoleModel,
			Parts: []*genai.Part{{
				FunctionCall: &genai.FunctionCall{
					ID:   "call-no-args",
					Name: "plugin_go_fmt",
					Args: nil,
				},
			}},
		}},
	})

	require.Len(t, messages, 1)

	raw, err := json.Marshal(messages[0])
	require.NoError(t, err)

	var payload struct {
		Content []struct {
			Type  string         `json:"type"`
			Input map[string]any `json:"input"`
		} `json:"content"`
	}
	require.NoError(t, json.Unmarshal(raw, &payload))
	require.Len(t, payload.Content, 1)
	assert.Equal(t, "tool_use", payload.Content[0].Type)
	assert.Equal(t, map[string]any{}, payload.Content[0].Input)
}
