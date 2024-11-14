package bindings

import (
	"context"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"

	"github.com/Southclaws/storyden/app/resources/account/account_querier"
	"github.com/Southclaws/storyden/app/resources/account/authentication"
	"github.com/Southclaws/storyden/app/resources/account/email"
	"github.com/Southclaws/storyden/app/resources/settings"
	auth_svc "github.com/Southclaws/storyden/app/services/authentication"
	"github.com/Southclaws/storyden/app/services/authentication/email_verify"
	"github.com/Southclaws/storyden/app/services/authentication/provider/email/email_only"
	"github.com/Southclaws/storyden/app/services/authentication/provider/email/email_password"
	"github.com/Southclaws/storyden/app/services/authentication/provider/password"
	session1 "github.com/Southclaws/storyden/app/transports/http/middleware/session"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/config"
)

type Authentication struct {
	settings     *settings.SettingsRepository
	p            *password.Provider
	ep           *email_only.Provider
	epp          *email_password.Provider
	cj           *session1.Jar
	accountQuery *account_querier.Querier
	er           email.EmailRepo
	am           *auth_svc.Manager
	ev           email_verify.Verifier
}

func NewAuthentication(
	settings *settings.SettingsRepository,
	cfg config.Config,
	p *password.Provider,
	ep *email_only.Provider,
	epp *email_password.Provider,
	accountQuery *account_querier.Querier,
	er email.EmailRepo,
	sm *session1.Jar,
	am *auth_svc.Manager,
	ev email_verify.Verifier,
) Authentication {
	return Authentication{settings, p, ep, epp, sm, accountQuery, er, am, ev}
}

func (o *Authentication) AuthProviderList(ctx context.Context, request openapi.AuthProviderListRequestObject) (openapi.AuthProviderListResponseObject, error) {
	settings, err := o.settings.Get(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	providers, err := o.am.GetProviderList(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	list, err := dt.MapErr(providers, serialiseAuthProvider)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	mode := settings.AuthenticationMode.Or(authentication.ModeHandle)

	return openapi.AuthProviderList200JSONResponse{
		AuthProviderListOKJSONResponse: openapi.AuthProviderListOKJSONResponse{
			Providers: list,
			Mode:      openapi.AuthMode(mode.String()),
		},
	}, nil
}

func (a *Authentication) AuthProviderLogout(ctx context.Context, request openapi.AuthProviderLogoutRequestObject) (openapi.AuthProviderLogoutResponseObject, error) {
	return openapi.AuthProviderLogout200Response{
		Headers: openapi.AuthProviderLogout200ResponseHeaders{
			SetCookie: a.cj.Destroy().String(),
		},
	}, nil
}

func serialiseAuthProvider(p auth_svc.Provider) (openapi.AuthProvider, error) {
	if op, ok := p.(auth_svc.OAuthProvider); ok {
		link, err := op.Link("/")
		if err != nil {
			return openapi.AuthProvider{}, fault.Wrap(err)
		}
		return openapi.AuthProvider{
			Provider: p.Provides().String(),
			Link:     link,
		}, nil
	}

	return openapi.AuthProvider{
		Provider: p.Provides().String(),
	}, nil
}
