package oauth

import (
	"context"

	"github.com/Southclaws/storyden/app/resources/account"
)

// Provider describes a type that can provide an OAuth2 authentication
// method for users.
//
// Link simply returns a URL to start the OAuth2 process.
//
// Login is called by the callback and handles the code/token exchange and
// returns a User object to the caller to be encoded into a cookie.
type Provider interface {
	// Enabled tells the OAuth component if the provider is enabled. Providers
	// can be enabled and disabled at initialisation time based on the presence
	// of environment variable configuration for each provider.
	Enabled() bool

	// Provider general information.
	Name() string
	ID() string
	LogoURL() string

	// Link is the actual login link to show to users.
	Link() string

	// Login is the callback function.
	Login(ctx context.Context, state, code string) (*account.Account, error)
}
