package bindings

import (
	"context"
	"errors"
	"net/mail"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"

	"github.com/Southclaws/storyden/app/services/authentication/provider/email/email_only"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
)

func (i *Authentication) AuthEmailSignup(ctx context.Context, request openapi.AuthEmailSignupRequestObject) (openapi.AuthEmailSignupResponseObject, error) {
	address, err := mail.ParseAddress(request.Body.Email)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	invitedBy, err := deserialiseInvitationID(request.Params.InvitationId)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	handle := opt.NewPtr(request.Body.Handle)

	acc, err := i.ep.Register(ctx, *address, handle, invitedBy)
	if err != nil {
		// SPEC: If the email exists, return a 422 response with no session.
		if errors.Is(err, email_only.ErrAccountAlreadyExists) {
			return openapi.AuthEmailSignup422Response{}, nil
		}

		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.AuthEmailSignup200JSONResponse{
		AuthSuccessOKJSONResponse: openapi.AuthSuccessOKJSONResponse{
			Body: openapi.AuthSuccessOK{Id: acc.ID.String()},
			Headers: openapi.AuthSuccessOKResponseHeaders{
				SetCookie: i.cj.Create(acc.ID.String()).String(),
			},
		},
	}, nil
}

func (i *Authentication) AuthEmailSignin(ctx context.Context, request openapi.AuthEmailSigninRequestObject) (openapi.AuthEmailSigninResponseObject, error) {
	// i.ep.Login(ctx, request.Body.Email)
	return nil, nil
}
