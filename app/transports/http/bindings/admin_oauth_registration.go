package bindings

import (
	"context"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/opt"

	oauthresource "github.com/Southclaws/storyden/app/resources/oauth"
	oauthservice "github.com/Southclaws/storyden/app/services/authentication/oauth"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
)

func (i *Admin) AdminOAuthDynamicRegistrationAccessTokenCreate(ctx context.Context, request openapi.AdminOAuthDynamicRegistrationAccessTokenCreateRequestObject) (openapi.AdminOAuthDynamicRegistrationAccessTokenCreateResponseObject, error) {
	if request.Body == nil {
		return openapi.AdminOAuthDynamicRegistrationAccessTokenCreate400Response{}, nil
	}

	creator, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	issued, err := i.oauth.CreateDynamicRegistrationAccessToken(ctx, oauthservice.DynamicRegistrationAccessTokenCreate{
		Label:            request.Body.Label,
		CreatorAccountID: creator,
		ExpiresAt:        request.Body.ExpiresAt,
		MaxRegistrations: opt.NewPtr(request.Body.MaxRegistrations).Or(1),
	})
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	cache := "no-store"
	return openapi.AdminOAuthDynamicRegistrationAccessTokenCreate201JSONResponse{
		Body: openapi.OAuthDynamicRegistrationAccessTokenIssued{
			Iat:   serialiseOAuthDynamicRegistrationAccessToken(issued.IAT),
			Token: issued.Token,
		},
		Headers: openapi.AdminOAuthDynamicRegistrationAccessTokenCreate201ResponseHeaders{CacheControl: &cache},
	}, nil
}

func (i *Admin) AdminOAuthDynamicRegistrationAccessTokenList(ctx context.Context, request openapi.AdminOAuthDynamicRegistrationAccessTokenListRequestObject) (openapi.AdminOAuthDynamicRegistrationAccessTokenListResponseObject, error) {
	result, err := i.oauth.ListDynamicRegistrationAccessTokens(ctx, deserialisePageParams(request.Params.Page, 50))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.AdminOAuthDynamicRegistrationAccessTokenList200JSONResponse{
		CurrentPage: result.CurrentPage,
		NextPage:    result.NextPage.Ptr(),
		PageSize:    result.Size,
		Results:     result.Results,
		TotalPages:  result.TotalPages,
		Iats:        dt.Map(result.Items, serialiseOAuthDynamicRegistrationAccessToken),
	}, nil
}

func (i *Admin) AdminOAuthDynamicRegistrationAccessTokenGet(ctx context.Context, request openapi.AdminOAuthDynamicRegistrationAccessTokenGetRequestObject) (openapi.AdminOAuthDynamicRegistrationAccessTokenGetResponseObject, error) {
	iat, err := i.oauth.GetDynamicRegistrationAccessToken(ctx, oauthresource.DynamicRegistrationAccessTokenID(deserialiseID(request.OauthDcrIatId)))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.AdminOAuthDynamicRegistrationAccessTokenGet200JSONResponse(serialiseOAuthDynamicRegistrationAccessToken(iat)), nil
}

func (i *Admin) AdminOAuthDynamicRegistrationAccessTokenRevoke(ctx context.Context, request openapi.AdminOAuthDynamicRegistrationAccessTokenRevokeRequestObject) (openapi.AdminOAuthDynamicRegistrationAccessTokenRevokeResponseObject, error) {
	if err := i.oauth.RevokeDynamicRegistrationAccessToken(ctx, oauthresource.DynamicRegistrationAccessTokenID(deserialiseID(request.OauthDcrIatId))); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.AdminOAuthDynamicRegistrationAccessTokenRevoke204Response{}, nil
}

func serialiseOAuthDynamicRegistrationAccessToken(in *oauthresource.DynamicRegistrationAccessToken) openapi.OAuthDynamicRegistrationAccessToken {
	return openapi.OAuthDynamicRegistrationAccessToken{
		Id:                openapi.Identifier(in.ID.XID().String()),
		CreatedAt:         openapi.CreatedAt(in.CreatedAt),
		Label:             in.Label,
		CreatorAccountId:  openapi.Identifier(in.CreatorAccountID.String()),
		ExpiresAt:         in.ExpiresAt,
		MaxRegistrations:  in.MaxRegistrations,
		RegistrationCount: in.RegistrationCount,
		RevokedAt:         in.RevokedAt.Ptr(),
	}
}
