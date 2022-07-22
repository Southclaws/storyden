package thread

import (
	"time"

	"4d63.com/optional"
	"github.com/samber/lo"

	"github.com/Southclaws/storyden/backend/internal/infrastructure/db/model"
	"github.com/Southclaws/storyden/backend/internal/utils"
	"github.com/Southclaws/storyden/backend/pkg/resources/category"
	"github.com/Southclaws/storyden/backend/pkg/resources/post"
)

type Thread struct {
	ID post.PostID

	Title     string
	Slug      string
	Pinned    bool
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt optional.Optional[time.Time]
	Category  category.Category
	Posts     []post.Post
}

func FromModel(m *model.Post) *Thread {
	return &Thread{
		ID: post.PostID(m.ID),

		Title:     m.Title,
		Slug:      m.Slug,
		Pinned:    m.Pinned,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		DeletedAt: optional.Of(m.DeletedAt),
		Category:  utils.Deref(category.FromModel(m.Edges.Category), 0),

		Posts: lo.Map(m.Edges.Posts, func(t *model.Post, i int) post.Post {
			p := post.FromModel(t)
			return *p
		}),
	}
}
