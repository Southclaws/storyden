package bindings

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"

	"github.com/Southclaws/storyden/internal/openapi"
)

func (i *Authentication) AuthPasswordSignin(ctx context.Context, request openapi.AuthPasswordSigninRequestObject) (openapi.AuthPasswordSigninResponseObject, error) {
	u, err := i.p.Login(ctx, request.Body.Identifier, request.Body.Token)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.AuthPasswordSignin200JSONResponse{
		AuthSuccessOKJSONResponse: openapi.AuthSuccessOKJSONResponse{
			Body: openapi.AuthSuccess{Id: u.ID.String()},
			Headers: openapi.AuthSuccessOKResponseHeaders{
				SetCookie: i.sm.Create(u.ID.String()).String(),
			},
		},
	}, nil
}

func (i *Authentication) AuthPasswordSignup(ctx context.Context, request openapi.AuthPasswordSignupRequestObject) (openapi.AuthPasswordSignupResponseObject, error) {
	u, err := i.p.Register(ctx, request.Body.Identifier, request.Body.Token)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.AuthPasswordSignup200JSONResponse{
		AuthSuccessOKJSONResponse: openapi.AuthSuccessOKJSONResponse{
			Body: openapi.AuthSuccessOK{Id: u.ID.String()},
			Headers: openapi.AuthSuccessOKResponseHeaders{
				SetCookie: i.sm.Create(u.ID.String()).String(),
			},
		},
	}, nil
}
