package tags

import (
	"net/http"

	"github.com/Southclaws/storyden/api/src/infra/web"
)

type TagsParams struct {
	Query string `qstring:"query"`
}

// @Summary  Get all tags
// @Tags     categories
// @Produce  json
// @Param    query  query     TagsParams  true  "Search parameters"
// @Success  200    {object}  []tag.Tag
// @Failure  500    {object}  web.Error
// @Router   /tags [get]
func (c *controller) get(w http.ResponseWriter, r *http.Request) {
	var p TagsParams
	if !web.ParseQuery(w, r, &p) {
		return
	}

	tags, err := c.tags.GetTags(r.Context(), p.Query)
	if err != nil {
		web.StatusInternalServerError(w, err)
		return
	}

	web.Write(w, tags)
}
