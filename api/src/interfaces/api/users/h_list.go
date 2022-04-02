package users

import (
	"net/http"

	"github.com/Southclaws/storyden/api/src/infra/web"
)

type queryParams struct {
	Sort string `qstring:"sort"`
	Max  int    `qstring:"max"`
	Skip int    `qstring:"skip"`
}

func (c *controller) list(w http.ResponseWriter, r *http.Request) {
	var p queryParams
	if !web.ParseQuery(w, r, &p) {
		return
	}

	if p.Sort == "" {
		p.Sort = "desc"
	}
	if p.Max == 0 {
		p.Max = 50
	} else if p.Max > 100 {
		p.Max = 100
	}

	users, err := c.repo.GetUsers(r.Context(), p.Sort, p.Max, p.Skip, true)
	if err != nil {
		web.StatusInternalServerError(w, err)
		return
	}

	web.Write(w, users)
}
