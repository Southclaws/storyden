package authentication

//go:generate go run github.com/Southclaws/enumerator

type serviceEnum string

const (
	servicePassword  serviceEnum = "password"   // Username + password (default)
	serviceEmailOnly serviceEnum = "email_only" // Email + verification code
	servicePhone     serviceEnum = "phone"      // Phone number + verification code

	// OAuth services
	serviceOAuthGoogle   serviceEnum = "oauth_google"   // Google
	serviceOAuthGitHub   serviceEnum = "oauth_github"   // GitHub
	serviceOAuthLinkedin serviceEnum = "oauth_linkedin" // LinkedIn
)
