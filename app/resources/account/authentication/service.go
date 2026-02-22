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

type Service struct {
	v       string
	builtIn BuiltInService
}

func NewService(v string) Service {
	bis, err := NewBuiltInService(v)
	if err != nil {
		return Service{v: v}
	}
	return Service{builtIn: bis, v: string(bis.v)}
}

func (s Service) String() string {
	if s.builtIn.v != "" {
		return string(s.builtIn.v)
	}
	return s.v
}

type tokenTypeEnum string

const (
	tokenTypeNone         tokenTypeEnum = "none"          // Authenticated by other means
	tokenTypePasswordHash tokenTypeEnum = "password_hash" // argon2 hashed password
	tokenTypeWebAuthn     tokenTypeEnum = "webauthn"      // WebAuthn token
	tokenTypeOAuth        tokenTypeEnum = "oauth"         // OAuth2 token
)
