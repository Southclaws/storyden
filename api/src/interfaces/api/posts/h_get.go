package posts

import (
	"errors"
	"html"
	"net/http"

	"github.com/Southclaws/storyden/api/src/infra/web"
	"github.com/Southclaws/storyden/api/src/services/authentication"
	"github.com/go-chi/chi"
)

type GetParams struct {
	Max  int `qstring:"max"  valid:"range(1|50)"`
	Skip int `qstring:"skip"`
}

// @Summary  Get a post and all its replies
// @Tags     posts
// @Produce  json
// @Param    slug    path      string     true   "Root post slug"
// @Param    params  query     GetParams  false  "Page params"
// @Success  200     {object}  []post.Post
// @Failure  400     {object}  web.Error
// @Failure  500     {object}  web.Error
// @Failure  404     {object}  web.Error
// @Router   /posts/{id} [get]
func (c *controller) get(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	var p GetParams
	if !web.ParseQuery(w, r, &p) {
		return
	}

	if p.Max == 0 {
		p.Max = 1000
	}

	isAdmin := authentication.IsRequestAdmin(r)

	posts, err := c.repo.GetPosts(r.Context(), slug, p.Max, p.Skip, isAdmin, isAdmin)
	if err != nil {
		web.StatusInternalServerError(w, err)
		return
	}
	if posts == nil {
		web.StatusNotFound(w, web.WithDescription(errors.New("not found"), "No posts were found with that ID"))
		return
	}

	// TODO: move this post body html escape elsewhere...
	for i, p := range posts {
		posts[i].Body = html.EscapeString(p.Body)
	}

	web.Write(w, posts)
}
