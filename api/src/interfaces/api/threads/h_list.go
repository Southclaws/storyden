package threads

import (
	"net/http"
	"time"

	"github.com/Southclaws/storyden/api/src/infra/web"
	"github.com/Southclaws/storyden/api/src/services/authentication"
)

type ListParams struct {
	Tags     []string  `qstring:"tags"`
	Category string    `qstring:"category"`
	Query    string    `qstring:"query"`
	Before   time.Time `qstring:"before"`
	Sort     string    `qstring:"sort"`
	Offset   int       `qstring:"offset"`
	Max      int       `qstring:"max"`
	Posts    bool      `qstring:"posts"`
}

// @Summary  Get, search and filter threads
// @Tags     threads
// @Produce  json
// @Param    params  query     ListParams  false  "Search, filteirng and pagination parameters"
// @Success  200     {object}  []post.Post
// @Failure  400     {object}  web.Error
// @Failure  500     {object}  web.Error
// @Router   /threads [get]
func (c *controller) list(w http.ResponseWriter, r *http.Request) {
	var p ListParams
	if !web.ParseQuery(w, r, &p) {
		return
	}

	if p.Before.IsZero() {
		p.Before = time.Now()
	}
	if p.Sort == "" {
		p.Sort = "desc"
	}
	if p.Max == 0 {
		p.Max = 20
	} else if p.Max > 20 {
		p.Max = 20
	}

	// Admins get to see deleted posts
	isAdmin := authentication.IsRequestAdmin(r)

	posts, err := c.threads.GetThreads(r.Context(), p.Tags, p.Category, p.Query, p.Before, p.Sort, p.Offset, p.Max, p.Posts, isAdmin, isAdmin)
	if err != nil {
		web.StatusInternalServerError(w, err)
		return
	}

	web.Write(w, posts)
}
