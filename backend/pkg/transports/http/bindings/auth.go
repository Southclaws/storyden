package bindings

import (
	"net/http"

	"github.com/Southclaws/storyden/backend/pkg/services/authentication"
	"github.com/Southclaws/storyden/backend/pkg/transports/http/openapi"
)

type Authentication struct{ s authentication.Service }

func NewAuthentication(s authentication.Service) Authentication { return Authentication{s} }

func (i *Authentication) GetV1AuthPassword(w http.ResponseWriter, r *http.Request, params openapi.GetV1AuthPasswordParams) {
	//
}

func (i *Authentication) middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if session, ok := i.s.DecodeSession(r); ok {
			ctx := authentication.AddUserToContext(r.Context(), session)
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			next.ServeHTTP(w, r)
		}
	})
}
