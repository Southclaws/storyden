package openai

import (
	"context"
	"iter"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/responses"
	"google.golang.org/adk/model"
	"google.golang.org/genai"
)

func (o *OpenAI) GenerateContent(ctx context.Context, req *model.LLMRequest, stream bool) iter.Seq2[*model.LLMResponse, error] {
	input := convertToOpenAIInput(req)
	tools := convertToOpenAITools(req)

	params := responses.ResponseNewParams{
		Model: openai.ChatModel(o.modelName),
		Input: responses.ResponseNewParamsInputUnion{
			OfInputItemList: input,
		},
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

func (o *OpenAI) generateContentSync(ctx context.Context, params responses.ResponseNewParams, yield func(*model.LLMResponse, error) bool) {
	res, err := o.client.Responses.New(ctx, params)
	if err != nil {
		yield(nil, fault.Wrap(mapError(err), fctx.With(ctx)))
		return
	}

	content := convertOpenAIResponseToGenaiContent(*res)

	yield(&model.LLMResponse{
		Content:      content,
		FinishReason: convertOpenAIResponseStatusToGenai(string(res.Status)),
		TurnComplete: true,
	}, nil)
}

func (o *OpenAI) generateContentStream(ctx context.Context, params responses.ResponseNewParams, yield func(*model.LLMResponse, error) bool) {
	streamReader := o.client.Responses.NewStreaming(ctx, params)

	var fullContent string

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		if !streamReader.Next() {
			break
		}

		event := streamReader.Current()

		if event.Type == "response.output_text.delta" && event.Delta != "" {
			fullContent += event.Delta
			if !yield(&model.LLMResponse{
				Content: &genai.Content{
					Role:  genai.RoleModel,
					Parts: []*genai.Part{{Text: event.Delta}},
				},
				Partial: true,
			}, nil) {
				return
			}
		}

		if event.Type == "response.completed" {
			yield(&model.LLMResponse{
				Content:      convertOpenAIResponseToGenaiContent(event.Response),
				FinishReason: convertOpenAIResponseStatusToGenai(string(event.Response.Status)),
				TurnComplete: true,
			}, nil)
			return
		}

		if event.Type == "response.failed" || event.Type == "response.incomplete" {
			yield(nil, fault.New("response ended with status: "+string(event.Response.Status)))
			return
		}
	}

	if err := streamReader.Err(); err != nil {
		yield(nil, fault.Wrap(mapError(err), fctx.With(ctx)))
		return
	}

	if fullContent != "" {
		content := &genai.Content{Role: genai.RoleModel, Parts: []*genai.Part{{Text: fullContent}}}
		yield(&model.LLMResponse{
			Content:      content,
			FinishReason: genai.FinishReasonStop,
			TurnComplete: true,
		}, nil)
	}
}
