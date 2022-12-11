package authentication

import (
	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/samber/lo"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/Southclaws/storyden/app/services/authentication/provider/oauth/github"
	"github.com/Southclaws/storyden/app/services/authentication/provider/oauth/linkedin"
	"github.com/Southclaws/storyden/app/services/authentication/provider/password"
	"github.com/Southclaws/storyden/app/services/authentication/provider/webauthn"
)

type Manager struct {
	providers map[string]Provider
}

var ErrInvalidProvider = fault.New("invalid provider")

// Adding a new OAuth2 provider?
//
// 1. Add the constructor to the `fx.Provide` call in the builder.
// 2. Add the instance to `allProviders` in the Manager constructor.
//
// See lines annotated with (1) and (2) below...
//
func Build() fx.Option {
	return fx.Options(
		fx.Provide(
			// (1)
			// All auth providers are initialised, those that fail are disabled.
			password.New,
			webauthn.New,
			github.New,
			linkedin.New,
		),

		fx.Provide(New),
	)
}

func New(
	l *zap.Logger,
	wa *webauthn.Provider,
	gh *github.GitHubProvider,
	li *linkedin.LinkedInProvider,
) *Manager {
	allProviders := []Provider{
		// (2)
		// All OAuth2 providers are statically added to this list regardless of
		// whether they are enabled or not. Disabled providers are filtered out.
		wa,
		gh,
		li,
	}

	// Filter out disabled providers.
	enabledProviders := lo.Filter(allProviders, func(p Provider, _ int) bool {
		return p.Enabled()
	})

	l.Info("initialised oauth providers",
		zap.Strings("all_providers", dt.Map(allProviders, name)),
		zap.Strings("enabled_providers", dt.Map(enabledProviders, name)),
	)

	return &Manager{
		providers: lo.KeyBy(enabledProviders, func(p Provider) string {
			return p.ID()
		}),
	}
}

func (oa *Manager) Providers() []Provider {
	return lo.Values(oa.providers)
}

func (oa *Manager) Provider(id string) (Provider, error) {
	p, ok := oa.providers[id]
	if !ok {
		return nil, fault.Wrap(ErrInvalidProvider)
	}

	return p, nil
}

func name(p Provider) string { return p.ID() }
