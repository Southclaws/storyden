package users

import (
	"errors"
	"net/http"

	"github.com/Southclaws/storyden/api/src/infra/web"
	"github.com/Southclaws/storyden/api/src/services/authentication"
	"github.com/go-chi/chi"
)

type patchPayload struct {
	Email *string `json:"email"`
	Name  *string `json:"name"`
	Bio   *string `json:"bio"`
}

func (c *controller) patch(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	ai, ok := authentication.GetAuthenticationInfo(w, r)
	if !ok {
		return
	}

	if ai.Cookie.UserID != id && !ai.Cookie.Admin {
		web.StatusUnauthorized(w, errors.New("not authorised to modify user"))
		return
	}

	var p patchPayload
	if !web.ParseBody(w, r, &p) {
		return
	}

	user, err := c.repo.UpdateUser(r.Context(), ai.Cookie.UserID, p.Email, p.Name, p.Bio)
	if err != nil {
		web.StatusInternalServerError(w, err)
		return
	}

	web.Write(w, user) //nolint:errcheck
}
