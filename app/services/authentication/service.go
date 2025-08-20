package authentication

import (
	"context"
	"log/slog"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/samber/lo"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account/authentication"
	"github.com/Southclaws/storyden/app/resources/settings"
	"github.com/Southclaws/storyden/app/services/authentication/email_verify"
	"github.com/Southclaws/storyden/app/services/authentication/provider/email_only"
	"github.com/Southclaws/storyden/app/services/authentication/provider/oauth/discord"
	"github.com/Southclaws/storyden/app/services/authentication/provider/oauth/github"
	"github.com/Southclaws/storyden/app/services/authentication/provider/oauth/google"
	"github.com/Southclaws/storyden/app/services/authentication/provider/oauth/keycloak"
	"github.com/Southclaws/storyden/app/services/authentication/provider/password"
	"github.com/Southclaws/storyden/app/services/authentication/provider/password/password_reset"
	"github.com/Southclaws/storyden/app/services/authentication/provider/phone"
	"github.com/Southclaws/storyden/app/services/authentication/provider/webauthn"
	"github.com/Southclaws/storyden/app/services/authentication/session"
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
			password.New,
			email_only.New,
			webauthn.New,
			google.New,
			github.New,
			discord.New,
			keycloak.New,
			phone.New,
		),
		fx.Provide(email_verify.New),
		fx.Provide(password_reset.NewTokenProvider, password_reset.NewEmailResetter),
		fx.Provide(New, session.NewValidator, session.NewIssuer),
	)
}

func New(
	logger *slog.Logger,
	settings *settings.SettingsRepository,

	pw *password.Provider,
	eo *email_only.Provider,
	wa *webauthn.Provider,
	gg *google.Provider,
	gh *github.Provider,
	dp *discord.Provider,
	kc *keycloak.Provider,
	pp *phone.Provider,
) *Manager {
	providers := []Provider{
		pw,
		eo,
		wa,
		gg,
		gh,
		dp,
		kc,
		pp,
	}

	logger.Debug("initialised auth providers",
		slog.Any("providers", dt.Map(providers, name)),
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
