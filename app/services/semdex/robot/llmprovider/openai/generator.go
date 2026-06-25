package openai

import (
	"context"
	"iter"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/openai/openai-go/v3"
	"google.golang.org/adk/model"
	"google.golang.org/genai"
)

func (o *OpenAI) GenerateContent(ctx context.Context, req *model.LLMRequest, stream bool) iter.Seq2[*model.LLMResponse, error] {
	messages := convertToOpenAIMessages(req)
	tools := convertToOpenAITools(req)

	params := openai.ChatCompletionNewParams{
		Model:    openai.ChatModel(o.modelName),
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
		yield(nil, fault.Wrap(mapError(err), fctx.With(ctx)))
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
		yield(nil, fault.Wrap(mapError(err), fctx.With(ctx)))
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
