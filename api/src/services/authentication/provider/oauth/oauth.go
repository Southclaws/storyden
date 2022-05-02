package oauth

import (
	"context"

	"go.uber.org/fx"

	"github.com/Southclaws/storyden/api/src/resources/user"
)

// OAuthProvider describes a type that can provide an OAuth2 authentication
// method for users.
//
// Link simply returns a URL to start the OAuth2 process.
//
// Login is called by the callback and handles the code/token exchange and
// returns a User object to the caller to be encoded into a cookie.
type OAuthProvider interface {
	Link() string
	Login(ctx context.Context, state, code string) (*user.User, error)
}

func Build() fx.Option {
	return fx.Options(
	// github.Build(),
	)
}
