package password

import (
	"github.com/Southclaws/fault"
	"github.com/Southclaws/storyden/app/resources/account/authentication"
)

var (
	ErrHandleRegistrationDisabled = fault.New("cannot register while in non-handle authentication mode")
	ErrAccountAlreadyExists       = fault.New("account already exists")
	ErrPasswordMismatch           = fault.New("password mismatch")
	ErrNoPassword                 = fault.New("password not enabled")
	ErrPasswordAlreadySet         = fault.New("password already enabled")
	ErrPasswordTooShort           = fault.New("password too short")
	ErrNotFound                   = fault.New("account not found")
)

var tokenType = authentication.TokenTypePassword

// We use a constant label for the identifier in order to ensure there's
// only a single auth method associated with an account that may hold a
// hashed password. There's a unique constraint across token type and
// identifier, in this case token type is password and identifier is const.
const authRecordIdentifier = "password"

// type Provider struct {
// 	logger       *zap.Logger
// 	settings     *settings.SettingsRepository
// 	auth         authentication.Repository
// 	accountQuery *account_querier.Querier
// 	register     *register.Registrar
// 	er           email.EmailRepo

// 	// TODO: Replace with an MQ message and sender job.
// 	sender email_verify.Verifier
// }

// func New(
// 	logger *zap.Logger,
// 	settings *settings.SettingsRepository,
// 	auth authentication.Repository,
// 	accountQuery *account_querier.Querier,
// 	er email.EmailRepo,
// 	register *register.Registrar,
// 	sender email_verify.Verifier,
// ) *Provider {
// 	return &Provider{
// 		logger:       logger,
// 		settings:     settings,
// 		auth:         auth,
// 		accountQuery: accountQuery,
// 		er:           er,
// 		register:     register,
// 		sender:       sender,
// 	}
// }
