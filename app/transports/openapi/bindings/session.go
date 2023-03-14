package bindings

import (
	"context"
	"net/http"

	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/services/authentication"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/securecookie"
)

const secureCookieName = "storyden-session"

type cookieJar struct {
	ss     *securecookie.Session
	domain string
}

func newCookieJar(cfg config.Config, ss *securecookie.Session) *cookieJar {
	return &cookieJar{domain: cfg.CookieDomain, ss: ss}
}

// Create an encrypted cookie from an account ID.
func (j *cookieJar) Create(accountID string) string {
	return (&http.Cookie{
		Name:     secureCookieName,
		Value:    j.ss.Encrypt(accountID),
		SameSite: http.SameSiteDefaultMode,
		Path:     "/",
		Domain:   j.domain,
		Secure:   true,
		HttpOnly: true,
	}).String()
}

// withSession checks the request for a session and drops it into a context.
func (j *cookieJar) withSession(r *http.Request) context.Context {
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
func (j *cookieJar) WithAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r.WithContext(j.withSession(r)))
	})
}
