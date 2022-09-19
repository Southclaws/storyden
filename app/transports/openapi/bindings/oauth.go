package bindings

import (
	"context"

	"github.com/Southclaws/storyden/app/transports/openapi/openapi"
)

type OAuth struct{}

func NewOAuth() OAuth {
	return OAuth{}
}

func (o *OAuth) AuthOAuthProviderLink(ctx context.Context, request openapi.AuthOAuthProviderLinkRequestObject) (openapi.AuthOAuthProviderLinkResponseObject, error) {
	return nil, nil
}

func (o *OAuth) AuthOAuthProviderCallback(ctx context.Context, request openapi.AuthOAuthProviderCallbackRequestObject) (openapi.AuthOAuthProviderCallbackResponseObject, error) {
	// provider, exists, err := o.oauth.LookupProvider(ctx, string(request.OauthProvider))
	// provider.Callback(ctx, ...)
	return nil, nil
}
