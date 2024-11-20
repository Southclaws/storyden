package password_reset

import (
	"context"
	"fmt"
	"net/mail"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/matcornic/hermes/v2"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/authentication"
	"github.com/Southclaws/storyden/app/resources/settings"
	"github.com/Southclaws/storyden/app/services/comms/mailqueue"
	"github.com/Southclaws/storyden/app/services/comms/mailtemplate"
	"github.com/Southclaws/storyden/internal/infrastructure/mailer"
)

var (
	ErrNotFound        = fault.New("no verification code match found", ftag.With(ftag.Unauthenticated))
	ErrAccountNotFound = fault.New("no account was found", ftag.With(ftag.NotFound))
)

type EmailResetter struct {
	tokenProvider *TokenProvider
	authRepo      authentication.Repository
	sender        mailer.Sender
	mailqueue     *mailqueue.Queuer
	settings      *settings.SettingsRepository
}

func NewEmailResetter(
	tokenProvider *TokenProvider,
	authRepo authentication.Repository,
	sender mailer.Sender,
	mailqueue *mailqueue.Queuer,
	settings *settings.SettingsRepository,
) *EmailResetter {
	return &EmailResetter{
		tokenProvider: tokenProvider,
		authRepo:      authRepo,
		sender:        sender,
		mailqueue:     mailqueue,
		settings:      settings,
	}
}

func (s *EmailResetter) SendPasswordReset(
	ctx context.Context,
	accountID account.AccountID,
	address mail.Address,
	lt LinkTemplate,
) error {
	token, err := s.tokenProvider.GetResetToken(ctx, accountID)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	link := lt.GetURL(token)

	return s.sendResetEmail(ctx, address, link)
}

func (s *EmailResetter) sendResetEmail(ctx context.Context, address mail.Address, link string) error {
	set, err := s.settings.Get(ctx)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	recipientName := address.Address
	instanceTitle := set.Title.Or(settings.DefaultTitle)
	welcome := fmt.Sprintf("Reset your password on %s!", instanceTitle)

	return s.mailqueue.Queue(ctx, address, recipientName, welcome, []string{welcome}, []mailtemplate.Action{
		{
			Instructions: "Click the link below to reset your password.",
			Button: hermes.Button{
				Text: "Reset password",
				Link: link,
			},
		},
	})
}

func (s *EmailResetter) Verify(ctx context.Context, token string) (account.AccountID, error) {
	accountID, err := s.tokenProvider.Validate(ctx, token)
	if err != nil {
		return account.AccountID{}, fault.Wrap(err, fctx.With(ctx))
	}

	return accountID, nil
}
