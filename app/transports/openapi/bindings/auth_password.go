package bindings

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"

	"github.com/Southclaws/storyden/internal/openapi"
)

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
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	cookie, err := i.sm.encodeSession(u.ID)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.AuthPasswordSignin200JSONResponse{
		AuthSuccessJSONResponse: openapi.AuthSuccessJSONResponse{
			Body:    openapi.AuthSuccess{Id: u.ID.String()},
			Headers: openapi.AuthSuccessResponseHeaders{SetCookie: cookie},
		},
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
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	cookie, err := i.sm.encodeSession(u.ID)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.AuthPasswordSignup200JSONResponse{
		AuthSuccessJSONResponse: openapi.AuthSuccessJSONResponse{
			Body:    openapi.AuthSuccess{Id: u.ID.String()},
			Headers: openapi.AuthSuccessResponseHeaders{SetCookie: cookie},
		},
	}, nil
}
