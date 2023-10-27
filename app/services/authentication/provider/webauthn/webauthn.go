package webauthn

import (
	"context"
	"encoding/base64"
	"encoding/json"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/fault/ftag"
	"github.com/go-webauthn/webauthn/webauthn"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/authentication"
	"github.com/Southclaws/storyden/app/services/authentication/register"
	"github.com/Southclaws/storyden/internal/ent"
)

var (
	ErrNoAuthRecord           = fault.New("webauthn does not match account")
	ErrExistsOnAnotherAccount = fault.New("webauthn id already bound to another account")
	ErrAccountExists          = fault.New("requester already has an account")
)

const (
	id   = "webauthn"
	name = "WebAuthn"
)

type Provider struct {
	auth_repo    authentication.Repository
	account_repo account.Repository
	reg          register.Service

	wa *webauthn.WebAuthn
}

func New(
	auth_repo authentication.Repository,
	account_repo account.Repository,
	reg register.Service,

	wa *webauthn.WebAuthn,
) (*Provider, error) {
	return &Provider{
		auth_repo:    auth_repo,
		account_repo: account_repo,
		reg:          reg,
		wa:           wa,
	}, nil
}

func (p *Provider) Enabled() bool { return true }
func (p *Provider) ID() string    { return id }
func (p *Provider) Name() string  { return name }

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
		return nil, fault.Wrap(ErrAccountExists,
			fctx.With(ctx),
			ftag.With(ftag.AlreadyExists),
			fmsg.WithDesc(
				"already exists",
				"An account with this handle has already been registered without a Passkey (WebAuthn) credential.",
			),
		)
	}

	acc, err = p.reg.Create(ctx, handle)
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
