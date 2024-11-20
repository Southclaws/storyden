package email_verify

import (
	"context"
	"fmt"
	"net/mail"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/matcornic/hermes/v2"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/email"
	"github.com/Southclaws/storyden/app/resources/settings"
	"github.com/Southclaws/storyden/app/services/comms/mailqueue"
)

var (
	ErrNotFound        = fault.New("no verification code match found", ftag.With(ftag.Unauthenticated))
	ErrAccountNotFound = fault.New("no account was found", ftag.With(ftag.NotFound))
)

type Verifier struct {
	emailRepo *email.Repository
	mailqueue *mailqueue.Queuer
	settings  *settings.SettingsRepository
}

func New(
	emailRepo *email.Repository,
	mailqueue *mailqueue.Queuer,
	settings *settings.SettingsRepository,
) *Verifier {
	return &Verifier{
		emailRepo: emailRepo,
		mailqueue: mailqueue,
		settings:  settings,
	}
}

// BeginEmailVerification adds an email record for the specified account, sets
// it to unverified and sends an email to the address with a verification code.
// You can also optionally supply an authentication record ID to link the email
// record to an authentication record to indicate the email is used for login.
func (s *Verifier) BeginEmailVerification(
	ctx context.Context,
	accountID account.AccountID,
	address mail.Address,
	code string,
) (*account.EmailAddress, error) {
	ae, err := s.emailRepo.Add(ctx, accountID, address, code)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return ae, s.sendVerification(ctx, address, code)
}

func (s *Verifier) sendVerification(ctx context.Context, address mail.Address, code string) error {
	set, err := s.settings.Get(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	recipientName := address.Address
	instanceTitle := set.Title.Or(settings.DefaultTitle)
	welcome := fmt.Sprintf("Welcome to %s!", instanceTitle)

	return s.mailqueue.Queue(ctx, address, recipientName, welcome, []string{welcome}, []hermes.Action{
		{
			Instructions: "Please use the following code to verify your account:",
			InviteCode:   code,
		},
	})
}

func (s *Verifier) ResendVerification(ctx context.Context, address mail.Address) error {
	code, err := s.emailRepo.GetCode(ctx, address)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return s.sendVerification(ctx, address, code)
}

func (s *Verifier) Verify(ctx context.Context, emailAddress mail.Address, code string) (*account.Account, error) {
	acc, exists, err := s.emailRepo.LookupCode(ctx, emailAddress, code)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}
	if !exists {
		return nil, fault.Wrap(ErrNotFound, fctx.With(ctx))
	}

	if err := acc.RejectSuspended(); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	err = s.emailRepo.Verify(ctx, acc.ID, emailAddress)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return acc, nil
}
