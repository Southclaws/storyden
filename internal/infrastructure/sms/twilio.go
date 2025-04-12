package sms

import (
	"context"
	"fmt"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
	"go.uber.org/zap"

	"github.com/Southclaws/storyden/internal/config"
)

var errFailedToSend = fault.New("failed to send sms")

type TwilioSender struct {
	client *twilio.RestClient
	number string
}

func newTwilio(l *zap.Logger, cfg config.Config) (Sender, error) {
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: cfg.TwilioAccountSID,
		Password: cfg.TwilioAuthToken,
	})

	acc, err := client.Api.FetchAccount(cfg.TwilioAccountSID)
	if err != nil {
		return nil, fault.Wrap(err)
	}

	l.Info("twilio enabled", zap.String("account_id", *acc.Sid))

	return &TwilioSender{
		client: client,
		number: cfg.TwilioPhoneNumber,
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
