package ai

import (
	"context"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/storyden/internal/config"
)

type OpenAI struct {
	client *openai.Client
}

func newOpenAI(cfg config.Config) (*OpenAI, error) {
	client := openai.NewClient(option.WithAPIKey(cfg.OpenAIKey))
	return &OpenAI{client: client}, nil
}

func (o *OpenAI) Prompt(ctx context.Context, input string) (*Result, error) {
	res, err := o.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model: openai.F(openai.ChatModelChatgpt4oLatest),
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(input),
		}),
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
