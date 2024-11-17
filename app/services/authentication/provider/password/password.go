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
	"go.uber.org/zap"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/account_querier"
	"github.com/Southclaws/storyden/app/resources/account/authentication"
	"github.com/Southclaws/storyden/app/resources/account/email"
	"github.com/Southclaws/storyden/app/resources/settings"
	"github.com/Southclaws/storyden/app/services/account/register"
	"github.com/Southclaws/storyden/app/services/authentication/email_verify"
	"github.com/Southclaws/storyden/app/services/system/instance_info"
	"github.com/Southclaws/storyden/internal/otp"
)

var (
	ErrHandleRegistrationDisabled = fault.New("cannot register while in non-handle authentication mode")
	ErrEmailRegistrationDisabled  = fault.New("cannot register while in non-email authentication mode")
	ErrAccountAlreadyExists       = fault.New("account already exists")
	ErrPasswordMismatch           = fault.New("password mismatch")
	ErrNoPassword                 = fault.New("password not enabled")
	ErrPasswordAlreadySet         = fault.New("password already enabled")
	ErrPasswordTooShort           = fault.New("password too short")
	ErrNotFound                   = fault.New("account not found")
)

var tokenType = authentication.TokenTypePasswordHash

type Provider struct {
	logger       *zap.Logger
	settings     *settings.SettingsRepository
	system       *instance_info.Provider
	auth         authentication.Repository
	accountQuery *account_querier.Querier
	er           email.EmailRepo
	register     *register.Registrar

	// TODO: Replace with an MQ message and sender job.
	sender email_verify.Verifier
}

var service = authentication.ServicePassword

func New(
	logger *zap.Logger,
	settings *settings.SettingsRepository,
	system *instance_info.Provider,
	auth authentication.Repository,
	accountQuery *account_querier.Querier,
	er email.EmailRepo,
	register *register.Registrar,
	sender email_verify.Verifier,
) *Provider {
	return &Provider{
		logger:       logger,
		settings:     settings,
		system:       system,
		auth:         auth,
		accountQuery: accountQuery,
		er:           er,
		register:     register,
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

	return acc, nil
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
	authRecord, err := p.addPasswordAuth(ctx, accountID, password)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to create account authentication instance"))
	}

	code, err := otp.Generate()
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	err = p.sender.BeginEmailVerification(ctx, accountID, email, code, opt.New(authRecord.ID))
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}
