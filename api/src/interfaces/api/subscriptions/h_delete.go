package subscriptions

import (
	"net/http"

	"github.com/Southclaws/storyden/api/src/infra/web"
	"github.com/Southclaws/storyden/api/src/services/authentication"
	"github.com/go-chi/chi"
)

type DeleteResponse struct {
	Deleted int `json:"deleted"`
}

// @Summary  Delete subscription
// @Tags     subscriptions
// @Accept   json
// @Produce  json
// @Param    id   path      string  true  "Subscription ID"
// @Success  200  {object}  DeleteResponse
// @Failure  500  {object}  web.Error
// @Router   /subscriptions/{id} [delete]
func (c *controller) delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	info, ok := authentication.GetAuthenticationInfo(w, r)
	if !ok {
		return
	}

	deleted, err := c.repo.Unsubscribe(r.Context(), info.Cookie.UserID, id)
	if err != nil {
		web.StatusInternalServerError(w, err)
	}

	web.Write(w, DeleteResponse{deleted})
}
