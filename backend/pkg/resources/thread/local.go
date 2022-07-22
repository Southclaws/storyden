package thread

import (
	"context"

	"github.com/google/uuid"

	"github.com/Southclaws/storyden/backend/pkg/resources/category"
	"github.com/Southclaws/storyden/backend/pkg/resources/post"
	"github.com/Southclaws/storyden/backend/pkg/resources/user"
)

type local struct {
	m map[post.PostID]Thread
}

func NewLocal() Repository {
	return &local{m: map[post.PostID]Thread{}}
}

func (l *local) CreateThread(
	ctx context.Context,
	title string,
	body string,
	authorID user.UserID,
	categoryID category.CategoryID,
	tags []string,
) (*Thread, error) {
	id := post.PostID(uuid.New())

	t := Thread{
		ID: id,

		Posts: []post.Post{
			{
				ID: id,
			},
		},
	}

	l.m[id] = t

	return &t, nil
}
