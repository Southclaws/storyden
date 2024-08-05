package bindings

import (
	"context"
	"net/url"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/account_querier"
	"github.com/Southclaws/storyden/app/services/account/account_auth"
	"github.com/Southclaws/storyden/app/services/account/account_update"
	"github.com/Southclaws/storyden/app/services/authentication"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/app/services/avatar"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
)

type Accounts struct {
	avatarService avatar.Service
	authManager   *authentication.Manager
	accountQuery  account_querier.Querier
	accountUpdate account_update.Updater
	accountAuth   account_auth.Manager
}

func NewAccounts(
	avatarService avatar.Service,
	authManager *authentication.Manager,
	accountQuery account_querier.Querier,
	accountUpdate account_update.Updater,
	accountAuth account_auth.Manager,
) Accounts {
	return Accounts{
		avatarService,
		authManager,
		accountQuery,
		accountUpdate,
		accountAuth,
	}
}

func (i *Accounts) AccountGet(ctx context.Context, request openapi.AccountGetRequestObject) (openapi.AccountGetResponseObject, error) {
	accountID, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	acc, err := i.accountQuery.GetByID(ctx, accountID)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.AccountGet200JSONResponse{
		AccountGetOKJSONResponse: openapi.AccountGetOKJSONResponse(serialiseAccount(acc)),
	}, nil
}

func (i *Accounts) AccountUpdate(ctx context.Context, request openapi.AccountUpdateRequestObject) (openapi.AccountUpdateResponseObject, error) {
	accountID, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	links, err := opt.MapErr(opt.NewPtr(request.Body.Links), deserialiseExternalLinkList)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	acc, err := i.accountUpdate.Update(ctx, accountID, account_update.Partial{
		Handle:    opt.NewPtrMap(request.Body.Handle, func(i openapi.AccountHandle) string { return string(i) }),
		Name:      opt.NewPtr(request.Body.Name),
		Bio:       opt.NewPtr(request.Body.Bio),
		Links:     links,
		Meta:      opt.NewPtr((*map[string]any)(request.Body.Meta)),
		Interests: opt.NewPtrMap(request.Body.Interests, tagsIDs),
	})
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.AccountUpdate200JSONResponse{
		AccountUpdateOKJSONResponse: openapi.AccountUpdateOKJSONResponse(serialiseAccount(acc)),
	}, nil
}

func deserialiseExternalLinkList(i openapi.ProfileExternalLinkList) ([]account.ExternalLink, error) {
	return dt.MapErr(i, deserialiseExternalLink)
}

func deserialiseExternalLink(l openapi.ProfileExternalLink) (account.ExternalLink, error) {
	u, err := url.Parse(string(l.Url))
	if err != nil {
		return account.ExternalLink{}, err
	}

	return account.ExternalLink{
		Text: string(l.Text),
		URL:  *u,
	}, nil
}

func (i *Accounts) AccountAuthProviderList(ctx context.Context, request openapi.AccountAuthProviderListRequestObject) (openapi.AccountAuthProviderListResponseObject, error) {
	accountID, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	authmethods, err := i.accountAuth.GetAuthMethods(ctx, accountID)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	available, err := dt.MapErr(i.authManager.Providers(), serialiseAuthProvider)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	active, err := dt.MapErr(authmethods, serialiseAuthMethod)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.AccountAuthProviderList200JSONResponse{
		AccountAuthProviderListOKJSONResponse: openapi.AccountAuthProviderListOKJSONResponse{
			Available: available,
			Active:    active,
		},
	}, nil
}

func (i *Accounts) AccountAuthMethodDelete(ctx context.Context, request openapi.AccountAuthMethodDeleteRequestObject) (openapi.AccountAuthMethodDeleteResponseObject, error) {
	accountID, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	id, err := xid.FromString(request.AuthMethodId)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), ftag.With(ftag.InvalidArgument))
	}

	err = i.accountAuth.DeleteAuthMethod(ctx, accountID, id)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	authmethods, err := i.accountAuth.GetAuthMethods(ctx, accountID)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	available, err := dt.MapErr(i.authManager.Providers(), serialiseAuthProvider)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	active, err := dt.MapErr(authmethods, serialiseAuthMethod)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.AccountAuthMethodDelete200JSONResponse{
		AccountAuthProviderListOKJSONResponse: openapi.AccountAuthProviderListOKJSONResponse{
			Available: available,
			Active:    active,
		},
	}, nil
}

func (i *Accounts) AccountGetAvatar(ctx context.Context, request openapi.AccountGetAvatarRequestObject) (openapi.AccountGetAvatarResponseObject, error) {
	id, err := openapi.ResolveHandle(ctx, i.accountQuery, request.AccountHandle)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	r, size, err := i.avatarService.Get(ctx, id)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.AccountGetAvatar200ImagepngResponse{
		AccountGetAvatarImagepngResponse: openapi.AccountGetAvatarImagepngResponse{
			Body:          r,
			ContentLength: size,
		},
	}, nil
}

func (i *Accounts) AccountSetAvatar(ctx context.Context, request openapi.AccountSetAvatarRequestObject) (openapi.AccountSetAvatarResponseObject, error) {
	accountID, err := session.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if err := i.avatarService.Set(ctx, accountID, request.Body, int64(request.Params.ContentLength)); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.AccountSetAvatar200Response{}, nil
}

func serialiseAuthMethod(in *account_auth.AuthMethod) (openapi.AccountAuthMethod, error) {
	p, err := serialiseAuthProvider(in.Provider)
	if err != nil {
		return openapi.AccountAuthMethod{}, fault.Wrap(err)
	}

	return openapi.AccountAuthMethod{
		Id:         in.Instance.ID.String(),
		CreatedAt:  in.Instance.Created,
		Name:       in.Instance.Name.Or("Unknown"),
		Identifier: in.Instance.Identifier,
		Provider:   p,
	}, nil
}
