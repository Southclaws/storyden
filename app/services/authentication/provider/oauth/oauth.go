package oauth

import (
	"errors"

	"github.com/Southclaws/dt"
	"github.com/samber/lo"
	"go.uber.org/zap"

	"github.com/Southclaws/storyden/app/services/authentication/provider/oauth/github"
)

var ErrInvalidProvider = errors.New("invalid provider")

type OAuth struct {
	providers map[string]Provider
}

func New(
	l *zap.Logger,
	gh *github.GitHubProvider,
) *OAuth {
	allProviders := []Provider{
		// All OAuth2 providers are statically added to this list regardless of
		// whether they are enabled or not. Disabled providers are filtered out.
		gh,
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
		return nil, ErrInvalidProvider
	}

	return p, nil
}

func name(p Provider) string { return p.ID() }
