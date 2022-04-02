package categories

import (
	"net/http"

	"github.com/Southclaws/storyden/api/src/infra/web"
	"github.com/Southclaws/storyden/api/src/resources/category"
)

type PatchAllBody []category.Category

// @Summary  Bulk update categories
// @Tags     categories
// @Accept   json
// @Produce  json
// @Param    body  body      PatchAllBody  true  "Updated categories"
// @Success  200   {object}  []category.Category
// @Failure  400   {object}  web.Error
// @Failure  500   {object}  web.Error
// @Router   /categories [patch]
func (c *controller) patchAll(w http.ResponseWriter, r *http.Request) {
	var p PatchAllBody
	if err := web.DecodeBody(r, &p); err != nil {
		web.StatusBadRequest(w, err)
		return
	}

	result := []category.Category{}
	for _, category := range p {
		newCategory, err := c.repo.UpdateCategory(r.Context(), category.ID, &category.Name, &category.Description, &category.Colour, &category.Sort, &category.Admin)
		if err != nil {
			web.StatusInternalServerError(w, err)
			return
		}
		result = append(result, *newCategory)
	}

	web.Write(w, result)
}
