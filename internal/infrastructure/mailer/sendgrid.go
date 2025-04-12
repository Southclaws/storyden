package mailer

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

var ErrSendgridFailed = fault.New("sendgrid responded with an unexpected status code")

type SendGrid struct {
	logger      *slog.Logger
	client      *sendgrid.Client
	fromName    string
	fromAddress string
}

const attachmentContentDisposition = "attachment"

func newSendgridMailer(logger *slog.Logger, cfg config.Config) (*SendGrid, error) {
	sg := &SendGrid{
		logger:      logger.With(slog.String("mailer", "sendgrid")),
		client:      sendgrid.NewSendClient(cfg.SendGridAPIKey),
		fromName:    cfg.SendGridFromName,
		fromAddress: cfg.SendGridFromAddress,
	}

	return sg, nil
}

func (m *SendGrid) Send(
	ctx context.Context,
	msg Message,
) error {
	from := mail.NewEmail(m.fromName, m.fromAddress)
	to := mail.NewEmail(msg.Name, msg.Address.Address)
	message := mail.NewSingleEmail(from, msg.Subject, to, msg.Content.Plain, msg.Content.HTML)

	m.logger.Info("sending live email",
		slog.String("email", to.Address),
		slog.String("name", to.Name),
		slog.String("subject", msg.Subject),
	)

	res, err := m.client.SendWithContext(ctx, message)
	if err != nil {
		return fault.Wrap(err, fmsg.With(res.Body))
	}
	if res.StatusCode != http.StatusAccepted {
		return fault.Wrap(ErrSendgridFailed, fctx.With(ctx))
	}

	return nil
}
