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
	"go.uber.org/zap"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/account_querier"
	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/account/authentication"
	"github.com/Southclaws/storyden/app/resources/settings"
	"github.com/Southclaws/storyden/app/services/account/register"
	"github.com/Southclaws/storyden/app/services/authentication/provider"
)

type UsernamePasswordProvider struct {
	logger       *zap.Logger
	settings     *settings.SettingsRepository
	auth         authentication.Repository
	accountQuery *account_querier.Querier
	register     *register.Registrar
}

var (
	userRequiredMode = authentication.ModeHandle
	userService      = authentication.ServiceUsernamePassword
)

func NewUsernamePasswordProvider(
	logger *zap.Logger,
	settings *settings.SettingsRepository,
	auth authentication.Repository,
	accountQuery *account_querier.Querier,
	register *register.Registrar,
) *UsernamePasswordProvider {
	return &UsernamePasswordProvider{
		logger:       logger,
		settings:     settings,
		auth:         auth,
		accountQuery: accountQuery,
		register:     register,
	}
}

func (p *UsernamePasswordProvider) Service() authentication.Service { return userService }
func (p *UsernamePasswordProvider) Token() authentication.TokenType { return tokenType }

func (p *UsernamePasswordProvider) Enabled(ctx context.Context) (bool, error) {
	// Handle+password registration and login is always enabled.

	return true, nil
}

func (p *UsernamePasswordProvider) Register(ctx context.Context, handle string, password string, inviteCode opt.Optional[xid.ID]) (*account.Account, error) {
	if len(password) < 8 {
		return nil, fault.Wrap(ErrPasswordTooShort,
			fctx.With(ctx),
			ftag.With(ftag.InvalidArgument),
			fmsg.WithDesc("too short", "Password must be at least 8 characters."))
	}

	if err := provider.CheckMode(ctx, p.logger, p.settings, userRequiredMode); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	settings, err := p.settings.Get(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if settings.AuthenticationMode.Or(authentication.ModeHandle) != userRequiredMode {
		return nil, fault.Wrap(ErrHandleRegistrationDisabled, fctx.With(ctx))
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

	account, err := p.register.Create(ctx, handle, opts...)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to create account"))
	}

	if err := p.addPasswordAuth(ctx, account.ID, password); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return account, nil
}

func (b *UsernamePasswordProvider) Login(ctx context.Context, handle string, password string) (*account.Account, error) {
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

	a, exists, err := b.auth.LookupByTokenType(ctx, acc.ID, authentication.TokenTypePassword, authRecordIdentifier)
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

func (b *UsernamePasswordProvider) Create(ctx context.Context, aid account.AccountID, password string) (*account.Account, error) {
	if len(password) < 8 {
		return nil, fault.Wrap(ErrPasswordTooShort,
			fctx.With(ctx),
			ftag.With(ftag.InvalidArgument),
			fmsg.WithDesc("too short", "Password must be at least 8 characters."))
	}

	acc, err := b.accountQuery.GetByID(ctx, aid)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to get account"))
	}

	_, exists, err := b.auth.LookupByTokenType(ctx, acc.ID, authentication.TokenTypePassword, authRecordIdentifier)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if exists {
		return nil, fault.Wrap(ErrPasswordAlreadySet,
			fctx.With(ctx),
			ftag.With(ftag.InvalidArgument),
			fmsg.WithDesc("already has password", "The specified account already uses password authentication."))
	}

	if err := b.addPasswordAuth(ctx, acc.ID, password); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return acc, nil
}

func (b *UsernamePasswordProvider) Update(ctx context.Context, aid account.AccountID, oldpassword, newpassword string) (*account.Account, error) {
	if len(newpassword) < 8 {
		return nil, fault.Wrap(ErrPasswordTooShort,
			fctx.With(ctx),
			ftag.With(ftag.InvalidArgument),
			fmsg.WithDesc("too short", "Password must be at least 8 characters."))
	}

	a, err := b.accountQuery.GetByID(ctx, aid)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to get account"))
	}

	auth, exists, err := b.auth.LookupByTokenType(ctx, a.ID, authentication.TokenTypePassword, authRecordIdentifier)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if !exists {
		return nil, fault.Wrap(ErrNoPassword,
			fctx.With(ctx),
			ftag.With(ftag.InvalidArgument),
			fmsg.WithDesc("no password", "The specified account does not use password authentication. Please try a different method."))
	}

	if err := a.RejectSuspended(); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	match, _, err := argon2id.CheckHash(oldpassword, auth.Token)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to compare secure password hash"))
	}

	if !match {
		return nil, fault.Wrap(ErrPasswordMismatch,
			fctx.With(ctx),
			ftag.With(ftag.Unauthenticated),
			fmsg.WithDesc("mismatch", "The provided password did not match the account."))
	}

	hashed, err := argon2id.CreateHash(newpassword, argon2id.DefaultParams)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to create secure password hash"))
	}

	auth, err = b.auth.Update(ctx, auth.ID, authentication.WithToken(hashed))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return &auth.Account, nil
}

func (b *UsernamePasswordProvider) addPasswordAuth(ctx context.Context, accountID account.AccountID, password string) error {
	hashed, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to create secure password hash"))
	}

	_, err = b.auth.Create(ctx, accountID, userService, authentication.TokenTypePassword, authRecordIdentifier, string(hashed), nil)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to create account authentication instance"))
	}

	return nil
}
