package authentication

import (
	"errors"
	"net/http"

	"github.com/Southclaws/storyden/backend/internal/web"
)

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
