package bindings

import (
	"context"

	"github.com/Southclaws/fault/errctx"
	"github.com/pkg/errors"

	account_resource "github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/services/account"
	"github.com/Southclaws/storyden/app/services/authentication"
	"github.com/Southclaws/storyden/app/services/avatar"
	"github.com/Southclaws/storyden/app/transports/openapi/openapi"
	"github.com/Southclaws/storyden/internal/utils"
)

type Accounts struct {
	as account.Service
	av avatar.Service
}

func NewAccounts(as account.Service, av avatar.Service) Accounts { return Accounts{as, av} }

func (i *Accounts) AccountsGet(ctx context.Context, request openapi.AccountsGetRequestObject) (openapi.AccountsGetResponseObject, error) {
	accountID, err := authentication.GetAccountID(ctx)
	if err != nil {
		return nil, errctx.Wrap(err, ctx)
	}

	acc, err := i.as.Get(ctx, accountID)
	if err != nil {
		return nil, errctx.Wrap(errors.Wrap(err, "failed to get account"), ctx)
	}

	return openapi.AccountsGet200JSONResponse{
		Id:        openapi.Identifier(acc.ID.String()),
		Handle:    &acc.Handle,
		Name:      utils.Ref(acc.Name),
		Bio:       utils.Ref(acc.Bio.ElseZero()),
		CreatedAt: acc.CreatedAt,
		UpdatedAt: acc.UpdatedAt,
		DeletedAt: utils.OptionalToPointer(acc.DeletedAt),
	}, nil
}

func (i *Accounts) AccountsGetAvatar(ctx context.Context, request openapi.AccountsGetAvatarRequestObject) (openapi.AccountsGetAvatarResponseObject, error) {
	r, err := i.av.Get(ctx, account_resource.AccountID(request.AccountId.XID()))
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
