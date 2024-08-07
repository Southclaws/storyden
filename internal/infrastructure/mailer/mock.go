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
	mu   sync.Mutex
	sent []MockEmail
}

func (m *Mock) Send(
	ctx context.Context,
	address mail.Address,
	name string,
	subject string,
	html string,
	plain string,
) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	fmt.Printf(`Mock email sent to: %s <%s> '%s'
%s
`,
		name,
		address.String(),
		subject,
		plain,
	)

	m.sent = append(m.sent, MockEmail{
		Address: address,
		Name:    name,
		Subject: subject,
		Html:    html,
		Plain:   plain,
	})

	return nil
}

func (m *Mock) GetLast() MockEmail {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.sent[len(m.sent)-1]
}
