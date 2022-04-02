package subscriptions

import (
	"net/http"

	"github.com/Southclaws/storyden/api/src/infra/web"
	"github.com/Southclaws/storyden/api/src/resources/notification"
	"github.com/Southclaws/storyden/api/src/services/authentication"
)

type PostBody notification.Subscription

// @Summary  Create a subscription
// @Tags     subscriptions
// @Accept   json
// @Produce  json
// @Param    payload  body      PostBody  true  "Subscription data"
// @Success  200      {object}  notification.Subscription
// @Failure  400      {object}  web.Error
// @Failure  500      {object}  web.Error
// @Failure  404      {object}  web.Error
// @Router   /subscriptions [post]
func (c *controller) post(w http.ResponseWriter, r *http.Request) {
	info, ok := authentication.GetAuthenticationInfo(w, r)
	if !ok {
		return
	}

	var b PostBody
	if !web.ParseBody(w, r, &b) {
		return
	}

	sub, err := c.repo.Subscribe(r.Context(), info.Cookie.UserID, b.RefersType, b.RefersTo)
	if err != nil {
		web.StatusInternalServerError(w, err)
		return
	}

	web.Write(w, sub)
}
