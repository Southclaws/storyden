// Package session provides session handling primitives and middleware for the
// API. Sessions work by encrypting an account's ID inside a cookie value. This
// is read via a middleware and dropped into the request context for later use.
package session

import (
	"context"
	"net/http"
	"time"

	"github.com/rs/xid"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/internal/config"
	"github.com/Southclaws/storyden/internal/infrastructure/endec/securecookie"
)

// TODO: Allow changing this via config.
const (
	secureCookieName = "storyden-session"
	sameSiteMode     = http.SameSiteDefaultMode
	cookieLifespan   = time.Hour * 24 * 90
)

func expiryFunc() time.Time {
	return time.Now().Add(cookieLifespan)
}

type Jar struct {
	ss               *securecookie.Session
	domain           string
	secureCookieName string
}

func New(cfg config.Config, ss *securecookie.Session) *Jar {
	return &Jar{
		domain:           cfg.CookieDomain,
		ss:               ss,
		secureCookieName: secureCookieName,
	}
}

func (j *Jar) createWithValue(value string, expire time.Time) *http.Cookie {
	return &http.Cookie{
		Name:     secureCookieName,
		Value:    value,
		SameSite: sameSiteMode,
		Expires:  expire,
		Path:     "/",
		Domain:   j.domain,

		// Always secure, localhost is automatically excluded by browsers.
		Secure: true,

		// JS never needs to access these cookies.
		HttpOnly: true,
	}
}

func (j *Jar) Create(accountID string) *http.Cookie {
	return j.createWithValue(j.ss.Encrypt(accountID), expiryFunc())
}

func (j *Jar) Destroy() *http.Cookie {
	return j.createWithValue("", time.Now())
}

// WithSession checks the request for a session and drops it into a context.
func (j *Jar) WithSession(r *http.Request) context.Context {
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

	return session.WithAccountID(r.Context(), account.AccountID(id))
}

// WithAuth simply pulls out the session from the cookie and propagates it.
func (j *Jar) WithAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := j.WithSession(r)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (j *Jar) GetCookieName() string {
	return j.secureCookieName
}
