package account

import (
	"context"

	"github.com/pkg/errors"

	"github.com/Southclaws/storyden/app/resources/account"
)

func (s *service) Get(ctx context.Context, id account.AccountID) (*account.Account, error) {
	acc, err := s.account_repo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get account by ID")
	}

	return acc, nil
}
