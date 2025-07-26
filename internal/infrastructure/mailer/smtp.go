package mailer

import (
	"context"
	"crypto/tls"
	"fmt"
	"log/slog"
	"net/smtp"
	"strings"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/storyden/internal/config"
)

var ErrSMTPFailed = fault.New("SMTP sending failed")

type SMTP struct {
	logger      *slog.Logger
	host        string
	port        int
	username    string
	password    string
	fromName    string
	fromAddress string
	useTLS      bool
}

func newSMTPMailer(logger *slog.Logger, cfg config.Config) (*SMTP, error) {
	if cfg.SMTPHost == "" {
		return nil, fault.New("SMTP_HOST must be provided when using SMTP email provider")
	}
	if cfg.SMTPPort == 0 {
		return nil, fault.New("SMTP_PORT must be provided when using SMTP email provider")
	}
	if cfg.SMTPFromAddress == "" {
		return nil, fault.New("SMTP_FROM_ADDRESS must be provided when using SMTP email provider")
	}

	s := &SMTP{
		logger:      logger.With(slog.String("mailer", "smtp")),
		host:        cfg.SMTPHost,
		port:        cfg.SMTPPort,
		username:    cfg.SMTPUsername,
		password:    cfg.SMTPPassword,
		fromName:    cfg.SMTPFromName,
		fromAddress: cfg.SMTPFromAddress,
		useTLS:      cfg.SMTPUseTLS,
	}

	return s, nil
}

func (s *SMTP) Send(
	ctx context.Context,
	msg Message,
) error {
	s.logger.Info("sending SMTP email",
		slog.String("email", msg.Address.Address),
		slog.String("name", msg.Name),
		slog.String("subject", msg.Subject),
		slog.String("host", s.host),
		slog.Int("port", s.port),
	)

	// Build the email message
	from := s.fromAddress
	if s.fromName != "" {
		from = fmt.Sprintf("%s <%s>", s.fromName, s.fromAddress)
	}

	to := msg.Address.Address
	if msg.Name != "" {
		to = fmt.Sprintf("%s <%s>", msg.Name, msg.Address.Address)
	}

	// Determine which content to send (prefer HTML over plain text)
	var body string
	var contentType string
	if msg.Content.HTML != "" {
		body = msg.Content.HTML
		contentType = "text/html; charset=UTF-8"
	} else {
		body = msg.Content.Plain
		contentType = "text/plain; charset=UTF-8"
	}

	// Build the message
	message := []string{
		fmt.Sprintf("From: %s", from),
		fmt.Sprintf("To: %s", to),
		fmt.Sprintf("Subject: %s", msg.Subject),
		fmt.Sprintf("Content-Type: %s", contentType),
		"",
		body,
	}

	messageBytes := []byte(strings.Join(message, "\r\n"))

	// Connect to the SMTP server
	addr := fmt.Sprintf("%s:%d", s.host, s.port)
	
	if s.useTLS {
		// Use TLS connection
		err := s.sendWithTLS(addr, messageBytes, msg.Address.Address)
		if err != nil {
			return fault.Wrap(err, fctx.With(ctx))
		}
	} else {
		// Use plain connection
		err := s.sendPlain(addr, messageBytes, msg.Address.Address)
		if err != nil {
			return fault.Wrap(err, fctx.With(ctx))
		}
	}

	return nil
}

func (s *SMTP) sendWithTLS(addr string, messageBytes []byte, toAddr string) error {
	// Set up authentication if credentials are provided
	var auth smtp.Auth
	if s.username != "" && s.password != "" {
		auth = smtp.PlainAuth("", s.username, s.password, s.host)
	}

	// Connect to the server with TLS
	tlsConfig := &tls.Config{
		ServerName: s.host,
	}

	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		return fmt.Errorf("failed to connect with TLS: %w", err)
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, s.host)
	if err != nil {
		return fmt.Errorf("failed to create SMTP client: %w", err)
	}
	defer client.Quit()

	// Authenticate if needed
	if auth != nil {
		if err := client.Auth(auth); err != nil {
			return fmt.Errorf("authentication failed: %w", err)
		}
	}

	// Set the sender and recipient
	if err := client.Mail(s.fromAddress); err != nil {
		return fmt.Errorf("failed to set sender: %w", err)
	}

	if err := client.Rcpt(toAddr); err != nil {
		return fmt.Errorf("failed to set recipient: %w", err)
	}

	// Send the message
	writer, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to get data writer: %w", err)
	}

	_, err = writer.Write(messageBytes)
	if err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	return writer.Close()
}

func (s *SMTP) sendPlain(addr string, messageBytes []byte, toAddr string) error {
	// Set up authentication if credentials are provided
	var auth smtp.Auth
	if s.username != "" && s.password != "" {
		auth = smtp.PlainAuth("", s.username, s.password, s.host)
	}

	// Use the standard SendMail function for plain connections
	return smtp.SendMail(addr, auth, s.fromAddress, []string{toAddr}, messageBytes)
}