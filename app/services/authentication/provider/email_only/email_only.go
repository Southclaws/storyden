package email_only

import (
	"context"
	"net/mail"
	"strings"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"
	petname "github.com/dustinkirkland/golang-petname"
	"github.com/pkg/errors"
	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/authentication"
	"github.com/Southclaws/storyden/app/resources/account/email"
	"github.com/Southclaws/storyden/app/services/authentication/email_verifier"
	"github.com/Southclaws/storyden/app/services/authentication/register"
	"github.com/Southclaws/storyden/internal/otp"
)

var (
	ErrAccountAlreadyExists = errors.New("account already exists")
	ErrNotFound             = errors.New("account not found")
	ErrAuthMethodNotFound   = errors.New("authentication method not found")
	ErrTokenMismatch        = fault.New("token mismatch", ftag.With(ftag.Unauthenticated))
)

const (
	id   = "email_only"
	name = "Email"
)

type Provider struct {
	auth     authentication.Repository
	ar       account.Repository
	register register.Service
	er       email.EmailRepo

	// TODO: Replace with an MQ message and sender job.
	sender email_verifier.VerificationMailSender
}

func New(
	auth authentication.Repository,
	ar account.Repository,
	register register.Service,
	er email.EmailRepo,
	sender email_verifier.VerificationMailSender,
) *Provider {
	return &Provider{auth, ar, register, er, sender}
}

func (p *Provider) Enabled() bool { return true } // TODO: Allow disabling.
func (p *Provider) ID() string    { return id }
func (p *Provider) Name() string  { return name }

func (b *Provider) Register(ctx context.Context, email mail.Address, handle opt.Optional[string]) (*account.Account, error) {
	code, err := otp.Generate()
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	// NOTE: Do we want to duplicate the email here?
	_, exists, err := b.auth.LookupByIdentifier(ctx, id, email.Address)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to get account"))
	}

	if exists {
		return nil, fault.Wrap(ErrAccountAlreadyExists,
			fctx.With(ctx),
			ftag.With(ftag.AlreadyExists),
			fmsg.WithDesc("exists", "The specified email address has already been registered."))
	}

	// For direct email registration, we generate a random handle for the new
	// account which is a simple placeholder that the owner can overwrite later.
	identifier := handle.Or(petname.Generate(2, "-"))

	account, err := b.register.Create(ctx, identifier)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to create account"))
	}

	if err := b.addEmailAuth(ctx, account.ID, email, code); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	err = b.sender.SendVerificationEmail(ctx, email, code)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return account, nil
}

func (b *Provider) Link(_ string) (string, error) {
	// Password provider does not use external links.
	return "", nil
}

func (b *Provider) Login(ctx context.Context, accountID string, code string) (*account.Account, error) {
	aid, err := xid.FromString(accountID)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	acc, err := b.ar.GetByID(ctx, account.AccountID(aid))
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	a, exists, err := b.auth.LookupByHandle(ctx, id, acc.Handle)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}
	if !exists {
		return nil, fault.Wrap(ErrAuthMethodNotFound, fctx.With(ctx))
	}
	authRecordEmailAddress, err := mail.ParseAddress(a.Identifier)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if err := acc.RejectSuspended(); err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if strings.TrimSpace(a.Token) != strings.TrimSpace(code) {
		return nil, fault.Wrap(ErrTokenMismatch, fctx.With(ctx))
	}

	err = b.er.Verify(ctx, acc.ID, *authRecordEmailAddress)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	return acc, nil
}

func (b *Provider) addEmailAuth(ctx context.Context, accountID account.AccountID, email mail.Address, token string) error {
	em, err := b.er.Add(ctx, accountID, email, true)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	identifier := em.Email.Address

	_, err = b.auth.Create(ctx, accountID, id, identifier, token, nil)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to create account authentication instance"))
	}

	return nil
}
