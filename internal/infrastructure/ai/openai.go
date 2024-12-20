package ai

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/philippgille/chromem-go"

	"github.com/Southclaws/storyden/internal/config"
)

type OpenAI struct {
	client *openai.Client
	ef     func(ctx context.Context, text string) ([]float32, error)
}

func newOpenAI(cfg config.Config) (*OpenAI, error) {
	client := openai.NewClient(option.WithAPIKey(cfg.OpenAIKey))
	ef := chromem.NewEmbeddingFuncOpenAI(cfg.OpenAIKey, chromem.EmbeddingModelOpenAI3Large)
	return &OpenAI{client: client, ef: ef}, nil
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

func (o *OpenAI) PromptStream(ctx context.Context, input string) (chan string, chan error) {
	ch := make(chan string)
	errCh := make(chan error)

	go func() {
		defer close(ch)
		defer close(errCh)

		stream := o.client.Chat.Completions.NewStreaming(ctx, openai.ChatCompletionNewParams{
			Model: openai.F(openai.ChatModelChatgpt4oLatest),
			Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
				openai.UserMessage(input),
			}),
		})

		for {
			select {
			case <-ctx.Done():
				errCh <- ctx.Err()
				return
			default:
			}

			if !stream.Next() {
				break
			}

			chunk := stream.Current()

			if len(chunk.Choices) > 0 {
				ch <- chunk.Choices[0].Delta.Content
			}
		}

		if err := stream.Err(); err != nil {
			errCh <- err
		}
	}()

	return ch, errCh
}

func (o *OpenAI) EmbeddingFunc() func(ctx context.Context, text string) ([]float32, error) {
	return o.ef
}
