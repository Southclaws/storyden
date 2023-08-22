package password

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/fault/ftag"
	"github.com/alexedwards/argon2id"
	"github.com/pkg/errors"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/authentication"
	"github.com/Southclaws/storyden/app/services/authentication/register"
)

var (
	ErrAccountAlreadyExists = errors.New("account already exists")
	ErrPasswordMismatch     = errors.New("password mismatch")
	ErrNotFound             = errors.New("account not found")
)

const (
	id   = "password"
	name = "Password"
	logo = "" // TODO: a basic logo symbol for password based auth.
)

type Provider struct {
	auth     authentication.Repository
	register register.Service
}

func New(auth authentication.Repository, register register.Service) *Provider {
	return &Provider{auth, register}
}

func (p *Provider) Enabled() bool   { return true } // TODO: Allow disabling.
func (p *Provider) ID() string      { return id }
func (p *Provider) Name() string    { return name }
func (p *Provider) LogoURL() string { return logo }

func (b *Provider) Register(ctx context.Context, identifier string, password string) (*account.Account, error) {
	_, exists, err := b.auth.LookupByIdentifier(ctx, id, identifier)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to get account"))
	}

	if exists {
		return nil, fault.Wrap(ErrAccountAlreadyExists, fctx.With(ctx), ftag.With(ftag.AlreadyExists))
	}

	account, err := b.register.Create(ctx, identifier)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to create account"))
	}

	hashed, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to create secure password hash"))
	}

	_, err = b.auth.Create(ctx, account.ID, id, identifier, string(hashed), nil)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to create account authentication instance"))
	}

	return account, nil
}

func (b *Provider) Link() string {
	// Password provider does not use external links.
	return ""
}

func (b *Provider) Login(ctx context.Context, identifier string, password string) (*account.Account, error) {
	a, exists, err := b.auth.LookupByIdentifier(ctx, id, identifier)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to get account"))
	}

	if !exists {
		return nil, fault.Wrap(ErrNotFound, fctx.With(ctx), ftag.With(ftag.NotFound))
	}

	match, _, err := argon2id.CheckHash(password, a.Token)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to compare secure password hash"))
	}

	if !match {
		return nil, fault.Wrap(ErrPasswordMismatch, fctx.With(ctx), ftag.With(ftag.Unauthenticated))
	}

	return &a.Account, nil
}
