package authentication

import (
	"context"
	"net/http"
	"time"

	"github.com/gorilla/securecookie"
	"github.com/pkg/errors"

	"github.com/Southclaws/storyden/api/src/config"
	"github.com/Southclaws/storyden/api/src/infra/web"
	"github.com/Southclaws/storyden/api/src/resources/user"
)

// State stores state for performing authentication
type State struct {
	sc     *securecookie.SecureCookie
	domain string
	users  user.Repository
}

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

// New initialises a new authentication service
func New(cfg config.Config, users user.Repository) *State {
	a := &State{
		sc:     securecookie.New(cfg.HashKey, cfg.BlockKey),
		domain: cfg.CookieDomain,
		users:  users,
	}

	return a
}

// EncodeAuthCookie writes the secure user auth cookie to the response writer.
func (a *State) EncodeAuthCookie(w http.ResponseWriter, user user.User) {
	encoded, err := a.sc.Encode(secureCookieName, Cookie{
		UserID:  user.ID,
		Admin:   user.Admin,
		Created: time.Now(),
	})
	if err != nil {
		web.StatusUnauthorized(w, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     secureCookieName,
		Value:    encoded,
		Path:     "/",
		Domain:   a.domain,
		Secure:   true,
		HttpOnly: true,
	})
}

// GetAuthenticationInfo extracts auth info from a request context and, if not
// present, will write a 500 error to the response and return not-ok. In this
// failure case, the request should be immediately terminated.
func GetAuthenticationInfo(
	w http.ResponseWriter,
	r *http.Request,
) (*Info, bool) {
	if auth, ok := GetAuthenticationInfoFromContext(r.Context()); ok {
		return auth, true
	}
	web.StatusInternalServerError(w, web.WithSuggestion(
		errors.New("failed to extract auth context from request"),
		"Could not read session data from cookies.",
		"Try clearing your cookies and logging in to your account again."))
	return nil, false
}

// GetAuthenticationInfoFromContext pulls out auth data from a request context
func GetAuthenticationInfoFromContext(ctx context.Context) (*Info, bool) {
	if auth, ok := ctx.Value(contextKey).(Info); ok {
		return &auth, true
	}
	return nil, false
}

func IsRequestAdmin(r *http.Request) bool {
	info, _ := GetAuthenticationInfoFromContext(r.Context())
	if info == nil {
		return false
	}
	return info.Cookie.Admin
}

func (s *State) GetOrCreateFromContext(ctx context.Context, email, authMethod, username string) (*user.User, error) {
	if existing, ok := GetAuthenticationInfoFromContext(ctx); ok && existing.Authenticated {
		u, err := s.users.GetUser(ctx, existing.Cookie.UserID, false)
		if err != nil {
			return nil, errors.Wrap(err, "failed to find user account")
		}
		return u, nil
	} else {
		u, err := s.users.CreateUser(ctx, email, user.AuthMethod(authMethod), username)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create user account")
		}
		return u, nil
	}
}
