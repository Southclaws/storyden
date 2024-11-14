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
	"github.com/pkg/errors"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/account/authentication"
	"github.com/Southclaws/storyden/app/resources/account/email"
	"github.com/Southclaws/storyden/app/resources/settings"
	"github.com/Southclaws/storyden/app/services/account/register"
	"github.com/Southclaws/storyden/app/services/authentication/email_verify"
	"github.com/Southclaws/storyden/internal/otp"
)

var ErrAccountAlreadyExists = errors.New("account already exists")

var (
	requiredMode = authentication.ModeEmail
	provider     = authentication.ServiceEmailOnly
)

type Provider struct {
	settings *settings.SettingsRepository
	auth     authentication.Repository
	register *register.Registrar
	er       email.EmailRepo

	// TODO: Replace with an MQ message and sender job.
	sender email_verify.Verifier
}

func New(
	settings *settings.SettingsRepository,
	auth authentication.Repository,

	register *register.Registrar,
	er email.EmailRepo,
	sender email_verify.Verifier,
) *Provider {
	return &Provider{
		settings: settings,
		auth:     auth,
		register: register,
		er:       er,
		sender:   sender,
	}
}

func (p *Provider) Provides() authentication.Service { return provider }

func (p *Provider) Enabled(ctx context.Context) (bool, error) {
	settings, err := p.settings.Get(ctx)
	if err != nil {
		return false, fault.Wrap(err, fctx.With(ctx))
	}

	return settings.AuthenticationMode.Or(authentication.ModeHandle) == requiredMode, nil
}

func (b *Provider) Register(ctx context.Context, email mail.Address, handle opt.Optional[string], inviteCode opt.Optional[xid.ID]) (*account.Account, error) {
	_, exists, err := b.er.LookupAccount(ctx, email)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}
	if exists {
		// If they've already registered, resend the code.
		// TODO: Put this on a queue and ensure there's a sufficient rate limit.
		err := b.sender.ResendVerification(ctx, email)
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

	account, err := b.register.Create(ctx, identifier, opts...)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to create account"))
	}

	if err := b.addEmailAuth(ctx, account.ID, email); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return account, nil
}

func (b *Provider) Link(_ string) (string, error) {
	// Password provider does not use external links.
	return "", nil
}

func (b *Provider) Login(ctx context.Context, accountID string, code string) (*account.Account, error) {
	// NOTE: There's no login method for this, it uses the email.Verify method.
	return nil, nil
}

func (b *Provider) addEmailAuth(ctx context.Context, accountID account.AccountID, email mail.Address) error {
	code, err := otp.Generate()
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	// Email auth records do not hold tokens or identifiers. There's no password
	// hash and the verification code is held in the email resource.
	identifier := xid.New().String()
	token := xid.New().String()

	authRecord, err := b.auth.Create(ctx, accountID, provider, identifier, token, nil)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to create account authentication instance"))
	}

	err = b.sender.BeginEmailVerification(ctx, accountID, email, code, opt.New(authRecord.ID))
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}
