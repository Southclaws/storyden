package threads

import (
	"net/http"

	"github.com/go-chi/chi"

	"github.com/Southclaws/storyden/api/src/infra/web"
	"github.com/Southclaws/storyden/api/src/services/authentication"
)

type DeleteResponse struct {
	Count int `json:"count"`
}

// @Summary  Soft-delete all posts in a thread
// @Tags     threads
// @Produce  json
// @Param    id   path      string  true  "Root post ID"
// @Success  200  {object}  DeleteResponse
// @Failure  500  {object}  web.Error
// @Router   /threads/{id} [delete]
func (c *controller) delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	info, ok := authentication.GetAuthenticationInfo(w, r)
	if !ok {
		return
	}

	count, err := c.threads.Delete(r.Context(), id, info.Cookie.UserID)
	if err != nil {
		web.StatusInternalServerError(w, err)
		return
	}

	web.Write(w, DeleteResponse{count})
}
