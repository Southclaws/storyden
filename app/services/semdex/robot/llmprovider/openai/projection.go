package openai

import (
	"encoding/json"
	"strings"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fmsg"
	"github.com/openai/openai-go/v3/packages/param"
	"github.com/openai/openai-go/v3/responses"
	"google.golang.org/adk/model"
	"google.golang.org/genai"
)

const openAIReasoningItemKey = "openai_reasoning_item"

func convertToOpenAIInput(req *model.LLMRequest) []responses.ResponseInputItemUnionParam {
	var input []responses.ResponseInputItemUnionParam

	if req.Config != nil && req.Config.SystemInstruction != nil {
		text := extractAllText(req.Config.SystemInstruction.Parts)
		if text != "" {
			input = append(input, responses.ResponseInputItemParamOfMessage(text, responses.EasyInputMessageRoleSystem))
		}
	}

	for _, content := range req.Contents {
		if content == nil {
			continue
		}

		// Check for function responses first (they can appear in any role)
		functionResponses := extractFunctionResponses(content.Parts)
		if len(functionResponses) > 0 {
			input = append(input, functionResponses...)
			continue
		}

		switch content.Role {
		case genai.RoleUser:
			text := extractAllText(content.Parts)
			if text != "" {
				input = append(input, responses.ResponseInputItemParamOfMessage(text, responses.EasyInputMessageRoleUser))
			}

		case genai.RoleModel:
			input = appendOpenAIModelContent(input, content)
		}
	}

	return input
}

func extractFunctionResponses(parts []*genai.Part) []responses.ResponseInputItemUnionParam {
	var output []responses.ResponseInputItemUnionParam

	for _, part := range parts {
		if part == nil {
			continue
		}
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
			output = append(output, responses.ResponseInputItemParamOfFunctionCallOutput(id, resultJSON))
		}
	}

	return output
}

func appendOpenAIModelContent(input []responses.ResponseInputItemUnionParam, content *genai.Content) []responses.ResponseInputItemUnionParam {
	for _, part := range content.Parts {
		if part == nil {
			continue
		}

		if reasoning := openAIReasoningItemFromPart(part); reasoning != nil {
			input = append(input, responses.ResponseInputItemUnionParam{OfReasoning: reasoning})
		}
		if part.Text != "" {
			input = append(input, responses.ResponseInputItemParamOfMessage(part.Text, responses.EasyInputMessageRoleAssistant))
		}
		if part.FunctionCall != nil {
			argsJSON := "{}"
			if part.FunctionCall.Args != nil {
				if b, err := json.Marshal(part.FunctionCall.Args); err == nil {
					argsJSON = string(b)
				}
			}
			input = append(input, responses.ResponseInputItemParamOfFunctionCall(argsJSON, part.FunctionCall.ID, part.FunctionCall.Name))
		}
	}

	return input
}

func openAIReasoningItemFromPart(part *genai.Part) *responses.ResponseReasoningItemParam {
	if part.PartMetadata == nil {
		return nil
	}
	raw, ok := part.PartMetadata[openAIReasoningItemKey].(string)
	if !ok || raw == "" {
		return nil
	}

	var reasoning responses.ResponseReasoningItemParam
	if err := json.Unmarshal([]byte(raw), &reasoning); err != nil {
		return nil
	}
	return &reasoning
}

func convertToOpenAITools(req *model.LLMRequest) []responses.ToolUnionParam {
	if req.Config == nil || len(req.Config.Tools) == 0 {
		return nil
	}

	var tools []responses.ToolUnionParam

	for _, tool := range req.Config.Tools {
		if tool.FunctionDeclarations == nil {
			continue
		}

		for _, fn := range tool.FunctionDeclarations {
			schema := map[string]any{
				"type":       "object",
				"properties": map[string]any{},
			}

			// Try Parameters first (genai.Schema)
			if fn.Parameters != nil {
				if b, err := json.Marshal(fn.Parameters); err == nil {
					var parsed map[string]any
					if json.Unmarshal(b, &parsed) == nil && parsed != nil {
						schema = parsed
					}
				}
			}

			// Fallback to ParametersJsonSchema if Parameters is nil
			if fn.ParametersJsonSchema != nil {
				if schemaMap, ok := fn.ParametersJsonSchema.(map[string]interface{}); ok && schemaMap != nil {
					schema = schemaMap
				} else if b, err := json.Marshal(fn.ParametersJsonSchema); err == nil {
					var parsed map[string]any
					if json.Unmarshal(b, &parsed) == nil && parsed != nil {
						schema = parsed
					}
				}
			}

			tools = append(tools, responses.ToolUnionParam{OfFunction: &responses.FunctionToolParam{
				Name:        fn.Name,
				Description: param.NewOpt(fn.Description),
				Parameters:  schema,
				Strict:      param.NewOpt(false),
			}})
		}
	}

	return tools
}

func convertOpenAIResponseToGenaiContent(response responses.Response) (*genai.Content, error) {
	content := &genai.Content{
		Role:  genai.RoleModel,
		Parts: []*genai.Part{},
	}

	for _, item := range response.Output {
		switch output := item.AsAny().(type) {
		case responses.ResponseOutputMessage:
			for _, part := range output.Content {
				if part.Type == "output_text" {
					content.Parts = append(content.Parts, &genai.Part{Text: part.Text})
				}
			}
		case responses.ResponseFunctionToolCall:
			args := make(map[string]interface{})
			if output.Arguments != "" {
				if err := json.Unmarshal([]byte(output.Arguments), &args); err != nil {
					return nil, fault.Wrap(err, fmsg.Withf("failed to parse arguments for function call %q", output.Name))
				}
			}
			content.Parts = append(content.Parts, &genai.Part{FunctionCall: &genai.FunctionCall{
				ID: output.CallID, Name: output.Name, Args: args,
			}})
		case responses.ResponseReasoningItem:
			content.Parts = append(content.Parts, &genai.Part{PartMetadata: map[string]any{
				openAIReasoningItemKey: output.RawJSON(),
			}})
		}
	}

	return content, nil
}

func convertOpenAIResponseToGenaiFinishReason(response responses.Response) genai.FinishReason {
	switch response.Status {
	case responses.ResponseStatusCompleted:
		return genai.FinishReasonStop
	case responses.ResponseStatusIncomplete:
		switch response.IncompleteDetails.Reason {
		case "max_output_tokens":
			return genai.FinishReasonMaxTokens
		case "content_filter":
			return genai.FinishReasonSafety
		}
	default:
		return genai.FinishReasonUnspecified
	}

	return genai.FinishReasonUnspecified
}

func openAIResponseError(response responses.Response) error {
	if response.Error.Message != "" {
		return fault.New("openai response failed: " + string(response.Error.Code) + ": " + response.Error.Message)
	}
	return fault.New("openai response ended with status: " + string(response.Status))
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
