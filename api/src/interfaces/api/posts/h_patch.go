package posts

import (
	"net/http"

	"github.com/Southclaws/storyden/api/src/infra/web"
	"github.com/Southclaws/storyden/api/src/services/authentication"
	"github.com/go-chi/chi"
)

type PatchBody struct {
	Title string `json:"title" valid:"stringlength(1|64)"`
	Body  string `json:"body"  valid:"stringlength(1|65535)"`
}

// @Summary  Update a post
// @Tags     posts
// @Accept   json
// @Produce  json
// @Param    id    path      string     true  "Post ID"
// @Param    body  body      PatchBody  true  "Updated post contents"
// @Success  200   {object}  post.Post
// @Failure  400   {object}  web.Error
// @Failure  500   {object}  web.Error
// @Failure  404   {object}  web.Error
// @Router   /posts/{id} [patch]
func (c *controller) patch(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	info, ok := authentication.GetAuthenticationInfo(w, r)
	if !ok {
		return
	}

	var b PatchBody
	if !web.ParseBody(w, r, &b) {
		return
	}

	post, err := c.repo.EditPost(r.Context(), info.Cookie.UserID, id, &b.Title, &b.Body)
	if err != nil {
		web.StatusInternalServerError(w, err)
		return
	}
	if post == nil {
		web.StatusNotFound(w, nil)
		return
	}

	web.Write(w, post)
}
