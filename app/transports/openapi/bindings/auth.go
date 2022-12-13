package bindings

import (
	"context"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/deepmap/oapi-codegen/pkg/middleware"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/labstack/echo/v4"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/services/authentication"
	"github.com/Southclaws/storyden/app/services/authentication/provider/password"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/openapi"
)

type Authentication struct {
	p      *password.Password
	sm     *Session
	ar     account.Repository
	am     *authentication.Manager
	wa     *webauthn.WebAuthn
	domain string
}

func NewAuthentication(
	cfg config.Config,
	p *password.Password,
	ar account.Repository,
	sm *Session,
	am *authentication.Manager,
	wa *webauthn.WebAuthn,
) Authentication {
	return Authentication{p, sm, ar, am, wa, cfg.CookieDomain}
}

func (o *Authentication) AuthProviderList(ctx context.Context, request openapi.AuthProviderListRequestObject) (openapi.AuthProviderListResponseObject, error) {
	list := dt.Map(o.am.Providers(),
		func(p authentication.Provider) openapi.AuthProvider {
			return openapi.AuthProvider{
				Provider: p.ID(),
				Name:     p.Name(),
				LogoUrl:  p.LogoURL(),
				Link:     p.Link(),
			}
		},
	)

	return openapi.AuthProviderListJSONResponse(list), nil
}

func (i *Authentication) middleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		r := c.Request()
		ctx := r.Context()

		session, ok := i.sm.decodeSession(r)
		if ok {
			c.SetRequest(r.WithContext(authentication.WithAccountID(ctx, session.UserID)))
		}

		return next(c)
	}
}

func (i *Authentication) validator(ctx context.Context, ai *openapi3filter.AuthenticationInput) error {
	c := ctx.Value(middleware.EchoContextKey).(echo.Context)

	// first check if the middleware injected an account ID, if not, fail.
	aid, err := authentication.GetAccountID(c.Request().Context())
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.Unauthenticated))
	}

	// Then look up the account.
	// TODO: Cache this.
	_, err = i.ar.GetByID(ctx, aid)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}
