package password

import (
	"context"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"
	"github.com/alexedwards/argon2id"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/account/authentication"
	"github.com/Southclaws/storyden/app/services/authentication/provider"
)

func (p *Provider) RegisterWithHandle(ctx context.Context, handle string, password string, inviteCode opt.Optional[xid.ID]) (*account.Account, error) {
	if err := provider.CheckMode(ctx, p.logger, p.settings, authentication.ModeHandle); err != nil {
		return nil, err
	}

	if len(password) < 8 {
		return nil, fault.Wrap(ErrPasswordTooShort,
			fctx.With(ctx),
			ftag.With(ftag.InvalidArgument),
			fmsg.WithDesc("too short", "Password must be at least 8 characters."))
	}

	_, exists, err := p.accountQuery.LookupByHandle(ctx, handle)
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

	account, err := p.register.Create(ctx, opt.New(handle), opts...)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to create account"))
	}

	if _, err := p.addPasswordAuth(ctx, account.ID, password); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return account, nil
}

func (b *Provider) LoginWithHandle(ctx context.Context, handle string, password string) (*account.Account, error) {
	if len(password) < 8 {
		return nil, fault.Wrap(ErrPasswordTooShort,
			fctx.With(ctx),
			ftag.With(ftag.InvalidArgument),
			fmsg.WithDesc("too short", "Password must be at least 8 characters."))
	}

	acc, exists, err := b.accountQuery.LookupByHandle(ctx, handle)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to get account"))
	}

	if !exists {
		return nil, fault.Wrap(ErrNotFound,
			fctx.With(ctx),
			ftag.With(ftag.NotFound),
			fmsg.WithDesc("not found", "No account was found with the provided handle."))
	}

	a, exists, err := b.auth.LookupByTokenType(ctx, acc.ID, tokenType, acc.ID.String())
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if !exists {
		return nil, fault.Wrap(ErrNoPassword,
			fctx.With(ctx),
			ftag.With(ftag.InvalidArgument),
			fmsg.WithDesc("no password", "The specified account does not use password authentication. Please try a different method."))
	}

	if err := a.Account.RejectSuspended(); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
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
