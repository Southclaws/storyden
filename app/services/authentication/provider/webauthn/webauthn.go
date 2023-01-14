package webauthn

import (
	"context"
	"encoding/base64"
	"encoding/json"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/go-webauthn/webauthn/webauthn"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/authentication"
	"github.com/Southclaws/storyden/internal/ent"
)

var (
	ErrNoAuthRecord           = fault.New("webauthn does not match account")
	ErrExistsOnAnotherAccount = fault.New("webauthn id already bound to another account")
)

const (
	id   = "webauthn"
	name = "WebAuthn"
	logo = "https://www.yubico.com/wp-content/uploads/2021/02/illus-yubikey-fingerprint-password-dkteal-r4.svg" // todo; change this image
)

type Provider struct {
	auth_repo    authentication.Repository
	account_repo account.Repository

	wa *webauthn.WebAuthn
}

func New(
	auth_repo authentication.Repository,
	account_repo account.Repository,

	wa *webauthn.WebAuthn,
) (*Provider, error) {
	return &Provider{
		auth_repo:    auth_repo,
		account_repo: account_repo,
		wa:           wa,
	}, nil
}

func (p *Provider) Enabled() bool   { return true }
func (p *Provider) ID() string      { return id }
func (p *Provider) Name() string    { return name }
func (p *Provider) LogoURL() string { return logo }

func (p *Provider) Link() string {
	return ""
}

func (p *Provider) Login(ctx context.Context, handle, pubkey string) (*account.Account, error) {
	return nil, nil
}

func (p *Provider) register(ctx context.Context, handle string, credential *webauthn.Credential) (*account.Account, error) {
	// TODO: LookupByHandle returning (account, bool, error) to stop this mess.
	accfound := true
	acc, err := p.account_repo.GetByHandle(ctx, handle)
	if err != nil {
		if ent.IsNotFound(err) {
			accfound = false
		} else {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}
	}

	if accfound {
		return nil, fault.New("requester already has an account")
	}

	acc, err = p.account_repo.Create(ctx, handle)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	encoded, err := json.Marshal(credential)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	_, err = p.auth_repo.Create(ctx, acc.ID, id, base64.RawURLEncoding.EncodeToString(credential.ID), string(encoded), nil)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return acc, nil
}
