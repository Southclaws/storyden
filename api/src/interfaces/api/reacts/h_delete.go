package reacts

import (
	"net/http"

	"github.com/go-chi/chi"

	"github.com/Southclaws/storyden/api/src/infra/web"
	"github.com/Southclaws/storyden/api/src/services/authentication"
)

// @Summary  Delete a reaction
// @Tags     reacts
// @Accept   json
// @Produce  json
// @Param    id   path      string  true  "React ID"
// @Success  200  {object}  react.React
// @Failure  400  {object}  web.Error
// @Failure  500  {object}  web.Error
// @Failure  404  {object}  web.Error
// @Router   /reacts/{id} [delete]
func (c *controller) delete(w http.ResponseWriter, r *http.Request) {
	reactID := chi.URLParam(r, "react_id")
	info, ok := authentication.GetAuthenticationInfo(w, r)
	if !ok {
		return
	}

	reaction, err := c.reacts.Remove(r.Context(), info.Cookie.UserID, reactID)
	if err != nil {
		web.StatusInternalServerError(w, err)
		return
	}
	if reaction == nil {
		web.StatusNotFound(w, nil)
		return
	}

	web.Write(w, reaction)
}
