package bindings

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"

	"github.com/Southclaws/storyden/app/resources/account/authentication"
	auth_service "github.com/Southclaws/storyden/app/services/authentication"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
)

func (o *Authentication) OAuthProviderCallback(ctx context.Context, request openapi.OAuthProviderCallbackRequestObject) (openapi.OAuthProviderCallbackResponseObject, error) {
	service, err := authentication.NewService(request.OauthProvider)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	provider, err := o.authManager.Provider(service)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	op, ok := provider.(auth_service.OAuthProvider)
	if !ok {
		return nil, fault.New("provider is not an OAuth provider", fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	account, err := op.Login(ctx, request.Body.State, request.Body.Code)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	t, err := o.si.Issue(ctx, account.ID)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.OAuthProviderCallback200JSONResponse{
		AuthSuccessOKJSONResponse: openapi.AuthSuccessOKJSONResponse{
			Body: openapi.AuthSuccess{Id: account.ID.String()},
			Headers: openapi.AuthSuccessOKResponseHeaders{
				SetCookie: o.cj.Create(*t).String(),
			},
		},
	}, nil
}
