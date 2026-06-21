package anthropic

import (
	"encoding/json"
	"strings"

	"github.com/anthropics/anthropic-sdk-go"
	"google.golang.org/adk/model"
	"google.golang.org/genai"
)

func convertToAnthropicMessages(req *model.LLMRequest) []anthropic.MessageParam {
	var messages []anthropic.MessageParam

	for _, content := range req.Contents {
		if content == nil {
			continue
		}

		// Tool results come back as FunctionResponse parts regardless of role;
		// Anthropic expects them in user-turn messages as tool_result blocks.
		if results := extractAnthropicToolResults(content.Parts); len(results) > 0 {
			messages = append(messages, anthropic.NewUserMessage(results...))
			continue
		}

		switch content.Role {
		case genai.RoleUser:
			if text := extractAllText(content.Parts); text != "" {
				messages = append(messages, anthropic.NewUserMessage(anthropic.NewTextBlock(text)))
			}

		case genai.RoleModel:
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

			if len(blocks) > 0 {
				messages = append(messages, anthropic.NewAssistantMessage(blocks...))
			}
		}
	}

	return messages
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
	var schema map[string]any

	if fn.Parameters != nil {
		if b, err := json.Marshal(fn.Parameters); err == nil {
			json.Unmarshal(b, &schema)
		}
	}

	if len(schema) == 0 && fn.ParametersJsonSchema != nil {
		if b, err := json.Marshal(fn.ParametersJsonSchema); err == nil {
			json.Unmarshal(b, &schema)
		}
	}

	// genai schema types are uppercase ("STRING", "OBJECT") but Anthropic requires lowercase.
	normalizeSchemaTypeStrings(schema)

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

// normalizeSchemaTypeStrings recursively lowercases "type" values in a JSON schema map.
func normalizeSchemaTypeStrings(schema map[string]any) {
	if schema == nil {
		return
	}
	if t, ok := schema["type"].(string); ok {
		schema["type"] = strings.ToLower(t)
	}
	if props, ok := schema["properties"].(map[string]any); ok {
		for _, v := range props {
			if sub, ok := v.(map[string]any); ok {
				normalizeSchemaTypeStrings(sub)
			}
		}
	}
	if items, ok := schema["items"].(map[string]any); ok {
		normalizeSchemaTypeStrings(items)
	}
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
