package email_verify

import (
	"context"
	"fmt"
	"net/mail"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"
	"github.com/matcornic/hermes/v2"
	"github.com/rs/xid"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/authentication"
	"github.com/Southclaws/storyden/app/resources/account/email"
	"github.com/Southclaws/storyden/app/resources/mailtemplate"
	"github.com/Southclaws/storyden/app/resources/settings"
	"github.com/Southclaws/storyden/internal/mailer"
)

var (
	ErrNotFound        = fault.New("no verification code match found", ftag.With(ftag.Unauthenticated))
	ErrAccountNotFound = fault.New("no account was found", ftag.With(ftag.NotFound))
)

type Verifier struct {
	fx.In

	AccountRepo account.Repository
	AuthRepo    authentication.Repository
	EmailRepo   email.EmailRepo
	Sender      mailer.Sender
	Template    *mailtemplate.Builder
	Settings    settings.Repository
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
	authRecordID opt.Optional[xid.ID],
) error {
	_, err := s.EmailRepo.Add(ctx, accountID, address, code, authRecordID)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return s.sendVerification(ctx, address, code)
}

func (s *Verifier) sendVerification(ctx context.Context, address mail.Address, code string) error {
	settings, err := s.Settings.Get(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	recipientName := address.Address
	instanceTitle := settings.Title.Get()
	welcome := fmt.Sprintf("Welcome to %s!", instanceTitle)

	template, err := s.Template.Build(ctx, recipientName, []string{welcome}, []hermes.Action{
		{
			Instructions: "Please use the following code to verify your account:",
			InviteCode:   code,
		},
	})
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return s.Sender.Send(ctx, address, recipientName, welcome, template.HTML, template.Plain)
}

func (s *Verifier) ResendVerification(ctx context.Context, address mail.Address) error {
	code, err := s.EmailRepo.GetCode(ctx, address)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return s.sendVerification(ctx, address, code)
}

func (s *Verifier) Verify(ctx context.Context, emailAddress mail.Address, code string) (*account.Account, error) {
	acc, exists, err := s.EmailRepo.LookupCode(ctx, emailAddress, code)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}
	if !exists {
		return nil, fault.Wrap(ErrNotFound, fctx.With(ctx))
	}

	if err := acc.RejectSuspended(); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	err = s.EmailRepo.Verify(ctx, acc.ID, emailAddress)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return acc, nil
}
