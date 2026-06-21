package openai

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/openai/openai-go/v3"
	openaioption "github.com/openai/openai-go/v3/option"
	"google.golang.org/adk/model"

	"github.com/Southclaws/storyden/app/resources/robot/llm_provider"
	"github.com/Southclaws/storyden/app/resources/robot/model_ref"
)

var Provider = llm_provider.ProviderOpenAI

type OpenAI struct {
	mu        sync.RWMutex
	apiKey    string
	client    *openai.Client
	modelName string
}

func (p *OpenAI) Name() string {
	return "openai/" + p.modelName
}

func (m *OpenAI) Validate(ctx context.Context) error {
	_, err := m.client.Models.Get(ctx, m.modelName)
	if err != nil {
		return fmt.Errorf("model %q not available: %w", m.modelName, err)
	}
	return nil
}

func (*OpenAI) Provider() model_ref.Provider { return Provider }

func (*OpenAI) RequiresAPIKey() bool { return true }

func (p *OpenAI) Configure(config llm_provider.Config) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.apiKey = config.APIKey
}

func (p *OpenAI) ListModels(ctx context.Context) ([]model_ref.Info, error) {
	p.mu.RLock()
	apiKey := p.apiKey
	p.mu.RUnlock()

	client := openai.NewClient(openaioption.WithAPIKey(apiKey))
	page, err := client.Models.List(ctx)
	if err != nil {
		return nil, err
	}

	out := make([]model_ref.Info, 0, len(page.Data))
	for _, item := range page.Data {
		raw := map[string]any{}
		_ = json.Unmarshal([]byte(item.RawJSON()), &raw)

		out = append(out, model_ref.Info{
			Ref: model_ref.ModelRef{
				Provider: Provider,
				Model:    model_ref.NewModel(item.ID),
			},
			Raw: raw,
		})
	}

	return out, nil
}

func (p *OpenAI) GetADKModelLLM(ctx context.Context, ref model_ref.ModelRef) (model.LLM, error) {
	p.mu.RLock()
	apiKey := p.apiKey
	p.mu.RUnlock()

	client := openai.NewClient(openaioption.WithAPIKey(apiKey))

	return &OpenAI{
		apiKey:    apiKey,
		client:    &client,
		modelName: ref.Model.String(),
	}, nil
}

func (p *OpenAI) ValidateModel(ctx context.Context, ref model_ref.ModelRef) error {
	p.mu.RLock()
	apiKey := p.apiKey
	p.mu.RUnlock()

	client := openai.NewClient(openaioption.WithAPIKey(apiKey))
	_, err := client.Models.Get(ctx, ref.Model.String())
	if err != nil {
		return err
	}
	return nil
}
