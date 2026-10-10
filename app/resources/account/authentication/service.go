package authentication

//go:generate go run github.com/Southclaws/enumerator

type builtInServiceEnum string

const (
	servicePassword    builtInServiceEnum = "password"     // Password + either username or email
	serviceEmailVerify builtInServiceEnum = "email_verify" // Email + verification code
	servicePhoneVerify builtInServiceEnum = "phone_verify" // Phone number + verification code
	serviceWebAuthn    builtInServiceEnum = "webauthn"     // WebAuthn/Passkey
	serviceAccessKey   builtInServiceEnum = "access_key"   // API access key

	// OAuth services
	serviceOAuthGoogle   builtInServiceEnum = "oauth_google"   // Google
	serviceOAuthGitHub   builtInServiceEnum = "oauth_github"   // GitHub
	serviceOAuthDiscord  builtInServiceEnum = "oauth_discord"  // Discord
	serviceOAuthKeycloak builtInServiceEnum = "oauth_keycloak" // Keycloak
)

// Service-typed aliases for built-in services.
//
// Use these when APIs require a Service but the identifier is one of the
// built-in provider values.
var (
	AuthServicePassword    = Service{builtIn: ServicePassword}
	AuthServiceEmailVerify = Service{builtIn: ServiceEmailVerify}
	AuthServicePhoneVerify = Service{builtIn: ServicePhoneVerify}
	AuthServiceWebAuthn    = Service{builtIn: ServiceWebAuthn}
	AuthServiceAccessKey   = Service{builtIn: ServiceAccessKey}

	AuthServiceOAuthGoogle   = Service{builtIn: ServiceOAuthGoogle}
	AuthServiceOAuthGitHub   = Service{builtIn: ServiceOAuthGitHub}
	AuthServiceOAuthDiscord  = Service{builtIn: ServiceOAuthDiscord}
	AuthServiceOAuthKeycloak = Service{builtIn: ServiceOAuthKeycloak}
)

type Service struct {
	custom  string
	builtIn BuiltInService
}

func NewService(v string) Service {
	bis, err := NewBuiltInService(v)
	if err != nil {
		return Service{custom: v}
	}
	return Service{builtIn: bis}
}

func (s Service) String() string {
	if s.IsBuiltIn() {
		return string(s.builtIn.v)
	}
	return s.custom
}

func (s Service) IsBuiltIn() bool {
	return s.builtIn.v != ""
}

func (s Service) BuiltIn() (*BuiltInService, bool) {
	if !s.IsBuiltIn() {
		return nil, false
	}
	return &s.builtIn, true
}

type tokenTypeEnum string

const (
	tokenTypeNone         tokenTypeEnum = "none"          // Authenticated by other means
	tokenTypePasswordHash tokenTypeEnum = "password_hash" // argon2 hashed password
	tokenTypeWebAuthn     tokenTypeEnum = "webauthn"      // WebAuthn token
	tokenTypeOAuth        tokenTypeEnum = "oauth"         // OAuth2 token
)
