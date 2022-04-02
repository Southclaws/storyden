package subscriptions

import (
	"net/http"

	"github.com/Southclaws/storyden/api/src/infra/web"
	"github.com/Southclaws/storyden/api/src/services/authentication"
)

// @Summary  Get all subscriptions
// @Tags     subscriptions
// @Accept   json
// @Produce  json
// @Success  200  {object}  notification.Subscription
// @Failure  500  {object}  web.Error
// @Failure  404  {object}  web.Error
// @Router   /subscriptions [get]
func (c *controller) get(w http.ResponseWriter, r *http.Request) {
	info, ok := authentication.GetAuthenticationInfo(w, r)
	if !ok {
		return
	}

	notifications, err := c.repo.GetSubscriptionsForUser(r.Context(), info.Cookie.UserID)
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
