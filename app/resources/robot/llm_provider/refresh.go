package llm_provider

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"

	"github.com/Southclaws/storyden/app/resources/robot/model_ref"
)

func (f *Factory) RefreshProviderModels(ctx context.Context, providerName model_ref.Provider) ([]model_ref.Info, error) {
	_, impl, err := f.provider(ctx, providerName)
	if err != nil {
		return nil, err
	}

	return f.refreshProviderModelsWithProvider(ctx, providerName, impl)
}

func (f *Factory) RefreshProviderModelsWithKey(ctx context.Context, providerName model_ref.Provider, apiKey string) ([]model_ref.Info, error) {
	provider, err := f.providerImpl(providerName)
	if err != nil {
		return nil, err
	}

	if apiKey == "" && provider.RequiresAPIKey() {
		return nil, fmt.Errorf("robot model provider %q has no API key configured", providerName)
	}
	provider.Configure(Config{
		Enabled: true,
		APIKey:  apiKey,
	})

	return f.refreshProviderModelsWithProvider(ctx, providerName, provider)
}

func (f *Factory) refreshProviderModelsWithProvider(ctx context.Context, providerName model_ref.Provider, provider Provider) ([]model_ref.Info, error) {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	models, err := provider.ListModels(ctx)
	if err != nil {
		f.logger.Warn("failed to refresh robot provider models", slog.String("provider", providerName.String()), slog.String("error", err.Error()))
		if cacheErr := f.cache.SetError(context.WithoutCancel(ctx), providerName, err.Error()); cacheErr != nil {
			f.logger.Error("failed to record robot model refresh error", slog.String("provider", providerName.String()), slog.String("error", cacheErr.Error()))
		}
		return nil, fault.Wrap(err, fctx.With(ctx), fmsg.Withf("refresh %s models", providerName))
	}

	if err := f.cache.UpsertProviderModels(ctx, providerName, models); err != nil {
		return nil, err
	}

	if err := f.cache.SetRefreshed(ctx, providerName, time.Now()); err != nil {
		return nil, err
	}

	return f.cache.ListProviderModels(ctx, providerName)
}

func (f *Factory) refreshBestEffort(ctx context.Context, provider model_ref.Provider, apiKey string) {
	impl, err := f.providerImpl(provider)
	if err != nil || (apiKey == "" && impl.RequiresAPIKey()) {
		return
	}
	impl.Configure(Config{
		Enabled: true,
		APIKey:  apiKey,
	})

	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	models, err := impl.ListModels(ctx)
	if err != nil {
		if cacheErr := f.cache.SetError(context.WithoutCancel(ctx), provider, err.Error()); cacheErr != nil {
			f.logger.Error("failed to record robot model refresh error", slog.String("provider", provider.String()), slog.String("error", cacheErr.Error()))
		}
		return
	}

	if err := f.cache.UpsertProviderModels(ctx, provider, models); err != nil {
		f.logger.Error("failed to refresh robot model cache", slog.String("provider", provider.String()), slog.String("error", err.Error()))
		return
	}

	if err := f.cache.SetRefreshed(ctx, provider, time.Now()); err != nil {
		f.logger.Error("failed to record robot model refresh", slog.String("provider", provider.String()), slog.String("error", err.Error()))
	}
}
