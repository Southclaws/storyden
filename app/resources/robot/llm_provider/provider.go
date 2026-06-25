package llm_provider

import (
	"context"

	"google.golang.org/adk/model"

	"github.com/Southclaws/storyden/app/resources/robot/model_ref"
)

// built-in llm model providers
var (
	ProviderOpenAI    = model_ref.NewProvider("openai")
	ProviderAnthropic = model_ref.NewProvider("anthropic")
	ProviderMock      = model_ref.NewProvider("mock")
)

type Config struct {
	Enabled bool
	APIKey  string
}

type Provider interface {
	Provider() model_ref.Provider
	RequiresAPIKey() bool
	Configure(Config)
	ListModels(ctx context.Context) ([]model_ref.Info, error)
	GetADKModelLLM(ctx context.Context, ref model_ref.ModelRef) (model.LLM, error)
	ValidateModel(ctx context.Context, ref model_ref.ModelRef) error
}
