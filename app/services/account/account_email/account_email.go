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

func (m *Manager) AddUnverified(ctx context.Context, accountID account.AccountID, address mail.Address) (*account.EmailAddress, error) {
	code, err := otp.Generate()
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	err = m.profileCache.Invalidate(ctx, xid.ID(accountID))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	ae, err := m.emailRepo.Add(ctx, accountID, address, code)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	m.bus.Publish(ctx, &rpc.EventAccountUpdated{
		ID: accountID,
	})

	return ae, nil
}

func (m *Manager) LookupAccount(ctx context.Context, address mail.Address) (*account.AccountWithEdges, bool, error) {
	acc, exists, err := m.emailRepo.LookupAccount(ctx, address)
	if err != nil {
		return nil, false, fault.Wrap(err, fctx.With(ctx))
	}

	return acc, exists, nil
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

func (m *Manager) SetVerifiedStatus(ctx context.Context, accountID account.AccountID, id xid.ID, verified bool) (*account.EmailAddress, error) {
	err := m.profileCache.Invalidate(ctx, xid.ID(accountID))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	ae, err := m.emailRepo.SetVerifiedStatus(ctx, accountID, id, verified)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	m.bus.Publish(ctx, &rpc.EventAccountUpdated{
		ID: accountID,
	})

	return ae, nil
}
