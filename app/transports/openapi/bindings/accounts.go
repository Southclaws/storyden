package bindings

import (
	"context"

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

func (i *Accounts) AccountsGet(ctx context.Context, request openapi.AccountsGetRequestObject) (openapi.AccountsGetResponseObject, error) {
	accountID, err := authentication.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	acc, err := i.as.Get(ctx, accountID)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.AccountsGet200JSONResponse{
		AccountsGetSuccessJSONResponse: openapi.AccountsGetSuccessJSONResponse(serialiseAccount(acc)),
	}, nil
}

func (i *Accounts) AccountsUpdate(ctx context.Context, request openapi.AccountsUpdateRequestObject) (openapi.AccountsUpdateResponseObject, error) {
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

	return openapi.AccountsUpdate200JSONResponse{
		AccountsUpdateSuccessJSONResponse: openapi.AccountsUpdateSuccessJSONResponse(serialiseAccount(acc)),
	}, nil
}

func (i *Accounts) AccountsGetAvatar(ctx context.Context, request openapi.AccountsGetAvatarRequestObject) (openapi.AccountsGetAvatarResponseObject, error) {
	id, err := openapi.ResolveHandle(ctx, i.ar, request.AccountHandle)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	r, err := i.av.Get(ctx, id)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.AccountsGetAvatar200ImagepngResponse{
		AccountsGetAvatarImagepngResponse: openapi.AccountsGetAvatarImagepngResponse{
			Body: r,
		},
	}, nil
}

func (i *Accounts) AccountsSetAvatar(ctx context.Context, request openapi.AccountsSetAvatarRequestObject) (openapi.AccountsSetAvatarResponseObject, error) {
	accountID, err := authentication.GetAccountID(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if err := i.av.Set(ctx, accountID, request.Body); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return openapi.AccountsSetAvatar200Response{}, nil
}
