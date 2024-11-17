package email_only

import (
	"context"
	"net/mail"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"
	petname "github.com/dustinkirkland/golang-petname"
	"github.com/rs/xid"
	"go.uber.org/zap"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/account/authentication"
	"github.com/Southclaws/storyden/app/resources/account/email"
	"github.com/Southclaws/storyden/app/resources/settings"
	"github.com/Southclaws/storyden/app/services/account/register"
	"github.com/Southclaws/storyden/app/services/authentication/email_verify"
	"github.com/Southclaws/storyden/app/services/authentication/provider"
	"github.com/Southclaws/storyden/internal/otp"
)

var (
	ErrEmailRegistrationDisabled = fault.New("cannot register while in non-email authentication mode")
	ErrAccountAlreadyExists      = fault.New("account already exists")
	ErrEmailNotFound             = fault.New("email address not found")
	ErrAccountMismatch           = fault.New("account mismatch")
)

var (
	requiredMode = authentication.ModeEmail
	service      = authentication.ServiceEmailVerify
	tokenType    = authentication.TokenTypeNone
)

type Provider struct {
	logger   *zap.Logger
	settings *settings.SettingsRepository
	auth     authentication.Repository
	register *register.Registrar
	er       email.EmailRepo

	// TODO: Replace with an MQ message and sender job.
	sender email_verify.Verifier
}

func New(
	logger *zap.Logger,
	settings *settings.SettingsRepository,
	auth authentication.Repository,

	register *register.Registrar,
	er email.EmailRepo,
	sender email_verify.Verifier,
) *Provider {
	return &Provider{
		logger:   logger,
		settings: settings,
		auth:     auth,
		register: register,
		er:       er,
		sender:   sender,
	}
}

func (p *Provider) Service() authentication.Service { return service }
func (p *Provider) Token() authentication.TokenType { return tokenType }

func (p *Provider) Enabled(ctx context.Context) (bool, error) {
	settings, err := p.settings.Get(ctx)
	if err != nil {
		return false, fault.Wrap(err, fctx.With(ctx))
	}

	return settings.AuthenticationMode.Or(authentication.ModeHandle) == requiredMode, nil
}

func (p *Provider) Register(ctx context.Context, email mail.Address, handle opt.Optional[string], inviteCode opt.Optional[xid.ID]) (*account.Account, error) {
	if err := provider.CheckMode(ctx, p.logger, p.settings, requiredMode); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	_, exists, err := p.er.LookupAccount(ctx, email)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}
	if exists {
		// If they've already registered, resend the code.
		// TODO: Put this on a queue and ensure there's a sufficient rate limit.
		err := p.sender.ResendVerification(ctx, email)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}

		return nil, fault.Wrap(ErrAccountAlreadyExists,
			fctx.With(ctx),
			ftag.With(ftag.AlreadyExists),
			fmsg.WithDesc("exists", "The specified email address has already been registered."))
	}

	// For direct email registration, we generate a random handle for the new
	// account which is a simple placeholder that the owner can overwrite later.
	identifier := handle.Or(petname.Generate(2, "-"))

	opts := []account_writer.Option{}
	inviteCode.Call(func(id xid.ID) { opts = append(opts, account_writer.WithInvitedBy(id)) })

	account, err := p.register.Create(ctx, identifier, opts...)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to create account"))
	}

	if err := p.addEmailAuth(ctx, account.ID, email); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return account, nil
}

func (p *Provider) Login(ctx context.Context, email mail.Address) error {
	if err := provider.CheckMode(ctx, p.logger, p.settings, requiredMode); err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	_, exists, err := p.auth.LookupByEmail(ctx, email)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}
	if !exists {
		return fault.Wrap(ErrEmailNotFound,
			fctx.With(ctx),
			ftag.With(ftag.NotFound),
			fmsg.WithDesc("not found", "The specified email address is not associated with an account."))
	}

	err = p.sender.ResendVerification(ctx, email)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}

func (p *Provider) addEmailAuth(ctx context.Context, accountID account.AccountID, email mail.Address) error {
	code, err := otp.Generate()
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	// Email verification authentication does not use any form of token, however
	// there needs to be some value set so generate a random ID for each record.
	identifier := ""
	token := xid.New().String()

	authRecord, err := p.auth.Create(ctx, accountID, service, authentication.TokenTypeNone, identifier, token, nil)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to create account authentication instance"))
	}

	err = p.sender.BeginEmailVerification(ctx, accountID, email, code, opt.New(authRecord.ID))
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}
