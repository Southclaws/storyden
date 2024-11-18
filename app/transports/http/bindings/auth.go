package bindings

import (
	"context"
	"fmt"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"

	"github.com/Southclaws/storyden/app/resources/account/account_querier"
	"github.com/Southclaws/storyden/app/resources/account/authentication"
	"github.com/Southclaws/storyden/app/resources/account/email"
	"github.com/Southclaws/storyden/app/resources/settings"
	auth_svc "github.com/Southclaws/storyden/app/services/authentication"
	"github.com/Southclaws/storyden/app/services/authentication/email_verify"
	"github.com/Southclaws/storyden/app/services/authentication/provider/email_only"
	"github.com/Southclaws/storyden/app/services/authentication/provider/password"
	"github.com/Southclaws/storyden/app/transports/http/middleware/session"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
)

type Authentication struct {
	cj                            *session.Jar
	settings                      *settings.SettingsRepository
	passwordAuthProvider          *password.Provider
	emailVerificationAuthProvider *email_only.Provider
	accountQuery                  *account_querier.Querier
	emailRepo                     *email.Repository
	authManager                   *auth_svc.Manager
	emailVerifier                 *email_verify.Verifier
}

func NewAuthentication(
	cj *session.Jar,
	settings *settings.SettingsRepository,
	passwordAuthProvider *password.Provider,
	emailVerificationAuthProvider *email_only.Provider,
	accountQuery *account_querier.Querier,
	emailRepo *email.Repository,
	authManager *auth_svc.Manager,
	emailVerifier *email_verify.Verifier,
) Authentication {
	return Authentication{
		cj:                            cj,
		settings:                      settings,
		passwordAuthProvider:          passwordAuthProvider,
		emailVerificationAuthProvider: emailVerificationAuthProvider,
		accountQuery:                  accountQuery,
		emailRepo:                     emailRepo,
		authManager:                   authManager,
		emailVerifier:                 emailVerifier,
	}
}

func (o *Authentication) AuthProviderList(ctx context.Context, request openapi.AuthProviderListRequestObject) (openapi.AuthProviderListResponseObject, error) {
	settings, err := o.settings.Get(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	providers, err := o.authManager.GetProviderList(ctx)
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
			Provider: p.Service().String(),
			Name:     fmt.Sprintf("%v", p.Service()),
			Link:     &link,
		}, nil
	}

	return openapi.AuthProvider{
		Provider: p.Service().String(),
		Name:     fmt.Sprintf("%v", p.Service()),
	}, nil
}

func deserialiseAuthMode(in openapi.AuthMode) (authentication.Mode, error) {
	mode, err := authentication.NewMode(string(in))
	if err != nil {
		return authentication.Mode{}, fault.Wrap(err, ftag.With(ftag.InvalidArgument))
	}
	return mode, nil
}
