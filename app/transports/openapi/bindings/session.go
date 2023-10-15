package bindings

import (
	"context"
	"net/http"
	"time"

	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/services/authentication"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/securecookie"
)

const secureCookieName = "storyden-session"

type CookieJar struct {
	ss     *securecookie.Session
	domain string
}

func newCookieJar(cfg config.Config, ss *securecookie.Session) *CookieJar {
	return &CookieJar{domain: cfg.CookieDomain, ss: ss}
}

// Create an encrypted cookie from an account ID.
func (j *CookieJar) Create(accountID string) *http.Cookie {
	return &http.Cookie{
		Name:     secureCookieName,
		Value:    j.ss.Encrypt(accountID),
		SameSite: http.SameSiteDefaultMode,
		Expires:  time.Now().Add(time.Hour * 24 * 90),
		Path:     "/",
		Domain:   j.domain,
		Secure:   true,
		HttpOnly: true,
	}
}

// WithSession checks the request for a session and drops it into a context.
func (j *CookieJar) WithSession(r *http.Request) context.Context {
	cookie, err := r.Cookie(secureCookieName)
	if err != nil {
		return r.Context()
	}

	accountID, ok := j.ss.Decrypt(cookie.Value)
	if !ok {
		return r.Context()
	}

	id, err := xid.FromString(accountID)
	if err != nil {
		return r.Context()
	}

	return authentication.WithAccountID(r.Context(), account.AccountID(id))
}

// WithAuth simply pulls out the session from the cookie and propagates it.
func (j *CookieJar) WithAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r.WithContext(j.WithSession(r)))
	})
}
