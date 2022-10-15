package bindings

import (
	"context"

	"github.com/Southclaws/fault/errctx"
	"github.com/Southclaws/opt"
	"github.com/pkg/errors"

	account_repo "github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/services/account"
	"github.com/Southclaws/storyden/app/services/authentication"
	"github.com/Southclaws/storyden/app/services/avatar"
	"github.com/Southclaws/storyden/app/transports/openapi/openapi"
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
		return nil, errctx.Wrap(err, ctx)
	}

	acc, err := i.as.Get(ctx, accountID)
	if err != nil {
		return nil, errctx.Wrap(errors.Wrap(err, "failed to get account"), ctx)
	}

	return openapi.AccountsGet200JSONResponse(serialiseAccount(acc)), nil
}

func (i *Accounts) AccountsUpdate(ctx context.Context, request openapi.AccountsUpdateRequestObject) (openapi.AccountsUpdateResponseObject, error) {
	accountID, err := authentication.GetAccountID(ctx)
	if err != nil {
		return nil, errctx.Wrap(err, ctx)
	}

	acc, err := i.as.Update(ctx, accountID, account.Partial{
		Handle:    opt.NewPtrMap(request.Body.Handle, func(i openapi.AccountHandle) string { return string(i) }),
		Name:      opt.NewPtr(request.Body.Name),
		Bio:       opt.NewPtr(request.Body.Bio),
		Interests: opt.NewPtrMap(request.Body.Interests, tagsIDs),
	})
	if err != nil {
		return nil, errctx.Wrap(err, ctx)
	}

	return openapi.AccountsUpdateSuccessJSONResponse(serialiseAccount(acc)), nil
}

func (i *Accounts) AccountsGetAvatar(ctx context.Context, request openapi.AccountsGetAvatarRequestObject) (openapi.AccountsGetAvatarResponseObject, error) {
	id, err := request.AccountHandle.ID(ctx, i.ar)
	if err != nil {
		return nil, errctx.Wrap(err, ctx)
	}

	r, err := i.av.Get(ctx, id)
	if err != nil {
		return nil, errctx.Wrap(err, ctx)
	}

	return openapi.AccountsGetAvatarImagepngResponse{
		Body: r,
	}, nil
}

func (i *Accounts) AccountsSetAvatar(ctx context.Context, request openapi.AccountsSetAvatarRequestObject) (openapi.AccountsSetAvatarResponseObject, error) {
	accountID, err := authentication.GetAccountID(ctx)
	if err != nil {
		return nil, errctx.Wrap(err, ctx)
	}

	if err := i.av.Set(ctx, accountID, request.Body); err != nil {
		return nil, errctx.Wrap(err, ctx)
	}

	return openapi.AccountsSetAvatar200Response{}, nil
}
