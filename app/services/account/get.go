package account

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"

	"github.com/Southclaws/storyden/app/resources/account"
)

func (s *service) Get(ctx context.Context, id account.AccountID) (*account.Account, error) {
	acc, err := s.account_repo.GetByID(ctx, id)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to get account by ID"))
	}

	return acc, nil
}
