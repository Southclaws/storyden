package thread

import (
	"github.com/samber/lo"

	"github.com/Southclaws/storyden/api/src/infra/db/model"
	"github.com/Southclaws/storyden/api/src/resources/post"
)

type Thread struct {
	ID post.PostID

	Posts []post.Post
}

func FromModel(m *model.Post) *Thread {
	return &Thread{
		ID: post.PostID(m.ID),

		Posts: lo.Map(m.Edges.Posts, func(t *model.Post, i int) post.Post {
			p := post.FromModel(t)
			return *p
		}),
	}
}
