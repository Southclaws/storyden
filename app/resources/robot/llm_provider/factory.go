package llm_provider

import (
	"context"
	"fmt"
	"log/slog"
	"slices"
	"strings"
	"time"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/opt"
	"github.com/puzpuzpuz/xsync/v4"
	"google.golang.org/adk/model"

	"github.com/Southclaws/storyden/app/resources/robot/model_ref"
	"github.com/Southclaws/storyden/app/resources/robot/robot_model_cache"
	"github.com/Southclaws/storyden/app/resources/settings"
)

const (
	ModelCacheTTL = 24 * time.Hour
)

type Factory struct {
	logger    *slog.Logger
	settings  *settings.SettingsRepository
	cache     *robot_model_cache.Repository
	providers *xsync.Map[model_ref.Provider, Provider]
}

type RuntimeSettings struct {
	Enabled      bool
	DefaultModel opt.Optional[model_ref.ModelRef]
	Providers    map[model_ref.Provider]RuntimeProviderSettings
}

type RuntimeProviderSettings struct {
	Enabled bool
	APIKey  string
}

func New(
	logger *slog.Logger,
	settingsRepo *settings.SettingsRepository,
	modelCache *robot_model_cache.Repository,
) *Factory {
	return &Factory{
		logger:    logger,
		settings:  settingsRepo,
		cache:     modelCache,
		providers: xsync.NewMap[model_ref.Provider, Provider](),
	}
}

func (f *Factory) Put(provider Provider) {
	f.providers.Store(provider.Provider(), provider)
}

func (f *Factory) Delete(provider model_ref.Provider) {
	f.providers.Delete(provider)
}

func (f *Factory) Providers() []model_ref.Provider {
	providers := []model_ref.Provider{}
	f.providers.Range(func(provider model_ref.Provider, _ Provider) bool {
		providers = append(providers, provider)
		return true
	})
	slices.SortFunc(providers, func(a, b model_ref.Provider) int {
		return strings.Compare(a.String(), b.String())
	})
	return providers
}

func (f *Factory) HasProvider(provider model_ref.Provider) bool {
	_, ok := f.providers.Load(provider)
	return ok
}

func (f *Factory) RequiresAPIKey(provider model_ref.Provider) bool {
	impl, ok := f.providers.Load(provider)
	return !ok || impl.RequiresAPIKey()
}

func (f *Factory) RuntimeSettings(ctx context.Context) (RuntimeSettings, error) {
	set, err := f.settings.Get(ctx)
	if err != nil {
		return RuntimeSettings{}, fault.Wrap(err, fctx.With(ctx))
	}

	out := RuntimeSettings{
		Providers: map[model_ref.Provider]RuntimeProviderSettings{},
	}

	services, ok := set.Services.Get()
	if !ok {
		return out, nil
	}

	robots, ok := services.Robots.Get()
	if !ok {
		return out, nil
	}

	if enabled, ok := robots.Enabled.Get(); ok {
		out.Enabled = enabled
	}
	if defaultModel, ok := robots.DefaultModel.Get(); ok {
		ref, err := model_ref.ParseID(defaultModel)
		if err != nil {
			return RuntimeSettings{}, err
		}
		out.DefaultModel = opt.New(ref)
	}
	if providers, ok := robots.Providers.Get(); ok {
		for name, provider := range providers {
			providerName := model_ref.NewProvider(name)
			current := out.Providers[providerName]
			if enabled, ok := provider.Enabled.Get(); ok {
				current.Enabled = enabled
			}
			if key, ok := provider.APIKey.Get(); ok {
				current.APIKey = key
			}
			out.Providers[providerName] = current
		}
	}

	if defaultModel, ok := out.DefaultModel.Get(); ok {
		if _, ok := out.Providers[defaultModel.Provider]; !ok && f.HasProvider(defaultModel.Provider) {
			out.Providers[defaultModel.Provider] = RuntimeProviderSettings{Enabled: out.Enabled}
		}
	}

	return out, nil
}

func (f *Factory) DefaultModel(ctx context.Context) (model_ref.ModelRef, error) {
	settings, err := f.RuntimeSettings(ctx)
	if err != nil {
		return model_ref.ModelRef{}, err
	}
	if !settings.Enabled {
		return model_ref.ModelRef{}, fault.New("robots are disabled")
	}
	defaultModel, ok := settings.DefaultModel.Get()
	if !ok {
		return model_ref.ModelRef{}, fault.New("robots are enabled but no default model is configured")
	}
	return defaultModel, nil
}

func (f *Factory) GetADKModelLLM(ctx context.Context, ref model_ref.ModelRef) (model.LLM, error) {
	_, provider, err := f.provider(ctx, ref.Provider)
	if err != nil {
		return nil, err
	}

	if err := f.EnsureModelAvailable(ctx, ref); err != nil {
		return nil, err
	}

	return provider.GetADKModelLLM(ctx, ref)
}

func (f *Factory) EnsureModelAvailable(ctx context.Context, ref model_ref.ModelRef) error {
	provider, impl, err := f.provider(ctx, ref.Provider)
	if err != nil {
		return err
	}

	status, err := f.cache.GetStatus(ctx, ref.Provider)
	if err != nil {
		return err
	}

	stale := true
	if last, ok := status.LastRefreshedAt.Get(); ok {
		stale = time.Since(last) > ModelCacheTTL
	}

	if _, err := f.cache.GetModel(ctx, ref); err == nil {
		if stale {
			go f.refreshBestEffort(context.WithoutCancel(ctx), ref.Provider, provider.APIKey)
		}
		return nil
	}

	if err := impl.ValidateModel(ctx, ref); err == nil {
		return nil
	}

	if !stale {
		return fmt.Errorf("model %q is not available in the %s model cache", ref.Model, ref.Provider)
	}

	if _, err := f.RefreshProviderModels(ctx, ref.Provider); err != nil {
		return err
	}

	if _, err := f.cache.GetModel(ctx, ref); err != nil {
		return fmt.Errorf("model %q is not available for provider %s", ref.Model, ref.Provider)
	}

	return nil
}

func (f *Factory) ListCachedModels(ctx context.Context, providers []model_ref.Provider) ([]model_ref.Info, error) {
	return f.cache.ListActiveModels(ctx, providers)
}

func (f *Factory) ProviderStatus(ctx context.Context, provider model_ref.Provider) (model_ref.CacheStatus, []model_ref.Info, error) {
	status, err := f.cache.GetStatus(ctx, provider)
	if err != nil {
		return model_ref.CacheStatus{}, nil, err
	}

	models, err := f.cache.ListProviderModels(ctx, provider)
	if err != nil {
		return model_ref.CacheStatus{}, nil, err
	}

	if f.statusIsStale(status) {
		runtime, err := f.RuntimeSettings(ctx)
		if err == nil {
			impl, implErr := f.providerImpl(provider)
			if p, ok := runtime.Providers[provider]; ok && implErr == nil && p.Enabled && (p.APIKey != "" || !impl.RequiresAPIKey()) {
				go f.refreshBestEffort(context.WithoutCancel(ctx), provider, p.APIKey)
			}
		}
	}

	return status, models, nil
}

func (f *Factory) Validate(ctx context.Context) error {
	settings, err := f.RuntimeSettings(ctx)
	if err != nil {
		return err
	}
	defaultModel, ok := settings.DefaultModel.Get()
	if !settings.Enabled || !ok {
		return nil
	}

	return f.EnsureModelAvailable(ctx, defaultModel)
}

func (f *Factory) provider(ctx context.Context, providerName model_ref.Provider) (RuntimeProviderSettings, Provider, error) {
	settings, err := f.RuntimeSettings(ctx)
	if err != nil {
		return RuntimeProviderSettings{}, nil, err
	}
	if !settings.Enabled {
		return RuntimeProviderSettings{}, nil, fault.New("robots are disabled")
	}

	provider, ok := settings.Providers[providerName]
	if !ok {
		return RuntimeProviderSettings{}, nil, fmt.Errorf("unsupported robot model provider %q", providerName)
	}
	if !provider.Enabled {
		return RuntimeProviderSettings{}, nil, fmt.Errorf("robot model provider %q is disabled", providerName)
	}

	impl, err := f.providerImpl(providerName)
	if err != nil {
		return RuntimeProviderSettings{}, nil, fmt.Errorf("unsupported robot model provider %q", providerName)
	}
	if provider.APIKey == "" && impl.RequiresAPIKey() {
		return RuntimeProviderSettings{}, nil, fmt.Errorf("robot model provider %q has no API key configured", providerName)
	}

	impl.Configure(Config{
		Enabled: provider.Enabled,
		APIKey:  provider.APIKey,
	})

	return provider, impl, nil
}

func (f *Factory) providerImpl(provider model_ref.Provider) (Provider, error) {
	impl, ok := f.providers.Load(provider)
	if !ok {
		return nil, fmt.Errorf("unsupported robot model provider %q", provider)
	}
	return impl, nil
}

func (f *Factory) statusIsStale(status model_ref.CacheStatus) bool {
	last, ok := status.LastRefreshedAt.Get()
	return !ok || time.Since(last) > ModelCacheTTL
}

func ProviderSettingsFromOptional(provider settings.RobotProviderSettings) RuntimeProviderSettings {
	return RuntimeProviderSettings{
		Enabled: provider.Enabled.Or(false),
		APIKey:  provider.APIKey.Or(""),
	}
}

func OptionalProviderSettings(enabled opt.Optional[bool], apiKey opt.Optional[string]) settings.RobotProviderSettings {
	return settings.RobotProviderSettings{
		Enabled: enabled,
		APIKey:  apiKey,
	}
}
