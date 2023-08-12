package account

import (
	"context"

	"github.com/Southclaws/dt"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/samber/lo"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/services/authentication"
)

func (s *service) Get(ctx context.Context, id account.AccountID) (*account.Account, error) {
	acc, err := s.account_repo.GetByID(ctx, id)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to get account by ID"))
	}

	return acc, nil
}

func (s *service) GetAuthMethods(ctx context.Context, id account.AccountID) ([]*AuthMethod, error) {
	acc, err := s.account_repo.GetByID(ctx, id)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	ps := s.auth_svc.Providers()

	active := lo.SliceToMap(acc.Auths, func(s string) (string, bool) { return s, true })

	return dt.Map(ps, func(p authentication.Provider) *AuthMethod {
		return &AuthMethod{
			Provider: p,
			Active:   active[p.ID()],
		}
	}), nil
}
