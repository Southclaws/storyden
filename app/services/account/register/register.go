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
	"github.com/Southclaws/storyden/app/resources/mark"
	"github.com/Southclaws/storyden/app/services/authentication/email_verify"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/app/services/onboarding"
	"github.com/Southclaws/storyden/internal/infrastructure/pubsub"
	"github.com/Southclaws/storyden/internal/otp"
	"github.com/Southclaws/storyden/lib/plugin/rpc"
)

var (
	errEmailAlreadyRegistered  = fault.New("email already registered")
	errAccountMismatch         = fault.New("account mismatch")
	errEmailNotVerified        = fault.New("email not verified")
	errAuthMethodAlreadyLinked = fault.New("authentication method already linked to another account")
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
		return nil, fault.Wrap(err,
			fctx.With(ctx),
			fmsg.WithDesc("failed to check onboarding status", "Unable to verify system setup. Please try again or contact site administration."))
	}

	if status == &onboarding.StatusRequiresFirstAccount {
		// If we're doing first-time-setup then set the first account to admin.
		opts = append(opts, account_writer.WithAdmin(true))
	}

	// If no handle was given, generate one using adjective-animal.
	handleOrGenerated := handle.Or(petname.Generate(2, "-"))

	if err := account.ValidateHandle(ctx, handleOrGenerated); err != nil {
		return nil, err
	}

	acc, err := s.accountWriter.Create(ctx, handleOrGenerated, opts...)
	if err != nil {
		return nil, fault.Wrap(err,
			fctx.With(ctx),
			fmsg.WithDesc("failed to create account", "Unable to create your account."))
	}

	s.bus.Publish(ctx, &rpc.EventAccountCreated{
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

	// For logged-in accounts, ensure the auth method being used is not already
	// linked to a different account.
	if sessionAccount, ok := session.Get(); ok && authMethodExists {
		if sessionAccount.ID != authmethod.Account.ID {
			return nil, fault.Wrap(errAccountMismatch,
				fctx.With(ctx),
				fmsg.WithDesc("account mismatch", "This authentication method is linked to a different account."))
		}
	}

	emailOwner, emailExists, err := s.emailRepo.LookupAccount(ctx, email)
	if err != nil {
		return nil, fault.Wrap(err,
			fctx.With(ctx),
			fmsg.WithDesc("failed to lookup email address", "Unable to check if this email is already registered. Please try again."))
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

	// Normalize handle to slug format. We don't error here because the handle
	// is provided by an external provider, so it's not necessarily always in
	// the end-user's control. Erroring here would result in a dead-end UX flow.
	handle = mark.Slugify(handle)

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
				fmsg.WithDesc("email already in use by another account", "Unable to complete sign-in. Please contact support if this issue persists."),
			)
		}

		if sessionAccount, ok := session.Get(); ok {
			if sessionAccount.ID != emailOwner.ID {
				return nil, fault.Wrap(errAccountMismatch,
					fctx.With(ctx),
					fmsg.WithDesc("account mismatch", "This authentication method is linked to a different account."))
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

		err = s.linkAndVerifyEmail(ctx, authmethod.Account.ID, email)
		if err != nil {
			return nil, fault.Wrap(err,
				fctx.With(ctx),
				fmsg.WithDesc("failed to link and verify email", "Unable to link email address to your account. Please try again."))
		}

		logger.Info("get or create: auth method exists, email not recorded, linking new email to existing account")

		return &authmethod.Account, nil

	case !authMethodExists && emailExists:
		// Member has already registered with this email address, perhaps on
		// a different auth provider. We know it's the same email assuming
		// the caller has verified this via OAuth or similar. The email may
		// have also been added manually via a newsletter list. Link them.

		if !isVerified {
			return nil, fault.Wrap(errEmailNotVerified,
				fctx.With(ctx),
				fmsg.WithDesc("email not verified", "Unable to complete sign-in. Please contact support if this issue persists."),
			)
		}

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
				return nil, fault.Wrap(err,
					fctx.With(ctx),
					fmsg.WithDesc("failed to link and verify email", "Unable to send verification email. Please try again or contact site administration."))
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

	if sessionAccount, ok := session.Get(); ok && authMethodExists {
		if sessionAccount.ID != authmethod.Account.ID {
			return nil, fault.Wrap(errAccountMismatch,
				fctx.With(ctx),
				fmsg.WithDesc("account mismatch", "This authentication method is linked to a different account."))
		}
	}

	handleOwner, handleExists, err := s.accountQuerier.LookupByHandle(ctx, handle)
	if err != nil {
		return nil, fault.Wrap(err,
			fctx.With(ctx),
			fmsg.WithDesc("failed to lookup username", "Unable to check if this username is already registered. Please try again."))
	}

	logger := s.logger.With(
		slog.String("handle", handle),
		slog.String("name", name),
		slog.Bool("auth_method_exists", authMethodExists),
		slog.Bool("handle_exists", handleExists),
	)

	switch {
	case authMethodExists && handleExists:
		// Member has already registered with this handle and the same handle
		// exists and points to an account. Verify those accounts are the same.

		if authmethod.Account.ID != handleOwner.ID {
			// If the authentication method looked up via the OAuth provider
			// identifier exists, we have an account and it's verified to be
			// owned by someone due to the OAuth process. The handle passed in
			// is purely for registration purposes fetched from the provider,
			// but it's common for handles to change on both Storyden itself and
			// OAuth services, so we don't use the handle for matching only the
			// OAuth identifier. If they don't match, no problem, just return
			// the account linked to the auth method.

			logger.Info("get or create: different account already exists with handle, returning existing auth method linked account")

			return &authmethod.Account, nil
		}

		logger.Info("get or create: account already exists with handle")

		return &handleOwner.Account, nil

	case authMethodExists && !handleExists:
		// Member has already registered with this authentication method but
		// has changed their name, this is fine. Return the existing account.

		logger.Info("get or create: auth method exists, but account handle changed")

		return &authmethod.Account, nil

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
	_, authMethodExists, err := s.authRepo.LookupByIdentifier(ctx, service, identifier)
	if err != nil {
		return nil, fault.Wrap(err, fmsg.With("failed to lookup existing account"), fctx.With(ctx))
	}

	if authMethodExists {
		return nil, fault.Wrap(errAuthMethodAlreadyLinked,
			fctx.With(ctx),
			fmsg.WithDesc("authMethodExists",
				"This authentication provider has already been linked to another account."),
		)
	}

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
	_, authMethodExists, err := s.authRepo.LookupByIdentifier(ctx, service, identifier)
	if err != nil {
		return nil, fault.Wrap(err, fmsg.With("failed to lookup existing account"), fctx.With(ctx))
	}

	if authMethodExists {
		return nil, fault.Wrap(errAuthMethodAlreadyLinked,
			fctx.With(ctx),
			fmsg.WithDesc("authMethodExists",
				"This authentication provider has already been linked to another account."),
		)
	}

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
		return fault.Wrap(err,
			fctx.With(ctx),
			fmsg.WithDesc("failed to generate verification code", "Unable to create verification code. Please try again."))
	}

	_, err = s.emailVerify.BeginEmailVerification(ctx, accID, email, code)
	if err != nil {
		return fault.Wrap(err,
			fctx.With(ctx),
			fmsg.WithDesc("failed to begin email verification", "Unable to send verification email. Please try again or contact site administration."))
	}

	return nil
}
