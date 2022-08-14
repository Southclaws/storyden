package bindings

import (
	"context"
	"errors"
	"net/http"

	"github.com/deepmap/oapi-codegen/pkg/middleware"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/gorilla/securecookie"
	"github.com/labstack/echo/v4"

	"github.com/Southclaws/storyden/backend/internal/config"
	"github.com/Southclaws/storyden/backend/pkg/resources/user"
	"github.com/Southclaws/storyden/backend/pkg/services/authentication"
	"github.com/Southclaws/storyden/backend/pkg/services/authentication/provider/password"
	"github.com/Southclaws/storyden/backend/pkg/transports/http/openapi"
)

type Authentication struct {
	// s      authentication.Service
	p      *password.Password
	sc     *securecookie.SecureCookie
	domain string
}

func NewAuthentication(
	cfg config.Config,
	// s authentication.Service,
	p *password.Password,
	sc *securecookie.SecureCookie,
) Authentication {
	return Authentication{p, sc, cfg.CookieDomain}
}

func (i *Authentication) Signin(ctx context.Context, request openapi.SigninRequestObject) any {
	u, err := func() (*user.User, error) {
		if request.JSONBody != nil {
			return i.p.Login(ctx, request.JSONBody.Identifier, request.JSONBody.Token)
		} else if request.FormdataBody != nil {
			return i.p.Login(ctx, request.FormdataBody.Identifier, request.FormdataBody.Token)
		}
		return nil, errors.New("missing body")
	}()
	if err != nil {
		if errors.Is(err, password.ErrPasswordMismatch) {
			return openapi.Signin401Response{}
		}
		return openapi.Signin500JSONResponse{Error: err.Error()}
	}

	cookie, err := i.encodeSession(u.ID)
	if err != nil {
		return err
	}

	return openapi.Signin200JSONResponse{
		Body:    openapi.AuthenticationResponse{Id: u.ID.String()},
		Headers: openapi.Signin200ResponseHeaders{SetCookie: cookie},
	}
}

func (i *Authentication) Signup(ctx context.Context, request openapi.SignupRequestObject) any {
	u, err := func() (*user.User, error) {
		if request.JSONBody != nil {
			return i.p.Register(ctx, request.JSONBody.Identifier, request.JSONBody.Token)
		} else if request.FormdataBody != nil {
			return i.p.Register(ctx, request.FormdataBody.Identifier, request.FormdataBody.Token)
		}
		return nil, errors.New("missing body")
	}()
	if err != nil {
		return openapi.Signup500JSONResponse{Error: err.Error()}
	}

	cookie, err := i.encodeSession(u.ID)
	if err != nil {
		return err
	}

	return openapi.Signup200JSONResponse{
		Body:    openapi.AuthenticationResponse{Id: u.ID.String()},
		Headers: openapi.Signup200ResponseHeaders{SetCookie: cookie},
	}
}

func (i *Authentication) validator(ctx context.Context, ai *openapi3filter.AuthenticationInput) error {
	echo := ctx.Value(middleware.EchoContextKey).(echo.Context)
	req := echo.Request()

	session, ok := i.decodeSession(req)
	if !ok {
		return errors.New("no cookie found in request")
	}

	req.WithContext(authentication.AddUserToContext(req.Context(), session.UserID))

	return nil
}

const secureCookieName = "storyden-session"

type session struct {
	UserID user.UserID
}

func (i *Authentication) encodeSession(userID user.UserID) (string, error) {
	encoded, err := i.sc.Encode(secureCookieName, userID)
	if err != nil {
		return "", err
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
