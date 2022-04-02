package notifications

import (
	"net/http"

	"github.com/Southclaws/storyden/api/src/infra/web"
	"github.com/Southclaws/storyden/api/src/resources/notification"
	"github.com/Southclaws/storyden/api/src/services/authentication"
	"github.com/go-chi/chi"
)

type PatchBody notification.Notification

// @Summary  Update a notification's read state
// @Tags     notifications
// @Accept   json
// @Produce  json
// @Param    id       path      string     true  "Notification ID"
// @Param    payload  body      PatchBody  true  "Notification object with read state set"
// @Success  200      {object}  notification.Notification
// @Failure  400      {object}  web.Error
// @Failure  500      {object}  web.Error
// @Failure  404      {object}  web.Error
// @Router   /subscriptions/notifications/{id} [patch]
func (c *controller) patch(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	info, ok := authentication.GetAuthenticationInfo(w, r)
	if !ok {
		return
	}

	var p PatchBody
	if err := web.DecodeBody(r, &p); err != nil {
		web.StatusBadRequest(w, err)
		return
	}

	notifications, err := c.repo.SetReadState(r.Context(), info.Cookie.UserID, id, p.Read)
	if err != nil {
		web.StatusInternalServerError(w, err)
		return
	}
	if notifications == nil {
		web.StatusNotFound(w, err)
		return
	}

	web.Write(w, notifications)
}
