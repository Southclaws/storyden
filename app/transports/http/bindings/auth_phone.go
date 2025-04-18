package bindings

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"

	"github.com/Southclaws/storyden/app/services/authentication/provider/phone"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/app/transports/http/middleware/session_cookie"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
)

type PhoneAuth struct {
	pp *phone.Provider
	cj *session_cookie.Jar
	si *session.Issuer
}

func NewPhoneAuth(pp *phone.Provider, cj *session_cookie.Jar, si *session.Issuer) PhoneAuth {
	return PhoneAuth{pp, cj, si}
}

func (i *PhoneAuth) PhoneRequestCode(ctx context.Context, request openapi.PhoneRequestCodeRequestObject) (openapi.PhoneRequestCodeResponseObject, error) {
	invitedBy, err := deserialiseInvitationID(request.Params.InvitationId)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	acc, err := i.pp.Register(ctx, request.Body.Identifier, request.Body.PhoneNumber, invitedBy)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.PhoneRequestCode200JSONResponse{
		AuthSuccessOKJSONResponse: openapi.AuthSuccessOKJSONResponse{
			Body: openapi.AuthSuccess{Id: acc.ID.String()},
		},
	}, nil
}

func (i *PhoneAuth) PhoneSubmitCode(ctx context.Context, request openapi.PhoneSubmitCodeRequestObject) (openapi.PhoneSubmitCodeResponseObject, error) {
	acc, err := i.pp.Login(ctx, request.AccountHandle, request.Body.Code)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	t, err := i.si.Issue(ctx, acc.ID)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.PhoneSubmitCode200JSONResponse{
		AuthSuccessOKJSONResponse: openapi.AuthSuccessOKJSONResponse{
			Body: openapi.AuthSuccess{Id: acc.ID.String()},
			Headers: openapi.AuthSuccessOKResponseHeaders{
				SetCookie: i.cj.Create(*t).String(),
			},
		},
	}, nil
}
