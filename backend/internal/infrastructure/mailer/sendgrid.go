package mailer

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type SendGrid struct {
	fromname string
	fromaddr string
	client   *sendgrid.Client
}

type Config struct {
	SendgridFromName string `envconfig:"SENDGRID_FROM_NAME" required:"true"`
	SendgridFromAddr string `envconfig:"SENDGRID_FROM_ADDR" required:"true"`
	SendgridAPIKey   string `envconfig:"SENDGRID_API_KEY" required:"true"`
}

func NewSendGrid() (Mailer, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, err
	}

	return &SendGrid{
		cfg.SendgridFromName,
		cfg.SendgridFromAddr,
		sendgrid.NewSendClient(cfg.SendgridAPIKey),
	}, nil
}

func (m *SendGrid) Mail(toname, toaddr, subj, rich, text string) error {
	from := mail.NewEmail(m.fromname, m.fromaddr)
	to := mail.NewEmail(toname, toaddr)

	message := mail.NewSingleEmail(from, subj, to, text, rich)

	_, err := m.client.Send(message)
	if err != nil {
		return errors.Wrap(err, "failed to send email via sendgrid")
	}

	return nil
}
