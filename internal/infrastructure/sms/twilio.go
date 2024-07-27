package sms

import (
	"context"
	"fmt"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/kelseyhightower/envconfig"
	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
	"go.uber.org/zap"
)

var errFailedToSend = fault.New("failed to send sms")

type Configuration struct {
	Enabled     bool   `envconfig:"TWILIO_ENABLED"`
	AccountSID  string `envconfig:"TWILIO_ACCOUNT_SID"`
	AuthToken   string `envconfig:"TWILIO_AUTH_TOKEN"`
	PhoneNumber string `envconfig:"TWILIO_PHONE_NUMBER"`
}

type TwilioSender struct {
	client *twilio.RestClient
	number string
}

func newTwilio(l *zap.Logger) (Sender, error) {
	pc := Configuration{}
	if err := envconfig.Process("", &pc); err != nil {
		return nil, fault.Wrap(err)
	}

	if !pc.Enabled {
		l.Info("twilio disabled")
		return nil, nil
	}

	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: pc.AccountSID,
		Password: pc.AuthToken,
	})

	acc, err := client.Api.FetchAccount(pc.AccountSID)
	if err != nil {
		return nil, fault.Wrap(err)
	}

	l.Info("twilio enabled", zap.String("account_id", *acc.Sid))

	return &TwilioSender{
		client: client,
		number: pc.PhoneNumber,
	}, nil
}

func (s *TwilioSender) Send(ctx context.Context, phone string, message string) error {
	params := &twilioApi.CreateMessageParams{}

	params.SetFrom(s.number)
	params.SetTo(phone)
	params.SetBody(message)

	resp, err := s.client.Api.CreateMessage(params)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	if resp.ErrorCode != nil {
		return fault.Wrap(errFailedToSend, fmsg.With(fmt.Sprintf("error code: %d", *resp.ErrorCode)))
	}

	return nil
}
