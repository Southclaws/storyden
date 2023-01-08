package webauthn

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/fault/ftag"
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

func (p *Provider) getOrCreateAccount(ctx context.Context, handle, credentialID, pubkey string) (*account.Account, error) {
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

	// TODO: Don't do this, instead get ALL auth methods of type "webauthn" and
	// iterate through them calling Validate. Each webauthn should be converted
	// into some type that satisfies `webauthn.User` for `p.wa.ValidateLogin()`.
	authrecord, authfound, err := p.auth_repo.LookupByIdentifier(ctx, authentication.Service(id), pubkey)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	// if we found an account with the given handle
	// and
	// we found an authentication record with the public key
	if accfound && authfound {
		// and they are the same account
		if acc.ID == authrecord.Account.ID {
			// log the user in
			return acc, nil
		} else {
			// they are for different accounts...
			// TODO: informative error
			return nil, fault.New("requester already has an account")
		}
	}

	// if we found an account, but no authentication record
	// this requester does not have access
	if accfound && !authfound {
		return nil, fault.Wrap(ErrNoAuthRecord,
			fctx.With(ctx),
			ftag.With(ftag.Unauthenticated),
			fmsg.WithDesc("account already exists without webauthn",
				"This handle has already been registered without a Passkey. If this is your account, you need to sign in to the account via another method and add your Passkey via settings."),
		)
	}

	if !accfound && authfound {
		// the handle is unused but this webauthn ID is used on another account
		return nil, fault.Wrap(ErrExistsOnAnotherAccount, fctx.With(ctx), ftag.With(ftag.AlreadyExists))
	}

	// no account or auth record, create a new account.

	acc, err = p.account_repo.Create(ctx, handle)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	_, err = p.auth_repo.Create(ctx, acc.ID, id, pubkey, credentialID, nil)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return acc, nil
}
