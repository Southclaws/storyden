package mailer_test

import (
	"log/slog"
	"net/mail"
	"os"
	"testing"

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