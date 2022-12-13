package authentication

import (
	"context"

	"github.com/Southclaws/storyden/app/resources/account"
)

// Provider describes a type that can be used to authenticate people.
//
// Link simply returns a URL to start the authentication process.
//
// Login is called by the callback and handles the code/token exchange and
// returns a User object to the caller to be encoded into a cookie.
type Provider interface {
	// Enabled tells the OAuth component if the provider is enabled. Providers
	// can be enabled and disabled at initialisation time based on the presence
	// of environment variable configuration for each provider.
	Enabled() bool

	// Provider general information. This could be a struct but it's simpler as
	// part of the interface for now, despite bloating the interface a bit.
	Name() string
	ID() string
	LogoURL() string

	// Link will, for providers that support it, provide a URL to a third-party
	// authenticator. OAuth providers use this to start the authentication flow.
	Link() string

	// Login is a function that will validate and authenticate a user given that
	// the provider is happy with the input. The input format differs depending
	// on the provider. For example, an OAuth provider will use `state` and
	// `secret` as the `state` and `code` components of the OAuth2 specification
	// and a simple password-based provider may simply want the account handle
	// in the `state` and their password in the `secret`.
	Login(ctx context.Context, state, secret string) (*account.Account, error)
}
