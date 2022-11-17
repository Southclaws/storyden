package oauth

import (
	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/samber/lo"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/Southclaws/storyden/app/services/authentication/provider/oauth/github"
	"github.com/Southclaws/storyden/app/services/authentication/provider/oauth/linkedin"
)

var ErrInvalidProvider = fault.New("invalid provider")

type OAuth struct {
	providers map[string]Provider
}

// Adding a new OAuth2 provider?
//
// 1. Add the constructor to the `fx.Provide` call in the builder.
// 2. Add the instance to `allProviders` in the OAuth constructor.
//
// See lines annotated with (1) and (2) below...
//
func Build() fx.Option {
	return fx.Provide(
		New,
		// (1)
		// All OAuth2 providers are initialised, those that fail are disabled.
		github.New,
		linkedin.New,
	)
}

func New(
	l *zap.Logger,
	gh *github.GitHubProvider,
	li *linkedin.LinkedInProvider,
) *OAuth {
	allProviders := []Provider{
		// (2)
		// All OAuth2 providers are statically added to this list regardless of
		// whether they are enabled or not. Disabled providers are filtered out.
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

	return &OAuth{
		providers: lo.KeyBy(enabledProviders, func(p Provider) string {
			return p.ID()
		}),
	}
}

func (oa *OAuth) Providers() []Provider {
	return lo.Values(oa.providers)
}

func (oa *OAuth) Provider(id string) (Provider, error) {
	p, ok := oa.providers[id]
	if !ok {
		return nil, fault.Wrap(ErrInvalidProvider)
	}

	return p, nil
}

func name(p Provider) string { return p.ID() }
