package users

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi"

	"github.com/Southclaws/storyden/api/src/infra/db"
	"github.com/Southclaws/storyden/api/src/infra/web"
	"github.com/Southclaws/storyden/api/src/services/authentication"
)

func (c *controller) get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	// context is public, hide sensitive PII like email
	public := true
	if auth, ok := authentication.GetAuthenticationInfoFromContext(r.Context()); ok {
		if auth.Cookie.UserID == id || auth.Cookie.Admin {
			// context is owner/admin, reveal sensitive PII
			public = false
		}
	}

	user, err := c.repo.GetUser(r.Context(), id, public)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			web.StatusNotFound(w, err)
		} else {
			web.StatusInternalServerError(w, err)
		}
		return
	}

	web.Write(w, user) //nolint:errcheck
}
