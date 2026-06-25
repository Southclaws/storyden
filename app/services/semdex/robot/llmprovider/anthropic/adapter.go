package anthropic

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/anthropics/anthropic-sdk-go"
	anthropicoption "github.com/anthropics/anthropic-sdk-go/option"
	"google.golang.org/adk/model"

	"github.com/Southclaws/storyden/app/resources/robot/llm_provider"
	"github.com/Southclaws/storyden/app/resources/robot/model_ref"
)

var Provider = llm_provider.ProviderAnthropic

type Anthropic struct {
	mu        sync.RWMutex
	apiKey    string
	client    *anthropic.Client
	modelName string
}

func (m *Anthropic) Name() string {
	return "anthropic/" + m.modelName
}

func (m *Anthropic) Validate(ctx context.Context) error {
	_, err := m.client.Models.Get(ctx, m.modelName, anthropic.ModelGetParams{})
	if err != nil {
		return fmt.Errorf("model %q not available: %w", m.modelName, mapError(err))
	}
	return nil
}

func (*Anthropic) Provider() model_ref.Provider { return Provider }

func (*Anthropic) RequiresAPIKey() bool { return true }

func (p *Anthropic) Configure(config llm_provider.Config) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.apiKey = config.APIKey
}

func (p *Anthropic) ListModels(ctx context.Context) ([]model_ref.Info, error) {
	p.mu.RLock()
	apiKey := p.apiKey
	p.mu.RUnlock()

	client := anthropic.NewClient(anthropicoption.WithAPIKey(apiKey))
	page, err := client.Models.List(ctx, anthropic.ModelListParams{})
	if err != nil {
		return nil, mapError(err)
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

func (p *Anthropic) GetADKModelLLM(ctx context.Context, ref model_ref.ModelRef) (model.LLM, error) {
	p.mu.RLock()
	apiKey := p.apiKey
	p.mu.RUnlock()

	client := anthropic.NewClient(anthropicoption.WithAPIKey(apiKey))

	return &Anthropic{
		apiKey:    apiKey,
		client:    &client,
		modelName: ref.Model.String(),
	}, nil
}

func (p *Anthropic) ValidateModel(ctx context.Context, ref model_ref.ModelRef) error {
	p.mu.RLock()
	apiKey := p.apiKey
	p.mu.RUnlock()

	client := anthropic.NewClient(anthropicoption.WithAPIKey(apiKey))
	_, err := client.Models.Get(ctx, ref.Model.String(), anthropic.ModelGetParams{})
	if err != nil {
		return mapError(err)
	}
	return nil
}
