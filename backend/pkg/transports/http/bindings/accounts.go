package bindings

import (
	"context"
	"time"

	"4d63.com/optional"

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

	return openapi.GetAccount200JSONResponse(openapi.Account{
		Id:        openapi.UUID(acc.ID),
		Bio:       utils.Ref(acc.Bio.ElseZero()),
		Email:     utils.Ref(acc.Email),
		Name:      utils.Ref(acc.Name),
		CreatedAt: utils.Ref(acc.CreatedAt.Format(time.RFC3339)),
		UpdatedAt: utils.Ref(acc.UpdatedAt.Format(time.RFC3339)),
		DeletedAt: OptionalElsePtr(acc.DeletedAt, FormatISO),
	})
}

func OptionalToPointer[T any](o optional.Optional[T]) *T {
	if v, ok := o.Get(); ok {
		return &v
	}
	return nil
}

func OptionalElsePtr[T, R any](o optional.Optional[T], fn func(T) R) *R {
	if v, ok := o.Get(); ok {
		r := fn(v)
		return &r
	}
	return nil
}

func FormatISO(t time.Time) string { return t.Format(time.RFC3339) }
