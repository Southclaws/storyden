package users

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/Southclaws/storyden/api/src/infra/web"
)

type PatchBody struct {
	Status bool `json:"status"`
}

// @Summary  Set a user's admin status
// @Tags     users
// @Accept   json
// @Produce  json
// @Param    id    path      string     true  "User ID"
// @Param    body  body      PatchBody  true  "Operation parameters"
// @Success  200   {object}  user.User
// @Failure  400   {object}  web.Error
// @Failure  500   {object}  web.Error
// @Failure  404   {object}  web.Error
// @Router   /users/{id}/adminstatus [patch]
func (c *controller) patchAdmin(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var p PatchBody
	if !web.ParseBody(w, r, &p) {
		return
	}

	done, err := c.repo.SetAdmin(r.Context(), id, p.Status)
	if err != nil {
		web.StatusInternalServerError(w, err)
		return
	}

	if !done {
		web.StatusNotFound(w, errors.New("failed to update user"))
		return
	}
}
