package openai

import (
	"encoding/json"
	"strings"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/packages/param"
	"github.com/openai/openai-go/v3/shared"
	"google.golang.org/adk/model"
	"google.golang.org/genai"
)

func convertToOpenAIMessages(req *model.LLMRequest) []openai.ChatCompletionMessageParamUnion {
	var messages []openai.ChatCompletionMessageParamUnion

	if req.Config != nil && req.Config.SystemInstruction != nil {
		text := extractAllText(req.Config.SystemInstruction.Parts)
		if text != "" {
			messages = append(messages, openai.SystemMessage(text))
		}
	}

	for _, content := range req.Contents {
		if content == nil {
			continue
		}

		// Check for function responses first (they can appear in any role)
		functionResponses := extractFunctionResponses(content.Parts)
		if len(functionResponses) > 0 {
			for _, resp := range functionResponses {
				messages = append(messages, resp)
			}
			continue
		}

		switch content.Role {
		case genai.RoleUser:
			text := extractAllText(content.Parts)
			if text != "" {
				messages = append(messages, openai.UserMessage(text))
			}

		case genai.RoleModel:
			text := extractAllText(content.Parts)
			toolCalls := extractToolCalls(content.Parts)

			if len(toolCalls) > 0 {
				msg := openai.ChatCompletionAssistantMessageParam{
					ToolCalls: toolCalls,
				}
				if text != "" {
					msg.Content = openai.ChatCompletionAssistantMessageParamContentUnion{
						OfString: param.NewOpt(text),
					}
				}
				messages = append(messages, openai.ChatCompletionMessageParamUnion{
					OfAssistant: &msg,
				})
			} else if text != "" {
				messages = append(messages, openai.AssistantMessage(text))
			}
		}
	}

	return messages
}

func extractFunctionResponses(parts []*genai.Part) []openai.ChatCompletionMessageParamUnion {
	var responses []openai.ChatCompletionMessageParamUnion

	for _, part := range parts {
		if part.FunctionResponse != nil {
			// Validate the ID - OpenAI requires a valid tool_call_id
			id := strings.TrimSpace(part.FunctionResponse.ID)
			if id == "" || id == "{}" || id == "null" {
				// Skip invalid IDs - they would cause API errors
				continue
			}

			resultJSON := ""
			if part.FunctionResponse.Response != nil {
				if b, err := json.Marshal(part.FunctionResponse.Response); err == nil {
					resultJSON = string(b)
				}
			}
			// ToolMessage signature is: ToolMessage(content, toolCallID)
			responses = append(responses, openai.ToolMessage(resultJSON, id))
		}
	}

	return responses
}

func extractToolCalls(parts []*genai.Part) []openai.ChatCompletionMessageToolCallUnionParam {
	var toolCalls []openai.ChatCompletionMessageToolCallUnionParam

	for _, part := range parts {
		if part.FunctionCall != nil {
			argsJSON := "{}"
			if part.FunctionCall.Args != nil {
				if b, err := json.Marshal(part.FunctionCall.Args); err == nil {
					argsJSON = string(b)
				}
			}

			toolCalls = append(toolCalls, openai.ChatCompletionMessageToolCallUnionParam{
				OfFunction: &openai.ChatCompletionMessageFunctionToolCallParam{
					ID: part.FunctionCall.ID,
					Function: openai.ChatCompletionMessageFunctionToolCallFunctionParam{
						Name:      part.FunctionCall.Name,
						Arguments: argsJSON,
					},
				},
			})
		}
	}

	return toolCalls
}

func convertToOpenAITools(req *model.LLMRequest) []openai.ChatCompletionToolUnionParam {
	if req.Config == nil || len(req.Config.Tools) == 0 {
		return nil
	}

	var tools []openai.ChatCompletionToolUnionParam

	for _, tool := range req.Config.Tools {
		if tool.FunctionDeclarations == nil {
			continue
		}

		for _, fn := range tool.FunctionDeclarations {
			var schema map[string]interface{}

			// Try Parameters first (genai.Schema)
			if fn.Parameters != nil {
				if b, err := json.Marshal(fn.Parameters); err == nil {
					json.Unmarshal(b, &schema)
				}
			}

			// Fallback to ParametersJsonSchema if Parameters is nil
			if schema == nil && fn.ParametersJsonSchema != nil {
				if schemaMap, ok := fn.ParametersJsonSchema.(map[string]interface{}); ok {
					schema = schemaMap
				} else if b, err := json.Marshal(fn.ParametersJsonSchema); err == nil {
					json.Unmarshal(b, &schema)
				}
			}

			tools = append(tools, openai.ChatCompletionFunctionTool(shared.FunctionDefinitionParam{
				Name:        fn.Name,
				Description: param.NewOpt(fn.Description),
				Parameters:  shared.FunctionParameters(schema),
			}))
		}
	}

	return tools
}

func convertOpenAIMessageToGenaiContent(msg openai.ChatCompletionMessage) *genai.Content {
	content := &genai.Content{
		Role:  genai.RoleModel,
		Parts: []*genai.Part{},
	}

	if msg.Content != "" {
		content.Parts = append(content.Parts, &genai.Part{Text: msg.Content})
	}

	for _, tc := range msg.ToolCalls {
		args := make(map[string]interface{})
		if tc.Function.Arguments != "" {
			json.Unmarshal([]byte(tc.Function.Arguments), &args)
		}

		content.Parts = append(content.Parts, &genai.Part{
			FunctionCall: &genai.FunctionCall{
				ID:   tc.ID,
				Name: tc.Function.Name,
				Args: args,
			},
		})
	}

	return content
}

func buildFinalGenaiContent(text string, toolCalls []openai.ChatCompletionMessageToolCallUnion) *genai.Content {
	content := &genai.Content{
		Role:  genai.RoleModel,
		Parts: []*genai.Part{},
	}

	if text != "" {
		content.Parts = append(content.Parts, &genai.Part{Text: text})
	}

	for _, tc := range toolCalls {
		args := make(map[string]interface{})
		if tc.Function.Arguments != "" {
			json.Unmarshal([]byte(tc.Function.Arguments), &args)
		}

		content.Parts = append(content.Parts, &genai.Part{
			FunctionCall: &genai.FunctionCall{
				ID:   tc.ID,
				Name: tc.Function.Name,
				Args: args,
			},
		})
	}

	return content
}

func convertOpenAIFinishReasonToGenai(reason string) genai.FinishReason {
	switch reason {
	case "stop":
		return genai.FinishReasonStop
	case "length":
		return genai.FinishReasonMaxTokens
	case "tool_calls":
		return genai.FinishReasonStop
	case "content_filter":
		return genai.FinishReasonSafety
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
