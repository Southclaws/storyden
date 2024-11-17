package authentication

import (
	"context"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/samber/lo"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/Southclaws/storyden/app/resources/account/authentication"
	"github.com/Southclaws/storyden/app/resources/settings"
	"github.com/Southclaws/storyden/app/services/authentication/provider/email_only"
	"github.com/Southclaws/storyden/app/services/authentication/provider/oauth/github"
	"github.com/Southclaws/storyden/app/services/authentication/provider/oauth/google"
	"github.com/Southclaws/storyden/app/services/authentication/provider/oauth/linkedin"
	"github.com/Southclaws/storyden/app/services/authentication/provider/password"
	"github.com/Southclaws/storyden/app/services/authentication/provider/phone"
	"github.com/Southclaws/storyden/app/services/authentication/provider/webauthn"
)

type Manager struct {
	settings  *settings.SettingsRepository
	providers map[authentication.Service]Provider
}

var ErrInvalidProvider = fault.New("invalid provider")

func Build() fx.Option {
	return fx.Options(
		fx.Provide(
			// All authentication provider services.
			password.NewEmailPasswordProvider,
			password.NewUsernamePasswordProvider,
			email_only.New,
			webauthn.New,
			google.New,
			github.New,
			linkedin.New,
			phone.New,
		),

		fx.Provide(New),
	)
}

func New(
	l *zap.Logger,
	settings *settings.SettingsRepository,

	pw *password.UsernamePasswordProvider,
	ep *password.EmailPasswordProvider,
	eo *email_only.Provider,
	wa *webauthn.Provider,
	gg *google.Provider,
	gh *github.Provider,
	li *linkedin.Provider,
	pp *phone.Provider,
) *Manager {
	providers := []Provider{
		pw,
		ep,
		eo,
		wa,
		gg,
		gh,
		li,
		pp,
	}

	l.Debug("initialised auth providers",
		zap.Strings("providers", dt.Map(providers, name)),
	)

	return &Manager{
		settings: settings,
		providers: lo.KeyBy(providers, func(p Provider) authentication.Service {
			return p.Service()
		}),
	}
}

func (oa *Manager) GetProviderList(ctx context.Context) ([]Provider, error) {
	providerList := lo.Values(oa.providers)

	filtered, err := dt.FilterErr(providerList, func(p Provider) (bool, error) {
		enabled, err := p.Enabled(ctx)
		if err != nil {
			return false, fault.Wrap(err, fctx.With(ctx))
		}

		return enabled, nil
	})
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return filtered, nil
}

func (oa *Manager) Provider(id authentication.Service) (Provider, error) {
	p, ok := oa.providers[id]
	if !ok {
		return nil, fault.Wrap(ErrInvalidProvider)
	}

	return p, nil
}

func service(p Provider) authentication.Service { return p.Service() }

func name(p Provider) string { return p.Service().String() }
