package password

import (
	"context"
	"net/mail"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"
	"github.com/alexedwards/argon2id"
	petname "github.com/dustinkirkland/golang-petname"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/account_writer"
)

func (p *Provider) RegisterWithEmail(ctx context.Context, email mail.Address, password string, handle opt.Optional[string], inviteCode opt.Optional[xid.ID]) (*account.Account, error) {
	enabled, err := p.isEmailAvailable(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}
	if !enabled {
		return nil, ErrEmailRegistrationDisabled
	}

	// TODO: Check if email capability is enabled

	if len(password) < 8 {
		return nil, fault.Wrap(ErrPasswordTooShort,
			fctx.With(ctx),
			ftag.With(ftag.InvalidArgument),
			fmsg.WithDesc("too short", "Password must be at least 8 characters."))
	}

	identifier := handle.Or(petname.Generate(2, "-"))

	_, exists, err := p.accountQuery.LookupByHandle(ctx, identifier)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to get account"))
	}

	if exists {
		return nil, fault.Wrap(ErrAccountAlreadyExists,
			fctx.With(ctx),
			ftag.With(ftag.AlreadyExists),
			fmsg.WithDesc("exists", "The specified handle has already been registered."))
	}

	opts := []account_writer.Option{}
	inviteCode.Call(func(id xid.ID) { opts = append(opts, account_writer.WithInvitedBy(id)) })

	account, err := p.register.Create(ctx, identifier, opts...)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to create account"))
	}

	if err := p.addPasswordAuthWithEmail(ctx, account.ID, email, password); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return account, nil
}

func (p *Provider) LoginWithEmail(ctx context.Context, emailAddress mail.Address, password string) (*account.Account, error) {
	enabled, err := p.isEmailAvailable(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}
	if !enabled {
		return nil, ErrEmailRegistrationDisabled
	}

	if len(password) < 8 {
		return nil, fault.Wrap(ErrPasswordTooShort,
			fctx.With(ctx),
			ftag.With(ftag.InvalidArgument),
			fmsg.WithDesc("too short", "Password must be at least 8 characters."))
	}

	acc, exists, err := p.er.LookupAccount(ctx, emailAddress)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to get account"))
	}

	if !exists {
		return nil, fault.Wrap(ErrNotFound,
			fctx.With(ctx),
			ftag.With(ftag.NotFound),
			fmsg.WithDesc("not found", "No account was found with the provided email address."))
	}

	if err := acc.RejectSuspended(); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	// Get the auth record for this account, it must be a password-based method.
	a, exists, err := p.auth.LookupByEmail(ctx, emailAddress)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if !exists {
		return nil, fault.Wrap(ErrNoPassword,
			fctx.With(ctx),
			ftag.With(ftag.InvalidArgument),
			fmsg.WithDesc("no password", "The specified account does not use email-password authentication. Please try a different method."))
	}

	match, _, err := argon2id.CheckHash(password, a.Token)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to compare secure password hash"))
	}

	if !match {
		return nil, fault.Wrap(ErrPasswordMismatch,
			fctx.With(ctx),
			ftag.With(ftag.Unauthenticated),
			fmsg.WithDesc("mismatch", "The provided password did not match the account."))
	}

	return &a.Account, nil
}
