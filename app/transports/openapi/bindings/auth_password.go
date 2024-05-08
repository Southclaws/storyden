package bindings

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"

	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/app/transports/openapi"
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

func (i *Authentication) AuthPasswordCreate(ctx context.Context, request openapi.AuthPasswordCreateRequestObject) (openapi.AuthPasswordCreateResponseObject, error) {
	id, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	u, err := i.p.Create(ctx, id, request.Body.Password)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.AuthPasswordCreate200JSONResponse{
		AuthSuccessOKJSONResponse: openapi.AuthSuccessOKJSONResponse{
			Body: openapi.AuthSuccessOK{Id: u.ID.String()},
			Headers: openapi.AuthSuccessOKResponseHeaders{
				SetCookie: i.sm.Create(u.ID.String()).String(),
			},
		},
	}, nil
}

func (i *Authentication) AuthPasswordUpdate(ctx context.Context, request openapi.AuthPasswordUpdateRequestObject) (openapi.AuthPasswordUpdateResponseObject, error) {
	id, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	u, err := i.p.Update(ctx, id, request.Body.Old, request.Body.New)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.AuthPasswordUpdate200JSONResponse{
		AuthSuccessOKJSONResponse: openapi.AuthSuccessOKJSONResponse{
			Body: openapi.AuthSuccessOK{Id: u.ID.String()},
			Headers: openapi.AuthSuccessOKResponseHeaders{
				SetCookie: i.sm.Create(u.ID.String()).String(),
			},
		},
	}, nil
}
