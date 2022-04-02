package categories

import (
	"net/http"

	"github.com/Southclaws/storyden/api/src/infra/web"
	"github.com/Southclaws/storyden/api/src/resources/category"
)

type PatchBody category.Category

// @Summary  Update a category
// @Tags     categories
// @Accept   json
// @Produce  json
// @Param    body  body      PatchBody  true  "Updated category"
// @Success  200   {object}  category.Category
// @Failure  400   {object}  web.Error
// @Failure  500   {object}  web.Error
// @Failure  404   {object}  web.Error
// @Router   /categories/{id} [patch]
func (c *controller) patch(w http.ResponseWriter, r *http.Request) {
	var p PatchBody
	if err := web.DecodeBody(r, &p); err != nil {
		web.StatusBadRequest(w, err)
		return
	}

	updated, err := c.repo.UpdateCategory(r.Context(), p.ID, &p.Name, &p.Description, &p.Colour, &p.Sort, &p.Admin)
	if err != nil {
		web.StatusInternalServerError(w, err)
		return
	}

	if updated == nil {
		web.StatusNotFound(w, nil)
		return
	}

	web.Write(w, updated)
}
