package anthropic

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/anthropics/anthropic-sdk-go"
	"google.golang.org/adk/v2/model"
	"google.golang.org/genai"

	"github.com/Southclaws/storyden/app/resources/asset"
	robotresource "github.com/Southclaws/storyden/app/resources/robot"
	"github.com/Southclaws/storyden/app/services/semdex/robot/llmprovider/toolschema"
	"github.com/Southclaws/storyden/app/services/semdex/robot/model_media"
)

func convertToAnthropicMessages(ctx context.Context, req *model.LLMRequest, media model_media.ImageResolver) ([]anthropic.MessageParam, error) {
	var messages []anthropic.MessageParam

	for _, content := range req.Contents {
		if content == nil {
			continue
		}

		message, ok, err := convertAnthropicContent(ctx, content, media)
		if err != nil {
			return nil, err
		}
		if ok {
			messages = append(messages, message)
		}
	}

	return messages, nil
}

func convertAnthropicContent(ctx context.Context, content *genai.Content, media model_media.ImageResolver) (anthropic.MessageParam, bool, error) {
	imageIDs, err := robotresource.ImageAssetIDs(content)
	if err != nil {
		return anthropic.MessageParam{}, false, err
	}

	if results := extractAnthropicToolResults(content.Parts); len(results) > 0 {
		if len(imageIDs) > 0 {
			return anthropic.MessageParam{}, false, fmt.Errorf("image attachments cannot be combined with tool results")
		}
		return anthropic.NewUserMessage(results...), true, nil
	}
	if content.Role == genai.RoleUser {
		return convertAnthropicUserContent(ctx, content, imageIDs, media)
	}
	if len(imageIDs) > 0 {
		return anthropic.MessageParam{}, false, fmt.Errorf("image attachments are only supported on user messages")
	}
	if content.Role == genai.RoleModel {
		return convertAnthropicModelContent(content)
	}
	return anthropic.MessageParam{}, false, nil
}

func convertAnthropicUserContent(ctx context.Context, content *genai.Content, imageIDs []asset.AssetID, media model_media.ImageResolver) (anthropic.MessageParam, bool, error) {
	var blocks []anthropic.ContentBlockParamUnion
	if text := extractAllText(content.Parts); text != "" {
		blocks = append(blocks, anthropic.NewTextBlock(text))
	}
	for _, id := range imageIDs {
		if media == nil {
			return anthropic.MessageParam{}, false, fmt.Errorf("image asset resolver is not configured")
		}
		image, err := media.ResolveImage(ctx, id)
		if err != nil {
			return anthropic.MessageParam{}, false, err
		}
		blocks = append(blocks, anthropic.NewImageBlock(anthropic.URLImageSourceParam{URL: image.URL}))
	}
	if len(blocks) == 0 {
		return anthropic.MessageParam{}, false, nil
	}
	return anthropic.NewUserMessage(blocks...), true, nil
}

func convertAnthropicModelContent(content *genai.Content) (anthropic.MessageParam, bool, error) {
	var blocks []anthropic.ContentBlockParamUnion
	if text := extractAllText(content.Parts); text != "" {
		blocks = append(blocks, anthropic.NewTextBlock(text))
	}
	for _, part := range content.Parts {
		if part == nil || part.FunctionCall == nil {
			continue
		}
		blocks = append(blocks, anthropic.NewToolUseBlock(
			part.FunctionCall.ID,
			anthropicToolInput(part.FunctionCall.Args),
			part.FunctionCall.Name,
		))
	}
	if len(blocks) == 0 {
		return anthropic.MessageParam{}, false, nil
	}
	return anthropic.NewAssistantMessage(blocks...), true, nil
}

func anthropicToolInput(args map[string]any) map[string]any {
	if args == nil {
		return map[string]any{}
	}
	return args
}

func extractAnthropicToolResults(parts []*genai.Part) []anthropic.ContentBlockParamUnion {
	var results []anthropic.ContentBlockParamUnion
	for _, part := range parts {
		if part == nil || part.FunctionResponse == nil {
			continue
		}

		id := strings.TrimSpace(part.FunctionResponse.ID)
		if id == "" || id == "{}" || id == "null" {
			continue
		}

		content := ""
		if part.FunctionResponse.Response != nil {
			if b, err := json.Marshal(part.FunctionResponse.Response); err == nil {
				content = string(b)
			}
		}

		results = append(results, anthropic.NewToolResultBlock(id, content, false))
	}
	return results
}

func convertToAnthropicTools(req *model.LLMRequest) []anthropic.ToolUnionParam {
	if req.Config == nil || len(req.Config.Tools) == 0 {
		return nil
	}

	var tools []anthropic.ToolUnionParam

	for _, tool := range req.Config.Tools {
		if tool.FunctionDeclarations == nil {
			continue
		}

		for _, fn := range tool.FunctionDeclarations {
			tools = append(tools, anthropic.ToolUnionParam{
				OfTool: &anthropic.ToolParam{
					Name:        fn.Name,
					Description: anthropic.String(fn.Description),
					InputSchema: buildToolInputSchema(fn),
				},
			})
		}
	}

	return tools
}

func buildToolInputSchema(fn *genai.FunctionDeclaration) anthropic.ToolInputSchemaParam {
	schema := toolschema.FromFunctionDeclaration(fn)

	result := anthropic.ToolInputSchemaParam{
		Properties: map[string]any{},
	}

	if props, ok := schema["properties"]; ok && props != nil {
		result.Properties = props
	}

	if req, ok := schema["required"].([]any); ok {
		for _, r := range req {
			if s, ok := r.(string); ok {
				result.Required = append(result.Required, s)
			}
		}
	}

	return result
}

func convertAnthropicMessageToGenai(blocks []anthropic.ContentBlockUnion) *genai.Content {
	content := &genai.Content{
		Role:  genai.RoleModel,
		Parts: []*genai.Part{},
	}

	for _, block := range blocks {
		switch block.Type {
		case "text":
			if block.Text != "" {
				content.Parts = append(content.Parts, &genai.Part{Text: block.Text})
			}
		case "tool_use":
			args := make(map[string]any)
			if len(block.Input) > 0 {
				json.Unmarshal(block.Input, &args)
			}
			content.Parts = append(content.Parts, &genai.Part{
				FunctionCall: &genai.FunctionCall{
					ID:   block.ID,
					Name: block.Name,
					Args: args,
				},
			})
		}
	}

	return content
}

func convertAnthropicStopReasonToGenai(reason anthropic.StopReason) genai.FinishReason {
	switch reason {
	case anthropic.StopReasonEndTurn, anthropic.StopReasonToolUse, anthropic.StopReasonStopSequence:
		return genai.FinishReasonStop
	case anthropic.StopReasonMaxTokens:
		return genai.FinishReasonMaxTokens
	default:
		return genai.FinishReasonUnspecified
	}
}

func extractAllText(parts []*genai.Part) string {
	var result string
	for _, part := range parts {
		if part != nil && part.Text != "" {
			if result != "" {
				result += "\n\n"
			}
			result += part.Text
		}
	}
	return result
}
