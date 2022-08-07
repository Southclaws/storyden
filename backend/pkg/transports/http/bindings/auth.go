package bindings

import (
	"context"
	"errors"
	"net/http"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/kr/pretty"

	"github.com/Southclaws/storyden/backend/pkg/resources/user"
	"github.com/Southclaws/storyden/backend/pkg/services/authentication"
	"github.com/Southclaws/storyden/backend/pkg/services/authentication/provider/password"
	"github.com/Southclaws/storyden/backend/pkg/transports/http/openapi"
)

type Authentication struct {
	s authentication.Service
	p *password.Password
}

func NewAuthentication(s authentication.Service, p *password.Password) Authentication {
	return Authentication{s, p}
}

func (i *Authentication) Signin(ctx context.Context, request openapi.SigninRequestObject) any {
	u, err := func() (*user.User, error) {
		if request.JSONBody != nil {
			return i.p.Register(ctx, request.JSONBody.Identifier, request.JSONBody.Token)
		} else if request.FormdataBody != nil {
			return i.p.Register(ctx, request.FormdataBody.Identifier, request.FormdataBody.Token)
		}
		return nil, errors.New("missing body")
	}()
	if err != nil {
		return openapi.Signin500JSONResponse{Error: err.Error()}
	}

	return openapi.Signin200JSONResponse{Id: u.ID.String()}
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

	return openapi.Signup200JSONResponse{Id: u.ID.String()}
}

func (i *Authentication) middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if session, ok := i.s.DecodeSession(r); ok {
			ctx := authentication.AddUserToContext(r.Context(), session)
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

func (i *Authentication) validator(ctx context.Context, ai *openapi3filter.AuthenticationInput) error {
	pretty.Println(ai.SecurityScheme)
	return errors.New("not allowed")
}
