package categories

import (
	"net/http"

	"github.com/go-chi/chi"

	"github.com/Southclaws/storyden/api/src/infra/web"
)

type DeleteBody struct {
	MoveTo string `json:"moveTo" valid:"required"`
}

// @Summary  Delete a category
// @Tags     categories
// @Accept   json
// @Produce  json
// @Param    id      path      string  true  "Category ID"
// @Param    moveTo  body      string  true  "New category to move posts to"
// @Success  200     {object}  category.Category
// @Failure  400     {object}  web.Error
// @Failure  500     {object}  web.Error
// @Failure  404     {object}  web.Error
// @Router   /categories/{id} [delete]
func (c *controller) delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var p DeleteBody
	if err := web.DecodeBody(r, &p); err != nil {
		web.StatusBadRequest(w, err)
		return
	}

	deleted, err := c.repo.DeleteCategory(r.Context(), id, p.MoveTo)
	if err != nil {
		web.StatusInternalServerError(w, err)
		return
	}
	if deleted == nil {
		web.StatusNotFound(w, nil)
		return
	}

	web.Write(w, deleted)
}
