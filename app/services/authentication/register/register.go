package register

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/services/onboarding"
)

type Service interface {
	Create(ctx context.Context, handle string, opts ...account.Option) (*account.Account, error)
}

func New(ar account.Repository, os onboarding.Service) Service {
	return &service{ar, os}
}

type service struct {
	ar account.Repository
	os onboarding.Service
}

func (s *service) Create(ctx context.Context, handle string, opts ...account.Option) (*account.Account, error) {
	status, err := s.os.GetOnboardingStatus(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if status == &onboarding.StatusRequiresFirstAccount {
		// If we're doing first-time-setup then set the first account to admin.
		opts = append(opts, account.WithAdmin(true))
	}

	acc, err := s.ar.Create(ctx, handle, opts...)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return acc, nil
}
