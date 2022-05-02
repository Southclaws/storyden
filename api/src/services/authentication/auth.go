package authentication

import (
	"net/http"

	"github.com/gorilla/securecookie"

	"github.com/Southclaws/storyden/api/src/config"
	"github.com/Southclaws/storyden/api/src/resources/authentication"
	"github.com/Southclaws/storyden/api/src/resources/user"
)

var contextKey = struct{}{}

const secureCookieName = "storyden-session"

// CookieAuth implements Authentication for cookies
type CookieAuth struct {
	repo   authentication.Repository
	sc     *securecookie.SecureCookie
	domain string
}

// New initialises a new authentication service
func NewCookieAuth(cfg config.Config, repo authentication.Repository) CookieAuth {
	return CookieAuth{
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
