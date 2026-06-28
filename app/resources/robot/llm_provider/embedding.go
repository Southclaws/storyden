package llm_provider

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/Southclaws/fault"

	"github.com/Southclaws/storyden/app/resources/robot/model_ref"
)

type EmbedFunc func(ctx context.Context, text string) ([]float32, error)

type EmbeddingProvider interface {
	SupportsEmbeddings() bool
	EmbedText(ctx context.Context, text string) ([]float32, error)
}

func (f *Factory) GetEmbedder(ctx context.Context) (EmbedFunc, error) {
	settings, err := f.RuntimeSettings(ctx)
	if err != nil {
		return nil, err
	}
	if !settings.Enabled {
		return nil, fault.New("semdex requires robots model providers to be enabled")
	}

	var candidateErr error
	for _, provider := range f.embeddingCandidates(settings) {
		config, ok := settings.Providers[provider]
		if !ok || !config.Enabled {
			continue
		}

		impl, err := f.providerImpl(provider)
		if err != nil {
			continue
		}

		embedder, ok := impl.(EmbeddingProvider)
		if !ok || !embedder.SupportsEmbeddings() {
			continue
		}

		_, configured, err := f.provider(ctx, provider)
		if err != nil {
			candidateErr = err
			continue
		}

		embedder, ok = configured.(EmbeddingProvider)
		if !ok || !embedder.SupportsEmbeddings() {
			return nil, fmt.Errorf("robot model provider %q no longer supports embeddings", provider)
		}

		return embedder.EmbedText, nil
	}

	if candidateErr != nil {
		return nil, fmt.Errorf("no enabled robot model provider with usable embeddings: %w", candidateErr)
	}

	return nil, fault.New("no enabled robot model provider supports embeddings")
}

func (f *Factory) embeddingCandidates(settings RuntimeSettings) []model_ref.Provider {
	providers := make([]model_ref.Provider, 0, len(settings.Providers))
	seen := map[model_ref.Provider]struct{}{}

	if defaultModel, ok := settings.DefaultModel.Get(); ok {
		providers = append(providers, defaultModel.Provider)
		seen[defaultModel.Provider] = struct{}{}
	}

	for provider := range settings.Providers {
		if _, ok := seen[provider]; ok {
			continue
		}
		providers = append(providers, provider)
	}

	slices.SortFunc(providers, func(a, b model_ref.Provider) int {
		_, aDefault := seen[a]
		_, bDefault := seen[b]
		if aDefault != bDefault {
			if aDefault {
				return -1
			}
			return 1
		}
		return strings.Compare(a.String(), b.String())
	})

	return providers
}
