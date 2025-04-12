package instance_info

import (
	"context"
	"log/slog"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/opt"

	"github.com/Southclaws/storyden/app/resources/account/authentication"
	"github.com/Southclaws/storyden/app/resources/settings"
	"github.com/Southclaws/storyden/app/services/onboarding"
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
	Settings         *settings.Settings
	OnboardingStatus *onboarding.Status
	Capabilities     Capabilities
}

func (p *Provider) Get(ctx context.Context) (*Info, error) {
	settings, err := p.settings.Get(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	settings, err = p.selfHeal(ctx, settings)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to self-heal malformed settings"))
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

	return &Info{
		Settings:         settings,
		OnboardingStatus: status,
		Capabilities:     caps,
	}, nil
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
