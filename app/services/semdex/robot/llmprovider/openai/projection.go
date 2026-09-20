package openai

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fmsg"
	"github.com/openai/openai-go/v3/packages/param"
	"github.com/openai/openai-go/v3/responses"
	"google.golang.org/adk/v2/model"
	"google.golang.org/genai"

	"github.com/Southclaws/storyden/app/resources/asset"
	robotresource "github.com/Southclaws/storyden/app/resources/robot"
	"github.com/Southclaws/storyden/app/services/semdex/robot/llmprovider/toolschema"
	"github.com/Southclaws/storyden/app/services/semdex/robot/model_media"
)

const openAIReasoningItemKey = "openai_reasoning_item"

func convertToOpenAIInput(ctx context.Context, req *model.LLMRequest, media model_media.ImageResolver) ([]responses.ResponseInputItemUnionParam, error) {
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

		items, err := convertOpenAIContent(ctx, content, media)
		if err != nil {
			return nil, err
		}
		input = append(input, items...)
	}

	return input, nil
}

func convertOpenAIContent(ctx context.Context, content *genai.Content, media model_media.ImageResolver) ([]responses.ResponseInputItemUnionParam, error) {
	imageIDs, err := robotresource.ImageAssetIDs(content)
	if err != nil {
		return nil, err
	}

	if functionResponses := extractFunctionResponses(content.Parts); len(functionResponses) > 0 {
		if len(imageIDs) > 0 {
			return nil, fmt.Errorf("image attachments cannot be combined with tool results")
		}
		return functionResponses, nil
	}

	if content.Role == genai.RoleUser {
		message, ok, err := convertOpenAIUserContent(ctx, content, imageIDs, media)
		if !ok || err != nil {
			return nil, err
		}
		return []responses.ResponseInputItemUnionParam{message}, nil
	}
	if len(imageIDs) > 0 {
		return nil, fmt.Errorf("image attachments are only supported on user messages")
	}
	if content.Role == genai.RoleModel {
		return appendOpenAIModelContent(nil, content), nil
	}
	return nil, nil
}

func convertOpenAIUserContent(ctx context.Context, content *genai.Content, imageIDs []asset.AssetID, media model_media.ImageResolver) (responses.ResponseInputItemUnionParam, bool, error) {
	var blocks responses.ResponseInputMessageContentListParam
	if text := extractAllText(content.Parts); text != "" {
		blocks = append(blocks, responses.ResponseInputContentParamOfInputText(text))
	}
	for _, id := range imageIDs {
		if media == nil {
			return responses.ResponseInputItemUnionParam{}, false, fmt.Errorf("image asset resolver is not configured")
		}
		image, err := media.ResolveImage(ctx, id)
		if err != nil {
			return responses.ResponseInputItemUnionParam{}, false, err
		}
		block := responses.ResponseInputContentParamOfInputImage(responses.ResponseInputImageDetailAuto)
		block.OfInputImage.ImageURL = param.NewOpt(image.URL)
		blocks = append(blocks, block)
	}
	if len(blocks) == 0 {
		return responses.ResponseInputItemUnionParam{}, false, nil
	}
	return responses.ResponseInputItemParamOfMessage(blocks, responses.EasyInputMessageRoleUser), true, nil
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
			tools = append(tools, responses.ToolUnionParam{OfFunction: &responses.FunctionToolParam{
				Name:        fn.Name,
				Description: param.NewOpt(fn.Description),
				Parameters:  toolschema.FromFunctionDeclaration(fn),
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
