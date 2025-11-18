package ai

import (
	"context"
	"encoding/json"
	"iter"
	"strings"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/packages/param"
	"github.com/openai/openai-go/v3/shared"
	"github.com/philippgille/chromem-go"
	"google.golang.org/adk/model"
	"google.golang.org/genai"

	"github.com/Southclaws/storyden/internal/config"
)

type OpenAI struct {
	client *openai.Client
	ef     func(ctx context.Context, text string) ([]float32, error)
}

func newOpenAI(cfg config.Config) (*OpenAI, error) {
	client := openai.NewClient(option.WithAPIKey(cfg.OpenAIKey))
	ef := chromem.NewEmbeddingFuncOpenAI(cfg.OpenAIKey, chromem.EmbeddingModelOpenAI3Large)
	return &OpenAI{client: &client, ef: ef}, nil
}

func (o *OpenAI) Prompt(ctx context.Context, input string) (*Result, error) {
	res, err := o.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model: openai.ChatModelGPT5,
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(input),
		},
	})
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if len(res.Choices) == 0 {
		return nil, fault.New("result is empty")
	}

	return &Result{
		Answer: res.Choices[0].Message.Content,
	}, nil
}

func (o *OpenAI) PromptStream(ctx context.Context, input string) (func(yield func(string, error) bool), error) {
	iter := func(yield func(string, error) bool) {
		stream := o.client.Chat.Completions.NewStreaming(ctx, openai.ChatCompletionNewParams{
			Model: openai.ChatModelGPT5,
			Messages: []openai.ChatCompletionMessageParamUnion{
				openai.UserMessage(input),
			},
		})

		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			if !stream.Next() {
				break
			}

			chunk := stream.Current()

			if len(chunk.Choices) > 0 {
				if !yield(chunk.Choices[0].Delta.Content, nil) {
					return
				}
			}
		}

		if err := stream.Err(); err != nil {
			yield("", fault.Wrap(err, fctx.With(ctx)))
			return
		}
	}

	return iter, nil
}

func (o *OpenAI) EmbeddingFunc() func(ctx context.Context, text string) ([]float32, error) {
	return o.ef
}

func (o *OpenAI) GenerateContent(ctx context.Context, req *model.LLMRequest, stream bool) iter.Seq2[*model.LLMResponse, error] {
	messages := convertToOpenAIMessages(req)
	tools := convertToOpenAITools(req)

	params := openai.ChatCompletionNewParams{
		Model:    openai.ChatModelGPT5,
		Messages: messages,
	}

	if len(tools) > 0 {
		params.Tools = tools
	}

	return func(yield func(*model.LLMResponse, error) bool) {
		if stream {
			o.generateContentStream(ctx, params, yield)
		} else {
			o.generateContentSync(ctx, params, yield)
		}
	}
}

func (o *OpenAI) generateContentSync(ctx context.Context, params openai.ChatCompletionNewParams, yield func(*model.LLMResponse, error) bool) {
	res, err := o.client.Chat.Completions.New(ctx, params)
	if err != nil {
		yield(nil, fault.Wrap(err, fctx.With(ctx)))
		return
	}

	if len(res.Choices) == 0 {
		yield(nil, fault.New("no choices in response"))
		return
	}

	choice := res.Choices[0]
	content := convertOpenAIMessageToGenaiContent(choice.Message)
	finishReason := convertOpenAIFinishReasonToGenai(string(choice.FinishReason))

	yield(&model.LLMResponse{
		Content:      content,
		FinishReason: finishReason,
		TurnComplete: true,
	}, nil)
}

func (o *OpenAI) generateContentStream(ctx context.Context, params openai.ChatCompletionNewParams, yield func(*model.LLMResponse, error) bool) {
	streamReader := o.client.Chat.Completions.NewStreaming(ctx, params)

	var fullContent string
	var collectedToolCalls []openai.ChatCompletionMessageToolCallUnion

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		if !streamReader.Next() {
			break
		}

		chunk := streamReader.Current()
		if len(chunk.Choices) == 0 {
			continue
		}

		choice := chunk.Choices[0]
		delta := choice.Delta

		if delta.Content != "" {
			fullContent += delta.Content
			if !yield(&model.LLMResponse{
				Content: &genai.Content{
					Role:  genai.RoleModel,
					Parts: []*genai.Part{{Text: delta.Content}},
				},
				Partial: true,
			}, nil) {
				return
			}
		}

		if len(delta.ToolCalls) > 0 {
			for _, tc := range delta.ToolCalls {
				idx := int(tc.Index)
				for len(collectedToolCalls) <= idx {
					collectedToolCalls = append(collectedToolCalls, openai.ChatCompletionMessageToolCallUnion{})
				}
				if tc.ID != "" {
					collectedToolCalls[idx].ID = tc.ID
				}
				if tc.Function.Name != "" {
					collectedToolCalls[idx].Function.Name = tc.Function.Name
				}
				if tc.Function.Arguments != "" {
					collectedToolCalls[idx].Function.Arguments += tc.Function.Arguments
				}
			}
		}

		if choice.FinishReason != "" {
			content := buildFinalGenaiContent(fullContent, collectedToolCalls)
			finishReason := convertOpenAIFinishReasonToGenai(string(choice.FinishReason))

			yield(&model.LLMResponse{
				Content:      content,
				FinishReason: finishReason,
				TurnComplete: true,
			}, nil)
			return
		}
	}

	if err := streamReader.Err(); err != nil {
		yield(nil, fault.Wrap(err, fctx.With(ctx)))
		return
	}

	if fullContent != "" || len(collectedToolCalls) > 0 {
		content := buildFinalGenaiContent(fullContent, collectedToolCalls)
		yield(&model.LLMResponse{
			Content:      content,
			FinishReason: genai.FinishReasonStop,
			TurnComplete: true,
		}, nil)
	}
}

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
			argsJSON := ""
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
			schema := make(map[string]interface{})
			if fn.Parameters != nil {
				if b, err := json.Marshal(fn.Parameters); err == nil {
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
				result += " "
			}
			result += part.Text
		}
	}
	return result
}
