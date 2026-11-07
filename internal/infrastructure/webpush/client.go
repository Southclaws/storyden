package webpush

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/mail"
	"net/url"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/infrastructure/httpsafe"
)

var (
	ErrDisabled = errors.New("web push is disabled")
	ErrExpired  = errors.New("push subscription expired")
)

type Client struct {
	http       *http.Client
	privateKey *ecdsa.PrivateKey
	publicKey  string
	subject    string
}

func New(cfg config.Config) (*Client, error) {
	if !cfg.WebPushEnabled {
		return &Client{}, nil
	}

	private, err := base64.RawURLEncoding.Strict().DecodeString(cfg.WebPushVAPIDPrivateKey)
	if err != nil {
		return nil, errors.New("WEB_PUSH_VAPID_PRIVATE_KEY must be unpadded base64url")
	}
	key, err := ecdsa.ParseRawPrivateKey(elliptic.P256(), private)
	if err != nil {
		return nil, errors.New("WEB_PUSH_VAPID_PRIVATE_KEY must be a valid 32-byte P-256 private key")
	}

	public, err := base64.RawURLEncoding.Strict().DecodeString(cfg.WebPushVAPIDPublicKey)
	if err != nil {
		return nil, errors.New("WEB_PUSH_VAPID_PUBLIC_KEY must be unpadded base64url")
	}

	expected, err := key.PublicKey.Bytes()
	if err != nil {
		return nil, err
	}
	if !bytes.Equal(public, expected) {
		return nil, errors.New("WEB_PUSH_VAPID_PUBLIC_KEY does not match the private key")
	}

	if err := validateSubject(cfg.WebPushVAPIDSubject); err != nil {
		return nil, err
	}

	return &Client{
		http: httpsafe.NewClient(httpsafe.Config{
			Timeout:      10 * time.Second,
			MaxRedirects: -1,
		}),
		privateKey: key,
		publicKey:  cfg.WebPushVAPIDPublicKey,
		subject:    cfg.WebPushVAPIDSubject,
	}, nil
}

func validateSubject(subject string) error {
	u, err := url.Parse(subject)
	if err == nil && u.User == nil && u.Fragment == "" {
		switch u.Scheme {
		case "https":
			if u.Hostname() != "" {
				return nil
			}
		case "mailto":
			if _, err := mail.ParseAddress(u.Opaque); err == nil {
				return nil
			}
		}
	}
	return errors.New("WEB_PUSH_VAPID_SUBJECT must be a mailto: or https: contact URI")
}

func (c *Client) Enabled() bool {
	return c.privateKey != nil
}

func (c *Client) Send(ctx context.Context, subscription Subscription, payload []byte) error {
	if !c.Enabled() {
		return ErrDisabled
	}

	body, err := encrypt(subscription, payload)
	if err != nil {
		return fmt.Errorf("encrypt push payload: %w", err)
	}

	endpoint := subscription.Endpoint
	host := strings.ToLower(endpoint.Hostname())
	if port := endpoint.Port(); port != "" && port != "443" {
		host = net.JoinHostPort(host, port)
	} else if strings.Contains(host, ":") {
		host = "[" + host + "]"
	}

	aud := "https://" + host

	token, err := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"aud": aud,
		"exp": time.Now().Add(12 * time.Hour).Unix(),
		"sub": c.subject,
	}).SignedString(c.privateKey)
	if err != nil {
		return fmt.Errorf("sign push authentication: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint.String(), bytes.NewReader(body))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "vapid t="+token+", k="+c.publicKey)
	req.Header.Set("Content-Encoding", "aes128gcm")
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("TTL", "86400")

	response, err := c.http.Do(req)
	if err != nil {
		var urlError *url.Error
		if errors.As(err, &urlError) {
			err = urlError.Err
		}
		return fmt.Errorf("send push: %w", err)
	}
	defer response.Body.Close()

	_, _ = io.Copy(io.Discard, io.LimitReader(response.Body, 4096))

	switch {
	case response.StatusCode >= 200 && response.StatusCode < 300:
		return nil
	case response.StatusCode == http.StatusNotFound || response.StatusCode == http.StatusGone:
		return ErrExpired
	default:
		return fmt.Errorf("push service returned HTTP %d", response.StatusCode)
	}
}
