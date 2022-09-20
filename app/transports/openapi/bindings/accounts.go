package bindings

import (
	"context"
	"time"

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
		Bio:       utils.Ref(acc.Bio.ElseZero()),
		Email:     utils.Ref(acc.Email),
		Name:      utils.Ref(acc.Name),
		CreatedAt: utils.Ref(acc.CreatedAt.Format(time.RFC3339)),
		UpdatedAt: utils.Ref(acc.UpdatedAt.Format(time.RFC3339)),
		DeletedAt: utils.OptionalElsePtr(acc.DeletedAt, utils.FormatISO),
	}, nil
}
