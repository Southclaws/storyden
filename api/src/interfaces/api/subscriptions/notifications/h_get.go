package notifications

import (
	"net/http"
	"time"

	"github.com/Southclaws/qstring"
	"github.com/Southclaws/storyden/api/src/infra/web"
	"github.com/Southclaws/storyden/api/src/services/authentication"
)

type GetParams struct {
	Read  bool      `qstring:"read" json:"read"`
	After time.Time `qstring:"after" json:"after"`
}

// @Summary  Get all notifications
// @Tags     notifications
// @Accept   json
// @Produce  json
// @Param    query  query     GetParams  true  "Query parameters for paging and read state"
// @Success  200    {object}  []notification.Notification
// @Failure  400    {object}  web.Error
// @Failure  500    {object}  web.Error
// @Failure  404    {object}  web.Error
// @Router   /subscriptions/notifications [get]
func (c *controller) get(w http.ResponseWriter, r *http.Request) {
	info, ok := authentication.GetAuthenticationInfo(w, r)
	if !ok {
		return
	}

	var p GetParams
	if err := qstring.Unmarshal(r.URL.Query(), &p); err != nil {
		web.StatusBadRequest(w, err)
		return
	}

	notifications, err := c.repo.GetNotifications(r.Context(), info.Cookie.UserID, p.Read, p.After)
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
