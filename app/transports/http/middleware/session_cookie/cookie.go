// Package session provides session handling primitives and middleware for the
// API. Sessions work by encrypting an account's ID inside a cookie value. This
// is read via a middleware and dropped into the request context for later use.
package session_cookie

import (
	"context"
	"net/http"
	"time"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/storyden/app/resources/account/token"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/internal/config"
)

// TODO: Allow changing this via config.
const (
	secureCookieName = "storyden-session"
	sameSiteMode     = http.SameSiteLaxMode
	cookieLifespan   = time.Hour * 24 * 90
)

func expiryFunc() time.Time {
	return time.Now().Add(cookieLifespan)
}

type Jar struct {
	validator        *session.Validator
	issuer           *session.Issuer
	domain           string
	secureCookieName string
}

func New(cfg config.Config, v *session.Validator) (*Jar, error) {
	domain, err := getCookieDomain(cfg.PublicAPIAddress, cfg.PublicWebAddress)
	if err != nil {
		return nil, fault.Wrap(err, fmsg.With("failed to parse domain from public API address"))
	}

	return &Jar{
		domain:           domain,
		validator:        v,
		secureCookieName: secureCookieName,
	}, nil
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

func (j *Jar) Create(t token.Token) *http.Cookie {
	return j.createWithValue(t.String(), expiryFunc())
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

	ctx, err := j.validator.Validate(r.Context(), cookie.Value)
	if err != nil {
		return r.Context()
	}

	return ctx
}

// WithAuth simply pulls out the session from the cookie and propagates it.
func (j *Jar) WithAuth() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := j.WithSession(r)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func (j *Jar) GetCookieName() string {
	return j.secureCookieName
}
