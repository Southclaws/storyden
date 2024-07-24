package email_only

import (
	"context"
	"fmt"
	"net/mail"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/matcornic/hermes/v2"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/app/resources/mailtemplate"
	"github.com/Southclaws/storyden/app/resources/settings"
	"github.com/Southclaws/storyden/internal/mailer"
)

type VerificationMailSender struct {
	fx.In

	Sender   mailer.Sender
	Template *mailtemplate.Builder
	Settings settings.Repository
}

func (s *VerificationMailSender) SendVerificationEmail(ctx context.Context, address mail.Address, code string) error {
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
