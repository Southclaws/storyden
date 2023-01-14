package bindings

import (
	"net/http"

	"github.com/Southclaws/fault"
	"github.com/gorilla/securecookie"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/internal/config"
)

type Session struct {
	sc     *securecookie.SecureCookie
	domain string
}

func NewSessionManager(cfg config.Config, sc *securecookie.SecureCookie) Session {
	return Session{sc, cfg.CookieDomain}
}

const secureCookieName = "storyden-session"

type session struct {
	UserID account.AccountID
}

func (s *Session) encodeSession(userID account.AccountID) (string, error) {
	encoded, err := s.sc.Encode(secureCookieName, session{userID})
	if err != nil {
		return "", fault.Wrap(err)
	}

	cookie := &http.Cookie{
		Name:     secureCookieName,
		Value:    encoded,
		SameSite: http.SameSiteNoneMode,
		Path:     "/",
		Domain:   s.domain,
		Secure:   false,
		HttpOnly: true,
	}

	return cookie.String(), nil
}

func (s *Session) decodeSession(r *http.Request) (*session, bool) {
	cookie, err := r.Cookie(secureCookieName)
	if err != nil {
		return nil, false
	}

	u := session{}

	if err = s.sc.Decode(secureCookieName, cookie.Value, &u); err != nil {
		return nil, false
	}

	return &u, true
}
