package posts

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"

	"github.com/Southclaws/storyden/api/src/infra/web"
	"github.com/Southclaws/storyden/api/src/resources/notification"
	"github.com/Southclaws/storyden/api/src/services/authentication"
)

type PostBody struct {
	Body    string `json:"body"   valid:"required,stringlength(1|65535)"`
	ReplyTo string `json:"replyTo"`
}

// @Summary  Post a reply to a thread
// @Tags     posts
// @Accepts  json
// @Produce  json
// @Param    thread_id  path      string    true  "Thread ID to post under"
// @Param    body       body      PostBody  true  "Post contents and optional reply-to ID"
// @Success  200        {object}  post.Post
// @Failure  400        {object}  web.Error
// @Failure  500        {object}  web.Error
// @Failure  404        {object}  web.Error
// @Router   /posts/{id} [post]
func (c *controller) post(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	info, ok := authentication.GetAuthenticationInfo(w, r)
	if !ok {
		return
	}

	var b PostBody
	if !web.ParseBody(w, r, &b) {
		web.StatusBadRequest(w, nil)
		return
	}

	post, err := c.repo.CreatePost(r.Context(), b.Body, info.Cookie.UserID, id, b.ReplyTo)
	if err != nil {
		web.StatusInternalServerError(w, err)
		return
	}
	if post == nil {
		web.StatusNotFound(w, nil)
		return
	}

	var link string
	if post.Slug != nil {
		link = c.publicAddress + "/discussion/" + *post.Slug
	}

	c.notifications.Notify(
		r.Context(),
		notification.NotificationTypeForumPostResponse,
		id,
		"Reply",
		fmt.Sprintf("%s: %s", post.Author.Name, post.Short),
		link)

	web.Write(w, post)
}
