package authentication

import (
	"context"
	"net/http"

	"github.com/gorilla/securecookie"
	"github.com/pkg/errors"

	"github.com/Southclaws/storyden/api/src/config"
	"github.com/Southclaws/storyden/api/src/infra/web"
	"github.com/Southclaws/storyden/api/src/resources/authentication"
	"github.com/Southclaws/storyden/api/src/resources/user"
)

// CookieAuth implements Authentication for cookies
type CookieAuth struct {
	repo   authentication.Repository
	sc     *securecookie.SecureCookie
	domain string
}

// New initialises a new authentication service
func NewCookieAuth(cfg config.Config, repo authentication.Repository) Contract {
	return &CookieAuth{
		repo,
		securecookie.New(cfg.HashKey, cfg.BlockKey),
		cfg.CookieDomain,
	}
}

func (a *CookieAuth) Encode(w http.ResponseWriter, user user.User) error {
	encoded, err := a.sc.Encode(secureCookieName, user)
	if err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     secureCookieName,
		Value:    encoded,
		Path:     "/",
		Domain:   a.domain,
		Secure:   true,
		HttpOnly: true,
	})

	return nil
}

func (a *CookieAuth) Decode(r *http.Request) (*user.User, error) {
	cookie, err := r.Cookie(secureCookieName)
	if err != nil {
		return nil, err
	}

	u := user.User{}

	if err = a.sc.Decode(secureCookieName, cookie.Value, &u); err != nil {
		return nil, err
	}

	return &u, nil
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
