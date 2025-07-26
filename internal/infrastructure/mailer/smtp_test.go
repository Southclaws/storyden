package mailer_test

import (
	"context"
	"log/slog"
	"net/mail"
	"os"
	"testing"
	"time"

	smtpmock "github.com/mocktools/go-smtp-mock/v2"
	"github.com/stretchr/testify/require"

	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/infrastructure/mailer"
)

func TestSMTPMailer_ValidConfiguration(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	
	tests := []struct {
		name        string
		cfg         config.Config
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid configuration",
			cfg: config.Config{
				EmailProvider:     "smtp",
				SMTPHost:         "smtp.example.com",
				SMTPPort:         587,
				SMTPUsername:     "user@example.com",
				SMTPPassword:     "password",
				SMTPFromAddress:  "noreply@example.com",
				SMTPFromName:     "Test Service",
				SMTPUseTLS:       true,
				JWTSecret:        []byte("test-secret-key"),
			},
			expectError: false,
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sender, err := mailer.NewMailer(logger, tt.cfg)
			
			if tt.expectError {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.errorMsg)
				require.Nil(t, sender)
			} else {
				require.NoError(t, err)
				require.NotNil(t, sender)
			}
		})
	}
}

func TestSMTPMailer_MessageConstruction(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	
	cfg := config.Config{
		EmailProvider:   "smtp",
		SMTPHost:        "smtp.example.com",
		SMTPPort:        587,
		SMTPFromAddress: "noreply@example.com",
		SMTPFromName:    "Test Service",
		SMTPUseTLS:      true,
		JWTSecret:       []byte("test-secret-key"),
	}

	sender, err := mailer.NewMailer(logger, cfg)
	require.NoError(t, err)
	require.NotNil(t, sender)

	// Test message creation
	addr := mail.Address{
		Name:    "Test User",
		Address: "test@example.com",
	}

	content, err := mailer.NewContent("<h1>Hello</h1>", "Hello")
	require.NoError(t, err)

	message, err := mailer.NewMessage(addr, "Test User", "Test Subject", *content)
	require.NoError(t, err)

	// This would normally send an email, but since we don't have a real SMTP server
	// we'll just verify the message is constructed properly
	require.Equal(t, "test@example.com", message.Address.Address)
	require.Equal(t, "Test User", message.Name)
	require.Equal(t, "Test Subject", message.Subject)
	require.Equal(t, "<h1>Hello</h1>", message.Content.HTML)
	require.Equal(t, "Hello", message.Content.Plain)
}

func TestSMTPMailer_EndToEnd(t *testing.T) {
	// Start mock SMTP server
	server := smtpmock.New(smtpmock.ConfigurationAttr{
		LogToStdout:       false,
		LogServerActivity: false,
		PortNumber:        2525, // Use non-standard port to avoid conflicts
	})

	err := server.Start()
	require.NoError(t, err, "Failed to start mock SMTP server")
	defer func() {
		err := server.Stop()
		require.NoError(t, err, "Failed to stop mock SMTP server")
	}()

	// Wait a bit for server to start
	time.Sleep(100 * time.Millisecond)

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	
	cfg := config.Config{
		EmailProvider:   "smtp",
		SMTPHost:        "localhost",
		SMTPPort:        2525,
		SMTPUsername:    "", // Mock server doesn't require auth
		SMTPPassword:    "",
		SMTPFromAddress: "sender@example.com",
		SMTPFromName:    "Test Sender",
		SMTPUseTLS:      false, // Mock server doesn't use TLS
		JWTSecret:       []byte("test-secret-key"),
	}

	sender, err := mailer.NewMailer(logger, cfg)
	require.NoError(t, err)
	require.NotNil(t, sender)

	// Create test message
	addr := mail.Address{
		Name:    "Test Recipient", 
		Address: "recipient@example.com",
	}

	content, err := mailer.NewContent("<h1>Test Email</h1><p>This is a test email with HTML content.</p>", "Test Email\n\nThis is a test email with plain text content.")
	require.NoError(t, err)

	message, err := mailer.NewMessage(addr, "Test Recipient", "End-to-End Test Subject", *content)
	require.NoError(t, err)

	// Send the email
	ctx := context.Background()
	err = sender.Send(ctx, *message)
	require.NoError(t, err, "Failed to send email")

	// Wait a bit for message to be processed
	time.Sleep(100 * time.Millisecond)

	// Verify the email was received by the mock server
	messages := server.Messages()
	require.Len(t, messages, 1, "Expected exactly one message to be sent")

	receivedMsg := messages[0]
	
	// Verify the message is consistent (all SMTP commands were successful)
	require.True(t, receivedMsg.IsConsistent(), "Expected message to be consistent")
	
	// Verify sender from MAILFROM command
	mailfromRequest := receivedMsg.MailfromRequest()
	require.Contains(t, mailfromRequest, "sender@example.com", "Expected sender address in MAILFROM")
	
	// Verify recipient from RCPTTO command
	rcpttoResponses := receivedMsg.RcpttoRequestResponse()
	require.Len(t, rcpttoResponses, 1, "Expected exactly one RCPTTO command")
	require.Contains(t, rcpttoResponses[0][0], "recipient@example.com", "Expected recipient address in RCPTTO")
	
	// Verify message content
	msgContent := receivedMsg.MsgRequest()
	require.Contains(t, msgContent, "From: Test Sender <sender@example.com>")
	require.Contains(t, msgContent, "To: Test Recipient <recipient@example.com>") 
	require.Contains(t, msgContent, "Subject: End-to-End Test Subject")
	require.Contains(t, msgContent, "Content-Type: text/html; charset=UTF-8")
	require.Contains(t, msgContent, "<h1>Test Email</h1><p>This is a test email with HTML content.</p>")
}

func TestSMTPMailer_EndToEndWithAuthentication(t *testing.T) {
	// This test verifies that the SMTP mailer can be configured with authentication credentials
	// and will attempt to authenticate if credentials are provided.
	// Since go-smtp-mock doesn't support enforcing authentication, we'll test that the mailer
	// handles the case where the server doesn't support AUTH gracefully.
	
	// Start mock SMTP server
	server := smtpmock.New(smtpmock.ConfigurationAttr{
		LogToStdout:       false,
		LogServerActivity: false,
		PortNumber:        2526, // Use different port
	})

	err := server.Start()
	require.NoError(t, err, "Failed to start mock SMTP server")
	defer func() {
		err := server.Stop()
		require.NoError(t, err, "Failed to stop mock SMTP server")
	}()

	// Wait a bit for server to start
	time.Sleep(100 * time.Millisecond)

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	
	cfg := config.Config{
		EmailProvider:   "smtp",
		SMTPHost:        "localhost",
		SMTPPort:        2526,
		SMTPUsername:    "testuser",
		SMTPPassword:    "testpass",
		SMTPFromAddress: "auth-sender@example.com",
		SMTPFromName:    "Authenticated Sender",
		SMTPUseTLS:      false, // Mock server doesn't use TLS
		JWTSecret:       []byte("test-secret-key"),
	}

	sender, err := mailer.NewMailer(logger, cfg)
	require.NoError(t, err)
	require.NotNil(t, sender)

	// Create test message
	addr := mail.Address{
		Name:    "Auth Test Recipient",
		Address: "auth-recipient@example.com",
	}

	content, err := mailer.NewContent("", "This is a plain text email for authentication testing.")
	require.NoError(t, err)

	message, err := mailer.NewMessage(addr, "Auth Test Recipient", "Authentication Test Subject", *content)
	require.NoError(t, err)

	// Send the email - this should fail because mock server doesn't support AUTH
	ctx := context.Background()
	err = sender.Send(ctx, *message)
	require.Error(t, err, "Expected error because mock server doesn't support AUTH")
	require.Contains(t, err.Error(), "server doesn't support AUTH", "Error should mention AUTH not supported")

	// Verify that while SMTP conversation may have started, no complete message was sent
	// (some SMTP commands like EHLO might be recorded even when auth fails)
	messages := server.Messages()
	// Either no messages, or if messages exist, they should not be consistent (complete)
	for _, msg := range messages {
		require.False(t, msg.IsConsistent(), "Expected incomplete message due to AUTH failure")
	}
}

func TestSMTPMailer_EndToEndWithoutAuthentication(t *testing.T) {
	// This test verifies that the SMTP mailer works correctly when no authentication
	// credentials are provided (common for local development or trusted networks)
	
	// Start mock SMTP server
	server := smtpmock.New(smtpmock.ConfigurationAttr{
		LogToStdout:       false,
		LogServerActivity: false,
		PortNumber:        2527, // Use different port
	})

	err := server.Start()
	require.NoError(t, err, "Failed to start mock SMTP server")
	defer func() {
		err := server.Stop()
		require.NoError(t, err, "Failed to stop mock SMTP server")
	}()

	// Wait a bit for server to start
	time.Sleep(100 * time.Millisecond)

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	
	cfg := config.Config{
		EmailProvider:   "smtp",
		SMTPHost:        "localhost",
		SMTPPort:        2527,
		SMTPUsername:    "", // No authentication
		SMTPPassword:    "",
		SMTPFromAddress: "no-auth-sender@example.com",
		SMTPFromName:    "No Auth Sender",
		SMTPUseTLS:      false,
		JWTSecret:       []byte("test-secret-key"),
	}

	sender, err := mailer.NewMailer(logger, cfg)
	require.NoError(t, err)
	require.NotNil(t, sender)

	// Create test message
	addr := mail.Address{
		Name:    "No Auth Recipient",
		Address: "no-auth-recipient@example.com",
	}

	content, err := mailer.NewContent("<p>HTML content without auth</p>", "Plain text content without auth")
	require.NoError(t, err)

	message, err := mailer.NewMessage(addr, "No Auth Recipient", "No Auth Test Subject", *content)
	require.NoError(t, err)

	// Send the email
	ctx := context.Background()
	err = sender.Send(ctx, *message)
	require.NoError(t, err, "Failed to send email without authentication")

	// Wait a bit for message to be processed
	time.Sleep(100 * time.Millisecond)

	// Verify the email was received by the mock server
	messages := server.Messages()
	require.Len(t, messages, 1, "Expected exactly one message to be sent")

	receivedMsg := messages[0]
	
	// Verify the message is consistent (all SMTP commands were successful)
	require.True(t, receivedMsg.IsConsistent(), "Expected message to be consistent")
	
	// Verify sender from MAILFROM command
	mailfromRequest := receivedMsg.MailfromRequest()
	require.Contains(t, mailfromRequest, "no-auth-sender@example.com", "Expected sender address in MAILFROM")
	
	// Verify recipient from RCPTTO command
	rcpttoResponses := receivedMsg.RcpttoRequestResponse()
	require.Len(t, rcpttoResponses, 1, "Expected exactly one RCPTTO command")
	require.Contains(t, rcpttoResponses[0][0], "no-auth-recipient@example.com", "Expected recipient address in RCPTTO")
	
	// Verify message content (should prefer HTML over plain text)
	msgContent := receivedMsg.MsgRequest()
	require.Contains(t, msgContent, "From: No Auth Sender <no-auth-sender@example.com>")
	require.Contains(t, msgContent, "To: No Auth Recipient <no-auth-recipient@example.com>")
	require.Contains(t, msgContent, "Subject: No Auth Test Subject")
	require.Contains(t, msgContent, "Content-Type: text/html; charset=UTF-8")
	require.Contains(t, msgContent, "<p>HTML content without auth</p>")
}