package bindings

import (
	"context"
	"net/http"

	"github.com/deepmap/oapi-codegen/pkg/middleware"
	"github.com/duo-labs/webauthn/webauthn"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/gorilla/securecookie"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/services/authentication"
	"github.com/Southclaws/storyden/app/services/authentication/provider/oauth"
	"github.com/Southclaws/storyden/app/services/authentication/provider/password"
	"github.com/Southclaws/storyden/app/transports/openapi/openapi"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/errctx"
	"github.com/Southclaws/storyden/internal/errtag"
)

type Authentication struct {
	p      *password.Password
	sc     *securecookie.SecureCookie
	ar     account.Repository
	oa     *oauth.OAuth
	wa     *webauthn.WebAuthn
	domain string
}

func NewAuthentication(
	cfg config.Config,
	p *password.Password,
	ar account.Repository,
	sc *securecookie.SecureCookie,
	oa *oauth.OAuth,
	wa *webauthn.WebAuthn,
) Authentication {
	return Authentication{p, sc, ar, oa, wa, cfg.CookieDomain}
}

func (i *Authentication) AuthPasswordSignin(ctx context.Context, request openapi.AuthPasswordSigninRequestObject) (openapi.AuthPasswordSigninResponseObject, error) {
	params := func() openapi.AuthPassword {
		if request.JSONBody != nil {
			return *request.JSONBody
		} else {
			return *request.FormdataBody
		}
	}()

	u, err := i.p.Login(ctx, params.Identifier, params.Token)
	if err != nil {
		return nil, errors.Wrap(err, "login failed")
	}

	cookie, err := i.encodeSession(u.ID)
	if err != nil {
		return nil, err
	}

	return openapi.AuthPasswordSignin200JSONResponse{
		Body:    openapi.AuthSuccess{Id: u.ID.String()},
		Headers: openapi.AuthSuccessResponseHeaders{SetCookie: cookie},
	}, nil
}

func (i *Authentication) AuthPasswordSignup(ctx context.Context, request openapi.AuthPasswordSignupRequestObject) (openapi.AuthPasswordSignupResponseObject, error) {
	params := func() openapi.AuthPassword {
		if request.JSONBody != nil {
			return *request.JSONBody
		} else {
			return *request.FormdataBody
		}
	}()

	u, err := i.p.Register(ctx, params.Identifier, params.Token)
	if err != nil {
		return nil, errctx.Wrap(err, ctx, "identifier", params.Identifier)
	}

	cookie, err := i.encodeSession(u.ID)
	if err != nil {
		return nil, err
	}

	return openapi.AuthPasswordSignup200JSONResponse{
		Body:    openapi.AuthSuccess{Id: u.ID.String()},
		Headers: openapi.AuthSuccessResponseHeaders{SetCookie: cookie},
	}, nil
}

func (o *Authentication) AuthOAuthProviderList(ctx context.Context, request openapi.AuthOAuthProviderListRequestObject) (openapi.AuthOAuthProviderListResponseObject, error) {
	list := dt.Map(o.oa.Providers(),
		func(p oauth.Provider) openapi.AuthOAuthProvider {
			return openapi.AuthOAuthProvider{
				Provider: "p.ID()",
				Name:     "p.Name()",
				LogoUrl:  "p.LogoURL()",
				Link:     p.Link(),
			}
		},
	)

	return openapi.AuthOAuthProviderListJSONResponse(list), nil
}

func (o *Authentication) AuthOAuthProviderCallback(ctx context.Context, request openapi.AuthOAuthProviderCallbackRequestObject) (openapi.AuthOAuthProviderCallbackResponseObject, error) {
	provider, err := o.oa.Provider(string(request.OauthProvider))
	if err != nil {
		return nil, errctx.Wrap(err, ctx)
	}

	account, err := provider.Login(ctx, request.Body.State, request.Body.Code)
	if err != nil {
		return nil, errctx.Wrap(err, ctx)
	}

	cookie, err := o.encodeSession(account.ID)
	if err != nil {
		return nil, err
	}

	return openapi.AuthPasswordSignin200JSONResponse{
		Body:    openapi.AuthSuccess{Id: account.ID.String()},
		Headers: openapi.AuthSuccessResponseHeaders{SetCookie: cookie},
	}, nil
}

func (i *Authentication) middleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		r := c.Request()
		ctx := r.Context()

		session, ok := i.decodeSession(r)
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
		return errtag.Wrap(err, errtag.Unauthenticated{})
	}

	// Then look up the account.
	// TODO: Cache this.
	_, err = i.ar.GetByID(ctx, aid)
	if err != nil {
		return errctx.Wrap(err, ctx)
	}

	return nil
}

const secureCookieName = "storyden-session"

type session struct {
	UserID account.AccountID
}

func (i *Authentication) encodeSession(userID account.AccountID) (string, error) {
	encoded, err := i.sc.Encode(secureCookieName, session{userID})
	if err != nil {
		return "", errtag.Wrap(err, errtag.Internal{})
	}

	cookie := &http.Cookie{
		Name:     secureCookieName,
		Value:    encoded,
		Path:     "/",
		Domain:   i.domain,
		Secure:   true,
		HttpOnly: true,
	}

	return cookie.String(), nil
}

func (i *Authentication) decodeSession(r *http.Request) (*session, bool) {
	cookie, err := r.Cookie(secureCookieName)
	if err != nil {
		return nil, false
	}

	u := session{}

	if err = i.sc.Decode(secureCookieName, cookie.Value, &u); err != nil {
		return nil, false
	}

	return &u, true
}
