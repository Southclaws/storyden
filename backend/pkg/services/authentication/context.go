package authentication

import (
	"context"
	"errors"
	"net/http"

	"github.com/Southclaws/storyden/backend/internal/web"
	"github.com/Southclaws/storyden/backend/pkg/resources/user"
)

var contextKey = struct{}{}

func AddUserToContext(ctx context.Context, u *user.User) context.Context {
	return context.WithValue(ctx, contextKey, u)
}

// GetUser extracts auth info from a request context and, if not present, will
// write a 500 error to the response and return not-ok. In this failure case,
// the request should be immediately terminated.
func GetUser(
	w http.ResponseWriter,
	r *http.Request,
) (*user.User, bool) {
	if auth, ok := GetUserFromContext(r.Context()); ok {
		return auth, true
	}

	web.StatusUnauthorized(w, web.WithSuggestion(
		errors.New("user not authenticated"),
		"The request did not have any authentication information with it.",
		"Ensure you are logged in, try logging out and back in again. If issues persist, please contact us.",
	))

	return nil, false
}

// GetUserFromContext pulls out auth data from a context.
func GetUserFromContext(ctx context.Context) (*user.User, bool) {
	if auth, ok := ctx.Value(contextKey).(user.User); ok {
		return &auth, true
	}

	return nil, false
}
