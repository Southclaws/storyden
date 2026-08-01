package ai

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/responses"
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
	return &OpenAI{client: &client, ef: ef}, nil
}

func (o *OpenAI) Prompt(ctx context.Context, input string) (*Result, error) {
	res, err := o.client.Responses.New(ctx, responses.ResponseNewParams{
		Model: openai.ChatModelGPT4_1,
		Input: responses.ResponseNewParamsInputUnion{
			OfString: openai.String(input),
		},
	})
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if res.OutputText() == "" {
		return nil, fault.New("result is empty")
	}

	return &Result{
		Answer: res.OutputText(),
	}, nil
}

func (o *OpenAI) PromptStream(ctx context.Context, input string) (func(yield func(string, error) bool), error) {
	iter := func(yield func(string, error) bool) {
		stream := o.client.Responses.NewStreaming(ctx, responses.ResponseNewParams{
			Model: openai.ChatModelGPT4_1,
			Input: responses.ResponseNewParamsInputUnion{
				OfString: openai.String(input),
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

			event := stream.Current()

			if event.Type == "response.output_text.delta" {
				if !yield(event.Delta, nil) {
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
