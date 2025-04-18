package bindings

import (
	"context"
	"errors"
	"net/mail"
	"strings"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"

	"github.com/Southclaws/storyden/app/services/authentication/provider/email_only"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
)

func (i *Authentication) AuthEmailSignup(ctx context.Context, request openapi.AuthEmailSignupRequestObject) (openapi.AuthEmailSignupResponseObject, error) {
	address, err := mail.ParseAddress(strings.ToLower(request.Body.Email))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	invitedBy, err := deserialiseInvitationID(request.Params.InvitationId)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	handle := opt.NewPtr(request.Body.Handle)

	acc, err := i.emailVerificationAuthProvider.Register(ctx, *address, handle, invitedBy)
	if err != nil {
		// SPEC: If the email exists, return a 422 response with no session.
		if errors.Is(err, email_only.ErrAccountAlreadyExists) {
			return openapi.AuthEmailSignup422Response{}, nil
		}

		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	t, err := i.si.Issue(ctx, acc.ID)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.AuthEmailSignup200JSONResponse{
		AuthSuccessOKJSONResponse: openapi.AuthSuccessOKJSONResponse{
			Body: openapi.AuthSuccessOK{Id: acc.ID.String()},
			Headers: openapi.AuthSuccessOKResponseHeaders{
				SetCookie: i.cj.Create(*t).String(),
			},
		},
	}, nil
}

func (i *Authentication) AuthEmailSignin(ctx context.Context, request openapi.AuthEmailSigninRequestObject) (openapi.AuthEmailSigninResponseObject, error) {
	address, err := mail.ParseAddress(strings.ToLower(request.Body.Email))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	err = i.emailVerificationAuthProvider.Login(ctx, *address)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.AuthEmailSignin200JSONResponse{
		AuthSuccessOKJSONResponse: openapi.AuthSuccessOKJSONResponse{},
	}, nil
}
