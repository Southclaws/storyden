package onboarding

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/internal/ent"
)

//go:generate go run -mod=mod github.com/Southclaws/enumerator

type Service interface {
	GetOnboardingStatus(ctx context.Context) (*Status, error)
}

type statusEnum string

const (
	statusRequiresFirstAccount statusEnum = `requires_first_account`
	statusRequiresCategory     statusEnum = `requires_category`
	statusRequiresMoreAccounts statusEnum = `requires_more_accounts`
	statusRequiresFirstPost    statusEnum = `requires_first_post`
	statusComplete             statusEnum = `complete`
)

type service struct {
	cachedStatus Status
	ec           *ent.Client
}

func Build() fx.Option {
	return fx.Provide(func(lc fx.Lifecycle, ec *ent.Client) Service {
		s := &service{
			cachedStatus: StatusRequiresFirstAccount,
			ec:           ec,
		}

		lc.Append(fx.Hook{
			OnStart: func(ctx context.Context) error {
				initial, err := s.GetOnboardingStatus(ctx)
				if err != nil {
					return fault.Wrap(err, fctx.With(ctx))
				}

				// NOTE: not thread safe but probably fine.
				s.cachedStatus = *initial

				return nil
			},
		})

		return s
	})
}

func (s *service) GetOnboardingStatus(ctx context.Context) (*Status, error) {
	if s.cachedStatus == StatusComplete {
		return &StatusComplete, nil
	}

	accounts, err := s.ec.Account.Query().Count(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if accounts == 0 {
		return &StatusRequiresFirstAccount, nil
	}

	categories, err := s.ec.Category.Query().Count(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if categories == 0 {
		return &StatusRequiresCategory, nil
	}

	if accounts == 1 {
		return &StatusRequiresMoreAccounts, nil
	}

	posts, err := s.ec.Post.Query().Count(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if posts == 0 {
		return &StatusRequiresFirstPost, nil
	}

	return &StatusComplete, nil
}
