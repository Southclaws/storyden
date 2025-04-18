package bindings

import (
	"context"
	"net/mail"
	"strings"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"

	"github.com/Southclaws/storyden/app/services/authentication/provider/password/password_reset"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
)

func (i *Authentication) AuthEmailPasswordSignup(ctx context.Context, request openapi.AuthEmailPasswordSignupRequestObject) (openapi.AuthEmailPasswordSignupResponseObject, error) {
	address, err := mail.ParseAddress(strings.ToLower(request.Body.Email))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	invitedBy, err := deserialiseInvitationID(request.Params.InvitationId)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	handle := opt.NewPtr(request.Body.Handle)

	acc, err := i.passwordAuthProvider.RegisterWithEmail(ctx, *address, request.Body.Password, handle, invitedBy)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	t, err := i.si.Issue(ctx, acc.ID)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.AuthEmailPasswordSignup200JSONResponse{
		AuthSuccessOKJSONResponse: openapi.AuthSuccessOKJSONResponse{
			Body: openapi.AuthSuccessOK{Id: acc.ID.String()},
			Headers: openapi.AuthSuccessOKResponseHeaders{
				SetCookie: i.cj.Create(*t).String(),
			},
		},
	}, nil
}

func (i *Authentication) AuthEmailPasswordSignin(ctx context.Context, request openapi.AuthEmailPasswordSigninRequestObject) (openapi.AuthEmailPasswordSigninResponseObject, error) {
	address, err := mail.ParseAddress(strings.ToLower(request.Body.Email))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	acc, err := i.passwordAuthProvider.LoginWithEmail(ctx, *address, request.Body.Password)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	t, err := i.si.Issue(ctx, acc.ID)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.AuthEmailPasswordSignin200JSONResponse{
		AuthSuccessOKJSONResponse: openapi.AuthSuccessOKJSONResponse{
			Body: openapi.AuthSuccess{Id: acc.ID.String()},
			Headers: openapi.AuthSuccessOKResponseHeaders{
				SetCookie: i.cj.Create(*t).String(),
			},
		},
	}, nil
}

func (i *Authentication) AuthPasswordResetRequestEmail(ctx context.Context, request openapi.AuthPasswordResetRequestEmailRequestObject) (openapi.AuthPasswordResetRequestEmailResponseObject, error) {
	address, err := mail.ParseAddress(strings.ToLower(request.Body.Email))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	lt, err := password_reset.NewLinkTemplate(request.Body.TokenUrl.Url, request.Body.TokenUrl.Query)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	err = i.passwordAuthProvider.RequestReset(ctx, *address, *lt)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.AuthPasswordResetRequestEmail200Response{}, nil
}
