package mailer

import (
	"context"
	"log/slog"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/resend/resend-go/v2"

	"github.com/Southclaws/storyden/internal/config"
)

var ErrResendFailed = fault.New("resend responded with an error")

type Resend struct {
	logger      *slog.Logger
	client      *resend.Client
	fromName    string
	fromAddress string
}

func newResendMailer(logger *slog.Logger, cfg config.Config) (*Resend, error) {
	client := resend.NewClient(cfg.ResendAPIKey)

	rs := &Resend{
		logger:      logger.With(slog.String("mailer", "resend")),
		client:      client,
		fromName:    cfg.ResendFromName,
		fromAddress: cfg.ResendFromAddress,
	}

	return rs, nil
}

func (m *Resend) Send(ctx context.Context, msg Message) error {
	params := &resend.SendEmailRequest{
		From:    m.fromName + " <" + m.fromAddress + ">",
		To:      []string{msg.Address.Address},
		Subject: msg.Subject,
	}

	if msg.Content.HTML != "" {
		params.Html = msg.Content.HTML
	}

	if msg.Content.Plain != "" {
		params.Text = msg.Content.Plain
	}

	m.logger.Info("sending live email",
		slog.String("email", msg.Address.Address),
		slog.String("name", msg.Name),
		slog.String("subject", msg.Subject),
	)

	sent, err := m.client.Emails.SendWithContext(ctx, params)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	m.logger.Debug("email sent successfully",
		slog.String("id", sent.Id),
		slog.String("email", msg.Address.Address),
	)

	return nil
}
