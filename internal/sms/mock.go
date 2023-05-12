package sms

import (
	"context"
	"fmt"

	"go.uber.org/zap"
)

type MockSender struct{}

func newMock(l *zap.Logger) (Sender, error) {
	l.Info("using mock sms sender - check the console for outgoing messages")
	return &MockSender{}, nil
}

func (s *MockSender) Send(ctx context.Context, phone string, message string) error {
	fmt.Printf(`
[MOCK SMS] to: "%s" message:

%s

`, phone, message)
	return nil
}
