package register

import (
	"context"
	"log/slog"
	"net/mail"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/opt"
	petname "github.com/dustinkirkland/golang-petname"
	"github.com/samber/lo"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/account/account_querier"
	"github.com/Southclaws/storyden/app/resources/account/account_writer"
	"github.com/Southclaws/storyden/app/resources/account/authentication"
	"github.com/Southclaws/storyden/app/resources/account/email"
	"github.com/Southclaws/storyden/app/resources/message"
	"github.com/Southclaws/storyden/app/services/authentication/email_verify"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/app/services/onboarding"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
	"github.com/Southclaws/storyden/internal/otp"
)

var (
	errEmailAlreadyRegistered = fault.New("email already registered")
	errAccountMismatch        = fault.New("account mismatch")
)

type Registrar struct {
	logger         *slog.Logger
	accountWriter  *account_writer.Writer
	accountQuerier *account_querier.Querier
	emailRepo      *email.Repository
	emailVerify    *email_verify.Verifier
	authRepo       authentication.Repository
	onboarding     onboarding.Service
	bus            *pubsub.Bus
}

func New(
	logger *slog.Logger,
	writer *account_writer.Writer,
	accountQuerier *account_querier.Querier,
	emailRepo *email.Repository,
	emailVerify *email_verify.Verifier,
	authRepo authentication.Repository,
	onboarding onboarding.Service,
	bus *pubsub.Bus,
) *Registrar {
	return &Registrar{
		logger:         logger,
		accountWriter:  writer,
		accountQuerier: accountQuerier,
		emailRepo:      emailRepo,
		emailVerify:    emailVerify,
		authRepo:       authRepo,
		onboarding:     onboarding,
		bus:            bus,
	}
}

func (s *Registrar) Create(ctx context.Context, handle opt.Optional[string], opts ...account_writer.Option) (*account.Account, error) {
	status, err := s.onboarding.GetOnboardingStatus(ctx)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	if status == &onboarding.StatusRequiresFirstAccount {
		// If we're doing first-time-setup then set the first account to admin.
		opts = append(opts, account_writer.WithAdmin(true))
	}

	// If no handle was given, generate one using adjective-animal.
	handleOrGenerated := handle.Or(petname.Generate(2, "-"))

	acc, err := s.accountWriter.Create(ctx, handleOrGenerated, opts...)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	s.bus.Publish(ctx, &message.EventAccountCreated{
		ID: acc.Account.ID,
	})

	return &acc.Account, nil
}

// GetOrCreateViaEmail is intended to be used for just OAuth2 providers. It will
// cover cases for existing auth records, existing emails and ensure accounts
// are linked correctly if necessary and also ensure mismatches are handled.
func (s *Registrar) GetOrCreateViaEmail(
	ctx context.Context,
	service authentication.Service,
	authName string,
	identifier string,
	token string,
	handle string,
	name string,
	email mail.Address,
) (*account.Account, error) {
	// Two key pieces of information here can point to an existing account. The
	// authentication record and the email address. In most cases, either both
	// will exist (a login) or neither will exist (a registration). However, in
	// some cases, a member may already have an account but not have an email.
	// Or somehow, their email is associated with a different account to the one
	// that their auth record points to. This function handles all these cases.

	// A session will be present if the user is attempting to link an account
	// to their Storyden account, rather than registering or logging in.
	session := session.GetOptAccount(ctx)

	authmethod, authMethodExists, err := s.authRepo.LookupByIdentifier(ctx, service, identifier)
	if err != nil {
		return nil, fault.Wrap(err, fmsg.With("failed to lookup existing account"), fctx.With(ctx))
	}

	emailOwner, emailExists, err := s.emailRepo.LookupAccount(ctx, email)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	isVerified := func() bool {
		if !emailExists {
			return false
		}

		current, found := lo.Find(emailOwner.EmailAddresses, func(e *account.EmailAddress) bool {
			return e.Email.Address == email.Address
		})
		if !found {
			return false
		}

		return current.Verified
	}()

	logger := s.logger.With(
		slog.String("handle", handle),
		slog.String("name", name),
		slog.String("email", email.Address),
		slog.Bool("auth_method_exists", authMethodExists),
		slog.Bool("email_exists", emailExists),
		slog.Bool("email_verified", isVerified),
	)

	switch {
	case authMethodExists && emailExists:
		// Member has already registered with this email address and the
		// same email exists and points to an account. Verify those accounts
		// are the same, if so, return the account. If not, bad state, error.

		if authmethod.Account.ID != emailOwner.ID {
			return nil, fault.Wrap(errEmailAlreadyRegistered,
				fctx.With(ctx),
				fmsg.WithDesc("email already in use by another account", "This email address is already in use by another account. Please use a different email address or log in with the existing account."),
			)
		}

		if sessionAccount, ok := session.Get(); ok {
			if sessionAccount.ID != emailOwner.ID {
				return nil, fault.Wrap(errAccountMismatch, fctx.With(ctx))
			}
		}

		logger.Info("get or create: account already exists")

		return &emailOwner.Account, nil

	case authMethodExists && !emailExists:
		// Member has already registered with this authentication method but
		// the email address associated with the authentication method has
		// not yet been recorded. This may happen if certain auth providers
		// are enabled at a later date, such as with OAuth providers.
		// Link the email address to the account and send a verification.

		err = s.linkAndVerifyEmail(ctx, emailOwner.ID, email)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}

		logger.Info("get or create: auth method exists, email not recorded, linking new email to existing account")

		return &emailOwner.Account, nil

	case !authMethodExists && emailExists:
		// Member has already registered with this email address, perhaps on
		// a different auth provider. We know it's the same email assuming
		// the caller has verified this via OAuth or similar. The email may
		// have also been added manually via a newsletter list. Link them.

		_, err = s.authRepo.Create(ctx, emailOwner.ID, service, authentication.TokenTypeOAuth, identifier, token, nil, authentication.WithName(authName))
		if err != nil {
			return nil, fault.Wrap(err, fmsg.With("failed to create new auth method for existing email"), fctx.With(ctx))
		}

		logger.Info("get or create: no auth record, email already points to existing account, linking new auth method to existing account")

		return &emailOwner.Account, nil

	case !authMethodExists && !emailExists:
		// Nothing exists for this member yet, create a new account.

		newAccount, err := s.CreateWithHandle(ctx, service, authName, identifier, token, name, handle)
		if err != nil {
			return nil, fault.Wrap(err, fmsg.With("failed to create new account"), fctx.With(ctx))
		}

		if !isVerified {
			err = s.linkAndVerifyEmail(ctx, newAccount.ID, email)
			if err != nil {
				return nil, fault.Wrap(err, fctx.With(ctx))
			}
		}

		logger.Info("get or create: no auth record, no email record, creating new account and verifying email")

		return newAccount, nil

	default:
		// switch block covers all cases.
		panic("unreachable")
	}
}

func (s *Registrar) GetOrCreateViaHandle(
	ctx context.Context,
	service authentication.Service,
	authName string,
	identifier string,
	token string,
	handle string,
	name string,
) (*account.Account, error) {
	// A session will be present if the user is attempting to link an account
	// to their Storyden account, rather than registering or logging in.
	session := session.GetOptAccount(ctx)

	authmethod, authMethodExists, err := s.authRepo.LookupByIdentifier(ctx, service, identifier)
	if err != nil {
		return nil, fault.Wrap(err, fmsg.With("failed to lookup existing account"), fctx.With(ctx))
	}

	handleOwner, handleExists, err := s.accountQuerier.LookupByHandle(ctx, handle)
	if err != nil {
		return nil, fault.Wrap(err, fctx.With(ctx))
	}

	logger := s.logger.With(
		slog.String("handle", handle),
		slog.String("name", name),
		slog.Bool("auth_method_exists", authMethodExists),
		slog.Bool("handle_exists", handleExists),
	)

	switch {
	case authMethodExists && handleExists:
		// Member has already registered with this email address and the
		// same email exists and points to an account. Verify those accounts
		// are the same, if so, return the account. If not, bad state, error.

		if authmethod.Account.ID != handleOwner.ID {
			logger.Info("get or create: different account already exists with handle, generating new handle for new account")

			return s.CreateWithRandomHandle(ctx, service, authName, identifier, token, name)
		}

		if sessionAccount, ok := session.Get(); ok {
			if sessionAccount.ID != handleOwner.ID {
				return nil, fault.Wrap(errAccountMismatch, fctx.With(ctx))
			}
		}

		logger.Info("get or create: account already exists with handle")

		return &handleOwner.Account, nil

	case authMethodExists && !handleExists:
		// Member has already registered with this authentication method but
		// has changed their name, this is fine. Return the existing account.

		logger.Info("get or create: auth method exists, but account handle changed")

		return &handleOwner.Account, nil

	case !authMethodExists && handleExists:
		// Member has already registered with this handle, we can only verify
		// ownership if there's an existing session, if not, create new account.

		if sessionAccount, ok := session.Get(); ok {
			_, err = s.authRepo.Create(ctx, sessionAccount.ID, service, authentication.TokenTypeOAuth, identifier, token, nil, authentication.WithName(authName))
			if err != nil {
				return nil, fault.Wrap(err, fmsg.With("failed to create new auth method for existing already logged-in account with same handle"), fctx.With(ctx))
			}

			return &sessionAccount, nil
		}

		logger.Info("get or create: no auth record, handle already points to existing account, creating new account with random handle")

		return s.CreateWithRandomHandle(ctx, service, authName, identifier, token, name)

	case !authMethodExists && !handleExists:
		// Nothing exists for this member yet, create a new account.

		logger.Info("get or create: no auth record, no email record, creating new account and verifying email")

		if sessionAccount, ok := session.Get(); ok {
			// They're already logged in but trying to link a new account. The
			// handle has not been registered yet so we can assume the handle
			// provided by the OAuth provider is different to the one they are
			// already using on Storyden. Link the auth method and return acc.
			_, err = s.authRepo.Create(ctx, sessionAccount.ID, service, authentication.TokenTypeOAuth, identifier, token, nil, authentication.WithName(authName))
			if err != nil {
				return nil, fault.Wrap(err, fmsg.With("failed to create new auth method for existing already logged-in account"), fctx.With(ctx))
			}

			return &sessionAccount, nil
		}

		return s.CreateWithHandle(ctx, service, authName, identifier, token, name, handle)

	default:
		// switch block covers all cases.
		panic("unreachable")
	}
}

func (s *Registrar) CreateWithRandomHandle(
	ctx context.Context,
	service authentication.Service,
	authName string,
	identifier string,
	token string,
	name string,
) (*account.Account, error) {
	randomHandle := petname.Generate(3, "-")

	newAccount, err := s.Create(ctx, opt.New(randomHandle),
		account_writer.WithName(name))
	if err != nil {
		return nil, fault.Wrap(err, fmsg.With("failed to create new account"), fctx.With(ctx))
	}

	_, err = s.authRepo.Create(ctx, newAccount.ID, service, authentication.TokenTypeOAuth, identifier, token, nil, authentication.WithName(authName))
	if err != nil {
		return nil, fault.Wrap(err, fmsg.With("failed to create new auth method for brand new account with random handle"), fctx.With(ctx))
	}

	return newAccount, nil
}

func (s *Registrar) CreateWithHandle(
	ctx context.Context,
	service authentication.Service,
	authName string,
	identifier string,
	token string,
	name string,
	handle string,
) (*account.Account, error) {
	newAccount, err := s.Create(ctx, opt.New(handle),
		account_writer.WithName(name))
	if err != nil {
		return nil, fault.Wrap(err, fmsg.With("failed to create new account"), fctx.With(ctx))
	}

	_, err = s.authRepo.Create(ctx, newAccount.ID, service, authentication.TokenTypeOAuth, identifier, token, nil, authentication.WithName(authName))
	if err != nil {
		return nil, fault.Wrap(err, fmsg.With("failed to create new auth method for brand new account"), fctx.With(ctx))
	}

	return newAccount, nil
}

func (s *Registrar) linkAndVerifyEmail(ctx context.Context, accID account.AccountID, email mail.Address) error {
	code, err := otp.Generate()
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	_, err = s.emailVerify.BeginEmailVerification(ctx, accID, email, code)
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	return nil
}
