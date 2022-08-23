package account

import (
	"context"
	"errors"

	"github.com/Southclaws/storyden/pkg/resources/account"
	"github.com/Southclaws/storyden/pkg/services/authentication"
)

var ErrNotAuthorised = errors.New("not authorised")

func (s *service) Get(ctx context.Context, id account.AccountID) (*account.Account, error) {
	subject, err := authentication.GetAccountID(ctx)
	if err != nil {
		return nil, err
	}

	if subject != id {
		return nil, ErrNotAuthorised
	}

	return s.account_repo.GetByID(ctx, id)
}
