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

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/services/authentication"
	"github.com/Southclaws/storyden/app/services/authentication/provider/email_only"
	"github.com/Southclaws/storyden/app/services/authentication/provider/password"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/app/transports/openapi"
	"github.com/Southclaws/storyden/internal/config"
)

type Authentication struct {
	p      *password.Provider
	ep     *email_only.Provider
	sm     *CookieJar
	ar     account.Repository
	am     *authentication.Manager
	domain string
}

func NewAuthentication(
	cfg config.Config,
	p *password.Provider,
	ep *email_only.Provider,
	ar account.Repository,
	sm *CookieJar,
	am *authentication.Manager,
) Authentication {
	return Authentication{p, ep, sm, ar, am, cfg.CookieDomain}
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
	a, err := i.ar.GetByID(ctx, aid)
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
