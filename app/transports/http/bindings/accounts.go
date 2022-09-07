package bindings

import (
	"context"
	"time"

	"github.com/pkg/errors"

	account_resource "github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/services/account"
	"github.com/Southclaws/storyden/app/transports/http/openapi"
	"github.com/Southclaws/storyden/internal/utils"
)

type Accounts struct {
	as account.Service
}

func NewAccounts(as account.Service) Accounts { return Accounts{as} }

func (i *Accounts) AccountsGet(ctx context.Context, request openapi.AccountsGetRequestObject) (openapi.AccountsGetResponseObject, error) {
	acc, err := i.as.Get(ctx, account_resource.AccountID(request.AccountId.XID()))
	if err != nil {
		return nil, errors.Wrap(err, "failed to get account")
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
