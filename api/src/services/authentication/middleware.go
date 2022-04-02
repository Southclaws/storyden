package authentication

import (
	"context"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/Southclaws/storyden/api/src/infra/web"
)

// Info represents data that is extracted from the path, validated
// against auth and stored in request context.
type Info struct {
	Authenticated bool
	Cookie        Cookie
}

// Cookie represents the data structure that is stored inside the securecookie
// value. It is sent to the client in an encrypted (HMAC) form and decrypted on
// any future requests and propagated to handlers.
type Cookie struct {
	UserID  string
	Admin   bool
	Created time.Time
}

var contextKey = struct{}{}

const secureCookieName = "openmultiplayer-session"

// WithAuthentication provides middleware for enforcing authentication
func (a *State) WithAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := Info{}

		if a.doCookieAuth(r, &auth) {
			auth.Authenticated = true
			// } else if a.doTokenAuth(r, &auth) {
			// 	auth.Authenticated = true
		} else {
			auth.Authenticated = false
		}

		// If the request contained a valid cookie, `auth.Authenticated` is now
		// true and `auth.Cookie` contains the user information.
		// Otherwise, `auth.Authenticated` is false and `auth.Cookie` is empty.

		next.ServeHTTP(w, r.WithContext(context.WithValue(
			r.Context(),
			contextKey,
			auth,
		)))
	})
}

func (a *State) doCookieAuth(r *http.Request, auth *Info) bool {
	cookie, err := r.Cookie(secureCookieName)
	if err != nil {
		return false
	}

	if err = a.sc.Decode(secureCookieName, cookie.Value, &auth.Cookie); err != nil {
		zap.L().Debug("failed to decode auth cookie", zap.Error(err))
		return false
	}

	return true
}

func MustBeAuthenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth, ok := GetAuthenticationInfo(w, r)
		if !ok {
			return
		}
		if !auth.Authenticated {
			web.StatusUnauthorized(w, web.WithSuggestion(
				errors.New("user not authenticated"),
				"The request did not have any authentication information with it.",
				"Ensure you are logged in, try logging out and back in again. If issues persist, please contact us.",
			))
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (a *State) MustBeAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth, ok := GetAuthenticationInfo(w, r)
		if !ok {
			return
		}

		if !auth.Cookie.Admin {
			web.StatusUnauthorized(w, errors.New("user is not an administrator"))
			return
		}

		next.ServeHTTP(w, r)
	})
}
