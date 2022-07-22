package authentication

import (
	"context"
	"errors"
	"net/http"

	"github.com/Southclaws/storyden/backend/internal/infrastructure/web"
)

// WithAuthentication provides middleware for extracting session data.
func (a *CookieAuth) WithAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := a.Decode(r)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		ctx := context.WithValue(
			r.Context(),
			contextKey,
			user,
		)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func MustBeAuthenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, ok := GetUser(w, r)
		if !ok {
			return
		}

		next.ServeHTTP(w, r)
	})
}

func MustBeAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := GetUser(w, r)
		if !ok {
			return
		}

		if !user.Admin {
			web.StatusUnauthorized(w, errors.New("user is not an administrator"))
			return
		}

		next.ServeHTTP(w, r)
	})
}
