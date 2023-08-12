package bindings

import (
	"context"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/opt"

	account_repo "github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/services/account"
	"github.com/Southclaws/storyden/app/services/authentication"
	"github.com/Southclaws/storyden/app/services/avatar"
	"github.com/Southclaws/storyden/internal/openapi"
)

type Accounts struct {
	as account.Service
	av avatar.Service
	ar account_repo.Repository
}

func NewAccounts(as account.Service, av avatar.Service, ar account_repo.Repository) Accounts {
	return Accounts{as, av, ar}
}

func (i *Accounts) AccountGet(ctx context.Context, request openapi.AccountGetRequestObject) (openapi.AccountGetResponseObject, error) {
	accountID, err := authentication.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	acc, err := i.as.Get(ctx, accountID)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.AccountGet200JSONResponse{
		AccountGetOKJSONResponse: openapi.AccountGetOKJSONResponse(serialiseAccount(acc)),
	}, nil
}

func (i *Accounts) AccountUpdate(ctx context.Context, request openapi.AccountUpdateRequestObject) (openapi.AccountUpdateResponseObject, error) {
	accountID, err := authentication.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	acc, err := i.as.Update(ctx, accountID, account.Partial{
		Handle:    opt.NewPtrMap(request.Body.Handle, func(i openapi.AccountHandle) string { return string(i) }),
		Name:      opt.NewPtr(request.Body.Name),
		Bio:       opt.NewPtr(request.Body.Bio),
		Interests: opt.NewPtrMap(request.Body.Interests, tagsIDs),
	})
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.AccountUpdate200JSONResponse{
		AccountUpdateOKJSONResponse: openapi.AccountUpdateOKJSONResponse(serialiseAccount(acc)),
	}, nil
}

func (i *Accounts) AccountAuthProviderList(ctx context.Context, request openapi.AccountAuthProviderListRequestObject) (openapi.AccountAuthProviderListResponseObject, error) {
	accountID, err := authentication.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	authmethods, err := i.as.GetAuthMethods(ctx, accountID)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.AccountAuthProviderList200JSONResponse{
		AccountAuthProviderListOKJSONResponse: openapi.AccountAuthProviderListOKJSONResponse{
			AuthMethods: dt.Map(authmethods, serialiseAuthMethod),
		},
	}, nil
}

func (i *Accounts) AccountGetAvatar(ctx context.Context, request openapi.AccountGetAvatarRequestObject) (openapi.AccountGetAvatarResponseObject, error) {
	id, err := openapi.ResolveHandle(ctx, i.ar, request.AccountHandle)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	r, err := i.av.Get(ctx, id)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.AccountGetAvatar200ImagepngResponse{
		AccountGetAvatarImagepngResponse: openapi.AccountGetAvatarImagepngResponse{
			Body: r,
		},
	}, nil
}

func (i *Accounts) AccountSetAvatar(ctx context.Context, request openapi.AccountSetAvatarRequestObject) (openapi.AccountSetAvatarResponseObject, error) {
	accountID, err := authentication.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if err := i.av.Set(ctx, accountID, request.Body); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.AccountSetAvatar200Response{}, nil
}

func serialiseAuthMethod(in *account.AuthMethod) openapi.AccountAuthMethod {
	return openapi.AccountAuthMethod{
		Provider: in.ID(),
		Name:     in.Name(),
		LogoUrl:  in.LogoURL(),
		Link:     in.Link(),
		Active:   in.Active,
	}
}
