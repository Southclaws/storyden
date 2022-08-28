package seed

import (
	"context"
	"fmt"

	"github.com/Southclaws/storyden/internal/utils"
	"github.com/Southclaws/storyden/pkg/resources/post"
	"github.com/Southclaws/storyden/pkg/resources/thread"
)

var (
	Post_01 = thread.Thread{
		ID:       post.PostID(id("00000000000000000010")),
		Title:    "Hello world!",
		Author:   thread.AuthorRef{ID: Account_000.ID},
		Category: Category_01_General,
	}
	Post_02 = thread.Thread{
		ID:       post.PostID(id("00000000000000000020")),
		Title:    "Hello, 世界",
		Author:   thread.AuthorRef{ID: Account_001.ID},
		Category: Category_01_General,
	}
)

func threads(r thread.Repository) {
	ctx := context.Background()

	for _, t := range []thread.Thread{
		Post_01,
		Post_02,
	} {
		utils.Must(r.Create(ctx, t.Title, t.Short, t.Author.ID, t.Category.ID, t.Tags, thread.WithID(t.ID)))
	}

	fmt.Println("created seed threads")
}
