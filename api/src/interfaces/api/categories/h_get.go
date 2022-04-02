package categories

import (
	"net/http"

	"github.com/Southclaws/storyden/api/src/infra/web"
	"github.com/Southclaws/storyden/api/src/services/authentication"
)

// @Summary  Get all categories
// @Tags     categories
// @Produce  json
// @Success  200  {object}  []category.Category
// @Failure  500  {object}  web.Error
// @Router   /categories [get]
func (c *controller) get(w http.ResponseWriter, r *http.Request) {
	categories, err := c.repo.GetCategories(r.Context(), authentication.IsRequestAdmin(r))
	if err != nil {
		web.StatusInternalServerError(w, err)
		return
	}

	web.Write(w, categories)
}
