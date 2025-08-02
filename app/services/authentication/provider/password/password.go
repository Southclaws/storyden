package password

import (
	"context"
	"log/slog"
	"net/mail"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/fault/ftag"
	"github.com/alexedwards/argon2id"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/account_querier"
	"github.com/Southclaws/storyden/app/resources/account/authentication"
	"github.com/Southclaws/storyden/app/resources/account/email"
	"github.com/Southclaws/storyden/app/resources/settings"
	"github.com/Southclaws/storyden/app/services/account/register"
	"github.com/Southclaws/storyden/app/services/authentication/email_verify"
	"github.com/Southclaws/storyden/app/services/authentication/provider/password/password_reset"
	"github.com/Southclaws/storyden/app/services/system/instance_info"
	"github.com/Southclaws/storyden/internal/otp"
)

var (
	ErrHandleRegistrationDisabled = fault.New("cannot register while in non-handle authentication mode", ftag.With(ftag.PermissionDenied))
	ErrEmailRegistrationDisabled  = fault.New("cannot register while in non-email authentication mode", ftag.With(ftag.PermissionDenied))
	ErrAccountAlreadyExists       = fault.New("account already exists", ftag.With(ftag.AlreadyExists))
	ErrPasswordMismatch           = fault.New("password mismatch", ftag.With(ftag.PermissionDenied))
	ErrNoPassword                 = fault.New("password not enabled", ftag.With(ftag.InvalidArgument))
	ErrPasswordAlreadySet         = fault.New("password already enabled", ftag.With(ftag.InvalidArgument))
	ErrPasswordTooShort           = fault.New("password too short", ftag.With(ftag.InvalidArgument))
	ErrNotFound                   = fault.New("account not found", ftag.With(ftag.NotFound))
)

var tokenType = authentication.TokenTypePasswordHash

type Provider struct {
	logger       *slog.Logger
	settings     *settings.SettingsRepository
	system       *instance_info.Provider
	auth         authentication.Repository
	accountQuery *account_querier.Querier
	er           *email.Repository
	register     *register.Registrar
	resetter     *password_reset.EmailResetter

	// TODO: Replace with an MQ message and sender job.
	sender *email_verify.Verifier
}

var service = authentication.ServicePassword

func New(
	logger *slog.Logger,
	settings *settings.SettingsRepository,
	system *instance_info.Provider,
	auth authentication.Repository,
	accountQuery *account_querier.Querier,
	er *email.Repository,
	register *register.Registrar,
	resetter *password_reset.EmailResetter,
	sender *email_verify.Verifier,
) *Provider {
	return &Provider{
		logger:       logger,
		settings:     settings,
		system:       system,
		auth:         auth,
		accountQuery: accountQuery,
		er:           er,
		register:     register,
		resetter:     resetter,
		sender:       sender,
	}
}

func (p *Provider) Service() authentication.Service { return service }
func (p *Provider) Token() authentication.TokenType { return tokenType }

func (p *Provider) Enabled(ctx context.Context) (bool, error) {
	// NOTE: password based auth and login is always enabled. This may change.
	return true, nil
}

func (b *Provider) AddPassword(ctx context.Context, aid account.AccountID, password string) (*account.Account, error) {
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

	_, exists, err := b.auth.LookupByTokenType(ctx, acc.ID, tokenType, acc.ID.String())
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if exists {
		return nil, fault.Wrap(ErrPasswordAlreadySet,
			fctx.With(ctx),
			ftag.With(ftag.InvalidArgument),
			fmsg.WithDesc("already has password", "The specified account already uses password authentication."))
	}

	if _, err := b.addPasswordAuth(ctx, acc.ID, password); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return &acc.Account, nil
}

func (b *Provider) UpdatePassword(ctx context.Context, aid account.AccountID, oldpassword, newpassword string) (*account.Account, error) {
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

	auth, exists, err := b.auth.LookupByTokenType(ctx, a.ID, tokenType, a.ID.String())
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

func (p *Provider) ResetPassword(ctx context.Context, token string, newpassword string) (*account.Account, error) {
	accountID, err := p.resetter.Verify(ctx, token)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	auth, exists, err := p.auth.LookupByTokenType(ctx, accountID, tokenType, accountID.String())
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}
	if !exists {
		return nil, fault.Wrap(ErrNoPassword,
			fctx.With(ctx),
			ftag.With(ftag.InvalidArgument),
			fmsg.WithDesc("no password", "The specified account does not use password authentication. Please try a different method."))
	}

	hashed, err := argon2id.CreateHash(newpassword, argon2id.DefaultParams)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to create secure password hash"))
	}

	auth, err = p.auth.Update(ctx, auth.ID, authentication.WithToken(hashed))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return &auth.Account, nil
}

func (p *Provider) isEmailAvailable(ctx context.Context) (bool, error) {
	info, err := p.system.Get(ctx)
	if err != nil {
		return false, fault.Wrap(err, fctx.With(ctx))
	}

	return info.Capabilities.Has(instance_info.CapabilityEmailClient), nil
}

func (b *Provider) addPasswordAuth(ctx context.Context, accountID account.AccountID, password string) (*authentication.Authentication, error) {
	hashed, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to create secure password hash"))
	}

	authRecord, err := b.auth.Create(ctx, accountID, service, tokenType, accountID.String(), string(hashed), nil)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to create account authentication instance"))
	}

	return authRecord, nil
}

func (p *Provider) addPasswordAuthWithEmail(ctx context.Context, accountID account.AccountID, email mail.Address, password string) error {
	code, err := otp.Generate()
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	_, err = p.sender.BeginEmailVerification(ctx, accountID, email, code)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	_, err = p.addPasswordAuth(ctx, accountID, password)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to create account authentication instance"))
	}

	return nil
}
