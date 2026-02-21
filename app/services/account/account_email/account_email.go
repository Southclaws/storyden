package account_email

import (
	"context"
	"net/mail"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/email"
	"github.com/Southclaws/storyden/app/resources/profile/profile_cache"
	"github.com/Southclaws/storyden/app/services/authentication/email_verify"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
	"github.com/Southclaws/storyden/internal/otp"
	"github.com/Southclaws/storyden/lib/plugin/rpc"
)

type Manager struct {
	emailRepo    *email.Repository
	verifier     *email_verify.Verifier
	profileCache *profile_cache.Cache
	bus          *pubsub.Bus
}

func New(emailRepo *email.Repository, verifier *email_verify.Verifier, profileCache *profile_cache.Cache, bus *pubsub.Bus) *Manager {
	return &Manager{emailRepo: emailRepo, verifier: verifier, profileCache: profileCache, bus: bus}
}

func (m *Manager) Add(ctx context.Context, accountID account.AccountID, address mail.Address) (*account.EmailAddress, error) {
	otp, err := otp.Generate()
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	err = m.profileCache.Invalidate(ctx, xid.ID(accountID))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	ae, err := m.verifier.BeginEmailVerification(ctx, accountID, address, otp)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	m.bus.Publish(ctx, &rpc.EventAccountUpdated{
		ID: accountID,
	})

	return ae, nil
}

func (m *Manager) Remove(ctx context.Context, accountID account.AccountID, id xid.ID) error {
	err := m.profileCache.Invalidate(ctx, xid.ID(accountID))
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	err = m.emailRepo.Remove(ctx, accountID, id)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	m.bus.Publish(ctx, &rpc.EventAccountUpdated{
		ID: accountID,
	})

	return nil
}
