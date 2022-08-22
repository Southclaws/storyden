package bindings

import (
	"context"
	"time"

	"github.com/Southclaws/storyden/backend/internal/utils"
	account_resource "github.com/Southclaws/storyden/backend/pkg/resources/account"
	"github.com/Southclaws/storyden/backend/pkg/services/account"
	"github.com/Southclaws/storyden/backend/pkg/transports/http/openapi"
)

type Accounts struct {
	as account.Service
}

func NewAccounts(as account.Service) Accounts { return Accounts{as} }

func (i *Accounts) GetAccount(ctx context.Context, request openapi.GetAccountRequestObject) any {
	acc, err := i.as.Get(ctx, account_resource.AccountID(request.Id))
	if err != nil {
		return err
	}

	return openapi.GetAccountSuccess{
		Id:        openapi.UUID(acc.ID),
		Bio:       utils.Ref(acc.Bio.ElseZero()),
		Email:     utils.Ref(acc.Email),
		Name:      utils.Ref(acc.Name),
		CreatedAt: utils.Ref(acc.CreatedAt.Format(time.RFC3339)),
		UpdatedAt: utils.Ref(acc.UpdatedAt.Format(time.RFC3339)),
		DeletedAt: utils.OptionalElsePtr(acc.DeletedAt, utils.FormatISO),
	}
}
