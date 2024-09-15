package register

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/services/onboarding"
)

type Registrar struct {
	writer     *account_writer.Writer
	onboarding onboarding.Service
}

func New(
	writer *account_writer.Writer,
	onboarding onboarding.Service,
) *Registrar {
	return &Registrar{
		writer:     writer,
		onboarding: onboarding,
	}
}

func (s *Registrar) Create(ctx context.Context, handle string, opts ...account_writer.Option) (*account.Account, error) {
	status, err := s.onboarding.GetOnboardingStatus(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if status == &onboarding.StatusRequiresFirstAccount {
		// If we're doing first-time-setup then set the first account to admin.
		opts = append(opts, account_writer.WithAdmin(true))
	}

	acc, err := s.writer.Create(ctx, handle, opts...)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return acc, nil
}
