package users

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/Southclaws/storyden/api/src/infra/web"
)

type BanstatusBody struct {
	Banned bool `json:"banned"`
}

// @Summary  Ban a user
// @Tags     users
// @Accept   json
// @Produce  json
// @Param    id    path      string         true  "User ID"
// @Param    body  body      BanstatusBody  true  "Ban status payload"
// @Success  200   {object}  user.User
// @Failure  400   {object}  web.Error
// @Failure  500   {object}  web.Error
// @Failure  404   {object}  web.Error
// @Router   /users/{id}/banstatus [patch]
func (c *controller) banstatus(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var p BanstatusBody
	if err := web.DecodeBody(r, &p); err != nil {
		web.StatusBadRequest(w, err)
		return
	}

	if p.Banned {
		user, err := c.repo.Ban(r.Context(), id)
		if err != nil {
			web.StatusInternalServerError(w, err)
			return
		}
		web.Write(w, user)
	} else {
		user, err := c.repo.Unban(r.Context(), id)
		if err != nil {
			web.StatusInternalServerError(w, err)
			return
		}
		web.Write(w, user)
	}
}
