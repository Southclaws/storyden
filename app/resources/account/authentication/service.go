package authentication

//go:generate go run github.com/Southclaws/enumerator

type serviceEnum string

const (
	serviceUsernamePassword serviceEnum = "username_password" // User/email + password
	serviceEmailPassword    serviceEnum = "email_password"    // Email + password
	serviceEmailVerify      serviceEnum = "email_verify"      // Email + verification code
	servicePhoneVerify      serviceEnum = "phone_verify"      // Phone number + verification code
	serviceWebAuthn         serviceEnum = "webauthn"          // WebAuthn/Passkey

	// OAuth services
	serviceOAuthGoogle   serviceEnum = "oauth_google"   // Google
	serviceOAuthGitHub   serviceEnum = "oauth_github"   // GitHub
	serviceOAuthLinkedin serviceEnum = "oauth_linkedin" // LinkedIn
)

type tokenTypeEnum string

const (
	tokenTypeNone     tokenTypeEnum = "none"     // Authenticated by other means
	tokenTypePassword tokenTypeEnum = "password" // argon2 hashed password
	tokenTypeWebAuthn tokenTypeEnum = "webauthn" // WebAuthn token
	tokenTypeOAuth    tokenTypeEnum = "oauth"    // OAuth2 token
)
