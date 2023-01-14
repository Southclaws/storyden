package bindings

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"

	"github.com/Southclaws/storyden/internal/openapi"
)

func (o *Authentication) AuthOAuthProviderCallback(ctx context.Context, request openapi.AuthOAuthProviderCallbackRequestObject) (openapi.AuthOAuthProviderCallbackResponseObject, error) {
	provider, err := o.am.Provider(string(request.OauthProvider))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	account, err := provider.Login(ctx, request.Body.State, request.Body.Code)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	cookie, err := o.sm.encodeSession(account.ID)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.AuthOAuthProviderCallback200JSONResponse{
		AuthSuccessJSONResponse: openapi.AuthSuccessJSONResponse{
			Body:    openapi.AuthSuccess{Id: account.ID.String()},
			Headers: openapi.AuthSuccessResponseHeaders{SetCookie: cookie},
		},
	}, nil
}
