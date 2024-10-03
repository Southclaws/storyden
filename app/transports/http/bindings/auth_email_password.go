package bindings

import (
	"context"
	"net/mail"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"

	"github.com/Southclaws/storyden/app/transports/http/openapi"
)

func (i *Authentication) AuthEmailPasswordSignup(ctx context.Context, request openapi.AuthEmailPasswordSignupRequestObject) (openapi.AuthEmailPasswordSignupResponseObject, error) {
	address, err := mail.ParseAddress(request.Body.Email)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	invitedBy, err := deserialiseInvitationID(request.Params.InvitationId)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	handle := opt.NewPtr(request.Body.Handle)

	acc, err := i.epp.Register(ctx, *address, request.Body.Password, handle, invitedBy)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.AuthEmailPasswordSignup200JSONResponse{
		AuthSuccessOKJSONResponse: openapi.AuthSuccessOKJSONResponse{
			Body: openapi.AuthSuccessOK{Id: acc.ID.String()},
			Headers: openapi.AuthSuccessOKResponseHeaders{
				SetCookie: i.cj.Create(acc.ID.String()).String(),
			},
		},
	}, nil
}

func (i *Authentication) AuthEmailPasswordSignin(ctx context.Context, request openapi.AuthEmailPasswordSigninRequestObject) (openapi.AuthEmailPasswordSigninResponseObject, error) {
	u, err := i.epp.Login(ctx, request.Body.Email, request.Body.Password)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.AuthEmailPasswordSignin200JSONResponse{
		AuthSuccessOKJSONResponse: openapi.AuthSuccessOKJSONResponse{
			Body: openapi.AuthSuccess{Id: u.ID.String()},
			Headers: openapi.AuthSuccessOKResponseHeaders{
				SetCookie: i.cj.Create(u.ID.String()).String(),
			},
		},
	}, nil
}
