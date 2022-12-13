package webauthn

import (
	"context"

	"github.com/go-webauthn/webauthn/webauthn"

	"github.com/Southclaws/storyden/app/resources/account"
)

const (
	id   = "webauthn"
	name = "WebAuthn"
	logo = "https://www.yubico.com/wp-content/uploads/2021/02/illus-yubikey-fingerprint-password-dkteal-r4.svg" // todo; change this image
)

type Provider struct {
	wa *webauthn.WebAuthn
}

func New(wa *webauthn.WebAuthn) (*Provider, error) {
	return &Provider{wa}, nil
}

func (p *Provider) Enabled() bool   { return true }
func (p *Provider) ID() string      { return id }
func (p *Provider) Name() string    { return name }
func (p *Provider) LogoURL() string { return logo }

func (p *Provider) Link() string {
	return ""
}

func (p *Provider) Login(ctx context.Context, state, code string) (*account.Account, error) {
	return nil, nil
}
