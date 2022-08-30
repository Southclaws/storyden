package seed

import (
	"context"
	"fmt"

	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/thread"
	"github.com/Southclaws/storyden/internal/utils"
)

var (
	Post_01 = thread.Thread{
		ID:       post.PostID(id("00000000000000000010")),
		Title:    "Hello world!",
		Author:   thread.AuthorRef{ID: Account_000.ID},
		Category: Category_01_General,
		Posts: []*post.Post{
			{
				ID:         post.PostID(id("00000000000000001010")),
				Body:       "First reply",
				Short:      "First reply",
				RootPostID: post.PostID(id("00000000000000000010")),
				Author:     post.Author{ID: Account_003.ID},
			},
			{
				ID:         post.PostID(id("00000000000000002010")),
				Body:       "Second reply",
				Short:      "Second reply",
				RootPostID: post.PostID(id("00000000000000000010")),
				Author:     post.Author{ID: Account_004.ID},
			},
		},
	}
	Post_02 = thread.Thread{
		ID:       post.PostID(id("00000000000000000020")),
		Title:    "Hello, 世界",
		Author:   thread.AuthorRef{ID: Account_001.ID},
		Category: Category_01_General,
		Posts: []*post.Post{
			{
				ID:         post.PostID(id("00000000000000001020")),
				Body:       "First reply of other thread",
				Short:      "First reply of other thread",
				RootPostID: post.PostID(id("00000000000000000020")),
				Author:     post.Author{ID: Account_005.ID},
			},
			{
				ID:         post.PostID(id("00000000000000002020")),
				Body:       "Second reply of other thread",
				Short:      "Second reply of other thread",
				RootPostID: post.PostID(id("00000000000000000020")),
				Author:     post.Author{ID: Account_006.ID},
			},
		},
	}
)

func threads(tr thread.Repository, pr post.Repository) {
	ctx := context.Background()

	for _, t := range []thread.Thread{
		Post_01,
		Post_02,
	} {
		th := utils.Must(tr.Create(ctx, t.Title, t.Short, t.Author.ID, t.Category.ID, t.Tags, thread.WithID(t.ID)))

		for _, p := range t.Posts {
			utils.Must(pr.Create(ctx, p.Body, p.Author.ID, th.ID, nil))
		}
	}

	fmt.Println("created seed threads")
}
