package authentication

import (
	"net/http"

	"github.com/gorilla/securecookie"

	"github.com/Southclaws/storyden/backend/internal/config"
	"github.com/Southclaws/storyden/backend/pkg/resources/authentication"
	"github.com/Southclaws/storyden/backend/pkg/resources/user"
)

// cookie implements authentication service using encrypted HTTP cookies.
type cookie struct {
	auth_repo authentication.Repository
	sc        *securecookie.SecureCookie
	domain    string
}

func newCookie(cfg config.Config, auth_repo authentication.Repository) Service {
	return &cookie{
		auth_repo: auth_repo,
		sc:        securecookie.New(cfg.HashKey, cfg.BlockKey),
		domain:    cfg.CookieDomain,
	}
}

const secureCookieName = "storyden-session"

func (s *cookie) EncodeSession(w http.ResponseWriter, user user.User) error {
	encoded, err := s.sc.Encode(secureCookieName, user)
	if err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     secureCookieName,
		Value:    encoded,
		Path:     "/",
		Domain:   s.domain,
		Secure:   true,
		HttpOnly: true,
	})

	return nil
}

func (s *cookie) DecodeSession(r *http.Request) (*user.User, bool) {
	cookie, err := r.Cookie(secureCookieName)
	if err != nil {
		return nil, false
	}

	u := user.User{}

	if err = s.sc.Decode(secureCookieName, cookie.Value, &u); err != nil {
		return nil, false
	}

	return &u, true
}
