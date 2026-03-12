package mailer_test

import (
	"context"
	"io"
	"log/slog"
	"net/mail"
	"testing"
	"time"

	smtpmock "github.com/mocktools/go-smtp-mock/v2"
	"github.com/stretchr/testify/require"

	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/infrastructure/mailer"
)

func TestSMTPMailer_ValidConfiguration(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	tests := []struct {
		name        string
		cfg         config.Config
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid configuration",
			cfg: config.Config{
				EmailProvider:   "smtp",
				SMTPHost:        "smtp.example.com",
				SMTPPort:        587,
				SMTPFromAddress: "noreply@example.com",
				SMTPFromName:    "Test Service",
				SMTPUseTLS:      true,
				JWTSecret:       []byte("test-secret-key"),
			},
		},
		{
			name: "missing host",
			cfg: config.Config{
				EmailProvider:   "smtp",
				SMTPPort:        587,
				SMTPFromAddress: "noreply@example.com",
				JWTSecret:       []byte("test-secret-key"),
			},
			expectError: true,
			errorMsg:    "SMTP_HOST must be provided",
		},
		{
			name: "missing port",
			cfg: config.Config{
				EmailProvider:   "smtp",
				SMTPHost:        "smtp.example.com",
				SMTPFromAddress: "noreply@example.com",
				JWTSecret:       []byte("test-secret-key"),
			},
			expectError: true,
			errorMsg:    "SMTP_PORT must be provided",
		},
		{
			name: "invalid port",
			cfg: config.Config{
				EmailProvider:   "smtp",
				SMTPHost:        "smtp.example.com",
				SMTPPort:        70000,
				SMTPFromAddress: "noreply@example.com",
				JWTSecret:       []byte("test-secret-key"),
			},
			expectError: true,
			errorMsg:    "SMTP_PORT must be between 1 and 65535",
		},
		{
			name: "missing from address",
			cfg: config.Config{
				EmailProvider: "smtp",
				SMTPHost:      "smtp.example.com",
				SMTPPort:      587,
				JWTSecret:     []byte("test-secret-key"),
			},
			expectError: true,
			errorMsg:    "SMTP_FROM_ADDRESS must be provided",
		},
		{
			name: "username without password",
			cfg: config.Config{
				EmailProvider:   "smtp",
				SMTPHost:        "smtp.example.com",
				SMTPPort:        587,
				SMTPUsername:    "user@example.com",
				SMTPFromAddress: "noreply@example.com",
				SMTPUseTLS:      true,
				JWTSecret:       []byte("test-secret-key"),
			},
			expectError: true,
			errorMsg:    "SMTP_USERNAME and SMTP_PASSWORD must both be provided together",
		},
		{
			name: "password without username",
			cfg: config.Config{
				EmailProvider:   "smtp",
				SMTPHost:        "smtp.example.com",
				SMTPPort:        587,
				SMTPPassword:    "password",
				SMTPFromAddress: "noreply@example.com",
				SMTPUseTLS:      true,
				JWTSecret:       []byte("test-secret-key"),
			},
			expectError: true,
			errorMsg:    "SMTP_USERNAME and SMTP_PASSWORD must both be provided together",
		},
		{
			name: "credentials without TLS",
			cfg: config.Config{
				EmailProvider:   "smtp",
				SMTPHost:        "smtp.example.com",
				SMTPPort:        587,
				SMTPUsername:    "user@example.com",
				SMTPPassword:    "password",
				SMTPFromAddress: "noreply@example.com",
				SMTPUseTLS:      false,
				JWTSecret:       []byte("test-secret-key"),
			},
			expectError: true,
			errorMsg:    "SMTP authentication requires TLS",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sender, err := mailer.NewMailer(logger, tt.cfg)

			if tt.expectError {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.errorMsg)
				require.Nil(t, sender)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, sender)
		})
	}
}

func TestSMTPMailer_EndToEndWithoutAuthentication(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	server := newMockServer(t)

	cfg := smtpConfig(server)
	sender, err := mailer.NewMailer(logger, cfg)
	require.NoError(t, err)

	addr := mail.Address{Name: "No Auth Recipient", Address: "no-auth-recipient@example.com"}
	content, err := mailer.NewContent("", "Plain text content without auth")
	require.NoError(t, err)
	message, err := mailer.NewMessage(addr, "No Auth Recipient", "No Auth Test Subject", *content)
	require.NoError(t, err)

	err = sender.Send(context.Background(), *message)
	require.NoError(t, err)

	messages, err := server.WaitForMessages(1, 2*time.Second)
	require.NoError(t, err)
	require.Len(t, messages, 1)
	require.True(t, messages[0].IsConsistent())
	require.Contains(t, messages[0].MailfromRequest(), "sender@example.com")
	rcpt := messages[0].RcpttoRequestResponse()
	require.Len(t, rcpt, 1)
	require.Contains(t, rcpt[0][0], "no-auth-recipient@example.com")
	msg := messages[0].MsgRequest()
	require.Contains(t, msg, "From: \"Test Sender\" <sender@example.com>")
	require.Contains(t, msg, "To: \"No Auth Recipient\" <no-auth-recipient@example.com>")
	require.Contains(t, msg, "Subject: No Auth Test Subject")
	require.Contains(t, msg, "Content-Type: text/plain; charset=UTF-8")
	require.Contains(t, msg, "Plain text content without auth")
}

func TestSMTPMailer_EndToEndMultipartAlternative(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	server := newMockServer(t)

	cfg := smtpConfig(server)
	sender, err := mailer.NewMailer(logger, cfg)
	require.NoError(t, err)

	addr := mail.Address{Name: "Recipient", Address: "recipient@example.com"}
	content, err := mailer.NewContent("<h1>Test Email</h1><p>HTML body</p>", "Test Email\n\nPlain body")
	require.NoError(t, err)
	message, err := mailer.NewMessage(addr, "Recipient", "Multipart Subject", *content)
	require.NoError(t, err)

	err = sender.Send(context.Background(), *message)
	require.NoError(t, err)

	messages, err := server.WaitForMessages(1, 2*time.Second)
	require.NoError(t, err)
	require.Len(t, messages, 1)
	require.True(t, messages[0].IsConsistent())
	msg := messages[0].MsgRequest()
	require.Contains(t, msg, "Content-Type: multipart/alternative; boundary=")
	require.Contains(t, msg, "Content-Type: text/plain; charset=UTF-8")
	require.Contains(t, msg, "Content-Type: text/html; charset=UTF-8")
	require.Contains(t, msg, "Test Email\r\n\r\nPlain body")
	require.Contains(t, msg, "<h1>Test Email</h1><p>HTML body</p>")
}

func TestSMTPMailer_RejectsHeaderInjection(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	server := newMockServer(t)

	cfg := smtpConfig(server)
	sender, err := mailer.NewMailer(logger, cfg)
	require.NoError(t, err)

	addr := mail.Address{Name: "Recipient", Address: "recipient@example.com"}
	content, err := mailer.NewContent("", "body")
	require.NoError(t, err)
	message, err := mailer.NewMessage(addr, "Recipient", "safe\r\nBcc: attacker@example.com", *content)
	require.NoError(t, err)

	err = sender.Send(context.Background(), *message)
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid SMTP header value")
	require.Len(t, server.Messages(), 0)
}

func TestSMTPMailer_RejectsInvalidFromAddressConfiguration(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	cfg := config.Config{
		EmailProvider:   "smtp",
		SMTPHost:        "smtp.example.com",
		SMTPPort:        587,
		SMTPFromAddress: "sender@example.com\r\nBcc: attacker@example.com",
		SMTPUseTLS:      true,
		JWTSecret:       []byte("test-secret-key"),
	}

	sender, err := mailer.NewMailer(logger, cfg)
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid SMTP header value")
	require.Nil(t, sender)
}

func TestSMTPMailer_TLSRequiresStartTLS(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	server := newMockServer(t)

	cfg := smtpConfig(server)
	cfg.SMTPUseTLS = true
	sender, err := mailer.NewMailer(logger, cfg)
	require.NoError(t, err)

	addr := mail.Address{Name: "Recipient", Address: "recipient@example.com"}
	content, err := mailer.NewContent("", "body")
	require.NoError(t, err)
	message, err := mailer.NewMessage(addr, "Recipient", "subject", *content)
	require.NoError(t, err)

	err = sender.Send(context.Background(), *message)
	require.Error(t, err)
	require.Contains(t, err.Error(), "STARTTLS")
	for _, msg := range server.Messages() {
		require.False(t, msg.IsConsistent())
	}
}

func TestSMTPMailer_RejectsAuthenticationWithoutTLS(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	cfg := config.Config{
		EmailProvider:   "smtp",
		SMTPHost:        "smtp.example.com",
		SMTPPort:        587,
		SMTPUsername:    "testuser",
		SMTPPassword:    "testpass",
		SMTPFromAddress: "sender@example.com",
		SMTPUseTLS:      false,
		JWTSecret:       []byte("test-secret-key"),
	}
	sender, err := mailer.NewMailer(logger, cfg)
	require.Error(t, err)
	require.Contains(t, err.Error(), "requires TLS")
	require.Nil(t, sender)
}

func TestSMTPMailer_RespectsCanceledContext(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	server := newMockServer(t)

	cfg := smtpConfig(server)
	sender, err := mailer.NewMailer(logger, cfg)
	require.NoError(t, err)

	addr := mail.Address{Name: "Recipient", Address: "recipient@example.com"}
	content, err := mailer.NewContent("", "body")
	require.NoError(t, err)
	message, err := mailer.NewMessage(addr, "Recipient", "subject", *content)
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err = sender.Send(ctx, *message)
	require.Error(t, err)
	require.Contains(t, err.Error(), "context canceled")
	require.Len(t, server.Messages(), 0)
}

func smtpConfig(server *smtpmock.Server) config.Config {
	return config.Config{
		EmailProvider:   "smtp",
		SMTPHost:        "localhost",
		SMTPPort:        server.PortNumber(),
		SMTPFromAddress: "sender@example.com",
		SMTPFromName:    "Test Sender",
		SMTPUseTLS:      false,
		JWTSecret:       []byte("test-secret-key"),
	}
}

func newMockServer(t *testing.T) *smtpmock.Server {
	t.Helper()

	server := smtpmock.New(smtpmock.ConfigurationAttr{})
	err := server.Start()
	require.NoError(t, err)

	t.Cleanup(func() {
		err := server.Stop()
		require.NoError(t, err)
	})

	return server
}
