package email_only

import (
	"context"
	"log/slog"
	"net/mail"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/account_querier"
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
	logger       *slog.Logger
	settings     *settings.SettingsRepository
	accountQuery *account_querier.Querier
	auth         authentication.Repository
	register     *register.Registrar
	er           *email.Repository

	// TODO: Replace with an MQ message and sender job.
	sender *email_verify.Verifier
}

func New(
	logger *slog.Logger,
	settings *settings.SettingsRepository,
	accountQuery *account_querier.Querier,
	auth authentication.Repository,
	register *register.Registrar,
	er *email.Repository,
	sender *email_verify.Verifier,
) *Provider {
	return &Provider{
		logger:       logger,
		settings:     settings,
		accountQuery: accountQuery,
		auth:         auth,
		register:     register,
		er:           er,
		sender:       sender,
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

	if h, ok := handle.Get(); ok {
		_, exists, err := p.accountQuery.LookupByHandle(ctx, h)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}

		if exists {
			return nil, fault.Wrap(ErrAccountAlreadyExists,
				fctx.With(ctx),
				ftag.With(ftag.AlreadyExists),
				fmsg.WithDesc("exists", "The specified handle has already been registered."))
		}
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

	opts := []account_writer.Option{}
	inviteCode.Call(func(id xid.ID) { opts = append(opts, account_writer.WithInvitedBy(id)) })

	account, err := p.register.Create(ctx, handle, opts...)
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

	acc, exists, err := p.er.LookupAccount(ctx, email)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}
	if !exists {
		return fault.Wrap(ErrEmailNotFound,
			fctx.With(ctx),
			ftag.With(ftag.NotFound),
			fmsg.WithDesc("not found", "The specified email address is not associated with an account."))
	}

	_, exists, err = p.auth.LookupByTokenType(ctx, acc.ID, tokenType, acc.ID.String())
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
	token := xid.New().String()

	_, err = p.auth.Create(ctx, accountID, service, authentication.TokenTypeNone, accountID.String(), token, nil)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to create account authentication instance"))
	}

	_, err = p.sender.BeginEmailVerification(ctx, accountID, email, code)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}
