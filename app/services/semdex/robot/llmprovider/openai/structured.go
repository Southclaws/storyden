package openai

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/openai/openai-go/v3"
	openaioption "github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/shared"

	"github.com/Southclaws/storyden/app/resources/robot/llm_provider"
	"github.com/Southclaws/storyden/app/resources/robot/model_ref"
)

func (*OpenAI) SupportsStructuredOutput() bool { return true }

func (p *OpenAI) StructuredPrompt(ctx context.Context, ref model_ref.ModelRef, request llm_provider.StructuredPromptRequest) (string, error) {
	p.mu.RLock()
	apiKey := p.apiKey
	p.mu.RUnlock()

	client := openai.NewClient(openaioption.WithAPIKey(apiKey))

	res, err := client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model: openai.ChatModel(ref.Model.String()),
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(request.Input),
		},
		ResponseFormat: openai.ChatCompletionNewParamsResponseFormatUnion{
			OfJSONSchema: &shared.ResponseFormatJSONSchemaParam{
				JSONSchema: shared.ResponseFormatJSONSchemaJSONSchemaParam{
					Name:        "json_schema",
					Strict:      openai.Bool(true),
					Description: openai.String(request.Description),
					Schema:      request.Schema,
				},
			},
		},
	})
	if err != nil {
		return "", fault.Wrap(mapError(err), fctx.With(ctx))
	}
	if len(res.Choices) == 0 {
		return "", fault.New("structured output response choices are empty")
	}

	content := res.Choices[0].Message.Content
	if content == "" {
		return "", fault.New("structured output response is empty")
	}

	return content, nil
}
