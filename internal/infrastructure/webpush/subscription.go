package webpush

import (
	"crypto/ecdh"
	"encoding/base64"
	"errors"
	"net/url"

	"github.com/Southclaws/storyden/internal/infrastructure/httpsafe"
)

type SubscriptionInput struct {
	Endpoint string
	P256DH   string
	Auth     string
}

type Subscription struct {
	Endpoint url.URL
	P256DH   *ecdh.PublicKey
	Auth     []byte
}

func ParseSubscription(input SubscriptionInput) (*Subscription, error) {
	endpoint, err := ParseEndpoint(input.Endpoint)
	if err != nil {
		return nil, err
	}

	public, err := base64.RawURLEncoding.Strict().DecodeString(input.P256DH)
	if err != nil {
		return nil, errors.New("p256dh must be unpadded base64url")
	}

	key, err := ecdh.P256().NewPublicKey(public)
	if err != nil {
		return nil, errors.New("p256dh must be a valid uncompressed P-256 public key")
	}

	auth, err := base64.RawURLEncoding.Strict().DecodeString(input.Auth)
	if err != nil || len(auth) != 16 {
		return nil, errors.New("auth must be an unpadded base64url-encoded 16-byte secret")
	}

	return &Subscription{
		Endpoint: *endpoint,
		P256DH:   key,
		Auth:     auth,
	}, nil
}

func ParseEndpoint(endpoint string) (*url.URL, error) {
	if len(endpoint) > 4096 {
		return nil, errors.New("push endpoint must not exceed 4096 bytes")
	}

	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, errors.New("push endpoint must be a valid URL")
	}

	if u.Scheme != "https" || u.Hostname() == "" || u.User != nil || u.Fragment != "" {
		return nil, errors.New("push endpoint must be an HTTPS URL without credentials or a fragment")
	}

	if httpsafe.IsDisallowedHost(u.Hostname()) {
		return nil, errors.New("push endpoint must use a public host")
	}

	return u, nil
}
