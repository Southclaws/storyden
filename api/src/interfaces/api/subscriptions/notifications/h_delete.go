package notifications

import (
	"net/http"

	"github.com/Southclaws/storyden/api/src/infra/web"
	"github.com/Southclaws/storyden/api/src/services/authentication"
	"github.com/go-chi/chi"
)

// @Summary  Delete notification
// @Tags     notifications
// @Accept   json
// @Produce  json
// @Param    id   path      string  true  "Notification ID"
// @Success  200  {object}  notification.Notification
// @Failure  500  {object}  web.Error
// @Failure  400  {object}  web.Error
// @Router   /subscriptions/notifications/{id} [delete]
func (c *controller) delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	info, ok := authentication.GetAuthenticationInfo(w, r)
	if !ok {
		return
	}

	notification, err := c.repo.Delete(r.Context(), info.Cookie.UserID, id)
	if err != nil {
		web.StatusInternalServerError(w, err)
		return
	}
	if notification == nil {
		web.StatusNotFound(w, err)
		return
	}

	web.Write(w, notification)
}
