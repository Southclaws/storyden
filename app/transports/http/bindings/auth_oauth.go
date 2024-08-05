package bindings

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"

	"github.com/Southclaws/storyden/app/transports/http/openapi"
)

func (o *Authentication) OAuthProviderCallback(ctx context.Context, request openapi.OAuthProviderCallbackRequestObject) (openapi.OAuthProviderCallbackResponseObject, error) {
	provider, err := o.am.Provider(string(request.OauthProvider))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	account, err := provider.Login(ctx, request.Body.State, request.Body.Code)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.OAuthProviderCallback200JSONResponse{
		AuthSuccessOKJSONResponse: openapi.AuthSuccessOKJSONResponse{
			Body: openapi.AuthSuccess{Id: account.ID.String()},
			Headers: openapi.AuthSuccessOKResponseHeaders{
				SetCookie: o.cj.Create(account.ID.String()).String(),
			},
		},
	}, nil
}
