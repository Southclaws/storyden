package openai

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/openai/openai-go/v3"
	openaioption "github.com/openai/openai-go/v3/option"
)

const embeddingModel = openai.EmbeddingModelTextEmbedding3Large

func (*OpenAI) SupportsEmbeddings() bool { return true }

func (p *OpenAI) EmbedText(ctx context.Context, text string) ([]float32, error) {
	p.mu.RLock()
	apiKey := p.apiKey
	p.mu.RUnlock()

	client := openai.NewClient(openaioption.WithAPIKey(apiKey))

	res, err := client.Embeddings.New(ctx, openai.EmbeddingNewParams{
		Input: openai.EmbeddingNewParamsInputUnion{
			OfString: openai.String(text),
		},
		Model:          embeddingModel,
		EncodingFormat: openai.EmbeddingNewParamsEncodingFormatFloat,
	})
	if err != nil {
		return nil, fault.Wrap(mapError(err), fctx.With(ctx))
	}
	if len(res.Data) == 0 {
		return nil, fault.New("embedding response is empty")
	}

	embedding := make([]float32, len(res.Data[0].Embedding))
	for i, v := range res.Data[0].Embedding {
		embedding[i] = float32(v)
	}

	return embedding, nil
}
