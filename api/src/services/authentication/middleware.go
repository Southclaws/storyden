package authentication

import (
	"context"
	"errors"
	"net/http"

	"github.com/Southclaws/storyden/api/src/infra/web"
)

// Info represents data that is extracted from the path, validated
// against auth and stored in request context.
type Info struct {
	Authenticated bool
	Cookie        Cookie
}

var contextKey = struct{}{}

const secureCookieName = "storyden-session"

// WithAuthentication provides middleware for enforcing authentication
func (a *State) WithAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := Info{}

		if a.doCookieAuth(r, &auth) {
			auth.Authenticated = true
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
