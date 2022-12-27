package thread

import (
	"time"

	"4d63.com/optional"
	"github.com/Southclaws/dt"

	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/resources/category"
	"github.com/Southclaws/storyden/app/resources/post"
	"github.com/Southclaws/storyden/app/resources/react"
	"github.com/Southclaws/storyden/internal/ent"
	"github.com/Southclaws/storyden/internal/utils"
)

type Thread struct {
	ID        post.PostID
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt optional.Optional[time.Time]

	Title    string
	Slug     string
	Short    string
	Pinned   bool
	Author   post.Author
	Tags     []string
	Category category.Category
	Posts    []*post.Post
	Reacts   []*react.React
	Meta     map[string]any
}

const Name = "Thread"

func (*Thread) GetRole() string { return Name }

func (*Thread) GetResourceName() string { return Name }

func FromModel(m *ent.Post) *Thread {
	// Thread data structure will always contain one post: itself in post form.
	posts := []*post.Post{
		post.FromModel(m),
	}

	if p, err := m.Edges.PostsOrErr(); err == nil && len(p) > 0 {
		posts = append(posts, dt.Map(p, post.FromModel)...)
	}

	return &Thread{
		ID:        post.PostID(m.ID),
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		DeletedAt: optional.OfPtr(m.DeletedAt),

		Title:  m.Title,
		Slug:   m.Slug,
		Short:  m.Short,
		Pinned: m.Pinned,
		Author: post.Author{
			ID:     account.AccountID(m.Edges.Author.ID),
			Handle: m.Edges.Author.Handle,
			Name:   m.Edges.Author.Name,
		},
		Tags:     dt.Map(m.Edges.Tags, func(t *ent.Tag) string { return t.Name }),
		Category: utils.Deref(category.FromModel(m.Edges.Category)),
		Posts:    posts,
		Reacts:   dt.Map(m.Edges.Reacts, react.FromModel),
		Meta:     m.Metadata,
	}
}
