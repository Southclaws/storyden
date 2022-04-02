package threads

import (
	"errors"
	"net/http"

	"github.com/Southclaws/storyden/api/src/infra/web"
	"github.com/Southclaws/storyden/api/src/resources/notification"
	"github.com/Southclaws/storyden/api/src/resources/post"
	"github.com/Southclaws/storyden/api/src/services/authentication"
)

type PostBody struct {
	Title    string   `json:"title"      valid:"required,stringlength(1|64)"`
	Body     string   `json:"body"       valid:"required,stringlength(1|65535)"`
	Tags     []string `json:"tags"       valid:"required"`
	Category string   `json:"category"   valid:"required"`
}

// @Summary  Create a new thread by creating a root post under a category.
// @Tags     threads
// @Accepts  json
// @Produce  json
// @Param    body  body      PostBody  true  "Post contents with title etc."
// @Success  200   {object}  post.Post
// @Failure  400   {object}  web.Error
// @Failure  500   {object}  web.Error
// @Failure  404   {object}  web.Error
// @Router   /threads [post]
func (c *controller) post(w http.ResponseWriter, r *http.Request) {
	info, ok := authentication.GetAuthenticationInfo(w, r)
	if !ok {
		return
	}

	var b PostBody
	if !web.ParseBody(w, r, &b) {
		return
	}

	p, err := c.threads.CreateThread(r.Context(), b.Title, b.Body, info.Cookie.UserID, b.Category, b.Tags)
	if err != nil {
		if errors.Is(err, post.ErrTagNameTooLong) {
			web.StatusBadRequest(w, web.WithSuggestion(err,
				"The name of one of the tags is too long",
				"The character limit for a tag is 24, delete any tags longer than this and retry."))
		} else {
			web.StatusInternalServerError(w, err)
		}
		return
	}

	c.notifications.Subscribe(r.Context(), info.Cookie.UserID, notification.NotificationTypeForumPostResponse, p.ID)

	web.Write(w, p)
}
