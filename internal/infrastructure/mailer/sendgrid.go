package mailer

import (
	"context"
	"net/http"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/kelseyhightower/envconfig"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"go.uber.org/zap"
)

var ErrSendgridFailed = fault.New("sendgrid responded with an unexpected status code")

type SendGrid struct {
	l *zap.Logger

	client      *sendgrid.Client
	fromName    string
	fromAddress string
}

type Configuration struct {
	FromName    string `envconfig:"SENDGRID_FROM_NAME"    required:"true"`
	FromAddress string `envconfig:"SENDGRID_FROM_ADDRESS" required:"true"`
	APIKey      string `envconfig:"SENDGRID_API_KEY"      required:"true"`
}

const attachmentContentDisposition = "attachment"

func newSendgridMailer(l *zap.Logger) (*SendGrid, error) {
	pc := Configuration{}
	if err := envconfig.Process("", &pc); err != nil {
		return nil, fault.Wrap(err)
	}

	sg := &SendGrid{
		l:           l.With(zap.String("mailer", "sendgrid")),
		client:      sendgrid.NewSendClient(pc.APIKey),
		fromName:    pc.FromName,
		fromAddress: pc.FromAddress,
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

	m.l.Info("sending live email",
		zap.String("email", to.Address),
		zap.String("name", to.Name),
		zap.String("subject", msg.Subject),
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
