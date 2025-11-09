package account_email

import (
	"context"
	"net/mail"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/fault/ftag"
	"github.com/rs/xid"
	"github.com/samber/lo"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/account_querier"
	authentication_repo "github.com/Southclaws/storyden/app/resources/account/authentication"
	"github.com/Southclaws/storyden/app/resources/account/email"
	"github.com/Southclaws/storyden/app/resources/message"
	"github.com/Southclaws/storyden/app/services/authentication/email_verify"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
	"github.com/Southclaws/storyden/internal/otp"
)

type Manager struct {
	emailRepo      *email.Repository
	verifier       *email_verify.Verifier
	bus            *pubsub.Bus
	accountQuerier *account_querier.Querier
	authRepo       authentication_repo.Repository
}

func New(
	emailRepo *email.Repository,
	verifier *email_verify.Verifier,
	bus *pubsub.Bus,
	accountQuerier *account_querier.Querier,
	authRepo authentication_repo.Repository,
) *Manager {
	return &Manager{
		emailRepo:      emailRepo,
		verifier:       verifier,
		bus:            bus,
		accountQuerier: accountQuerier,
		authRepo:       authRepo,
	}
}

func (m *Manager) Add(ctx context.Context, accountID account.AccountID, address mail.Address) (*account.EmailAddress, error) {
	otp, err := otp.Generate()
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	ae, err := m.verifier.BeginEmailVerification(ctx, accountID, address, otp)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	m.bus.Publish(ctx, &message.EventAccountUpdated{
		ID: accountID,
	})

	return ae, nil
}

func (m *Manager) Remove(ctx context.Context, accountID account.AccountID, id xid.ID) error {
	// Get account with all email addresses and auth methods
	acc, err := m.accountQuerier.GetByID(ctx, accountID)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	// Get all authentication methods for this account
	authMethods, err := m.authRepo.GetAuthMethods(ctx, accountID)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	// Check if account has OAuth authentication (Discord or GitHub)
	hasOAuthAuth := lo.ContainsBy(authMethods, func(auth *authentication_repo.Authentication) bool {
		return auth.Service == authentication_repo.ServiceOAuthDiscord ||
			auth.Service == authentication_repo.ServiceOAuthGitHub ||
			auth.Service == authentication_repo.ServiceOAuthGoogle
	})

	// Count verified emails
	verifiedEmails := lo.Filter(acc.EmailAddresses, func(e *account.EmailAddress, _ int) bool {
		return e.Verified
	})

	// Find the email being deleted
	emailToDelete, found := lo.Find(acc.EmailAddresses, func(e *account.EmailAddress) bool {
		return xid.ID(e.ID) == id
	})

	if !found {
		return fault.New("email address not found",
			fctx.With(ctx),
			ftag.With(ftag.NotFound),
		)
	}

	// If trying to delete a verified email and it's the last verified one
	if emailToDelete.Verified && len(verifiedEmails) == 1 {
		// Check if there are OAuth connections
		if !hasOAuthAuth {
			return fault.New("cannot delete last verified email address",
				fctx.With(ctx),
				ftag.With(ftag.PermissionDenied),
				fmsg.WithDesc(
					"last verified email",
					"You cannot delete your last verified email address without linking an OAuth account (Discord, GitHub, or Google) first. Please add an OAuth authentication method before removing this email.",
				),
			)
		}
	}

	err = m.emailRepo.Remove(ctx, accountID, id)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	m.bus.Publish(ctx, &message.EventAccountUpdated{
		ID: accountID,
	})

	return nil
}
