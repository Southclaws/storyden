package bindings

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"

	"github.com/Southclaws/storyden/app/services/authentication/provider/phone"
	"github.com/Southclaws/storyden/app/transports/http/cookie"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
)

type PhoneAuth struct {
	pp *phone.Provider
	cj *cookie.Jar
}

func NewPhoneAuth(pp *phone.Provider, cj *cookie.Jar) PhoneAuth {
	return PhoneAuth{pp, cj}
}

func (i *PhoneAuth) PhoneRequestCode(ctx context.Context, request openapi.PhoneRequestCodeRequestObject) (openapi.PhoneRequestCodeResponseObject, error) {
	acc, err := i.pp.Register(ctx, request.Body.Identifier, request.Body.PhoneNumber)
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

	return openapi.PhoneSubmitCode200JSONResponse{
		AuthSuccessOKJSONResponse: openapi.AuthSuccessOKJSONResponse{
			Body: openapi.AuthSuccess{Id: acc.ID.String()},
			Headers: openapi.AuthSuccessOKResponseHeaders{
				SetCookie: i.cj.Create(acc.ID.String()).String(),
			},
		},
	}, nil
}
