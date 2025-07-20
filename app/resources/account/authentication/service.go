package authentication

//go:generate go run github.com/Southclaws/enumerator

type serviceEnum string

const (
	servicePassword    serviceEnum = "password"     // Password + either username or email
	serviceEmailVerify serviceEnum = "email_verify" // Email + verification code
	servicePhoneVerify serviceEnum = "phone_verify" // Phone number + verification code
	serviceWebAuthn    serviceEnum = "webauthn"     // WebAuthn/Passkey
	serviceAccessKey   serviceEnum = "access_key"   // API access key

	// OAuth services
	serviceOAuthGoogle  serviceEnum = "oauth_google"  // Google
	serviceOAuthGitHub  serviceEnum = "oauth_github"  // GitHub
	serviceOAuthDiscord serviceEnum = "oauth_discord" // Discord
	serviceOAuthKeycloak serviceEnum = "oauth_keycloak" // Keycloak
)

type tokenTypeEnum string

const (
	tokenTypeNone         tokenTypeEnum = "none"          // Authenticated by other means
	tokenTypePasswordHash tokenTypeEnum = "password_hash" // argon2 hashed password
	tokenTypeWebAuthn     tokenTypeEnum = "webauthn"      // WebAuthn token
	tokenTypeOAuth        tokenTypeEnum = "oauth"         // OAuth2 token
)
