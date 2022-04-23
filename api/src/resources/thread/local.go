package thread

import (
	"context"

	"github.com/Southclaws/storyden/api/src/resources/post"
	"github.com/Southclaws/storyden/api/src/resources/user"
	"github.com/google/uuid"
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
	categoryName string,
	tags []string,
) (*Thread, error) {
	id := post.PostID(uuid.New())

	l.m[id] = Thread{
		ID: id,

		Posts: []post.Post{
			{
				ID: id,
			},
		},
	}

	return nil, nil
}
