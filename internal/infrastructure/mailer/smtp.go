package mailer

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/mail"
	"net/smtp"
	"strconv"
	"strings"
	"time"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/storyden/internal/config"
)

var (
	ErrSMTPFailed          = fault.New("SMTP sending failed")
	ErrInvalidSMTPHeader   = fault.New("invalid SMTP header value")
	ErrSMTPAuthRequiresTLS = fault.New("SMTP authentication requires TLS; set SMTP_USE_TLS=true or remove SMTP credentials")
)

const (
	implicitTLSPort      = 465
	defaultSMTPTimeout   = 30 * time.Second
	defaultSMTPBoundary  = "storyden-multipart"
	defaultSMTPPlainType = "text/plain; charset=UTF-8"
	defaultSMTPHTMLType  = "text/html; charset=UTF-8"
)

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
	if cfg.SMTPPort < 1 || cfg.SMTPPort > 65535 {
		return nil, fault.New("SMTP_PORT must be between 1 and 65535 when using SMTP email provider")
	}
	if cfg.SMTPFromAddress == "" {
		return nil, fault.New("SMTP_FROM_ADDRESS must be provided when using SMTP email provider")
	}
	if !cfg.SMTPUseTLS && (cfg.SMTPUsername != "" || cfg.SMTPPassword != "") {
		return nil, ErrSMTPAuthRequiresTLS
	}
	if (cfg.SMTPUsername == "") != (cfg.SMTPPassword == "") {
		return nil, fault.New("SMTP_USERNAME and SMTP_PASSWORD must both be provided together when using SMTP authentication")
	}

	fromAddress, err := normalizeEnvelopeAddress(cfg.SMTPFromAddress)
	if err != nil {
		return nil, fault.Wrap(err)
	}

	s := &SMTP{
		logger:      logger.With(slog.String("mailer", "smtp")),
		host:        cfg.SMTPHost,
		port:        cfg.SMTPPort,
		username:    cfg.SMTPUsername,
		password:    cfg.SMTPPassword,
		fromName:    cfg.SMTPFromName,
		fromAddress: fromAddress,
		useTLS:      cfg.SMTPUseTLS,
	}

	return s, nil
}

func (s *SMTP) Send(
	ctx context.Context,
	msg Message,
) error {
	if err := ctx.Err(); err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	s.logger.Info("sending live email",
		slog.String("name", msg.Name),
		slog.String("subject", msg.Subject),
	)

	from, err := formatAddressHeader(s.fromName, s.fromAddress)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	to, err := formatAddressHeader(msg.Name, msg.Address.Address)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	subject, err := sanitizeHeaderValue(msg.Subject)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	toAddr, err := normalizeEnvelopeAddress(msg.Address.Address)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	messageBytes := buildMessage(from, to, subject, msg.Content)

	addr := smtpServerAddress(s.host, s.port)

	if s.useTLS {
		err = s.sendWithTLS(ctx, addr, messageBytes, toAddr)
		if err != nil {
			return fault.Wrap(err, fctx.With(ctx))
		}
	} else {
		err = s.sendWithoutTLS(ctx, addr, messageBytes, toAddr)
		if err != nil {
			return fault.Wrap(err, fctx.With(ctx))
		}
	}

	return nil
}

func buildMessage(from string, to string, subject string, content Content) []byte {
	headers := []string{
		fmt.Sprintf("From: %s", from),
		fmt.Sprintf("To: %s", to),
		fmt.Sprintf("Subject: %s", subject),
		"MIME-Version: 1.0",
	}

	if content.HTML != "" && content.Plain != "" {
		boundary := fmt.Sprintf("%s-%d", defaultSMTPBoundary, time.Now().UnixNano())
		lines := append(headers,
			fmt.Sprintf("Content-Type: multipart/alternative; boundary=%q", boundary),
			"",
			fmt.Sprintf("--%s", boundary),
			fmt.Sprintf("Content-Type: %s", defaultSMTPPlainType),
			"",
			content.Plain,
			fmt.Sprintf("--%s", boundary),
			fmt.Sprintf("Content-Type: %s", defaultSMTPHTMLType),
			"",
			content.HTML,
			fmt.Sprintf("--%s--", boundary),
			"",
		)
		return []byte(strings.Join(lines, "\r\n"))
	}

	body := content.Plain
	contentType := defaultSMTPPlainType
	if content.HTML != "" {
		body = content.HTML
		contentType = defaultSMTPHTMLType
	}

	lines := append(headers,
		fmt.Sprintf("Content-Type: %s", contentType),
		"",
		body,
	)
	return []byte(strings.Join(lines, "\r\n"))
}

func (s *SMTP) sendWithTLS(ctx context.Context, addr string, messageBytes []byte, toAddr string) error {
	if useImplicitTLS(s.port) {
		return s.sendWithImplicitTLS(ctx, addr, messageBytes, toAddr)
	}
	return s.sendWithSTARTTLS(ctx, addr, messageBytes, toAddr)
}

func useImplicitTLS(port int) bool {
	return port == implicitTLSPort
}

func smtpServerAddress(host string, port int) string {
	return net.JoinHostPort(host, strconv.Itoa(port))
}

func (s *SMTP) sendWithSTARTTLS(ctx context.Context, addr string, messageBytes []byte, toAddr string) error {
	session, err := s.openPlainSession(ctx, addr)
	if err != nil {
		return err
	}
	defer session.close()

	if ok, _ := session.client.Extension("STARTTLS"); !ok {
		return errors.New("SMTP server does not support STARTTLS")
	}

	if err := session.client.StartTLS(s.tlsConfig()); err != nil {
		return fmt.Errorf("failed to upgrade SMTP connection with STARTTLS: %w", err)
	}

	if err := s.authenticate(session.client); err != nil {
		return err
	}

	if err := s.sendMessage(session.client, messageBytes, toAddr); err != nil {
		return err
	}

	if err := session.client.Quit(); err != nil {
		s.logger.Warn("failed to close SMTP connection", slog.String("error", err.Error()))
	}

	return nil
}

func (s *SMTP) sendWithImplicitTLS(ctx context.Context, addr string, messageBytes []byte, toAddr string) error {
	session, err := s.openImplicitTLSSession(ctx, addr)
	if err != nil {
		return err
	}
	defer session.close()

	if err := s.authenticate(session.client); err != nil {
		return err
	}

	if err := s.sendMessage(session.client, messageBytes, toAddr); err != nil {
		return err
	}

	if err := session.client.Quit(); err != nil {
		s.logger.Warn("failed to close SMTP connection", slog.String("error", err.Error()))
	}

	return nil
}

func (s *SMTP) sendWithoutTLS(ctx context.Context, addr string, messageBytes []byte, toAddr string) error {
	if s.username != "" || s.password != "" {
		return ErrSMTPAuthRequiresTLS
	}

	session, err := s.openPlainSession(ctx, addr)
	if err != nil {
		return err
	}
	defer session.close()

	if err := s.sendMessage(session.client, messageBytes, toAddr); err != nil {
		return err
	}

	if err := session.client.Quit(); err != nil {
		s.logger.Warn("failed to close SMTP connection", slog.String("error", err.Error()))
	}

	return nil
}

func (s *SMTP) sendMessage(client *smtp.Client, messageBytes []byte, toAddr string) error {
	if err := client.Mail(s.fromAddress); err != nil {
		return fmt.Errorf("failed to set sender: %w", err)
	}

	if err := client.Rcpt(toAddr); err != nil {
		return fmt.Errorf("failed to set recipient: %w", err)
	}

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

func (s *SMTP) authenticate(client *smtp.Client) error {
	if s.username == "" && s.password == "" {
		return nil
	}

	if ok, _ := client.Extension("AUTH"); !ok {
		return errors.New("SMTP server doesn't support AUTH")
	}

	auth := smtp.PlainAuth("", s.username, s.password, s.host)
	if err := client.Auth(auth); err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	return nil
}

func (s *SMTP) openPlainSession(ctx context.Context, addr string) (*smtpSession, error) {
	conn, err := (&net.Dialer{Timeout: defaultSMTPTimeout}).DialContext(ctx, "tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to SMTP server: %w", err)
	}

	stop := watchConnectionContext(ctx, conn)
	client, err := smtp.NewClient(conn, s.host)
	if err != nil {
		stop()
		_ = conn.Close()
		return nil, fmt.Errorf("failed to create SMTP client: %w", err)
	}

	return &smtpSession{client: client, conn: conn, stopContextWatch: stop}, nil
}

func (s *SMTP) openImplicitTLSSession(ctx context.Context, addr string) (*smtpSession, error) {
	conn, err := (&tls.Dialer{
		NetDialer: &net.Dialer{Timeout: defaultSMTPTimeout},
		Config:    s.tlsConfig(),
	}).DialContext(ctx, "tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to SMTP server with implicit TLS: %w", err)
	}

	stop := watchConnectionContext(ctx, conn)
	client, err := smtp.NewClient(conn, s.host)
	if err != nil {
		stop()
		_ = conn.Close()
		return nil, fmt.Errorf("failed to create SMTP client: %w", err)
	}

	return &smtpSession{client: client, conn: conn, stopContextWatch: stop}, nil
}

func (s *SMTP) tlsConfig() *tls.Config {
	return &tls.Config{
		ServerName: s.host,
		MinVersion: tls.VersionTLS12,
	}
}

func watchConnectionContext(ctx context.Context, conn net.Conn) func() {
	deadline := time.Now().Add(defaultSMTPTimeout)
	if ctxDeadline, ok := ctx.Deadline(); ok && ctxDeadline.Before(deadline) {
		deadline = ctxDeadline
	}
	_ = conn.SetDeadline(deadline)

	done := make(chan struct{})
	go func() {
		select {
		case <-ctx.Done():
			_ = conn.SetDeadline(time.Now())
		case <-done:
		}
	}()

	return func() {
		close(done)
	}
}

type smtpSession struct {
	client           *smtp.Client
	conn             net.Conn
	stopContextWatch func()
}

func (s *smtpSession) close() {
	if s.stopContextWatch != nil {
		s.stopContextWatch()
	}
	_ = s.client.Close()
	_ = s.conn.Close()
}

func sanitizeHeaderValue(value string) (string, error) {
	if strings.ContainsAny(value, "\r\n") {
		return "", errors.Join(ErrInvalidSMTPHeader, fmt.Errorf("header contains line breaks"))
	}

	return value, nil
}

func formatAddressHeader(name, address string) (string, error) {
	sanitizedName, err := sanitizeHeaderValue(name)
	if err != nil {
		return "", err
	}

	sanitizedAddress, err := sanitizeHeaderValue(address)
	if err != nil {
		return "", err
	}

	parsedAddress, err := mail.ParseAddress(sanitizedAddress)
	if err != nil {
		return "", errors.Join(ErrInvalidSMTPHeader, fmt.Errorf("invalid email address: %w", err))
	}

	mailAddress := &mail.Address{Name: sanitizedName, Address: parsedAddress.Address}

	return mailAddress.String(), nil
}

func normalizeEnvelopeAddress(address string) (string, error) {
	sanitizedAddress, err := sanitizeHeaderValue(address)
	if err != nil {
		return "", err
	}

	parsedAddress, err := mail.ParseAddress(sanitizedAddress)
	if err != nil {
		return "", errors.Join(ErrInvalidSMTPHeader, fmt.Errorf("invalid email address: %w", err))
	}

	return parsedAddress.Address, nil
}
