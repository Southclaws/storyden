package sms

import (
	"context"
	"fmt"
	"log/slog"
)

type MockSender struct{}

func newMock(l *slog.Logger) (Sender, error) {
	l.Debug("using mock sms sender - check the console for outgoing messages")
	return &MockSender{}, nil
}

func (s *MockSender) Send(ctx context.Context, phone string, message string) error {
	fmt.Printf(`
[MOCK SMS] to: "%s" message:

%s

`, phone, message)
	return nil
}
