package anthropic

import (
	"context"
	"iter"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/anthropics/anthropic-sdk-go"
	"google.golang.org/adk/v2/model"
	"google.golang.org/genai"
)

func (m *Anthropic) GenerateContent(ctx context.Context, req *model.LLMRequest, stream bool) iter.Seq2[*model.LLMResponse, error] {
	messages := convertToAnthropicMessages(req)
	tools := convertToAnthropicTools(req)

	params := anthropic.MessageNewParams{
		Model:     anthropic.Model(m.modelName),
		MaxTokens: 8096,
		Messages:  messages,
	}

	if req.Config != nil && req.Config.SystemInstruction != nil {
		if text := extractAllText(req.Config.SystemInstruction.Parts); text != "" {
			params.System = []anthropic.TextBlockParam{{Text: text}}
		}
	}

	if len(tools) > 0 {
		params.Tools = tools
		params.ToolChoice = anthropic.ToolChoiceUnionParam{
			OfAuto: &anthropic.ToolChoiceAutoParam{},
		}
	}

	return func(yield func(*model.LLMResponse, error) bool) {
		if stream {
			m.generateContentStream(ctx, params, yield)
		} else {
			m.generateContentSync(ctx, params, yield)
		}
	}
}

func (m *Anthropic) generateContentSync(ctx context.Context, params anthropic.MessageNewParams, yield func(*model.LLMResponse, error) bool) {
	msg, err := m.client.Messages.New(ctx, params)
	if err != nil {
		yield(nil, fault.Wrap(mapError(err), fctx.With(ctx)))
		return
	}

	yield(&model.LLMResponse{
		Content:      convertAnthropicMessageToGenai(msg.Content),
		FinishReason: convertAnthropicStopReasonToGenai(msg.StopReason),
		TurnComplete: true,
	}, nil)
}

func (m *Anthropic) generateContentStream(ctx context.Context, params anthropic.MessageNewParams, yield func(*model.LLMResponse, error) bool) {
	stream := m.client.Messages.NewStreaming(ctx, params)

	var accumulated anthropic.Message
	for stream.Next() {
		event := stream.Current()
		if err := accumulated.Accumulate(event); err != nil {
			yield(nil, fault.Wrap(err, fctx.With(ctx)))
			return
		}

		if delta, ok := event.AsAny().(anthropic.ContentBlockDeltaEvent); ok {
			if textDelta, ok := delta.Delta.AsAny().(anthropic.TextDelta); ok && textDelta.Text != "" {
				if !yield(&model.LLMResponse{
					Content: &genai.Content{
						Role:  genai.RoleModel,
						Parts: []*genai.Part{{Text: textDelta.Text}},
					},
					Partial: true,
				}, nil) {
					return
				}
			}
		}
	}

	if err := stream.Err(); err != nil {
		yield(nil, fault.Wrap(mapError(err), fctx.With(ctx)))
		return
	}

	yield(&model.LLMResponse{
		Content:      convertAnthropicMessageToGenai(accumulated.Content),
		FinishReason: convertAnthropicStopReasonToGenai(accumulated.StopReason),
		TurnComplete: true,
	}, nil)
}
