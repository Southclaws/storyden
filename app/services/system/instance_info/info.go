package instance_info

import (
	"context"
	"log/slog"
	"net/url"
	"strings"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/opt"

	"github.com/Southclaws/storyden/app/resources/account/authentication"
	"github.com/Southclaws/storyden/app/resources/asset"
	"github.com/Southclaws/storyden/app/resources/settings"
	"github.com/Southclaws/storyden/app/services/onboarding"
	"github.com/Southclaws/storyden/app/services/plugin/plugin_runner"
	"github.com/Southclaws/storyden/internal/config"
)

type Provider struct {
	logger     *slog.Logger
	config     config.Config
	settings   *settings.SettingsRepository
	onboarding onboarding.Service
}

func New(
	logger *slog.Logger,
	config config.Config,
	settings *settings.SettingsRepository,
	onboarding onboarding.Service,
) *Provider {
	return &Provider{
		logger:     logger,
		config:     config,
		settings:   settings,
		onboarding: onboarding,
	}
}

type Info struct {
	Media            MediaServiceInfo
	Settings         *settings.Settings
	OnboardingStatus *onboarding.Status
	Capabilities     Capabilities
	WebAddress       url.URL
	APIAddress       url.URL
}

type InfoWithAssets struct {
	*Info
	BannerAsset opt.Optional[asset.Asset]
}

func (p *Provider) GetWithAssets(ctx context.Context) (*InfoWithAssets, error) {
	expanded, err := p.settings.GetWithAssets(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}
	info, err := p.get(ctx, expanded.Settings)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}
	return &InfoWithAssets{Info: info, BannerAsset: expanded.BannerAsset}, nil
}

func (p *Provider) Get(ctx context.Context) (*Info, error) {
	settings, err := p.settings.Get(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return p.get(ctx, settings)
}

func (p *Provider) get(ctx context.Context, settings *settings.Settings) (*Info, error) {
	settings, err := p.selfHeal(ctx, settings)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to self-heal malformed settings"))
	}

	limits, err := settings.ImageLimits()
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	status, err := p.onboarding.GetOnboardingStatus(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	caps := Capabilities{}

	if p.config.LanguageModelProvider != "" {
		caps = append(caps, CapabilityGenAI)
	}

	if p.config.SemdexProvider != "" {
		caps = append(caps, CapabilitySemdex)
	}

	if p.config.EmailProvider != "" {
		caps = append(caps, CapabilityEmailClient)
	}

	if p.config.SMSProvider != "" {
		caps = append(caps, CapabilitySMSClient)
	}

	if p.config.OAuthEnabled {
		caps = append(caps, CapabilityOAuth)
	}

	if services, ok := settings.Services.Get(); ok {
		if robots, ok := services.Robots.Get(); ok && robotsAvailable(robots) {
			caps = append(caps, CapabilityRobots)
		}
	}

	if p.config.PluginRuntimeProvider != plugin_runner.RuntimeProviderNone.String() {
		caps = append(caps, CapabilityPlugins)
	}

	return &Info{
		Media: MediaServiceInfo{
			ImageMemoryBudget:    p.config.ImageMemoryBudget,
			ImageEstimatedMemory: estimateImageMemory(limits),
		},
		Settings:         settings,
		OnboardingStatus: status,
		Capabilities:     caps,
		WebAddress:       p.config.PublicWebAddress,
		APIAddress:       p.config.PublicAPIAddress,
	}, nil
}

func robotsAvailable(robots settings.RobotServiceSettings) bool {
	defaultModel, ok := robots.DefaultModel.Get()
	if !ok || defaultModel == "" {
		return false
	}

	providers, ok := robots.Providers.Get()
	if !ok {
		return false
	}

	providerName, _, ok := strings.Cut(defaultModel, "/")
	if !ok || providerName == "" {
		return false
	}

	provider, ok := providers[providerName]
	return ok && provider.Enabled.Or(false)
}

func (p Provider) selfHeal(ctx context.Context, set *settings.Settings) (*settings.Settings, error) {
	var err error
	switch set.AuthenticationMode.OrZero() {
	case authentication.ModeEmail:
		if p.config.EmailProvider == "" {
			p.logger.Warn("Email authentication mode is enabled but no email provider is configured - resetting auth mode to handle")
			set.AuthenticationMode = opt.New(authentication.ModeHandle)
			set, err = p.settings.Set(ctx, *set)
		}

	case authentication.ModePhone:
		// TODO: Implement phone auth check
		// if p.config.SMSProvider == "" {
		// 	p.logger.Warn("Phone authentication mode is enabled but no SMS provider is configured - resetting auth mode to handle")
		// 	current.AuthenticationMode = opt.New(authentication.ModeHandle)
		// 	current, err = p.settings.Set(ctx, *current)
		// }
	}

	return set, err
}
