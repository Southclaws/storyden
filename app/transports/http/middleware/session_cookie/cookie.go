// Package session provides session handling primitives and middleware for the
// API. Sessions work by encrypting an account's ID inside a cookie value. This
// is read via a middleware and dropped into the request context for later use.
package session_cookie

import (
	"context"
	"net/http"
	"net/url"
	"time"

	"github.com/rs/xid"
	"golang.org/x/net/publicsuffix"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fmsg"
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

func New(cfg config.Config, ss *securecookie.Session) (*Jar, error) {
	domain, err := getDomain(cfg.PublicAPIAddress)
	if err != nil {
		return nil, fault.Wrap(err, fmsg.With("failed to parse domain from public API address"))
	}

	return &Jar{
		domain:           domain,
		ss:               ss,
		secureCookieName: secureCookieName,
	}, nil
}

func getDomain(address url.URL) (string, error) {
	// We want to use the site's domain, not the API's subdomain to ensure that
	// cookies can be used in both the frontend and for the API. This assumption
	// is based on the idea that Storyden must be hosted on a single domain with
	// the API and frontend on different subdomains. For example, if your site
	// was "www.cats.com" and the API was "api.cats.com", then the domain config
	// would be set up with `www.cats.com` as the `PUBLIC_WEB_ADDRESS` and then
	// `api.cats.com` as the `PUBLIC_API_ADDRESS` and then this code would use
	// the API address hostname to parse `cats.com` as the actual cookie domain.
	// The reason for this is that it makes SSR frontends trivial to implement.

	hostname := address.Hostname()

	if hostname == "localhost" {
		return hostname, nil
	}

	domain, err := publicsuffix.EffectiveTLDPlusOne(hostname)
	if err != nil {
		return "", fault.Wrap(err, fmsg.With("failed to parse domain from public API address"))
	}

	return domain, nil
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
