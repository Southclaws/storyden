package email_password

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
	"github.com/pkg/errors"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/account_querier"
	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/account/authentication"
	"github.com/Southclaws/storyden/app/resources/account/email"
	"github.com/Southclaws/storyden/app/services/account/register"
	"github.com/Southclaws/storyden/app/services/authentication/email_verify"
	"github.com/Southclaws/storyden/internal/otp"
)

var (
	ErrAccountAlreadyExists = errors.New("account already exists")
	ErrPasswordMismatch     = errors.New("password mismatch")
	ErrNoPassword           = errors.New("password not enabled")
	ErrPasswordTooShort     = errors.New("password too short")
	ErrNotFound             = errors.New("account not found")
)

const (
	id   = "email_password"
	name = "Email and Password"
)

type Provider struct {
	auth         authentication.Repository
	accountQuery *account_querier.Querier
	er           email.EmailRepo
	register     *register.Registrar

	// TODO: Replace with an MQ message and sender job.
	sender email_verify.Verifier
}

func New(
	auth authentication.Repository,
	accountQuery *account_querier.Querier,
	er email.EmailRepo,
	register *register.Registrar,
	sender email_verify.Verifier,
) *Provider {
	return &Provider{auth, accountQuery, er, register, sender}
}

func (p *Provider) Enabled() bool { return true } // TODO: Allow disabling.
func (p *Provider) ID() string    { return id }
func (p *Provider) Name() string  { return name }

func (b *Provider) Register(ctx context.Context, email mail.Address, password string, handle opt.Optional[string], inviteCode opt.Optional[xid.ID]) (*account.Account, error) {
	if len(password) < 8 {
		return nil, fault.Wrap(ErrPasswordTooShort,
			fctx.With(ctx),
			ftag.With(ftag.InvalidArgument),
			fmsg.WithDesc("too short", "Password must be at least 8 characters."))
	}

	identifier := handle.Or(petname.Generate(2, "-"))

	_, exists, err := b.accountQuery.LookupByHandle(ctx, identifier)
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

	account, err := b.register.Create(ctx, identifier, opts...)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to create account"))
	}

	if err := b.addEmailPasswordAuth(ctx, account.ID, email, password); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return account, nil
}

func (b *Provider) Link(_ string) (string, error) {
	// Password provider does not use external links.
	return "", nil
}

func (b *Provider) Login(ctx context.Context, email string, password string) (*account.Account, error) {
	if len(password) < 8 {
		return nil, fault.Wrap(ErrPasswordTooShort,
			fctx.With(ctx),
			ftag.With(ftag.InvalidArgument),
			fmsg.WithDesc("too short", "Password must be at least 8 characters."))
	}

	emailAddress, err := mail.ParseAddress(email)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	acc, exists, err := b.er.LookupAccount(ctx, *emailAddress)
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

	// Get the auth record for this email address
	a, exists, err := b.auth.LookupByHandle(ctx, id, acc.Handle)
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

func (b *Provider) Create(ctx context.Context, aid account.AccountID, email mail.Address, password string) (*account.Account, error) {
	// TODO: Add an email-password auth record for an existing account.
	return nil, nil
}

func (b *Provider) Update(ctx context.Context, aid account.AccountID, email mail.Address, oldpassword, newpassword string) (*account.Account, error) {
	// TODO: Update password for an email-password auth type.
	return nil, nil
}

func (b *Provider) addEmailPasswordAuth(ctx context.Context, accountID account.AccountID, email mail.Address, password string) error {
	code, err := otp.Generate()
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	hashed, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to create secure password hash"))
	}

	// Email auth records do not hold identifiers. However the auth record does
	// use the token field for the password hash.
	identifier := xid.New().String()

	authRecord, err := b.auth.Create(ctx, accountID, id, identifier, string(hashed), nil)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to create account authentication instance"))
	}

	err = b.sender.BeginEmailVerification(ctx, accountID, email, code, opt.New(authRecord.ID))
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}
