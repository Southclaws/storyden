package bindings

import (
	"context"
	"net/mail"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/app/transports/openapi"
)

func (i *Authentication) AuthEmailSignup(ctx context.Context, request openapi.AuthEmailSignupRequestObject) (openapi.AuthEmailSignupResponseObject, error) {
	address, err := mail.ParseAddress(request.Body.Email)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	acc, err := i.ep.Register(ctx, *address)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.AuthEmailSignup200JSONResponse{
		AuthSuccessOKJSONResponse: openapi.AuthSuccessOKJSONResponse{
			Body: openapi.AuthSuccessOK{Id: acc.ID.String()},
			Headers: openapi.AuthSuccessOKResponseHeaders{
				SetCookie: i.sm.Create(acc.ID.String()).String(),
			},
		},
	}, nil
}

func (i *Authentication) AuthEmailSignin(ctx context.Context, request openapi.AuthEmailSigninRequestObject) (openapi.AuthEmailSigninResponseObject, error) {
	// i.ep.Login(ctx, request.Body.Email)
	return nil, nil
}

func (i *Authentication) AuthEmailVerify(ctx context.Context, request openapi.AuthEmailVerifyRequestObject) (openapi.AuthEmailVerifyResponseObject, error) {
	id, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	acc, err := i.ep.Login(ctx, id.String(), request.Body.Code)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.AuthEmailVerify200JSONResponse{
		AuthSuccessOKJSONResponse: openapi.AuthSuccessOKJSONResponse{
			Body: openapi.AuthSuccessOK{Id: acc.ID.String()},
			Headers: openapi.AuthSuccessOKResponseHeaders{
				SetCookie: i.sm.Create(acc.ID.String()).String(),
			},
		},
	}, nil
}
