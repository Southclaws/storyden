package anthropic

import (
	"context"
	"encoding/json"
	"testing"

	anthropicapi "github.com/anthropics/anthropic-sdk-go"
	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/adk/v2/model"
	"google.golang.org/genai"

	"github.com/Southclaws/storyden/app/resources/asset"
	robotresource "github.com/Southclaws/storyden/app/resources/robot"
	"github.com/Southclaws/storyden/app/services/semdex/robot/model_media"
)

type imageResolver map[asset.AssetID]string

func (r imageResolver) ResolveImage(_ context.Context, id asset.AssetID) (model_media.Image, error) {
	return model_media.Image{AssetID: id, URL: r[id], MIME: "image/png"}, nil
}

func TestBuildToolInputSchemaSerializesEmptyParametersAsObject(t *testing.T) {
	inputSchema := buildToolInputSchema(&genai.FunctionDeclaration{
		Name: "robot_create",
	})

	raw, err := json.Marshal(anthropicapi.ToolParam{
		Name:        "robot_create",
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
	messages, err := convertToAnthropicMessages(context.Background(), &model.LLMRequest{
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
	}, nil)
	require.NoError(t, err)

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

func TestConvertToAnthropicMessagesProjectsMultipleStorydenImagesInOrder(t *testing.T) {
	first := asset.AssetID(xid.New())
	second := asset.AssetID(xid.New())
	resolver := imageResolver{
		first:  "https://storyden.example/api/assets/first.png",
		second: "https://storyden.example/api/assets/second.webp",
	}
	req := &model.LLMRequest{Contents: []*genai.Content{{
		Role: genai.RoleUser,
		Parts: []*genai.Part{
			{Text: "Compare these images."},
			robotresource.NewImageAssetPart(first),
			robotresource.NewImageAssetPart(second),
			robotresource.NewImageAssetPart(first),
		},
	}}}

	messages, err := convertToAnthropicMessages(context.Background(), req, resolver)
	require.NoError(t, err)
	require.Len(t, messages, 1)

	raw, err := json.Marshal(messages[0])
	require.NoError(t, err)
	assert.JSONEq(t, `{
		"role":"user",
		"content":[
			{"type":"text","text":"Compare these images."},
			{"type":"image","source":{"type":"url","url":"https://storyden.example/api/assets/first.png"}},
			{"type":"image","source":{"type":"url","url":"https://storyden.example/api/assets/second.webp"}},
			{"type":"image","source":{"type":"url","url":"https://storyden.example/api/assets/first.png"}}
		]
	}`, string(raw))
}

func TestConvertToAnthropicMessagesRejectsExternalMediaParts(t *testing.T) {
	req := &model.LLMRequest{Contents: []*genai.Content{{
		Role:  genai.RoleUser,
		Parts: []*genai.Part{{InlineData: &genai.Blob{MIMEType: "image/png", Data: []byte("untrusted")}}},
	}}}

	_, err := convertToAnthropicMessages(context.Background(), req, nil)
	require.Error(t, err)
	assert.ErrorContains(t, err, "must reference a Storyden asset ID")
}
