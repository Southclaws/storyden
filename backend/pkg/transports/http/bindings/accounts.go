package bindings

import (
	"context"

	"github.com/Southclaws/storyden/backend/pkg/transports/http/openapi"
)

type Accounts struct{}

func NewAccounts() Accounts { return Accounts{} }

func (i *Accounts) GetAccount(ctx context.Context, request openapi.GetAccountRequestObject) any {
	return nil
}
