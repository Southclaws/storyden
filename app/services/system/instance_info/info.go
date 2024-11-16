package instance_info

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"

	"github.com/Southclaws/storyden/app/resources/settings"
	"github.com/Southclaws/storyden/app/services/onboarding"
	"github.com/Southclaws/storyden/internal/config"
)

type Provider struct {
	config     config.Config
	settings   *settings.SettingsRepository
	onboarding onboarding.Service
}

func New(
	config config.Config,
	settings *settings.SettingsRepository,
	onboarding onboarding.Service,
) *Provider {
	return &Provider{
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

	status, err := p.onboarding.GetOnboardingStatus(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	caps := Capabilities{}

	if p.config.SemdexEnabled {
		caps = append(caps, CapabilitySemdex)
	}

	if p.config.EmailProvider != "" {
		caps = append(caps, CapabilityEmailClient)
	}

	// TODO: SMS client check

	return &Info{
		Settings:         settings,
		OnboardingStatus: status,
		Capabilities:     caps,
	}, nil
}
