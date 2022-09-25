package bindings

import (
	"context"

	"github.com/pkg/errors"

	"github.com/Southclaws/storyden/app/services/account"
	"github.com/Southclaws/storyden/app/services/authentication"
	"github.com/Southclaws/storyden/app/transports/openapi/openapi"
	"github.com/Southclaws/storyden/internal/errctx"
	"github.com/Southclaws/storyden/internal/utils"
)

type Accounts struct {
	as account.Service
}

func NewAccounts(as account.Service) Accounts { return Accounts{as} }

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
