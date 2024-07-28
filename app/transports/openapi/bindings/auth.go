package bindings

import (
	"context"
	"net/http"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/oapi-codegen/echo-middleware"

	"github.com/Southclaws/storyden/app/resources/account/account_querier"
	"github.com/Southclaws/storyden/app/resources/account/email"
	"github.com/Southclaws/storyden/app/services/authentication"
	"github.com/Southclaws/storyden/app/services/authentication/email_verify"
	"github.com/Southclaws/storyden/app/services/authentication/provider/email/email_only"
	"github.com/Southclaws/storyden/app/services/authentication/provider/email/email_password"
	"github.com/Southclaws/storyden/app/services/authentication/provider/password"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/app/transports/openapi"
	"github.com/Southclaws/storyden/internal/config"
)

type Authentication struct {
	p            *password.Provider
	ep           *email_only.Provider
	epp          *email_password.Provider
	sm           *CookieJar
	accountQuery account_querier.Querier
	er           email.EmailRepo
	am           *authentication.Manager
	ev           email_verify.Verifier
	domain       string
}

func NewAuthentication(
	cfg config.Config,
	p *password.Provider,
	ep *email_only.Provider,
	epp *email_password.Provider,
	accountQuery account_querier.Querier,
	er email.EmailRepo,
	sm *CookieJar,
	am *authentication.Manager,
	ev email_verify.Verifier,
) Authentication {
	return Authentication{p, ep, epp, sm, accountQuery, er, am, ev, cfg.CookieDomain}
}

func (o *Authentication) AuthProviderList(ctx context.Context, request openapi.AuthProviderListRequestObject) (openapi.AuthProviderListResponseObject, error) {
	list, err := dt.MapErr(o.am.Providers(), serialiseAuthProvider)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.AuthProviderList200JSONResponse{
		AuthProviderListOKJSONResponse: openapi.AuthProviderListOKJSONResponse{
			Providers: list,
		},
	}, nil
}

func (a *Authentication) AuthProviderLogout(ctx context.Context, request openapi.AuthProviderLogoutRequestObject) (openapi.AuthProviderLogoutResponseObject, error) {
	return openapi.AuthProviderLogout200Response{
		Headers: openapi.AuthProviderLogout200ResponseHeaders{
			SetCookie: (&http.Cookie{
				Name:     secureCookieName,
				Value:    "",
				SameSite: http.SameSiteDefaultMode,
				Path:     "/",
				Domain:   a.domain,
				Secure:   true,
				HttpOnly: true,
			}).String(),
		},
	}, nil
}

func (i *Authentication) validator(ctx context.Context, ai *openapi3filter.AuthenticationInput) error {
	// security scheme name from openapi.yaml
	if ai.SecuritySchemeName != "browser" {
		return nil
	}

	c := ctx.Value(echomiddleware.EchoContextKey).(echo.Context)

	// first check if the middleware injected an account ID, if not, fail.
	aid, err := session.GetAccountID(c.Request().Context())
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Unauthenticated))
	}

	// Then look up the account.
	// TODO: Cache this.
	a, err := i.accountQuery.GetByID(ctx, aid)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	// Reject any requests from suspended accounts.
	if err := a.RejectSuspended(); err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}

func serialiseAuthProvider(p authentication.Provider) (openapi.AuthProvider, error) {
	link, err := p.Link("/")
	if err != nil {
		return openapi.AuthProvider{}, fault.Wrap(err)
	}

	return openapi.AuthProvider{
		Provider: p.ID(),
		Name:     p.Name(),
		Link:     link,
	}, nil
}
