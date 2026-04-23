package mailer

import (
	"context"
	"fmt"
	"net/mail"
	"sync"
)

type MockEmail struct {
	Address mail.Address
	Name    string
	Subject string
	Html    string
	Plain   string
}

type Mock struct {
	mu      sync.Mutex
	sent    []MockEmail
	sendErr error
}

func (m *Mock) Send(
	ctx context.Context,
	msg Message,
) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.sendErr != nil {
		return m.sendErr
	}

	fmt.Printf(`Mock email sent to: %s <%s> '%s'
%s
`,
		msg.Name,
		msg.Address.String(),
		msg.Subject,
		msg.Content.Plain,
	)

	m.sent = append(m.sent, MockEmail{
		Address: msg.Address,
		Name:    msg.Name,
		Subject: msg.Subject,
		Html:    msg.Content.HTML,
		Plain:   msg.Content.Plain,
	})

	return nil
}

func (m *Mock) GetLast() MockEmail {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.sent[len(m.sent)-1]
}

func (m *Mock) Count() int {
	m.mu.Lock()
	defer m.mu.Unlock()

	return len(m.sent)
}

func (m *Mock) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.sent = nil
	m.sendErr = nil
}

func (m *Mock) SetSendError(err error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.sendErr = err
}

func (m *Mock) ClearSendError() {
	m.SetSendError(nil)
}
