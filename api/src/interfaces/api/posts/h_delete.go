package posts

import (
	"net/http"

	"github.com/Southclaws/storyden/api/src/infra/web"
	"github.com/Southclaws/storyden/api/src/services/authentication"
	"github.com/go-chi/chi"
)

// @Summary  Soft-delete a single post
// @Tags     posts
// @Produce  json
// @Param    id   path      string  true  "Post ID"
// @Success  200  {object}  post.Post
// @Failure  500  {object}  web.Error
// @Failure  404  {object}  web.Error
// @Router   /posts/{id} [delete]
func (c *controller) delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	info, ok := authentication.GetAuthenticationInfo(w, r)
	if !ok {
		return
	}

	post, err := c.repo.DeletePost(r.Context(), info.Cookie.UserID, id, info.Cookie.Admin)
	if err != nil {
		web.StatusInternalServerError(w, err)
	}
	if post == nil {
		web.StatusNotFound(w, nil)
		return
	}

	web.Write(w, post)
}
